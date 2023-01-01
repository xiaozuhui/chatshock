package repositories

import (
	"chatshock/configs"
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/models"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

type ChatRoomRepo struct {
}

func (c ChatRoomRepo) FindAllChatRoom() ([]*entities.ChatRoom, error) {
	var chatRooms []models.ChatroomModel
	var users []models.UserModel
	var avatars []models.FileModel

	err := configs.DBEngine.Find(&chatRooms).Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]uuid.UUID, 0, 0)
	avatarIDs := make([]uuid.UUID, 0, 0)
	userMap := make(map[uuid.UUID]*entities.UserEntity)
	avatarMap := make(map[uuid.UUID]models.FileModel)
	for _, chatRoom := range chatRooms {
		userIDMap := make(map[uuid.UUID]struct{}, 0)
		err := json.Unmarshal(chatRoom.Users, &userIDMap)
		if err != nil {
			return nil, err
		}
		for uID := range userIDMap {
			userIDs = append(userIDs, uID)
		}
		avatarIDs = append(avatarIDs, chatRoom.ChatRoomAvatar)
	}
	err = configs.DBEngine.Where("uuid IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap[user.UUID] = user.ModelToEntity().(*entities.UserEntity)
	}
	err = configs.DBEngine.Where("uuid IN (?)", avatarIDs).Find(&avatars).Error
	if err != nil {
		return nil, err
	}
	for _, avatar := range avatars {
		avatarMap[avatar.UUID] = avatar
	}
	res := make([]*entities.ChatRoom, 0, 0)
	for _, chatRoom := range chatRooms {
		chatRoomEntity := chatRoom.ModelToEntity().(*entities.ChatRoom)
		chatRoomEntity.ChatRoomAvatar = avatarMap[chatRoom.UUID].ModelToEntity().(*entities.FileEntity)
		chatRoomEntity.Users = userMap
		chatRoomEntity.Master = userMap[chatRoom.Master]
		res = append(res, chatRoomEntity)
	}
	return res, nil
}

func (c ChatRoomRepo) FindChatRooms(IDs []uuid.UUID) ([]*entities.ChatRoom, error) {
	var chatRooms []models.ChatroomModel
	var users []models.UserModel
	var avatars []models.FileModel
	err := configs.DBEngine.Where("uuid IN (?)", IDs).Find(&chatRooms).Error
	if err != nil {
		return nil, err
	}
	userIDs := make([]uuid.UUID, 0, 0)
	avatarIDs := make([]uuid.UUID, 0, 0)
	userMap := make(map[uuid.UUID]*entities.UserEntity)
	avatarMap := make(map[uuid.UUID]models.FileModel)
	for _, chatRoom := range chatRooms {
		userIDMap := make(map[uuid.UUID]struct{}, 0)
		err := json.Unmarshal(chatRoom.Users, &userIDMap)
		if err != nil {
			return nil, err
		}
		for uID := range userIDMap {
			userIDs = append(userIDs, uID)
		}
		avatarIDs = append(avatarIDs, chatRoom.ChatRoomAvatar)
	}
	err = configs.DBEngine.Where("uuid IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap[user.UUID] = user.ModelToEntity().(*entities.UserEntity)
	}
	err = configs.DBEngine.Where("uuid IN (?)", avatarIDs).Find(&avatars).Error
	if err != nil {
		return nil, err
	}
	for _, avatar := range avatars {
		avatarMap[avatar.UUID] = avatar
	}
	res := make([]*entities.ChatRoom, 0, 0)
	for _, chatRoom := range chatRooms {
		chatRoomEntity := chatRoom.ModelToEntity().(*entities.ChatRoom)
		chatRoomEntity.ChatRoomAvatar = avatarMap[chatRoom.UUID].ModelToEntity().(*entities.FileEntity)
		chatRoomEntity.Users = userMap
		chatRoomEntity.Master = userMap[chatRoom.Master]
		res = append(res, chatRoomEntity)
	}
	return res, nil
}

