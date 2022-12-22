package tests

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-21 08:28:22
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-21 10:03:37
 * @Description:
 */

import (
	"chatshock/utils"
	"os"
	"testing"
)

// TestAvatar
/**
 * @description:
 * @param {*testing.T} t
 * @return {*}
 * @author: xiaozuhui
 */
func TestAvatar(t *testing.T) {
	avatar, err := utils.GenerateAvatar("xiaozuhui", "../")
	if err != nil {
		t.Error(err)
	}
	file, err := os.OpenFile("tmp.png", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	_, err = file.Write(avatar.Bytes())
	if err != nil {
		t.Error(err)
	}
}
