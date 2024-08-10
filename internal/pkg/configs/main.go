package configs

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	ServiceName = "d"
	Version     = "0.0.0" // the version is set via build flags
	Hostname    = ""      // set during initialization
)

// Init reads the given files and initiates the configuration (connections
// credentials), settings (a use cases parameters) and handlers.
func Init(fileName string) (*Root, error) {
	var err error
	if Hostname, err = os.Hostname(); err != nil {
		return nil, errors.Wrap(err, "hostname")
	}

	f, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	root := Root{}
	if err := yaml.Unmarshal(f, &root); err != nil {
		return nil, errors.Wrap(err, "parse")
	}

	if err := root.initiate(); err != nil {
		return nil, err
	}

	return &root, nil
}
