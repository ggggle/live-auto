package src

import (
	"encoding/xml"
	"fmt"
	"github.com/juju/errors"
	"io/ioutil"
)

var G_Config *Config

type Config struct {
	PieceSize int64  `xml:"piece_size"`
	Proxy     string `xml:"proxy"`
}

func ReadConfig(configFilePath string) error {
	content_b, err := ioutil.ReadFile(configFilePath)
	if nil != err {
		return errors.New(fmt.Sprintf("open[%s]file error, %s", configFilePath, err.Error()))
	}
	G_Config = new(Config)
	err = xml.Unmarshal(content_b, G_Config)
	if nil != err {
		return err
	}
	G_Config.PieceSize *= 1024 * 1024
	return nil
}
