package src

import (
	"fmt"
	"os"
)

func init() {
	if err := ReadConfig("./config.xml"); nil != err {
		fmt.Printf("配置文件./config.xml读取错误[%s]\n", err.Error())
		os.Exit(-1)
	}
	InitLogger()
	InitGDrive()
}
