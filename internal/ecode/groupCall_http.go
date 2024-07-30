package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// groupCall business-level http error codes.
// the groupCallNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	groupCallNO       = 41
	groupCallName     = "groupCall"
	groupCallBaseCode = errcode.HCode(groupCallNO)

	ErrCreateGroupCall     = errcode.NewError(groupCallBaseCode+1, "failed to create "+groupCallName)
	ErrDeleteByIDGroupCall = errcode.NewError(groupCallBaseCode+2, "failed to delete "+groupCallName)
	ErrUpdateByIDGroupCall = errcode.NewError(groupCallBaseCode+3, "failed to update "+groupCallName)
	ErrGetByIDGroupCall    = errcode.NewError(groupCallBaseCode+4, "failed to get "+groupCallName+" details")
	ErrListGroupCall       = errcode.NewError(groupCallBaseCode+5, "failed to list of "+groupCallName)

	ErrDeleteByIDsGroupCall    = errcode.NewError(groupCallBaseCode+6, "failed to delete by batch ids "+groupCallName)
	ErrGetByConditionGroupCall = errcode.NewError(groupCallBaseCode+7, "failed to get "+groupCallName+" details by conditions")
	ErrListByIDsGroupCall      = errcode.NewError(groupCallBaseCode+8, "failed to list by batch ids "+groupCallName)
	ErrListByLastIDGroupCall   = errcode.NewError(groupCallBaseCode+9, "failed to list by last id "+groupCallName)

	// error codes are globally unique, adding 1 to the previous error code
)
