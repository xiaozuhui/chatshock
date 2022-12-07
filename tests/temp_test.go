package tests

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-09 09:32:45
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:52:16
 * @Description:
 */

import (
	"chatshock/utils"
	"testing"
	"time"
)

func TestParseTemplate(t *testing.T) {
	emailStr, err := utils.ParseTemplate("register_code",
		utils.Options{FieldName: func() string {
			return "code"
		}, FieldValue: func() string {
			return "123456"
		}}, utils.Options{FieldName: func() string {
			return "sign_name"
		}, FieldValue: func() string {
			return "xiaozuhui"
		}}, utils.Options{FieldName: func() string {
			return "time_now"
		}, FieldValue: func() string {
			return time.Now().Format("2006年01月02日 15时04分")
		}})
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("emailStr: %s", *emailStr)
}
