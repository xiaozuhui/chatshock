/*
 * @Author: xiaozuhui
 * @Date: 2022-12-13 16:52:26
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:57:04
 * @Description:
 */
package tests

import (
	"chatshock/utils"
	"testing"
)

func TestSender(t *testing.T) {
	sender := utils.EmailAddress{EmailAddress: "xiaozuhui@outlook.com"}
	// err := sender.SendMessage("register_code", "注册", utils.Options{FieldName: func() string {
	// 	return "code"
	// }, FieldValue: func() string {
	// 	return "123456"
	// }}, utils.Options{FieldName: func() string {
	// 	return "sign_name"
	// }, FieldValue: func() string {
	// 	return "xiaozuhui"
	// }}, utils.Options{FieldName: func() string {
	// 	return "time_now"
	// }, FieldValue: func() string {
	// 	return time.Now().Format("2006年01月02日 15时04分")
	// }})
	// if err != nil {
	// 	t.Error(errors.WithStack(err))
	// }
	t.Log(sender)
}
