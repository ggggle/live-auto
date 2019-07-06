package main

import (
	"fmt"
	"live-auto/src"
)

func main() {
	/*	path := filepath.Join("./中文", "test")
		fmt.Println(path)
		os.MkdirAll(path, os.ModePerm)
		return*/
	fmt.Println(src.G_Config)
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
