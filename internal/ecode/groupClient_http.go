package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// groupClient business-level http error codes.
// the groupClientNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	groupClientNO       = 21
	groupClientName     = "groupClient"
	groupClientBaseCode = errcode.HCode(groupClientNO)

	ErrCreateGroupClient     = errcode.NewError(groupClientBaseCode+1, "failed to create "+groupClientName)
	ErrDeleteByIDGroupClient = errcode.NewError(groupClientBaseCode+2, "failed to delete "+groupClientName)
	ErrUpdateByIDGroupClient = errcode.NewError(groupClientBaseCode+3, "failed to update "+groupClientName)
	ErrGetByIDGroupClient    = errcode.NewError(groupClientBaseCode+4, "failed to get "+groupClientName+" details")
	ErrListGroupClient       = errcode.NewError(groupClientBaseCode+5, "failed to list of "+groupClientName)

	ErrDeleteByIDsGroupClient    = errcode.NewError(groupClientBaseCode+6, "failed to delete by batch ids "+groupClientName)
	ErrGetByConditionGroupClient = errcode.NewError(groupClientBaseCode+7, "failed to get "+groupClientName+" details by conditions")
	ErrListByIDsGroupClient      = errcode.NewError(groupClientBaseCode+8, "failed to list by batch ids "+groupClientName)
	ErrListByLastIDGroupClient   = errcode.NewError(groupClientBaseCode+9, "failed to list by last id "+groupClientName)

	// error codes are globally unique, adding 1 to the previous error code
)
