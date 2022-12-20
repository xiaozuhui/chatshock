package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:25:08
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 23:18:19
 * @Description:
 */

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/issue9/identicon/v2"
	"image/color"
	"image/png"
	"math/rand"
	"strings"
)

func GenerateAvatar(UUID uuid.UUID) (*bytes.Buffer, error) {
	ident := identicon.New(identicon.Style2, 128, color.NRGBA{}, color.NRGBA{}, color.NRGBA{}, color.NRGBA{})
	pnMD5 := fmt.Sprintf("%v", md5.Sum([]byte(UUID.String())))
	aNum := strings.Split(pnMD5, "a")
	img := ident.Rand(rand.New(rand.NewSource(int64(len(aNum) + 100))))
	img = ident.Make([]byte("192.168.1.1"))
	buff := bytes.Buffer{}
	err := png.Encode(&buff, img)
	if err != nil {
		return nil, err
	}
	return &buff, nil
}
