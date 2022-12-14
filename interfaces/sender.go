package interfaces

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-09 09:32:45
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 14:00:40
 * @Description:
 */

// ISender 发送工具的interfaces

type Options interface {
	GetKey() string
	GetValue() string
}

type ISender interface {
	SendMessage(st, subject string, options ...Options) error
	String() string
	Type() string
}
