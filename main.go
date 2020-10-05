package main

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"github.com/sirupsen/logrus"
	"live-auto/cfg"
	"live-auto/src"
	"strconv"
	"strings"
)

// 加载xml配置文件中的
func LoadCfgRecorder() {
	for _, rule := range cfg.G_Config.RecordRules.Rules {
		url_split := strings.Split(rule.Urls, ";")
		if 0 == len(url_split) {
			continue
		}
		rcd_cfg := src.RecordConfig{
			Loop:            rule.Loop,
			AutoRemove:      rule.EnableUploaders.AutoRemove,
			EnableUploaders: make([]src.UploaderType, 0),
		}
		uploader_split := strings.Split(rule.EnableUploaders.Uploaders, ";")
		for _, uploader := range uploader_split {
			if tmp_num, err := strconv.Atoi(uploader); nil == err && tmp_num < int(src.MAX_SUPPORT_UPLOADER) {
				rcd_cfg.EnableUploaders = append(rcd_cfg.EnableUploaders, src.UploaderType(tmp_num))
			}
		}
		for _, v := range url_split {
			go func(url string) {
				recorder, err := src.NewRecorder(url, rcd_cfg)
				if nil != err {
					src.Logger.WithFields(logrus.Fields{
						src.ERROR_CONTENT_DEF: err.Error(),
						"url":                 url,
					}).Warn("初始化recorder出错")
					return
				}
				if err = src.GetIRecorderMngr().AddRecorder(recorder); nil != err {
					src.Logger.WithFields(logrus.Fields{
						src.ERROR_CONTENT_DEF: err.Error(),
					}).Warn("初始化添加recorder出错")
					return
				}
				recorder.Start()
			}(v)
		}
	}
}

func main() {
	LoadCfgRecorder()
	cmd()
}

func cmd() {
	for {
		fmt.Println("1.add")
		fmt.Println("2.remove")
		fmt.Println("3.list")
		fmt.Print("input:")
		input := 0
		_, err := fmt.Scanf("%d\n", &input)
		if nil != err {
			continue
		}
		switch input {
		case 1:
			add()
		case 2:
			remove()
		case 3:
			list()
		default:
			continue
		}
	}
}

func add() {
	live_url := ""
	fmt.Print("input url:")
	_, err := fmt.Scanf("%s\n", &live_url)
	if nil != err {
		return
	}
	go func() {
		recorder, err := src.NewRecorder(live_url, src.RecordConfig{
			Loop:            true,
			AutoRemove:      true,
			EnableUploaders: []src.UploaderType{src.GDRIVE},
		})
		if nil != err {
			fmt.Println(err)
			return
		}
		err = src.GetIRecorderMngr().AddRecorder(recorder)
		if nil != err {
			fmt.Println(err)
			return
		}
		recorder.Start()
	}()
}

func remove() {
	live_id := ""
	fmt.Print("input live_id:")
	_, err := fmt.Scanf("%s\n", &live_id)
	if nil != err {
		return
	}
	err = src.GetIRecorderMngr().RemoveRecorder(api.LiveId(live_id))
	if nil != err {
		fmt.Println(err.Error())
	}
}

func list() {
	list := src.GetIRecorderMngr().GetAllRecorder()
	for _, rcd := range list {
		fmt.Printf("%s %s\n", rcd.Live.GetLiveId(), rcd.Live.GetRawUrl())
	}
}
