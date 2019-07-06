package src

const (
	GDRIVE UploaderType = iota // gdrive网盘
	YOUTUBE        // youtube
)

type UploaderType int

type Uploader interface {
	DoUpload(file_path string)
}

func NewUploader(_type UploaderType) Uploader {
	switch _type {
	case GDRIVE:
		return &GDriveUploader{}
	default:
		return nil
	}
}