package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateSmsRequest request params
type CreateSmsRequest struct {
	MachineCode string `json:"machineCode" binding:""`
	Address     string `json:"address" binding:""`
	Date        string `json:"date" binding:""`
	Body        string `json:"body" binding:""`
	SmsType     string `json:"smsType" binding:""`
	From        string `json:"from" binding:""`
}

// UpdateSmsByIDRequest request params
type UpdateSmsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	MachineCode string `json:"machineCode" binding:""`
	Address     string `json:"address" binding:""`
	Date        string `json:"date" binding:""`
	Body        string `json:"body" binding:""`
	SmsType     string `json:"smsType" binding:""`
}

// SmsObjDetail detail
type SmsObjDetail struct {
	ID string `json:"id"` // convert to string id

	MachineCode string    `json:"machineCode"`
	Address     string    `json:"address"`
	Date        string    `json:"date"`
	Body        string    `json:"body"`
	SmsType     string    `json:"smsType"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateSmsRespond only for api docs
type CreateSmsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateSmsByIDRespond only for api docs
type UpdateSmsByIDRespond struct {
	Result
}

// GetSmsByIDRespond only for api docs
type GetSmsByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Sms SmsObjDetail `json:"sms"`
	} `json:"data"` // return data
}

// DeleteSmsByIDRespond only for api docs
type DeleteSmsByIDRespond struct {
	Result
}

// DeleteSmssByIDsRespond only for api docs
type DeleteSmssByIDsRespond struct {
	Result
}

// ListSmssRequest request params
type ListSmssRequest struct {
	query.Params
}

// ListSmssRespond only for api docs
type ListSmssRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Smss []SmsObjDetail `json:"smss"`
	} `json:"data"` // return data
}

// DeleteSmssByIDsRequest request params
type DeleteSmssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetSmsByConditionRequest request params
type GetSmsByConditionRequest struct {
	query.Conditions
}

// GetSmsByConditionRespond only for api docs
type GetSmsByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Sms SmsObjDetail `json:"sms"`
	} `json:"data"` // return data
}

// ListSmssByIDsRequest request params
type ListSmssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListSmssByIDsRespond only for api docs
type ListSmssByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Smss []SmsObjDetail `json:"smss"`
	} `json:"data"` // return data
}
