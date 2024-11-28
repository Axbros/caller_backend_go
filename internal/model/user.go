package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type User struct {
	ggorm.Model `gorm:"embedded"` // embed id and time
	MachineCode string            `gorm:"column:machine_code;type:varchar(32)" json:"machineCode"`
	Sms         time.Time         `gorm:"column:sms;" json:"sms"`
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}
