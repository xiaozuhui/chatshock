package entities

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type FileEntity struct {
	UUID          uuid.UUID   `json:"uuid"`
	Bucket        string      `json:"bucket"`          // bucket名称
	FileName      string      `json:"filename"`        // 文件名称
	FileURL       string      `json:"file_url"`        // 文件可访问url
	URLExpireTime *time.Time  `json:"url_expire_time"` // url失效时间
	FileType      FileTypeStr `json:"file_type"`       // 文件类型
	ContentType   string      `json:"mime_type"`       // MIME类型
}

type FileType uint8
type FileTypeStr string

const (
	Photo FileType = iota
	Video
	Voice
	Document
	Binary

	PhotoStr    FileTypeStr = "photo"
	VideoStr                = "video"
	VoiceStr                = "voice"
	DocumentStr             = "document"
	BinaryStr               = "binary"
)

func (t FileType) Parse() FileTypeStr {
	var ft FileTypeStr
	switch t {
	case Photo:
		ft = PhotoStr
	case Video:
		ft = VideoStr
	case Voice:
		ft = VoiceStr
	case Document:
		ft = DocumentStr
	case Binary:
		ft = BinaryStr
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的文件类型: %v", t))))
	}
	return ft
}

func (t FileTypeStr) Parse() FileType {
	var ft FileType
	switch t {
	case PhotoStr:
		ft = Photo
	case VideoStr:
		ft = Video
	case VoiceStr:
		ft = Voice
	case DocumentStr:
		ft = Document
	case BinaryStr:
		ft = Binary
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的文件类型: %v", t))))
	}
	return ft
}

func ContentType2FileType(contentType string) (FileTypeStr, error) {
	splits := strings.Split(contentType, "/")
	if len(splits) < 2 {
		return "", errors.New(fmt.Sprintf("上传文件的ContentType[%s]不明", contentType))
	}
	var fileType FileTypeStr
	switch splits[0] {
	case "text":
		fileType = DocumentStr
	case "video":
		fileType = VideoStr
	case "audio":
		fileType = VoiceStr
	case "image":
		fileType = PhotoStr
	case "application":
		fileType = BinaryStr
	}
	return fileType, nil
}
