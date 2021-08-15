package common

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	WebServer struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"web-server"`
	Features struct {
		EnableRedirector bool `yaml:"enable-redirector"`
		EnableImages     bool `yaml:"enable-images"`
		EnableText       bool `yaml:"enable-text"`
		EnableFiles      bool `yaml:"enable-files"`
		API              struct {
			EnableAPI         bool `yaml:"manage-redirectors"`
			ManageRedirectors bool `yaml:"manage-redirectors"`
			ManageImages      bool `yaml:"manage-images"`
			ManageText        bool `yaml:"manage-text"`
			ManageFiles       bool `yaml:"manage-files"`
		} `yaml:"api"`
	} `yaml:"features"`
	MongoDB struct {
		URI        string `yaml:"uri"`
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		UseSRV     bool   `yaml:"use-srv"`
		User       string `yaml:"user"`
		Pass       string `yaml:"pass"`
		AuthSource string `yaml:"auth-source"`
		UseAuth    bool   `yaml:"use-auth"`
		DB         string `yaml:"db"`
	} `yaml:"mongodb"`
}

func LoadConfiguration(fileLocation string) (conf Configuration, err error) {
	file, err := os.Open(fileLocation)
	if err != nil {
		err = fmt.Errorf("encountered error while opening configuration file\n%v", err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("encountered error while reading configuration file\n%v", err)
		return
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		err = fmt.Errorf("encountered error while unmarshaling configuration data\n%v", err)
		return
	}

	return
}

func (conf *Configuration) GetURI() string {
	mdb := conf.MongoDB

	if mdb.URI != "" {
		return conf.MongoDB.URI
	}

	if mdb.UseSRV {
		if mdb.UseAuth {
			return fmt.Sprintf("mongodb+srv://%s:%s@%s:%d/?authSource=%s", mdb.User, mdb.Pass, mdb.Host, mdb.Port, mdb.AuthSource)
		}
		return fmt.Sprintf("mongodb+srv://%s:%d/?authSource=%s", mdb.Host, mdb.Port, mdb.AuthSource)
	}

	if mdb.UseAuth {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=%s", mdb.User, mdb.Pass, mdb.Host, mdb.Port, mdb.AuthSource)
	}

	return fmt.Sprintf("mongodb://%s:%d/?authSource=%s", mdb.Host, mdb.Port, mdb.AuthSource)
}
