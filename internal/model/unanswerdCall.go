package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type UnanswerdCall struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MachineId    string `gorm:"column:client_id;type:varchar(32)" json:"clientMachineCode"`
	MobileNumber string `gorm:"column:mobile_number;type:varchar(11)" json:"mobileNumber"`
	Location     string `gorm:"column:location;type:varchar(32)" json:"location"`
	Type         string `gorm:"column:type;type:varchar(12)" json:"type"`
}

// TableName table name
func (m *UnanswerdCall) TableName() string {
	return "callLog"
}
