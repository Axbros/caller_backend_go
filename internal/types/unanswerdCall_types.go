package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUnanswerdCallRequest request params
type CreateUnanswerdCallRequest struct {
	MachineId    string `json:"machine_id" binding:""`
	MobileNumber string `json:"number" binding:""`
	Type         string `json:"type" binding:""`
}

// UpdateUnanswerdCallByIDRequest request params
type UpdateUnanswerdCallByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	ClientMachineCode string `json:"clientMachineCode" binding:""`
	MobileNumber      string `json:"mobileNumber" binding:""`
}

// UnanswerdCallObjDetail detail
type UnanswerdCallObjDetail struct {
	ID string `json:"id"` // convert to string id

	ClientMachineCode string    `json:"clientMachineCode"`
	MobileNumber      string    `json:"mobileNumber"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// CreateUnanswerdCallRespond only for api docs
type CreateUnanswerdCallRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUnanswerdCallByIDRespond only for api docs
type UpdateUnanswerdCallByIDRespond struct {
	Result
}

// GetUnanswerdCallByIDRespond only for api docs
type GetUnanswerdCallByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UnanswerdCall UnanswerdCallObjDetail `json:"callLog"`
	} `json:"data"` // return data
}

// DeleteUnanswerdCallByIDRespond only for api docs
type DeleteUnanswerdCallByIDRespond struct {
	Result
}

// DeleteUnanswerdCallsByIDsRespond only for api docs
type DeleteUnanswerdCallsByIDsRespond struct {
	Result
}

// ListUnanswerdCallsRequest request params
type ListUnanswerdCallsRequest struct {
	query.Params
}

// ListUnanswerdCallsRespond only for api docs
type ListUnanswerdCallsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UnanswerdCalls []UnanswerdCallObjDetail `json:"callLogs"`
	} `json:"data"` // return data
}

// DeleteUnanswerdCallsByIDsRequest request params
type DeleteUnanswerdCallsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUnanswerdCallByConditionRequest request params
type GetUnanswerdCallByConditionRequest struct {
	query.Conditions
}

// GetUnanswerdCallByConditionRespond only for api docs
type GetUnanswerdCallByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UnanswerdCall UnanswerdCallObjDetail `json:"callLog"`
	} `json:"data"` // return data
}

// ListUnanswerdCallsByIDsRequest request params
type ListUnanswerdCallsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUnanswerdCallsByIDsRespond only for api docs
type ListUnanswerdCallsByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UnanswerdCalls []UnanswerdCallObjDetail `json:"callLogs"`
	} `json:"data"` // return data
}
