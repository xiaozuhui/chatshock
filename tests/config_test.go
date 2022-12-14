/*
 * @Author: xiaozuhui
 * @Date: 2022-12-09 13:34:44
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-09 13:57:12
 * @Description:
 */
package tests

import (
	"chatshock/configs"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestGetConfig(t *testing.T) {
	configs.Conf = &configs.Config{}
	configs.BaseDir = "../"
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(configs.BaseDir, "configs", "dev.yaml"))
	err := viper.ReadInConfig()
	if err != nil {
		t.Error(err.Error())
	}
	configs.Conf.Parse(viper.GetViper())
	t.Log(configs.Conf.PhoneConfig)
}
