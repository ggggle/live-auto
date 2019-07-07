package src

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"time"
)

var (
	GDRIVE_ROOT_DIR    = "LiveAuto" // 所有上传所在根目录
	GDRIVE_ROOT_DIR_ID = ""         // 查询出来的根目录ID
)

type GDriveUploader struct {
	parentDirID string
	Recorder    *Recorder
}

func (self *GDriveUploader) DoUpload(file_path string) {
	if "" == GDRIVE_ROOT_DIR_ID {
		return
	}
	if "" == self.parentDirID {
		dirName := self.generateDirName()
		self.parentDirID, _ = GetDirID(dirName)
		// 创建目录重试次数
		maxRetryTimes := 3
		for err := error(nil); nil != err && "" == self.parentDirID && maxRetryTimes > 0; maxRetryTimes-- {
			// 在ROOT目录基础上创建子目录
			self.parentDirID, err = MakeDir(dirName, GDRIVE_ROOT_DIR_ID)
			time.Sleep(time.Second)
		}
	}
	// 或者直接上传至根目录里？
	if "" == self.parentDirID {
		return
	}
	Logger.WithFields(logrus.Fields{
		"file_path": file_path,
	}).Info("gdrive开始上传")
	for uploadRetryTimes := 5; uploadRetryTimes >= 0; uploadRetryTimes-- {
		uploadArgs := []string{"upload", "-p", self.parentDirID}
		if self.Recorder.AutoRemove {
			uploadArgs = append(uploadArgs, "--delete")
		}
		uploadArgs = append(uploadArgs, file_path)
		cmd := exec.Command("gdrive", uploadArgs...)
		w := bytes.NewBuffer(nil)
		cmd.Stdout = w
		err := cmd.Run()
		if nil != err {
			Logger.WithFields(logrus.Fields{
				ERROR_CONTENT_DEF: err.Error(),
			}).Error("gdrive上传命令执行异常")
			return
		}
		uploadRet := string(w.Bytes())
		logEntry := Logger.WithFields(logrus.Fields{
			"return":     uploadRet,
			"left_times": uploadRetryTimes,
		})
		// 上传成功
		if strings.Contains(uploadRet, "Uploaded") {
			logEntry.Info("gdrive上传成功")
		} else {
			if 0 == uploadRetryTimes {
				logEntry.WithFields(logrus.Fields{
					"args": strings.Join(uploadArgs, ""),
				}).Info("gdrive上传重试截止")
				return
			}
			logEntry.Info("gdrive上传失败，稍后重试")
			select {
			case <-time.After(3 * time.Minute):
			}
		}
	}
}

func (self *GDriveUploader) generateDirName() string {
	// 平台-主播名
	return fmt.Sprintf("%s-%s",
		self.Recorder.Live.GetPlatformCNName(),
		self.Recorder.Live.GetCachedInfo().HostName)
}

func GetDirID(dir_name string) (dir_id string, err error) {
	queryArg := fmt.Sprintf("name='%s'", dir_name)
	// listCmdLine := fmt.Sprintf("gdrive list --query name='%s'", dir_name)
	cmd := exec.Command("gdrive", "list", "-q", queryArg)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if nil != err {
		Logger.WithFields(logrus.Fields{
			ERROR_CONTENT_DEF: err.Error(),
			"stderr":          string(stderr.String()),
			"cmd":             cmd.Args,
		}).Error("gdrive查询目录ID失败")
		return
	}
	lines := strings.Split(stdout.String(), "\n")
	if 1 == len(lines) {
		return
	}
	for _, oneLine := range lines {
		if strings.Contains(oneLine, " dir ") {
			lineSplit := strings.Split(oneLine, " ")
			dir_id = lineSplit[0]
			return
		}
	}
	return
}

/*调用gdrive创建目录
@dir_name  目录名
@parent_id 父目录ID

@dir_id  所创建目录的ID
*/
func MakeDir(dir_name string, parent_id ...string) (dir_id string, err error) {
	var args []string
	args = append(args, "mkdir")
	if 0 == len(parent_id) {
		args = append(args, dir_name)
	} else {
		args = append(args, "-p")
		args = append(args, parent_id[0])
		args = append(args, dir_name)
	}
	cmd := exec.Command("gdrive", args...)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if nil != err {
		Logger.WithFields(logrus.Fields{
			ERROR_CONTENT_DEF: err.Error(),
			"stderr":          stderr.String(),
			"args":            args,
		}).Error("gdrive创建目录失败")
	} else {
		if respStrSplit := strings.Split(stdout.String(), " "); len(respStrSplit) >= 2 {
			dir_id = respStrSplit[1]
		} else {
			Logger.WithFields(logrus.Fields{
				"resp": respStrSplit,
			}).Info("gdrive创建目录返回异常")
		}
	}
	return
}

func InitGDrive() {
	// 查询应用根目录ID
	go func() {
		GDRIVE_ROOT_DIR_ID, _ = GetDirID(GDRIVE_ROOT_DIR)
		maxRetryTimes := 3
		for err := error(nil); nil == err && "" == GDRIVE_ROOT_DIR_ID; maxRetryTimes-- {
			GDRIVE_ROOT_DIR_ID, err = MakeDir(GDRIVE_ROOT_DIR)
			time.Sleep(time.Second)
		}
		if "" == GDRIVE_ROOT_DIR_ID {
			Logger.Info("gdrive根目录ID获取失败，上传功能将不生效")
		} else {
			Logger.WithFields(logrus.Fields{
				"GDRIVE_DIR_ID": GDRIVE_ROOT_DIR_ID,
			}).Info("gdrive根目录ID获取成功")
		}
	}()
}
