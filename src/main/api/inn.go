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

func (j *jsonINN) getKey(p string) string {
	return genKey(p, j.ID)
}

func (j *jsonINN) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id", j.ID,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
		"slug", j.Slug,
	}
}

func (j *jsonINN) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id",      // 0
		"name_ru", // 1
		"name_ua", // 2
		"name_en", // 3
		"slug",    // 4
	}
}

func (j *jsonINN) setValues(v ...interface{}) bool {
	for i := range v {
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
	return j.ID != 0
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
	return makeINNs(v...)
}

func makeINNs(x ...int64) (jsonINNs, error) {
	v := make([]*jsonINN, len(x))
	for i := range v {
		v[i] = &jsonINN{ID: x[i]}
	}
	return jsonINNs(v), nil
}

func getINNXSync(h *dbxHelper, p string, d ...bool) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getINNX(h *dbxHelper, p string) (interface{}, error) {
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
