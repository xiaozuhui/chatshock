package tests

import (
	"chatshock/utils"
	"github.com/gofrs/uuid"
	"os"
	"testing"
)

func TestAvatar(t *testing.T) {
	v4, err := uuid.NewV4()
	if err != nil {
		t.Error(err)
	}
	avatar, err := utils.GenerateAvatar(v4)
	if err != nil {
		t.Error(err)
	}
	file, err := os.OpenFile("tmp.png", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Error(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Error(err)
		}
	}(file)
	_, err = file.Write(avatar.Bytes())
	if err != nil {
		t.Error(err)
	}
}
