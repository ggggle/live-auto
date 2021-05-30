package src

func init() {
	//if err := cfg.ReadConfig("./config.xml"); nil != err {
	//	fmt.Printf("配置文件./config.xml读取错误[%s]\n", err.Error())
	//	os.Exit(-1)
	//}
	InitLogger()
	InitGDrive()
}
