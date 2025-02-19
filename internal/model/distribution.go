package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Distribution struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	UserID    int    `gorm:"column:user_id;type:int(11)" json:"userId"`
	GroupName string `gorm:"column:group_name;type:varchar(255)" json:"groupName"`
}

// TableName table name
func (m *Distribution) TableName() string {
	return "distribution"
}
