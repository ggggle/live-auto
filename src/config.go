package src

import (
	"encoding/xml"
	"fmt"
	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var G_Config *Config

type Config struct {
	XMLName       xml.Name    `xml:"config"`
	PieceSize     int64       `xml:"piece_size"`
	Proxy         string      `xml:"proxy"`
	CheckInterval int         `xml:"check_interval"`
	OutPutDirPath string      `xml:"out_put_dir_path"`
	RecordRules   recordRules `xml:"record_rules"`
}

type recordRules struct {
	Rules []Rule `xml:"rule"`
}

type Rule struct {
	Urls            string          `xml:"urls"`
	EnableUploaders enableUploaders `xml:"enable_uploaders"`
	Loop            bool            `xml:"loop"`
	CheckInterval   int             `xml:"check_interval"`
}

type enableUploaders struct {
	AutoRemove bool   `xml:"auto_remove,attr"`
	Uploaders  string `xml:",chardata"`
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

// 加载xml配置文件中的
func LoadCfgRecorder() {
	for _, rule := range G_Config.RecordRules.Rules {
		url_split := strings.Split(rule.Urls, ";")
		if 0 == len(url_split) {
			continue
		}
		rcd_cfg := RecordConfig{
			Loop:            rule.Loop,
			AutoRemove:      rule.EnableUploaders.AutoRemove,
			EnableUploaders: make([]UploaderType, 0),
		}
		uploader_split := strings.Split(rule.EnableUploaders.Uploaders, ";")
		for _, uploader := range uploader_split {
			if tmp_num, err := strconv.Atoi(uploader); nil == err && tmp_num < int(MAX_SUPPORT_UPLOADER) {
				rcd_cfg.EnableUploaders = append(rcd_cfg.EnableUploaders, UploaderType(tmp_num))
			}
		}
		for _, v := range url_split {
			go func(url string) {
				recorder, err := NewRecorder(url, rcd_cfg)
				if nil != err {
					Logger.WithFields(logrus.Fields{
						ERROR_CONTENT_DEF: err.Error(),
						"url":             url,
					}).Warn("初始化recorder出错")
					return
				}
				if err = GetIRecorderMngr().AddRecorder(recorder); nil != err {
					Logger.WithFields(logrus.Fields{
						ERROR_CONTENT_DEF: err.Error(),
					}).Warn("初始化添加recorder出错")
					return
				}
				recorder.Start()
			}(v)
		}
	}
}
