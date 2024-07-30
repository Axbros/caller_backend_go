package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type GroupCall struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	GroupNumber      string `gorm:"column:group_number;type:varchar(4)" json:"groupNumber"`
	PhoneNumber      string `gorm:"column:phone_number;type:varchar(11)" json:"phoneNumber"`
	TransferClientID string `gorm:"column:transfer_client_id;type:varchar(32)" json:"transferClientId"`
}

// TableName table name
func (m *GroupCall) TableName() string {
	return "group_call"
}
