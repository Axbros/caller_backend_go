package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// unanswerdCall business-level http error codes.
// the unanswerdCallNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	unanswerdCallNO       = 47
	unanswerdCallName     = "unanswerdCall"
	unanswerdCallBaseCode = errcode.HCode(unanswerdCallNO)

	ErrCreateUnanswerdCall     = errcode.NewError(unanswerdCallBaseCode+1, "failed to create "+unanswerdCallName)
	ErrDeleteByIDUnanswerdCall = errcode.NewError(unanswerdCallBaseCode+2, "failed to delete "+unanswerdCallName)
	ErrUpdateByIDUnanswerdCall = errcode.NewError(unanswerdCallBaseCode+3, "failed to update "+unanswerdCallName)
	ErrGetByIDUnanswerdCall    = errcode.NewError(unanswerdCallBaseCode+4, "failed to get "+unanswerdCallName+" details")
	ErrListUnanswerdCall       = errcode.NewError(unanswerdCallBaseCode+5, "failed to list of "+unanswerdCallName)

	ErrDeleteByIDsUnanswerdCall    = errcode.NewError(unanswerdCallBaseCode+6, "failed to delete by batch ids "+unanswerdCallName)
	ErrGetByConditionUnanswerdCall = errcode.NewError(unanswerdCallBaseCode+7, "failed to get "+unanswerdCallName+" details by conditions")
	ErrListByIDsUnanswerdCall      = errcode.NewError(unanswerdCallBaseCode+8, "failed to list by batch ids "+unanswerdCallName)
	ErrListByLastIDUnanswerdCall   = errcode.NewError(unanswerdCallBaseCode+9, "failed to list by last id "+unanswerdCallName)

	// error codes are globally unique, adding 1 to the previous error code
)