func (c ChatRoomRepo) FindChatRoom(ID uuid.UUID) (*entities.ChatRoom, error) {
	var chatroom models.ChatroomModel
	var users []models.UserModel
	var avatar models.FileModel
	err := configs.DBEngine.First(&chatroom, "uuid = ?", ID).Error
	if err != nil {
		return nil, err
	}
	chatRoomEntity := chatroom.ModelToEntity().(*entities.ChatRoom)
	err = configs.DBEngine.First(&avatar, "uuid = ?", chatroom.ChatRoomAvatar).Error
	if err != nil {
		return nil, err
	}
	chatRoomEntity.ChatRoomAvatar = avatar.ModelToEntity().(*entities.FileEntity)
	userIDs := make([]uuid.UUID, 0, 0)
	userIDMap := make(map[uuid.UUID]struct{}, 0)
	err = json.Unmarshal(chatroom.Users, &userIDMap)
	if err != nil {
		return nil, err
	}
	for uID := range userIDMap {
		userIDs = append(userIDs, uID)
	}
	err = configs.DBEngine.Where("uuid IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	userEntitiesMap := make(map[uuid.UUID]*entities.UserEntity, 0)
	for _, u := range users {
		userEntitiesMap[u.UUID] = u.ModelToEntity().(*entities.UserEntity)
	}
	chatRoomEntity.Users = userEntitiesMap
	chatRoomEntity.Master = userEntitiesMap[chatroom.Master]
	return chatRoomEntity, nil
}

func (c ChatRoomRepo) FindChatRoomByMaster(masterID uuid.UUID) ([]*entities.ChatRoom, error) {
	var chatRooms []models.ChatroomModel
	var users []models.UserModel
	var avatars []models.FileModel
	err := configs.DBEngine.Where("master IN (?)", masterID).Find(&chatRooms).Error
	if err != nil {
		return nil, err
	}
	userIDs := make([]uuid.UUID, 0, 0)
	avatarIDs := make([]uuid.UUID, 0, 0)
	userMap := make(map[uuid.UUID]*entities.UserEntity)
	avatarMap := make(map[uuid.UUID]models.FileModel)
	for _, chatRoom := range chatRooms {
		userIDMap := make(map[uuid.UUID]struct{}, 0)
		err := json.Unmarshal(chatRoom.Users, &userIDMap)
		if err != nil {
			return nil, err
		}
		for uID := range userIDMap {
			userIDs = append(userIDs, uID)
		}
		avatarIDs = append(avatarIDs, chatRoom.ChatRoomAvatar)
	}
	err = configs.DBEngine.Where("uuid IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap[user.UUID] = user.ModelToEntity().(*entities.UserEntity)
	}
	err = configs.DBEngine.Where("uuid IN (?)", avatarIDs).Find(&avatars).Error
	if err != nil {
		return nil, err
	}
	for _, avatar := range avatars {
		avatarMap[avatar.UUID] = avatar
	}
	res := make([]*entities.ChatRoom, 0, 0)
	for _, chatRoom := range chatRooms {
		chatRoomEntity := chatRoom.ModelToEntity().(*entities.ChatRoom)
		chatRoomEntity.ChatRoomAvatar = avatarMap[chatRoom.UUID].ModelToEntity().(*entities.FileEntity)
		chatRoomEntity.Users = userMap
		chatRoomEntity.Master = userMap[masterID]
		res = append(res, chatRoomEntity)
	}
	return res, nil
}

func (c ChatRoomRepo) CreateChatRoom(chatRoomEntity *entities.ChatRoom) error {
	chatRoomModel, err := models.EntityToChatroomModel(chatRoomEntity)
	if err != nil {
		return err
	}
	baseModel, err := models.NewBaseModel(chatRoomEntity.UUID)
	if err != nil {
		return err
	}
	chatRoomModel.BaseModel = *baseModel
	err = configs.DBEngine.Model(&models.ChatroomModel{}).Create(chatRoomModel).Error
	if err != nil {
		return err
	}
	return nil
}

func (c ChatRoomRepo) DeleteChatRoom(ID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (c ChatRoomRepo) DeleteChatRoomsByMaster(masterID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (c ChatRoomRepo) FindChatRoomByUser(userID uuid.UUID) ([]*entities.ChatRoom, error) {
	var chatRooms []models.ChatroomModel
	var users []models.UserModel
	var avatars []models.FileModel
	err := configs.DBEngine.Find(&chatRooms, datatypes.JSONQuery("users").HasKey(userID.String())).Error
	if err != nil {
		return nil, err
	}
	userIDs := make([]uuid.UUID, 0, 0)
	avatarIDs := make([]uuid.UUID, 0, 0)
	userMap := make(map[uuid.UUID]*entities.UserEntity)
	avatarMap := make(map[uuid.UUID]models.FileModel)
	for _, chatRoom := range chatRooms {
		userIDMap := make(map[uuid.UUID]struct{}, 0)
		err := json.Unmarshal(chatRoom.Users, &userIDMap)
		if err != nil {
			return nil, err
		}
		for uID := range userIDMap {
			userIDs = append(userIDs, uID)
		}
		avatarIDs = append(avatarIDs, chatRoom.ChatRoomAvatar)
	}
	err = configs.DBEngine.Where("uuid IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap[user.UUID] = user.ModelToEntity().(*entities.UserEntity)
	}
	err = configs.DBEngine.Where("uuid IN (?)", avatarIDs).Find(&avatars).Error
	if err != nil {
		return nil, err
	}
	for _, avatar := range avatars {
		avatarMap[avatar.UUID] = avatar
	}
	res := make([]*entities.ChatRoom, 0, 0)
	for _, chatRoom := range chatRooms {
		chatRoomEntity := chatRoom.ModelToEntity().(*entities.ChatRoom)
		chatRoomEntity.ChatRoomAvatar = avatarMap[chatRoom.UUID].ModelToEntity().(*entities.FileEntity)
		chatRoomEntity.Users = userMap
		chatRoomEntity.Master = userMap[chatRoom.Master]
		res = append(res, chatRoomEntity)
	}
	return res, nil
}

func (c ChatRoomRepo) IntoChatRoom(userID, chatroomID uuid.UUID) error {
	var chatroom models.ChatroomModel
	err := configs.DBEngine.First(&chatroom, "uuid = ?", chatroomID).Error
	if err != nil {
		return err
	}
	userMap := make(map[uuid.UUID]struct{}, 0)
	err = json.Unmarshal(chatroom.Users, &userMap)
	if err != nil {
		return err
	}
	if _, ok := userMap[userID]; ok {
		return errors.WithStack(errors.New(fmt.Sprintf("用户[%v]已经存在于该聊天室中[%s]", userID, chatroom.Name)))
	}
	userMap[userID] = struct{}{}
	userMarshal, err := json.Marshal(userMap)
	if err != nil {
		return err
	}
	chatroom.Users = userMarshal
	err = configs.DBEngine.Save(chatroom).Error
	return err
}

func (c ChatRoomRepo) OutFromChatRoom(userID, chatroomID uuid.UUID) error {
	var chatroom models.ChatroomModel
	err := configs.DBEngine.First(&chatroom, "uuid = ?", chatroomID).Error
	if err != nil {
		return err
	}
	userMap := make(map[uuid.UUID]struct{}, 0)
	err = json.Unmarshal(chatroom.Users, &userMap)
	if err != nil {
		return err
	}
	if _, ok := userMap[userID]; !ok {
		return errors.WithStack(errors.New(fmt.Sprintf("用户[%v]不在该聊天室中[%s]", userID, chatroom.Name)))
	}
	// 判断是否是master用户
	if userID == chatroom.Master {
		return errors.WithStack(errors.New(fmt.Sprintf("用户[%v]不在该聊天室的群主，不能直接退出", userID)))
	}
	delete(userMap, userID) // 删除userID
	userMarshal, err := json.Marshal(userMap)
	if err != nil {
		return err
	}
	chatroom.Users = userMarshal
	err = configs.DBEngine.Save(chatroom).Error
	return err
}

var _ interfaces.IChatRoom = ChatRoomRepo{}
