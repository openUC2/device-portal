package machinename

import (
	"github.com/sargassum-world/godest/clientcache"
)

type Cache struct {
	Cache clientcache.Cache
}

// /machinename/name

func keyName() string {
	return "/machinename/name"
}

func (c *Cache) SetName(name string, costWeight float32) error {
	key := keyName()
	return c.Cache.SetEntry(key, name, costWeight, -1)
}

func (c *Cache) UnsetName() {
	key := keyName()
	c.Cache.UnsetEntry(key)
}

func (c *Cache) GetName() (string, bool, error) {
	key := keyName()
	var value string
	keyExists, valueExists, err := c.Cache.GetEntry(key, &value)
	if !keyExists || !valueExists || err != nil {
		return "", keyExists, err
	}

	return value, true, nil
}
