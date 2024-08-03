package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateGroupCallRequest request params
type CreateGroupCallRequest struct {
	GroupNumber      string `json:"groupNumber" binding:""`
	PhoneNumber      string `json:"phoneNumber" binding:""`
	TransferClientID string `json:"transferClientId" binding:""`
}

// UpdateGroupCallByIDRequest request params
type UpdateGroupCallByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	GroupNumber      string `json:"groupNumber" binding:""`
	PhoneNumber      string `json:"phoneNumber" binding:""`
	TransferClientID string `json:"transferClientId" binding:""`
}

// GroupCallObjDetail detail
type GroupCallObjDetail struct {
	ID string `json:"id"` // convert to string id

	GroupNumber      string    `json:"groupNumber"`
	PhoneNumber      string    `json:"phoneNumber"`
	TransferClientID string    `json:"transferClientId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// CreateGroupCallRespond only for api docs
type CreateGroupCallRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateGroupCallByIDRespond only for api docs
type UpdateGroupCallByIDRespond struct {
	Result
}

// GetGroupCallByIDRespond only for api docs
type GetGroupCallByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupCall GroupCallObjDetail `json:"groupCall"`
	} `json:"data"` // return data
}

// DeleteGroupCallByIDRespond only for api docs
type DeleteGroupCallByIDRespond struct {
	Result
}

// DeleteGroupCallsByIDsRespond only for api docs
type DeleteGroupCallsByIDsRespond struct {
	Result
}

// ListGroupCallsRequest request params
type ListGroupCallsRequest struct {
	query.Params
}

// ListGroupCallsRespond only for api docs
type ListGroupCallsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupCalls []GroupCallObjDetail `json:"groupCalls"`
	} `json:"data"` // return data
}

// DeleteGroupCallsByIDsRequest request params
type DeleteGroupCallsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetGroupCallByConditionRequest request params
type GetGroupCallByConditionRequest struct {
	query.Conditions
}

// GetGroupCallByConditionRespond only for api docs
type GetGroupCallByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupCall GroupCallObjDetail `json:"groupCall"`
	} `json:"data"` // return data
}

// ListGroupCallsByIDsRequest request params
type ListGroupCallsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListGroupCallsByIDsRespond only for api docs
type ListGroupCallsByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupCalls []GroupCallObjDetail `json:"groupCalls"`
	} `json:"data"` // return data
}
