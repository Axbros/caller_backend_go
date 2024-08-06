package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type GroupClient struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	GroupName string `gorm:"column:group_name;type:varchar(255);NOT NULL" json:"groupName"`
	ClientID  int    `gorm:"column:client_id;type:int(11);NOT NULL" json:"clientId"`
}

// TableName table name
func (m *GroupClient) TableName() string {
	return "group_client"
}
