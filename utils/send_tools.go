package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:48:58
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:18:01
 * @Description:
 */

import (
	"chatshock/configs"
	"chatshock/interfaces"
	"encoding/json"
	"errors"

	smsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"gopkg.in/gomail.v2"
)

type Options struct {
	FieldName  func() string
	FieldValue func() string
}

func (t Options) GetKey() string {
	return t.FieldName()
}

func (t Options) GetValue() string {
	return t.FieldValue()
}

const (
	Phone string = "phone"
	Email string = "email"
)

const (
	// 发送方式+放松目的
	PhoneRegister         string = "phone_register"
	EmailRegisterCode     string = "register_code"
	EmailRegisterCheckURL string = "register_check_url"
)

// SendPhoneMessage 发送手机消息
func SendPhoneMessage(phoneNumber, signName, templateCode string, content []byte) error {
	var err error
	sendSmsRequest := &smsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(string(content)),
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

// SendEmailMessage 发送email消息
func SendEmailMessage(toEmail, subject, HTMLStr string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(
		configs.Conf.EmailConfig.FromAddress,
		configs.Conf.EmailConfig.FromName))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", HTMLStr)
	d := gomail.NewDialer(configs.Conf.EmailConfig.EmailSmtp, 465,
		configs.Conf.EmailConfig.FromAddress,
		configs.Conf.EmailConfig.Secret)
	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

type EmailAddress struct {
	EmailAddress string `json:"email_address"`
}

func (s EmailAddress) String() string {
	return string(s.EmailAddress)
}

func (s EmailAddress) Type() string {
	return Email
}

func (s EmailAddress) SendMessage(st string, subject string, options ...interfaces.Options) error {
	// 获取模版
	tmpStr, err := ParseTemplate(string(st), options...)
	if err != nil {
		return err
	}
	toEmail := s.EmailAddress
	err = SendEmailMessage(toEmail, subject, *tmpStr)
	return err
}

type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}

func (s PhoneNumber) String() string {
	return string(s.PhoneNumber)
}

func (s PhoneNumber) Type() string {
	return Phone
}

func (s PhoneNumber) SendMessage(st string, signName string, options ...interfaces.Options) error {
	if _, ok := configs.Conf.PhoneConfig.SignTemplate[st]; !ok {
		return errors.New("发送类型{st[configs.SendType]}错误")
	}
	// 获取签名和模版code
	signDict := configs.Conf.PhoneConfig.SignTemplate[string(st)]
	// 获取模版字段
	contDict := make(map[string]string, 0)
	for _, option := range options {
		contDict[option.GetKey()] = option.GetValue()
	}
	content, err := json.Marshal(contDict)
	if err != nil {
		return err
	}
	for signName, tplCode := range signDict {
		err = SendPhoneMessage(s.String(), signName, tplCode, content)
	}
	return err
}
