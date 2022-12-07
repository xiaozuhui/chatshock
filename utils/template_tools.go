package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-09 09:32:45
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:50:00
 * @Description:
 */

import (
	"bytes"
	"chatshock/configs"
	"chatshock/interfaces"
	"html/template"
	"path/filepath"
)

// ParseTemplate 解析模版，将数据添加进模版，返回HTML字符串
func ParseTemplate(tmpName string, fields ...interfaces.Options) (*string, error) {
	buff := bytes.Buffer{}
	if configs.BaseDir == "" {
		configs.BaseDir = "../"
	}
	var tmpPath = filepath.Join(configs.BaseDir, "templates", tmpName+".tmpl")
	tpl, err := template.ParseFiles(tmpPath)
	if err != nil {
		return nil, err
	}
	contentMap := make(map[string]template.HTML, 0)
	for _, item := range fields {
		contentMap[item.GetKey()] = template.HTML(item.GetValue())
	}
	err = tpl.Execute(&buff, contentMap)
	if err != nil {
		return nil, err
	}
	s := buff.String()
	return &s, err
}
