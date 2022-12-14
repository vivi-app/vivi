package http

import (
	lua "github.com/awirix/lua"
	"net/http"
)

const headerTypeName = httpTypeName + "_header"

func pushHeader(L *lua.LState, header *http.Header) {
	ud := L.NewUserData()
	ud.Value = header
	L.SetMetatable(ud, L.GetTypeMetatable(headerTypeName))
	L.Push(ud)
}

func checkHeader(L *lua.LState, n int) *http.Header {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*http.Header); ok {
		return v
	}
	L.ArgError(n, "header expected")
	return nil
}

func newHeader(L *lua.LState) int {
	header := http.Header{}
	pushHeader(L, &header)
	return 1
}

func headerGet(L *lua.LState) int {
	header := checkHeader(L, 1)
	key := L.CheckString(2)

	value := header.Get(key)
	L.Push(lua.LString(value))
	return 1
}

func headerSet(L *lua.LState) int {
	header := checkHeader(L, 1)
	key := L.CheckString(2)
	value := L.CheckString(3)

	header.Set(key, value)
	return 0
}

func headerAdd(L *lua.LState) int {
	header := checkHeader(L, 1)
	key := L.CheckString(2)
	value := L.CheckString(3)

	header.Add(key, value)
	return 0
}

func headerDel(L *lua.LState) int {
	header := checkHeader(L, 1)
	key := L.CheckString(2)

	header.Del(key)
	return 0
}

func headerClone(L *lua.LState) int {
	header := checkHeader(L, 1)

	clone := header.Clone()
	pushHeader(L, &clone)
	return 1
}
