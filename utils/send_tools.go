package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:48:58
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-22 09:50:24
 * @Description:
 */

import (
	"chatshock/configs"
	"chatshock/interfaces"
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

type SendCode string

const (
	Unknown               SendCode = "unknown"
	EmailRegisterCode     SendCode = "register_code"
	EmailRegisterCheckURL SendCode = "register_check_url"
	EmailReBind           SendCode = "rebind_code"
)

func ParseSendCode(sendType string) SendCode {
	switch SendCode(sendType) {
	case EmailRegisterCode:
		return EmailRegisterCode
	case EmailRegisterCheckURL:
		return EmailRegisterCheckURL
	case EmailReBind:
		return EmailReBind
	default:
		return Unknown
	}
}

func (c SendCode) String() string {
	return string(c)
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
	tmpStr, err := ParseTemplate(st, options...)
	if err != nil {
		return err
	}
	toEmail := s.EmailAddress
	err = sendEmailMessage(toEmail, subject, *tmpStr)
	return err
}

// sendEmailMessage 发送email消息
func sendEmailMessage(toEmail, subject, HTMLStr string) error {
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
