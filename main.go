package main

import (
	"fmt"
	"live-auto/src"
	"time"
)

func main() {

	fmt.Println(src.G_Config)
	fmt.Println(src.GetDirID("liveAuto"))
	time.Sleep(time.Second)
	return
	// defer src.Logger.Exit(0)
	recorder, err := src.NewRecorder("https://www.douyu.com/196", src.RecordConfig{
		Loop:            true,
		AutoRemove:      true,
		EnableUploaders: []src.UploaderType{src.GDRIVE},
	})
	if nil != err {
		fmt.Println(err)
		return
	}
	recorder.Start()
	select {}
}
