/*
 * @Author: xiaozuhui
 * @Date: 2022-12-02 12:22:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:46:04
 * @Description:
 */
package interfaces

import (
	"chatshock/entities"

	"github.com/gofrs/uuid"
)

type IFile interface {
	// SaveFile 保存文件信息
	SaveFile(fileEntity *entities.FileEntity) (uuid.UUID, error)
	// SaveFiles 批量保存文件信息
	SaveFiles(fileEntities []*entities.FileEntity) ([]uuid.UUID, error)
	// GetFile 根据uuid获取文件信息
	GetFile(id uuid.UUID) (*entities.FileEntity, error)
	// GetFiles 批量获取文件信息
	GetFiles(ids []uuid.UUID) ([]*entities.FileEntity, error)
	// DeleteFile 删除文件
	DeleteFile(id uuid.UUID) error
	// DeleteFiles 批量删除文件
	DeleteFiles(ids []uuid.UUID) error
	// UpdateFile 更新文件信息，只能更新文件url和URL过期时间
	UpdateFile(id uuid.UUID, fileEntity *entities.FileEntity) (*entities.FileEntity, error)
	// UpdateFiles 批量更新
	UpdateFiles(files map[uuid.UUID]*entities.FileEntity) ([]*entities.FileEntity, error)
}
