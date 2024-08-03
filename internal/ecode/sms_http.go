package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// sms business-level http error codes.
// the smsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	smsNO       = 44
	smsName     = "sms"
	smsBaseCode = errcode.HCode(smsNO)

	ErrCreateSms     = errcode.NewError(smsBaseCode+1, "failed to create "+smsName)
	ErrDeleteByIDSms = errcode.NewError(smsBaseCode+2, "failed to delete "+smsName)
	ErrUpdateByIDSms = errcode.NewError(smsBaseCode+3, "failed to update "+smsName)
	ErrGetByIDSms    = errcode.NewError(smsBaseCode+4, "failed to get "+smsName+" details")
	ErrListSms       = errcode.NewError(smsBaseCode+5, "failed to list of "+smsName)

	ErrDeleteByIDsSms    = errcode.NewError(smsBaseCode+6, "failed to delete by batch ids "+smsName)
	ErrGetByConditionSms = errcode.NewError(smsBaseCode+7, "failed to get "+smsName+" details by conditions")
	ErrListByIDsSms      = errcode.NewError(smsBaseCode+8, "failed to list by batch ids "+smsName)
	ErrListByLastIDSms   = errcode.NewError(smsBaseCode+9, "failed to list by last id "+smsName)

	// error codes are globally unique, adding 1 to the previous error code
)
