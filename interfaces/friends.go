package interfaces

import (
	"chatshock/entities"

	"github.com/gofrs/uuid"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-04 16:42:40
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-04 16:54:09
 * @Description:
 */

type IFriend interface {
	// GetFriends 获取用户所关注的其他用户
	GetFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error)
	// GetBindFriends 获取该用户所双向绑定的其他用户
	GetBindFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error)
	// GetUnBindFriends 获取该用户的好友中已经删除该用户的其他用户
	GetUnBindFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error)
	// IsBindFriend 判断某个用户是否是该用户的绑定用户
	IsBindFriend(userID, otherID uuid.UUID) (bool, error)
	// AddFriend 添加用户（比如该用户删了对方，但对方没有删，可以用这个加回来）
	AddFriend(userID, otherID uuid.UUID) error
	// DeleteFriend 单方面删除other用户的好友
	DeleteFriend(userID, otherID uuid.UUID) error
}
