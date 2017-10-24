package api

import (
	"encoding/json"
	"time"
)

const (
	prefixSearchSPC = "search:spc"
	prefixSearchINN = "search:inn"
	prefixSearchORG = "search:org"
	prefixSearchATC = "search:atc"
	prefixSearchACT = "search:act"
)

type spec struct {
	ID   int64   `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Sale float64 `json:"sale,omitempty"`
	Text bool    `json:"text,omitempty"`
}

type item struct {
	ID   int64   `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Spec []*spec `json:"spec,omitempty"`
}

type result struct {
	INF []*item `json:"inf,omitempty"`
	INN []*item `json:"inn,omitempty"`
	ATC []*item `json:"atc,omitempty"`
	ACT []*item `json:"act,omitempty"`
	ORG []*item `json:"org,omitempty"`
}

func jsonToString(data []byte) (string, error) {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return "", err
	}
	return v, nil
}

func heatSearch(h *dbxHelper) (interface{}, error) {
	start := time.Now()

	c := h.getConn()
	defer h.delConn(c)

	var ids []int64
	var err error

	ids, err = loadSyncIDs(c, prefixClassATC, 0)
	if err != nil {
		return nil, err
	}
	atc, err := makeClasses(ids...)
	err = loadHashers(c, prefixClassATC, atc)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixClassATC, atc)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixINN, 0)
	if err != nil {
		return nil, err
	}
	inn, err := makeINNs(ids...)
	err = loadHashers(c, prefixINN, inn)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixINN, inn)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixMaker, 0)
	if err != nil {
		return nil, err
	}
	org, err := makeMakers(ids...)
	err = loadHashers(c, prefixMaker, org)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixMaker, org)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixSpecACT, 0)
	if err != nil {
		return nil, err
	}
	act, err := makeSpecs(ids...)
	err = loadHashers(c, prefixSpecACT, act)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecACT, act)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixSpecINF, 0)
	if err != nil {
		return nil, err
	}
	inf, err := makeSpecs(ids...)
	err = loadHashers(c, prefixSpecINF, inf)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecINF, inf)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixSpecDEC, 0)
	if err != nil {
		return nil, err
	}
	dec, err := makeSpecs(ids...)
	err = loadHashers(c, prefixSpecDEC, dec)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecDEC, dec)
	if err != nil {
		return nil, err
	}

	return time.Since(start).String(), nil
}

func listSugg(h *dbxHelper) (interface{}, error) {
	s, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	//var wg sync.WaitGroup
	//wg.Add(1)
	//wg.Wait()

	return s, nil
}

func findSugg(h *dbxHelper) (interface{}, error) {
	s, err := jsonToString(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return s, nil
}
