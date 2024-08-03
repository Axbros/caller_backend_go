package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type UnanswerdCall struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	ClientMachineCode string `gorm:"column:client_machine_code;type:varchar(32)" json:"clientMachineCode"`
	MobileNumber      string `gorm:"column:mobile_number;type:varchar(11)" json:"mobileNumber"`
}

// TableName table name
func (m *UnanswerdCall) TableName() string {
	return "unanswerd_call"
}
