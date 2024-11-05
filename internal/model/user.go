package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type User struct {
	ggorm.Model `gorm:"embedded"` // embed id and time
	MachineCode string            `gorm:"column:machine_code;type:varchar(32)" json:"machineCode"`
	Sms         string            `gorm:"column:sms;type:tinyint(1)" json:"sms"`
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}
