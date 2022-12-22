package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:25:08
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-21 10:00:27
 * @Description:
 */

import (
	"bytes"
	"chatshock/configs"
	"image/png"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/disintegration/letteravatar"
	"github.com/golang/freetype"
)

// /** GenerateAvatar (Deprecated)
//  * @description: 创建头像
//  * @param {uuid.UUID} UUID
//  * @return {*}
//  * @author: xiaozuhui
//  * @deprecated: true
//  */
// func GenerateAvatar(UUID uuid.UUID) (*bytes.Buffer, error) {
// 	ident := identicon.New(identicon.Style2, 128, color.NRGBA64{}, color.NRGBA64{}, color.NRGBA64{}, color.NRGBA64{})
// 	pnMD5 := fmt.Sprintf("%v", md5.Sum([]byte(UUID.String())))
// 	fmt.Println(pnMD5)
// 	img := ident.Make([]byte(pnMD5))
// 	buff := bytes.Buffer{}
// 	err := png.Encode(&buff, img)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &buff, nil
// }

func GenerateAvatar(nickname string, baseDir ...string) (*bytes.Buffer, error) {
	if len(baseDir) == 0 {
		baseDir = append(baseDir, configs.BaseDir)
	}
	fontFile, _ := os.ReadFile(filepath.Join(baseDir[0], "configs/statics", "Alimama_ShuHeiTi_Bold.ttf"))
	font, _ := freetype.ParseFont(fontFile)
	options := &letteravatar.Options{
		Font: font,
	}
	firstLetter, _ := utf8.DecodeRuneInString(nickname)
	img, err := letteravatar.Draw(200, firstLetter, options)
	if err != nil {
		return nil, err
	}
	buff := bytes.Buffer{}
	err = png.Encode(&buff, img)
	if err != nil {
		return nil, err
	}
	return &buff, nil
}
