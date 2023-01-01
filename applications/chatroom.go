package applications

import (
	"chatshock/entities"
	"chatshock/services"
	"chatshock/utils"
	"github.com/gofrs/uuid"
)

type ChatRoomApplication struct {
	UserService     services.UserService
	FileService     services.FileService
	ChatRoomService services.ChatRoomService
}

func NewChatRoomApplication() ChatRoomApplication {
	return ChatRoomApplication{
		UserService:     services.UserFactory(),
		FileService:     services.FileFactory(),
		ChatRoomService: services.ChatRoomFactory(),
	}
}

func (a ChatRoomApplication) CreateChatRoom(userID uuid.UUID, crName, crDescription string) (*entities.ChatRoom, error) {
	user, err := a.UserService.UserRepo.FindUser(userID)
	if err != nil {
		return nil, err
	}
	avatar, err := utils.GenerateAvatar(crName)
	if err != nil {
		return nil, err
	}
	charRoomUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	err = utils.MakeBucket(charRoomUUID.String())
	if err != nil {
		return nil, err
	}
	// 保存文件
	uploadInfo, err := utils.UploadImage(charRoomUUID.String(), charRoomUUID.String()+"_avatar.png", avatar)
	if err != nil {
		return nil, err
	}
	// 保存文件信息
	fileEntity, err := a.FileService.SaveFile(uploadInfo, entities.PhotoStr, "application/png")
	if err != nil {
		return nil, err
	}
	// 创建entity
	chatRoomEntity := entities.ChatRoom{
		BaseEntity:     *entities.NewBaseEntity(charRoomUUID),
		Name:           crName,
		Description:    crDescription,
		Master:         user,
		Users:          map[uuid.UUID]*entities.UserEntity{user.UUID: user},
		ChatRoomAvatar: fileEntity,
	}
	room, err := a.ChatRoomService.CreateChatRoom(&chatRoomEntity)
	if err != nil {
		return nil, err
	}
	return room, nil
}
