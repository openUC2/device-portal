package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PlanktoScope/machine-name/pkg/haikunator"
	"github.com/PlanktoScope/machine-name/pkg/wordlists"
	"github.com/pkg/errors"
	"github.com/sargassum-world/godest/env"
)

type MachineNameConfig struct {
	Lang         string
	SerialNumber uint32
	MachineName  string
}

func getMachineNameConfig() (c MachineNameConfig, err error) {
	c.Lang, err = getLangConfig()
	if err != nil {
		return MachineNameConfig{}, errors.Wrap(err, "couldn't determine lang")
	}
	c.SerialNumber, err = getSerialNumberConfig()
	validSerialNumber := true
	if err != nil {
		// FIXME: this should be a logger warning, not a printf statement
		fmt.Printf(
			"Warning: couldn't determine serial number, so the machine name will be 'invalid-name': %s\n",
			err,
		)
		validSerialNumber = false
	}
	if !validSerialNumber {
		c.MachineName = "invalid-name"
	} else {
		c.MachineName, err = generateMachineName(c.Lang, c.SerialNumber)
		if err != nil {
			return MachineNameConfig{}, errors.Wrap(err, "couldn't determine machine name")
		}
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

func getSerialNumberConfig() (sn uint32, err error) {
	rawSN := env.GetString("SERIAL_NUMBER", "")
	if rawSN == "" {
		// This is a file path specific to the Raspberry Pi
		const defaultSNFilePath = "/sys/firmware/devicetree/base/serial-number"
		snFilePath := env.GetString("SERIAL_NUMBER_FILE", defaultSNFilePath)
		rawFile, err := os.ReadFile(filepath.Clean(snFilePath))
		if err != nil {
			return 0, errors.Wrapf(err, "couldn't read serial number from file '%s'", snFilePath)
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
