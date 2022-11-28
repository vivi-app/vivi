package html

import (
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	registerDocumentType(L)
	registerSelectionType(L)

	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"parse": parse,
	})
}
