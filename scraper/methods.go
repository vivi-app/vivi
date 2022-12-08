package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func (s *Scraper) checkMedia() (*Media, error) {
	ret := s.state.Get(-1)
	s.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value: expected 'table' got '%s'", ret.Type().String())
	}

	media, err := newMedia(table)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s *Scraper) checkMediaSlice() ([]*Media, error) {
	ret := s.state.Get(-1)
	s.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value: expected 'table' got '%s'", ret.Type().String())
	}

	var (
		items  = make([]lua.LValue, 0)
		medias = make([]*Media, 0)
	)

	table.ForEach(func(_, value lua.LValue) {
		items = append(items, value)
	})

	for _, item := range items {
		table, ok := item.(*lua.LTable)
		if !ok {
			return nil, fmt.Errorf("invalid value in returned table: expected 'table' got '%s'", item.Type().String())
		}

		media, err := newMedia(table)
		if err != nil {
			return nil, err
		}

		medias = append(medias, media)
	}

	return medias, nil
}

func (s *Scraper) Search(query string) ([]*Media, error) {
	if !s.HasSearch() {
		panic("scraper does not have a search function")
	}

	err := s.state.CallByParam(lua.P{
		Fn:      s.functionSearch,
		NRet:    1,
		Protect: true,
	}, s.progress, lua.LString(query))

	if err != nil {
		return nil, err
	}

	return s.checkMediaSlice()
}

func (s *Scraper) Explore(media *Media) ([]*Media, error) {
	if !s.HasExplore() {
		panic("scraper does not have an explore function")
	}

	err := s.state.CallByParam(lua.P{
		Fn:      s.functionExplore,
		NRet:    1,
		Protect: true,
	}, s.progress, media.Value())

	if err != nil {
		return nil, err
	}

	return s.checkMediaSlice()
}

func (s *Scraper) Prepare(media *Media) (*Media, error) {
	err := s.state.CallByParam(lua.P{
		Fn:      s.functionPrepare,
		NRet:    1,
		Protect: true,
	}, s.progress, media.Value())

	if err != nil {
		return nil, err
	}

	return s.checkMedia()
}

func (s *Scraper) Play(media *Media) error {
	err := s.state.CallByParam(lua.P{
		Fn:      s.functionStream,
		NRet:    1,
		Protect: true,
	}, s.progress, media.Value())

	if err != nil {
		return err
	}

	return nil
}

func (s *Scraper) Download(media *Media) error {
	err := s.state.CallByParam(lua.P{
		Fn:      s.functionDownload,
		NRet:    1,
		Protect: true,
	}, s.progress, media.Value())

	if err != nil {
		return err
	}

	return nil
}
