package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// callHistory business-level http error codes.
// the callHistoryNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	callHistoryNO       = 36
	callHistoryName     = "callHistory"
	callHistoryBaseCode = errcode.HCode(callHistoryNO)

	ErrCreateCallHistory     = errcode.NewError(callHistoryBaseCode+1, "failed to create "+callHistoryName)
	ErrDeleteByIDCallHistory = errcode.NewError(callHistoryBaseCode+2, "failed to delete "+callHistoryName)
	ErrUpdateByIDCallHistory = errcode.NewError(callHistoryBaseCode+3, "failed to update "+callHistoryName)
	ErrGetByIDCallHistory    = errcode.NewError(callHistoryBaseCode+4, "failed to get "+callHistoryName+" details")
	ErrListCallHistory       = errcode.NewError(callHistoryBaseCode+5, "failed to list of "+callHistoryName)

	ErrDeleteByIDsCallHistory    = errcode.NewError(callHistoryBaseCode+6, "failed to delete by batch ids "+callHistoryName)
	ErrGetByConditionCallHistory = errcode.NewError(callHistoryBaseCode+7, "failed to get "+callHistoryName+" details by conditions")
	ErrListByIDsCallHistory      = errcode.NewError(callHistoryBaseCode+8, "failed to list by batch ids "+callHistoryName)
	ErrListByLastIDCallHistory   = errcode.NewError(callHistoryBaseCode+9, "failed to list by last id "+callHistoryName)

	// error codes are globally unique, adding 1 to the previous error code
)
