package machinename

import (
	"fmt"

	"github.com/PlanktoScope/machine-name/pkg/wordlists"
	"github.com/pkg/errors"
	"github.com/sargassum-world/godest/env"
)

const envPrefix = "MACHINENAME_"

type Config struct {
	Lang   string
	SNFile string

	CacheCost float32
}

func GetConfig() (c Config, err error) {
	c.Lang, err = getLangConfig()
	if err != nil {
		return Config{}, errors.Wrap(err, "couldn't make lang config")
	}

	// This is a file path specific to the Raspberry Pi
	const defaultSNFilePath = "/sys/firmware/devicetree/base/serial-number"
	c.SNFile = env.GetString(envPrefix+"SNFILE", defaultSNFilePath)

	const defaultCacheCost = 1.0
	c.CacheCost, err = env.GetFloat32(envPrefix+"CACHE_COST", defaultCacheCost)
	if err != nil {
		return Config{}, errors.Wrap(err, "couldn't make cache cost config")
	}
	return c, nil
}

func getLangConfig() (lang string, err error) {
	const defaultLang = "en_US.UTF-8"
	lang = env.GetString("LANG", defaultLang)

	allowed, err := wordlists.ListLanguages(wordlists.FS)
	if err != nil {
		return "", errors.Wrap(err, "couldn't determine the list of supported languages")
	}
	if _, ok := allowed[lang]; !ok {
		// FIXME: log this warning properly
		fmt.Printf("Warning: language '%s' is not supported, reverting to '%s'\n", lang, defaultLang)
		lang = defaultLang
	}
	return lang, nil
}
