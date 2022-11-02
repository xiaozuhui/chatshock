package utils
/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:25:08
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 23:18:19
 * @Description:
 */

import (
	"crypto/md5"
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/issue9/identicon/v2"
)

func GenerateAvatar(phoneNumber string) (*os.File, error) {
	ident := identicon.New(identicon.Style2, 128, color.NRGBA{}, color.NRGBA{})
	pnMD5 := fmt.Sprintf("%v", md5.Sum([]byte(phoneNumber)))
	img := ident.Make([]byte(pnMD5))
	tmpFile, err := os.Create("../tmp/tmp.png")
	if err != nil {
		return nil, err
	}
	err = png.Encode(tmpFile, img)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}
