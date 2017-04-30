package errortree_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/speijnik/go-errortree"
)

// ErrOptionMissing indicates that a configuration option is missing
var ErrOptionMissing = errors.New("Configuration option is missing")

type NetworkConfiguration struct {
	ListenAddress string
	MaxClients    uint
}

// Validate validates the network configuration
func (c *NetworkConfiguration) Validate() (err error) {
	if c.ListenAddress == "" {
		err = errortree.Add(err, "ListenAddress", ErrOptionMissing)
	} else if _, _, splitErr := net.SplitHostPort(c.ListenAddress); splitErr != nil {
		err = errortree.Add(err, "ListenAddress", errors.New("Must be in host:port format"))
	}

	if c.MaxClients < 1 {
		err = errortree.Add(err, "MaxClients", errors.New("Must be at least 1"))
	}

	return
}

type StorageConfiguration struct {
	DataDirectory string
}

func (c *StorageConfiguration) Validate() (err error) {
	if c.DataDirectory == "" {
		err = errortree.Add(err, "DataDirectory", ErrOptionMissing)
	} else if fileInfo, statErr := os.Stat(c.DataDirectory); statErr != nil {
		err = errortree.Add(err, "DataDirectory", errors.New("Directory does not exist"))
	} else if !fileInfo.IsDir() {
		err = errortree.Add(err, "DataDirectory", errors.New("Not a directory"))
	}
	return
}

// Configuration represents a configuration struct
// which may be filled from a file.
//
// It provides a Validate function which ensures that the configuration
// is correct and can be used.
type Configuration struct {
	// Network configuration
	Network NetworkConfiguration

	// Storage configuration
	Storage StorageConfiguration
}

func (c *Configuration) Validate() (err error) {
	err = errortree.Add(err, "Network", c.Network.Validate())
	err = errortree.Add(err, "Storage", c.Storage.Validate())

	return
}

func Example() {
	// Initialize an empty configuration. Validating this configuration should give us three errors:
	c := Configuration{}

	fmt.Println("[0] " + c.Validate().Error())

	// Now set some values and validate again
	c.Network.ListenAddress = "[[:8080"
	c.Network.MaxClients = 1
	c.Storage.DataDirectory = "/non-existing"

	fmt.Println("[1] " + c.Validate().Error())

	// Set the data directory to a temporary file (!) and validate again
	f, err := ioutil.TempFile("", "go-errortree-test-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	c.Network.ListenAddress = "0.0.0.0:8080"
	c.Storage.DataDirectory = f.Name()

	fmt.Println("[2] " + c.Validate().Error())

	// Fix everything and run again
	tempDir, err := ioutil.TempDir("", "go-errortree-test-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempDir)

	c.Storage.DataDirectory = tempDir

	if err := c.Validate(); err == nil {
		fmt.Println("[3] Config OK")
	}

	// Output: [0] 3 errors occurred:
	//
	// * Network:ListenAddress: Configuration option is missing
	// * Network:MaxClients: Must be at least 1
	// * Storage:DataDirectory: Configuration option is missing
	// [1] 2 errors occurred:
	//
	// * Network:ListenAddress: Must be in host:port format
	// * Storage:DataDirectory: Directory does not exist
	// [2] 1 error occurred:
	//
	// * Storage:DataDirectory: Not a directory
	// [3] Config OK
}
