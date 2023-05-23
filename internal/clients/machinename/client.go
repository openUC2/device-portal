// Package machinename determines the machine name of the underlying compute platform based on
// environment variables or a file (defaulting to the serial number file in the Raspberry Pi's
// device tree)
package machinename

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PlanktoScope/machine-name/pkg/haikunator"
	"github.com/PlanktoScope/machine-name/pkg/wordlists"
	"github.com/pkg/errors"
	"github.com/sargassum-world/godest"
	"github.com/sargassum-world/godest/clientcache"
	"github.com/sargassum-world/godest/env"
)

type Client struct {
	Config Config
	Logger godest.Logger
	Cache  *Cache
}

func NewClient(c Config, cache clientcache.Cache, l godest.Logger) *Client {
	return &Client{
		Config: c,
		Logger: l,
		Cache: &Cache{
			Cache: cache,
		},
	}
}

func (c *Client) GetName() (string, error) {
	if name, cacheHit := c.getNameFromCache(); cacheHit {
		c.Logger.Debugf("machine name was loaded from cache as %s", name)
		return name, nil // empty name indicates nonexistent name
	}
	return c.getNameFromSystem()
}

func (c *Client) getNameFromCache() (string, bool) {
	name, cacheHit, err := c.Cache.GetName()
	if err != nil {
		// Log the error but return as a cache miss so we can manually generate the name
		c.Logger.Error(errors.Wrap(err, "couldn't get the cache entry for the machine name"))
		return "", false // treat an unparseable cache entry like a cache miss
	}
	return name, cacheHit
}

func (c *Client) getNameFromSystem() (string, error) {
	sn, err := c.getSerialNumber()
	if err != nil {
		return "", errors.Wrap(err, "couldn't determine serial number")
	}
	name, err := generateMachineName(c.Config.Lang, sn)
	if err != nil {
		return "", errors.Wrapf(
			err, "couldn't generate machine name from serial number %d in language %s",
			sn, c.Config.Lang,
		)
	}
	if err := c.Cache.SetName(name, c.Config.CacheCost); err != nil {
		return "", errors.Wrapf(err, "couldn't cache machine name %s", name)
	}
	if name != "" {
		c.Logger.Infof("machine name was determined to be %s", name)
	}
	return name, nil
}

func (c *Client) getSerialNumber() (sn uint32, err error) {
	rawSN := env.GetString(envPrefix+"SN", "")
	if rawSN == "" {
		rawFile, err := os.ReadFile(filepath.Clean(c.Config.SNFile))
		if err != nil {
			return 0, errors.Wrapf(err, "couldn't read serial number from file '%s'", c.Config.SNFile)
		}
		rawSN = strings.TrimSpace(string(rawFile))
	}

	return parseSerialNumber(rawSN)
}

func parseSerialNumber(raw string) (uint32, error) {
	const base = 16
	const parsedWidth = 32
	parsed64, err := strconv.ParseUint(strings.TrimPrefix(raw, "0x"), base, parsedWidth)
	return uint32(parsed64), errors.Wrapf(err, "couldn't parse serial number '%s'", raw)
}

func generateMachineName(lang string, sn uint32) (name string, err error) {
	first, second, err := wordlists.Load(wordlists.FS, lang)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't load naming wordlists for language '%s", lang)
	}
	return haikunator.SelectName(sn, first, second), nil
}
