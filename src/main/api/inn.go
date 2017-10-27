package api

import (
	"encoding/json"
	"net/http"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixINN = "inn"
)

type jsonINN struct {
	ID int64 `json:"id,omitempty"`

	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?

	Name   string `json:"name,omitempty"` // *
	NameRU string `json:"name_ru,omitempty"`
	NameUA string `json:"name_ua,omitempty"`
	NameEN string `json:"name_en,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

func (j *jsonINN) getID() int64 {
	return j.ID
}

func (j *jsonINN) getNameRU(_ string) string {
	return j.NameRU
}

func (j *jsonINN) getNameUA(_ string) string {
	return j.NameUA
}

func (j *jsonINN) getFields() []interface{} {
	return []interface{}{
		"id",      // 0
		"name_ru", // 1
		"name_ua", // 2
		"name_en", // 3
		"slug",    // 4
	}
}

func (j *jsonINN) getValues() []interface{} {
	return []interface{}{
		j.ID,     // 0
		j.NameRU, // 1
		j.NameUA, // 2
		j.NameEN, // 3
		j.Slug,   // 4
	}
}

func (j *jsonINN) setValues(v ...interface{}) {
	for i := range v {
		if v[i] == nil {
			continue
		}
		switch i {
		case 0:
			j.ID, _ = redis.Int64(v[i], nil)
		case 1:
			j.NameRU, _ = redis.String(v[i], nil)
		case 2:
			j.NameUA, _ = redis.String(v[i], nil)
		case 3:
			j.NameEN, _ = redis.String(v[i], nil)
		case 4:
			j.Slug, _ = redis.String(v[i], nil)
		}
	}
}

type jsonINNs []*jsonINN

func (j jsonINNs) len() int {
	return len(j)
}

func (j jsonINNs) elem(i int) interface{} {
	return j[i]
}

func (j jsonINNs) nill(i int) {
	j[i] = nil
}

func jsonToINNs(data []byte) (jsonINNs, error) {
	var v []*jsonINN
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonINNs(v), nil
}

func jsonToINNsFromIDs(data []byte) (jsonINNs, error) {
	v, err := jsonToIDs(data)
	if err != nil {
		return nil, err
	}
	return makeINNs(v...), nil
}

func makeINNs(x ...int64) jsonINNs {
	v := make([]*jsonINN, len(x))
	for i := range v {
		v[i] = &jsonINN{ID: x[i]}
	}
	return jsonINNs(v)
}

func loadINNLinks(c redis.Conn, p string, v []*jsonINN) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		v[i].IDSpecDEC, err = loadLinkIDs(c, p, prefixSpecDEC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDSpecINF, err = loadLinkIDs(c, p, prefixSpecINF, v[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getINNXSync(h *dbxHelper, p string, d ...bool) ([]int64, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getINNX(h *dbxHelper, p string) (jsonINNs, error) {
	v, err := jsonToINNsFromIDs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = loadHashers(c, p, v)
	if err != nil {
		return nil, err
	}

	err = loadINNLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func setINNX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToINNs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = saveHashers(c, p, v)
	if err != nil {
		return nil, err
	}

	err = saveSearchers(c, p, v)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func delINNX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToINNsFromIDs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = freeHashers(c, p, v)
	if err != nil {
		return nil, err
	}

	err = freeSearchers(c, p, v)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

// INN

func getINNSync(h *dbxHelper) (interface{}, error) {
	return getINNXSync(h, prefixINN)
}

func getINNSyncDel(h *dbxHelper) (interface{}, error) {
	return getINNXSync(h, prefixINN)
}

func getINN(h *dbxHelper) (interface{}, error) {
	return getINNX(h, prefixINN)
}

func setINN(h *dbxHelper) (interface{}, error) {
	return setINNX(h, prefixINN)
}

func delINN(h *dbxHelper) (interface{}, error) {
	return delINNX(h, prefixINN)
}
