package entities

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

type FileEntity struct {
	UUID          uuid.UUID   `json:"uuid"`
	Bucket        string      `json:"bucket"`          // bucket名称
	FileName      string      `json:"filename"`        // 文件名称
	FileURL       string      `json:"file_url"`        // 文件可访问url
	URLExpireTime *time.Time  `json:"url_expire_time"` // url失效时间
	FileType      FileTypeStr `json:"file_type"`       // 文件类型
	MIMEType      string      `json:"mime_type"`       // MIME类型
}

type FileType uint8
type FileTypeStr string

const (
	Photo FileType = iota
	Video
	Voice
	Document

	PhotoStr    FileTypeStr = "photo"
	VideoStr    FileTypeStr = "video"
	VoiceStr    FileTypeStr = "voice"
	DocumentStr FileTypeStr = "document"
)

func (t FileType) Parse() FileTypeStr {
	var ft FileTypeStr
	switch t {
	case Photo:
		ft = PhotoStr
		break
	case Video:
		ft = VideoStr
		break
	case Voice:
		ft = VoiceStr
		break
	case Document:
		ft = DocumentStr
		break
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
		break
	case VideoStr:
		ft = Video
		break
	case VoiceStr:
		ft = Voice
		break
	case DocumentStr:
		ft = Document
		break
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的文件类型: %v", t))))
	}
	return ft
}
