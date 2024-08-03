package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type CallHistory struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	RequestMachineCode string `gorm:"column:request_machine_code;type:varchar(32)" json:"requestMachineCode"`
	ClientMachineCode  string `gorm:"column:client_machine_code;type:varchar(32)" json:"clientMachineCode"`
	MobileNumber       string `gorm:"column:mobile_number;type:varchar(11)" json:"mobileNumber"`
	Instruction        string `gorm:"column:instruction;type:varchar(16)" json:"instruction"`
}

// TableName table name
func (m *CallHistory) TableName() string {
	return "call_history"
}
