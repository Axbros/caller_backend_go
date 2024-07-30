package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateClientsRequest request params
type CreateClientsRequest struct {
	MachineCode string `json:"machineCode" binding:""`
	IPAddress   string `json:"ipAddress" binding:""`
}

// UpdateClientsByIDRequest request params
type UpdateClientsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	MachineCode string `json:"machineCode" binding:""`
	IPAddress   string `json:"ipAddress" binding:""`
}

// ClientsObjDetail detail
type ClientsObjDetail struct {
	ID string `json:"id"` // convert to string id

	MachineCode string    `json:"machineCode"`
	IPAddress   string    `json:"ipAddress"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateClientsRespond only for api docs
type CreateClientsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateClientsByIDRespond only for api docs
type UpdateClientsByIDRespond struct {
	Result
}

// GetClientsByIDRespond only for api docs
type GetClientsByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Clients ClientsObjDetail `json:"clients"`
	} `json:"data"` // return data
}

// DeleteClientsByIDRespond only for api docs
type DeleteClientsByIDRespond struct {
	Result
}

// DeleteClientssByIDsRespond only for api docs
type DeleteClientssByIDsRespond struct {
	Result
}

// ListClientssRequest request params
type ListClientssRequest struct {
	query.Params
}

// ListClientssRespond only for api docs
type ListClientssRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Clientss []ClientsObjDetail `json:"clientss"`
	} `json:"data"` // return data
}

// DeleteClientssByIDsRequest request params
type DeleteClientssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetClientsByConditionRequest request params
type GetClientsByConditionRequest struct {
	query.Conditions
}

// GetClientsByConditionRespond only for api docs
type GetClientsByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Clients ClientsObjDetail `json:"clients"`
	} `json:"data"` // return data
}

// ListClientssByIDsRequest request params
type ListClientssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListClientssByIDsRespond only for api docs
type ListClientssByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Clientss []ClientsObjDetail `json:"clientss"`
	} `json:"data"` // return data
}
