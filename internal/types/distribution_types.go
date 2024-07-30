package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateDistributionRequest request params
type CreateDistributionRequest struct {
	UserID      int `json:"userId" binding:""`
	GroupCallID int `json:"groupCallId" binding:""`
}

// UpdateDistributionByIDRequest request params
type UpdateDistributionByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID      int `json:"userId" binding:""`
	GroupCallID int `json:"groupCallId" binding:""`
}

// DistributionObjDetail detail
type DistributionObjDetail struct {
	ID string `json:"id"` // convert to string id

	UserID      int       `json:"userId"`
	GroupCallID int       `json:"groupCallId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateDistributionRespond only for api docs
type CreateDistributionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateDistributionByIDRespond only for api docs
type UpdateDistributionByIDRespond struct {
	Result
}

// GetDistributionByIDRespond only for api docs
type GetDistributionByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Distribution DistributionObjDetail `json:"distribution"`
	} `json:"data"` // return data
}

// DeleteDistributionByIDRespond only for api docs
type DeleteDistributionByIDRespond struct {
	Result
}

// DeleteDistributionsByIDsRespond only for api docs
type DeleteDistributionsByIDsRespond struct {
	Result
}

// ListDistributionsRequest request params
type ListDistributionsRequest struct {
	query.Params
}

// ListDistributionsRespond only for api docs
type ListDistributionsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Distributions []DistributionObjDetail `json:"distributions"`
	} `json:"data"` // return data
}

// DeleteDistributionsByIDsRequest request params
type DeleteDistributionsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetDistributionByConditionRequest request params
type GetDistributionByConditionRequest struct {
	query.Conditions
}

// GetDistributionByConditionRespond only for api docs
type GetDistributionByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Distribution DistributionObjDetail `json:"distribution"`
	} `json:"data"` // return data
}

// ListDistributionsByIDsRequest request params
type ListDistributionsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListDistributionsByIDsRespond only for api docs
type ListDistributionsByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Distributions []DistributionObjDetail `json:"distributions"`
	} `json:"data"` // return data
}
