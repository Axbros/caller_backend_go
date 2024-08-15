package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Sms struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MachineCode string `gorm:"column:machine_code;type:varchar(32)" json:"machineCode"`
	Address     string `gorm:"column:address;type:varchar(255)" json:"address"`
	Body        string `gorm:"column:body;type:text" json:"body"`
	SmsType     string `gorm:"column:sms_type;type:varchar(16)" json:"smsType"`
	From        string `gorm:"column:from;type:varchar(12)" json:"from"`
}
