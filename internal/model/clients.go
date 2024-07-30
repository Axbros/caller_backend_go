package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Clients struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MachineCode string `gorm:"column:machine_code;type:varchar(32)" json:"machineCode"`
	IPAddress   string `gorm:"column:ip_address;type:varchar(32)" json:"ipAddress"`
}
