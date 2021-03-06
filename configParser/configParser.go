package configParser

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Defaults
var (
	Config *ConfigFile = &ConfigFile{}
)

type ConfigFile struct {
	VoiceName         string `yaml:"voiceName"`
	LanguageCode      string `yaml:"languageCode"`
	Interface         string `yaml:"interface"`
	ArrivedSuffix     string `yaml:"arrivedSuffix"`
	LeftSuffix        string `yaml:"leftSuffix"`
	ArrivedPrefix     string `yaml:"arrivedPrefix"`
	LeftPrefix        string `yaml:"leftPrefix"`
	RegisteredDevices []struct {
		Name         string `yaml:"name"`
		Hostname     string `yaml:"hostname"`
		DefaultState string `yaml:"defaultState"`
	} `yaml:"registeredDevices"`
	ScanOnly  bool   `yaml:"scanOnly"`
	DNSServer string `yaml:"dnsServer"`
}

func ParseConfig() {
	f, err := ioutil.ReadFile(os.Getenv("CONFIG_PATH"))

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, Config)
	if err != nil {
		panic(err)
	}
}

// InterfaceToPrefix - Takes the `Interface` and gives the first 3 octets
func InterfaceToPrefix() string {
	split := strings.Split(Config.Interface, ".")
	prefix := strings.Join(split[:3], ".")
	return prefix
}
