package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUserRequest request params
type CreateUserRequest struct {
	MachineCode string `json:"machineCode" binding:""`
}

// UpdateUserByIDRequest request params
type UpdateUserByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	MachineCode string `json:"machineCode" binding:""`
}

// UserObjDetail detail
type UserObjDetail struct {
	ID string `json:"id"` // convert to string id

	MachineCode string    `json:"machineCode"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateUserRespond only for api docs
type CreateUserRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUserByIDRespond only for api docs
type UpdateUserByIDRespond struct {
	Result
}

// GetUserByIDRespond only for api docs
type GetUserByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		User UserObjDetail `json:"user"`
	} `json:"data"` // return data
}

// DeleteUserByIDRespond only for api docs
type DeleteUserByIDRespond struct {
	Result
}

// DeleteUsersByIDsRespond only for api docs
type DeleteUsersByIDsRespond struct {
	Result
}

// ListUsersRequest request params
type ListUsersRequest struct {
	query.Params
}

// ListUsersRespond only for api docs
type ListUsersRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users []UserObjDetail `json:"users"`
	} `json:"data"` // return data
}

// DeleteUsersByIDsRequest request params
type DeleteUsersByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUserByConditionRequest request params
type GetUserByConditionRequest struct {
	query.Conditions
}

// GetUserByConditionRespond only for api docs
type GetUserByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		User UserObjDetail `json:"user"`
	} `json:"data"` // return data
}

// ListUsersByIDsRequest request params
type ListUsersByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUsersByIDsRespond only for api docs
type ListUsersByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users []UserObjDetail `json:"users"`
	} `json:"data"` // return data
}
