package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// clients business-level http error codes.
// the clientsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	clientsNO       = 72
	clientsName     = "clients"
	clientsBaseCode = errcode.HCode(clientsNO)

	ErrCreateClients     = errcode.NewError(clientsBaseCode+1, "failed to create "+clientsName)
	ErrDeleteByIDClients = errcode.NewError(clientsBaseCode+2, "failed to delete "+clientsName)
	ErrUpdateByIDClients = errcode.NewError(clientsBaseCode+3, "failed to update "+clientsName)
	ErrGetByIDClients    = errcode.NewError(clientsBaseCode+4, "failed to get "+clientsName+" details")
	ErrListClients       = errcode.NewError(clientsBaseCode+5, "failed to list of "+clientsName)

	ErrDeleteByIDsClients    = errcode.NewError(clientsBaseCode+6, "failed to delete by batch ids "+clientsName)
	ErrGetByConditionClients = errcode.NewError(clientsBaseCode+7, "failed to get "+clientsName+" details by conditions")
	ErrListByIDsClients      = errcode.NewError(clientsBaseCode+8, "failed to list by batch ids "+clientsName)
	ErrListByLastIDClients   = errcode.NewError(clientsBaseCode+9, "failed to list by last id "+clientsName)

	// error codes are globally unique, adding 1 to the previous error code
)
