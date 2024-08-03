package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// distribution business-level http error codes.
// the distributionNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	distributionNO       = 73
	distributionName     = "distribution"
	distributionBaseCode = errcode.HCode(distributionNO)

	ErrCreateDistribution     = errcode.NewError(distributionBaseCode+1, "failed to create "+distributionName)
	ErrDeleteByIDDistribution = errcode.NewError(distributionBaseCode+2, "failed to delete "+distributionName)
	ErrUpdateByIDDistribution = errcode.NewError(distributionBaseCode+3, "failed to update "+distributionName)
	ErrGetByIDDistribution    = errcode.NewError(distributionBaseCode+4, "failed to get "+distributionName+" details")
	ErrListDistribution       = errcode.NewError(distributionBaseCode+5, "failed to list of "+distributionName)

	ErrDeleteByIDsDistribution    = errcode.NewError(distributionBaseCode+6, "failed to delete by batch ids "+distributionName)
	ErrGetByConditionDistribution = errcode.NewError(distributionBaseCode+7, "failed to get "+distributionName+" details by conditions")
	ErrListByIDsDistribution      = errcode.NewError(distributionBaseCode+8, "failed to list by batch ids "+distributionName)
	ErrListByLastIDDistribution   = errcode.NewError(distributionBaseCode+9, "failed to list by last id "+distributionName)

	// error codes are globally unique, adding 1 to the previous error code
)
