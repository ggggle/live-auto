package src

const (
	GDRIVE  UploaderType = iota // gdrive网盘
	YOUTUBE                     // youtube
)

type UploaderType int

type Uploader interface {
	DoUpload(file_path string)
}

func NewUploader(_type UploaderType, owner *Recorder) Uploader {
	switch _type {
	case GDRIVE:
		return &GDriveUploader{
			Recorder: owner,
		}
	default:
		return nil
	}
}
