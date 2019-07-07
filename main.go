package main

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/api"
	"live-auto/src"
)

func main() {
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

func list()  {
	list := src.GetIRecorderMngr().GetAllRecorder()
	for _, rcd := range list{
		fmt.Printf("%s %s\n", rcd.Live.GetLiveId(), rcd.Live.GetRawUrl())
	}
}