package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// callLog business-level http error codes.
// the callLogNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	callLogNO       = 47
	callLogName     = "callLog"
	callLogBaseCode = errcode.HCode(callLogNO)

	ErrCreateUnanswerdCall     = errcode.NewError(callLogBaseCode+1, "failed to create "+callLogName)
	ErrDeleteByIDUnanswerdCall = errcode.NewError(callLogBaseCode+2, "failed to delete "+callLogName)
	ErrUpdateByIDUnanswerdCall = errcode.NewError(callLogBaseCode+3, "failed to update "+callLogName)
	ErrGetByIDUnanswerdCall    = errcode.NewError(callLogBaseCode+4, "failed to get "+callLogName+" details")
	ErrListUnanswerdCall       = errcode.NewError(callLogBaseCode+5, "failed to list of "+callLogName)

	ErrDeleteByIDsUnanswerdCall    = errcode.NewError(callLogBaseCode+6, "failed to delete by batch ids "+callLogName)
	ErrGetByConditionUnanswerdCall = errcode.NewError(callLogBaseCode+7, "failed to get "+callLogName+" details by conditions")
	ErrListByIDsUnanswerdCall      = errcode.NewError(callLogBaseCode+8, "failed to list by batch ids "+callLogName)
	ErrListByLastIDUnanswerdCall   = errcode.NewError(callLogBaseCode+9, "failed to list by last id "+callLogName)

	// error codes are globally unique, adding 1 to the previous error code
)
