package src

import (
	"os"
	"path"
	"strings"
)

func createFile(filePath string) (file *os.File, err error) {
	// 创建目录
	err = os.MkdirAll(path.Dir(strings.ReplaceAll(filePath, "\\", "/")), os.ModePerm)
	if nil == err {
		// 创建文件
		return os.Create(filePath)
	}
	return
}
