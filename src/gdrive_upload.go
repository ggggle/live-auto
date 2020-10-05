package src

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"live-auto/cfg"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	GDRIVE_ROOT_DIR       = "LiveAuto" // 所有上传所在根目录
	GDRIVE_ROOT_DIR_ID    = ""         // 查询出来的根目录ID
)

const MIN_UPLOAD_FILE_SIZE int64 = 1024 * 256 // 最小上传大小

type GDriveUploader struct {
	parentDirID string
	Recorder    *Recorder
}

func (self *GDriveUploader) DoUpload(file_path string) {
	if file_info, err := os.Stat(file_path); nil != err || file_info.Size() < MIN_UPLOAD_FILE_SIZE {
		return
	}
	mv_dir := cfg.G_Config.RcloneMountDirPath + self.generateDirName()
	Logger.WithFields(logrus.Fields{
		"file_path": file_path,
	}).Info("gdrive开始上传")
	for uploadRetryTimes := 5; uploadRetryTimes >= 0; uploadRetryTimes-- {
		mkdir_args := []string{"-p", mv_dir}
		mkdir_cmd := exec.Command("mkdir", mkdir_args...)
		mkdir_cmd.Run()
		uploadArgs := []string{file_path, mv_dir}
		cmd := exec.Command("mv", uploadArgs...)
		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		if nil != err {
			Logger.WithFields(logrus.Fields{
				ERROR_CONTENT_DEF: err.Error(),
				"stderr":          stderr.String(),
			}).Error("gdrive上传命令执行异常")
			return
		} else {
			break
		}
	}
}

func (self *GDriveUploader) generateDirName() string {
	// 平台-主播名
	return fmt.Sprintf("%s-%s/%s/",
		self.Recorder.Live.GetPlatformCNName(),
		self.Recorder.Live.GetCachedInfo().HostName,
		time.Now().Format("2006-01"))
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
}
