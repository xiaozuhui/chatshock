/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:48:58
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 11:02:14
 * @Description:
 */
package utils

import (
	"chatshock/configs"
	"encoding/json"
	"fmt"

	smsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func SendValidMessage(phoneNumber string, content map[string]string) error {
	_cont, err := json.Marshal(content)
	if err != nil {
		fmt.Printf("序列化错误 err=%v", err)
		return err
	}
	signName := configs.Conf.PhoneValidConfig.SignName
	templateCode := configs.Conf.PhoneValidConfig.TemplateCode
	sendSmsRequest := &smsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(string(_cont)),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() error {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, err := configs.SMSClient.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		var smsError = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			smsError = _t
		} else {
			smsError.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, err := util.AssertAsString(smsError.Message)
		if err != nil {
			return err
		}
	}
	return err
}
