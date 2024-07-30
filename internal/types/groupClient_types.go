package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateGroupClientRequest request params
type CreateGroupClientRequest struct {
	GroupID  int `json:"groupId" binding:""`
	ClientID int `json:"clientId" binding:""`
}

// UpdateGroupClientByIDRequest request params
type UpdateGroupClientByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	GroupID  int `json:"groupId" binding:""`
	ClientID int `json:"clientId" binding:""`
}

// GroupClientObjDetail detail
type GroupClientObjDetail struct {
	ID string `json:"id"` // convert to string id

	GroupID   int       `json:"groupId"`
	ClientID  int       `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateGroupClientRespond only for api docs
type CreateGroupClientRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateGroupClientByIDRespond only for api docs
type UpdateGroupClientByIDRespond struct {
	Result
}

// GetGroupClientByIDRespond only for api docs
type GetGroupClientByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupClient GroupClientObjDetail `json:"groupClient"`
	} `json:"data"` // return data
}

// DeleteGroupClientByIDRespond only for api docs
type DeleteGroupClientByIDRespond struct {
	Result
}

// DeleteGroupClientsByIDsRespond only for api docs
type DeleteGroupClientsByIDsRespond struct {
	Result
}

// ListGroupClientsRequest request params
type ListGroupClientsRequest struct {
	query.Params
}

// ListGroupClientsRespond only for api docs
type ListGroupClientsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupClients []GroupClientObjDetail `json:"groupClients"`
	} `json:"data"` // return data
}

// DeleteGroupClientsByIDsRequest request params
type DeleteGroupClientsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetGroupClientByConditionRequest request params
type GetGroupClientByConditionRequest struct {
	query.Conditions
}

// GetGroupClientByConditionRespond only for api docs
type GetGroupClientByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupClient GroupClientObjDetail `json:"groupClient"`
	} `json:"data"` // return data
}

// ListGroupClientsByIDsRequest request params
type ListGroupClientsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListGroupClientsByIDsRespond only for api docs
type ListGroupClientsByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		GroupClients []GroupClientObjDetail `json:"groupClients"`
	} `json:"data"` // return data
}
