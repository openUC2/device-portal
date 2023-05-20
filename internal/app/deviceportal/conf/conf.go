// Package conf supports environment variable-based application configuration
package conf

import (
	"github.com/pkg/errors"
)

type Config struct {
	HTTP        HTTPConfig
	MachineName MachineNameConfig
}

func GetConfig() (c Config, err error) {
	c.HTTP, err = getHTTPConfig()
	if err != nil {
		return Config{}, errors.Wrap(err, "couldn't make http config")
	}
	c.MachineName, err = getMachineNameConfig()
	if err != nil {
		return Config{}, errors.Wrap(err, "couldn't make machine name config")
	}
	return c, nil
}
