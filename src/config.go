package src

import (
	"encoding/xml"
	"fmt"
	"github.com/juju/errors"
	"io/ioutil"
	"os"
)

var G_Config *Config

type Config struct {
	PieceSize     int64  `xml:"piece_size"`
	Proxy         string `xml:"proxy"`
	CheckInterval int    `xml:"check_interval"`
	OutPutDirPath string `xml:"out_put_dir_path"`
}

func ReadConfig(configFilePath string) error {
	content_b, err := ioutil.ReadFile(configFilePath)
	if nil != err {
		return errors.New(fmt.Sprintf("open[%s]file error, %s", configFilePath, err.Error()))
	}
	// 默认值
	G_Config = &Config{
		PieceSize:     2048,
		Proxy:         "",
		CheckInterval: 15,
		OutPutDirPath: fmt.Sprintf(".%c", os.PathSeparator),
	}
	err = xml.Unmarshal(content_b, G_Config)
	if nil != err {
		return err
	}
	G_Config.PieceSize *= 1024 * 1024
	return nil
}
