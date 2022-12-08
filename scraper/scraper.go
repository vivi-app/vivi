package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
)

type Scraper struct {
	state *lua.LState

	functionSearch   *lua.LFunction
	functionExplore  *lua.LFunction
	functionPrepare  *lua.LFunction
	functionStream   *lua.LFunction
	functionDownload *lua.LFunction
	progress         *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.functionSearch != nil
}

func (s *Scraper) HasExplore() bool {
	return s.functionExplore != nil
}

func (s *Scraper) SetProgress(progress func(string)) {
	s.progress = s.state.NewFunction(func(L *lua.LState) int {
		progress(L.ToString(1))
		return 0
	})
}

func New(L *lua.LState, r io.Reader) (*Scraper, error) {
	lfunc, err := L.Load(r, Module)
	if err != nil {
		return nil, err
	}

	L.Push(lfunc)

	err = L.CallByParam(lua.P{
		Fn:      lfunc,
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return nil, err
	}

	module := L.Get(-1)
	theScraper := &Scraper{}

	table, ok := module.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("scraper module must return a table, got %s", module.Type().String())
	}

	errorNotAFunction := func(name string, val lua.LValue) error {
		return fmt.Errorf("scraper module must return a function `%s`, got %s", name, val.Type().String())
	}

	functionSearch := table.RawGet(lua.LString(FunctionSearch))
	if functionSearch.Type() == lua.LTFunction {
		theScraper.functionSearch = functionSearch.(*lua.LFunction)
	} else if functionSearch.Type() != lua.LTNil {
		return nil, errorNotAFunction(FunctionSearch, functionSearch)
	}

	functionExplore := table.RawGet(lua.LString(FunctionExplore))
	if functionExplore.Type() == lua.LTFunction {
		theScraper.functionExplore = functionExplore.(*lua.LFunction)
	} else if functionExplore.Type() != lua.LTNil {
		return nil, errorNotAFunction(FunctionExplore, functionExplore)
	}

	if !theScraper.HasExplore() && !theScraper.HasSearch() {
		return nil, fmt.Errorf("scraper module must return at least one of the functions `%s` or `%s`", FunctionSearch, FunctionExplore)
	}

	functionPrepare := table.RawGet(lua.LString(FunctionPrepare))
	if functionPrepare.Type() == lua.LTFunction {
		theScraper.functionPrepare = functionPrepare.(*lua.LFunction)
	} else {
		return nil, errorNotAFunction(FunctionPrepare, functionPrepare)
	}

	functionStream := table.RawGet(lua.LString(FunctionStream))
	if functionStream.Type() == lua.LTFunction {
		theScraper.functionStream = functionStream.(*lua.LFunction)
	} else {
		return nil, errorNotAFunction(FunctionStream, functionStream)
	}

	functionDownload := table.RawGet(lua.LString(FunctionDownload))
	if functionDownload.Type() == lua.LTFunction {
		theScraper.functionDownload = functionDownload.(*lua.LFunction)
	} else {
		return nil, errorNotAFunction(FunctionDownload, functionDownload)
	}

	theScraper.state = L
	return theScraper, nil
}
