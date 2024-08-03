package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateCallHistoryRequest request params
type CreateCallHistoryRequest struct {
	RequestMachineCode string `json:"requestMachineCode" binding:""`
	ClientMachineCode  string `json:"clientMachineCode" binding:""`
	MobileNumber       string `json:"mobileNumber" binding:""`
	Instruction        string `json:"instruction" binding:""`
}

// UpdateCallHistoryByIDRequest request params
type UpdateCallHistoryByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	RequestMachineCode string `json:"requestMachineCode" binding:""`
	ClientMachineCode  string `json:"clientMachineCode" binding:""`
	MobileNumber       string `json:"mobileNumber" binding:""`
	Instruction        string `json:"instruction" binding:""`
}

// CallHistoryObjDetail detail
type CallHistoryObjDetail struct {
	ID string `json:"id"` // convert to string id

	RequestMachineCode string    `json:"requestMachineCode"`
	ClientMachineCode  string    `json:"clientMachineCode"`
	MobileNumber       string    `json:"mobileNumber"`
	Instruction        string    `json:"instruction"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// CreateCallHistoryRespond only for api docs
type CreateCallHistoryRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateCallHistoryByIDRespond only for api docs
type UpdateCallHistoryByIDRespond struct {
	Result
}

// GetCallHistoryByIDRespond only for api docs
type GetCallHistoryByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CallHistory CallHistoryObjDetail `json:"callHistory"`
	} `json:"data"` // return data
}

// DeleteCallHistoryByIDRespond only for api docs
type DeleteCallHistoryByIDRespond struct {
	Result
}

// DeleteCallHistorysByIDsRespond only for api docs
type DeleteCallHistorysByIDsRespond struct {
	Result
}

// ListCallHistorysRequest request params
type ListCallHistorysRequest struct {
	query.Params
}

// ListCallHistorysRespond only for api docs
type ListCallHistorysRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CallHistorys []CallHistoryObjDetail `json:"callHistorys"`
	} `json:"data"` // return data
}

// DeleteCallHistorysByIDsRequest request params
type DeleteCallHistorysByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetCallHistoryByConditionRequest request params
type GetCallHistoryByConditionRequest struct {
	query.Conditions
}

// GetCallHistoryByConditionRespond only for api docs
type GetCallHistoryByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CallHistory CallHistoryObjDetail `json:"callHistory"`
	} `json:"data"` // return data
}

// ListCallHistorysByIDsRequest request params
type ListCallHistorysByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListCallHistorysByIDsRespond only for api docs
type ListCallHistorysByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CallHistorys []CallHistoryObjDetail `json:"callHistorys"`
	} `json:"data"` // return data
}
