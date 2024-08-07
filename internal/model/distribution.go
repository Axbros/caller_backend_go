package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Distribution struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	UserID      int    `gorm:"column:user_id;type:int(11)" json:"userId"`
	GroupCallID uint64 `gorm:"column:group_call_id;type:int(11)" json:"groupCallId"`
}

// TableName table name
func (m *Distribution) TableName() string {
	return "distribution"
}
