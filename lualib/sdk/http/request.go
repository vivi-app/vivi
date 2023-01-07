package http

import (
	lua "github.com/awirix/lua"
	"io"
	"net/http"
	"strings"
)

const requestTypeName = httpTypeName + "_request"

func checkRequest(L *lua.LState, idx int) *http.Request {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*http.Request); ok {
		return v
	}
	L.ArgError(1, "request expected")
	return nil
}

func pushRequest(L *lua.LState, request *http.Request) {
	ud := L.NewUserData()
	ud.Value = request
	L.SetMetatable(ud, L.GetTypeMetatable(requestTypeName))
	L.Push(ud)
}

func newRequest(L *lua.LState) int {
	method := L.CheckString(1)
	url := L.CheckString(2)
	body := L.OptString(3, "")

	var reqBody io.Reader

	if body != "" {
		reqBody = strings.NewReader(body)
	}

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushRequest(L, request)
	return 1
}

func requestSetHeader(L *lua.LState) int {
	request := checkRequest(L, 1)
	header := checkHeader(L, 2)
	request.Header = *header
	return 0
}

func requestGetHeader(L *lua.LState) int {
	request := checkRequest(L, 1)
	pushHeader(L, &request.Header)
	return 1
}

func requestSetBasicAuth(L *lua.LState) int {
	request := checkRequest(L, 1)
	username := L.CheckString(2)
	password := L.CheckString(3)
	request.SetBasicAuth(username, password)
	return 0
}
