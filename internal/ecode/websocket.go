package ecode

import "github.com/zhufuyi/sponge/pkg/errcode"

var (
	websocketNo      = 88
	websockeName     = "websocket"
	websockeBaseCode = errcode.HCode(websocketNo)

	ErrClientNotFound = errcode.NewError(websockeBaseCode+1, "current user's client is not found")
)
