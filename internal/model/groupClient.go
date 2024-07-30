package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type GroupClient struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	GroupID  int `gorm:"column:group_id;type:int(11);NOT NULL" json:"groupId"`
	ClientID int `gorm:"column:client_id;type:int(11);NOT NULL" json:"clientId"`
}

// TableName table name
func (m *GroupClient) TableName() string {
	return "group_client"
}
