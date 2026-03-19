// Package conf supports environment variable-based application configuration
package conf

import (
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"
)

type Config struct {
	Cache ristretto.Config
	HTTP  HTTPConfig
}

type HTTPConfig struct {
	Port      int
	BasePath  string
	GzipLevel int
}

func GetConfig() (c Config, err error) {
	c.Cache, err = getCacheConfig()
	if err != nil {
		return Config{}, errors.Wrap(err, "couldn't make cache config")
	}

	return c, nil
}
