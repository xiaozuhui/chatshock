package entities

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-02 12:22:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-07 16:42:54
 * @Description:
 */

import (
	"fmt"

	"github.com/pkg/errors"
)

type MessageEntity struct {
}

type MessageType uint8
type MessageTypeStr string

/*
*
消息类型：

	1、文字消息
	2、图片消息
	3、动图消息（表情）
	4、视频消息
	5、语音消息
	6、文件消息
	7、复合消息（文字消息和图片组合在一起的消息类型）
*/
const (
	MtText MessageType = iota
	MtPhoto
	MtDynamicPhoto
	MtVideo
	MtVoice
	MtFile
	MtTextAndPhoto

	MtTextStr         MessageTypeStr = "text"
	MtPhotoStr        MessageTypeStr = "photo"
	MtDynamicPhotoStr MessageTypeStr = "dynamic_photo"
	MtVideoStr        MessageTypeStr = "video"
	MtVoiceStr        MessageTypeStr = "voice"
	MtFileStr         MessageTypeStr = "file"
	MtTextAndPhotoStr MessageTypeStr = "text_and_photo"
)

// Parse 将类型MessageType转为MessageTypeStr类型
func (t MessageType) Parse() MessageTypeStr {
	var typ MessageTypeStr
	switch t {
	case MtText:
		typ = MtTextStr
	case MtPhoto:
		typ = MtPhotoStr
	case MtVideo:
		typ = MtVideoStr
	case MtVoice:
		typ = MtVoiceStr
	case MtFile:
		typ = MtFileStr
	case MtTextAndPhoto:
		typ = MtTextAndPhotoStr
	case MtDynamicPhoto:
		typ = MtDynamicPhotoStr
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的类型：%v", t))))
	}
	return typ
}

// Parse 将类型MessageTypeStr转为类型MessageType
func (t MessageTypeStr) Parse() MessageType {
	var typ MessageType
	switch t {
	case MtTextStr:
		typ = MtText
	case MtPhotoStr:
		typ = MtPhoto
	case MtVideoStr:
		typ = MtVideo
	case MtVoiceStr:
		typ = MtVoice
	case MtFileStr:
		typ = MtFile
	case MtTextAndPhotoStr:
		typ = MtTextAndPhoto
	case MtDynamicPhotoStr:
		typ = MtDynamicPhoto
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的类型：%v", t))))
	}
	return typ
}
