package entities

import "github.com/gofrs/uuid"

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-04 09:48:24
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-04 09:51:31
 * @Description:
 */

type FriendsEntity struct {
	BaseEntity
	UserUUID  uuid.UUID `json:"user_id"`
	OtherUUID uuid.UUID `json:"other_id"`
}
