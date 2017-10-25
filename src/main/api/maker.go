package api

import (
	"encoding/json"
	"net/http"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixMaker = "maker"
)

type jsonMaker struct {
	ID     int64 `json:"id,omitempty"`
	IDNode int64 `json:"id_node,omitempty"`

	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?

	Name   string `json:"name,omitempty"` // *
	NameRU string `json:"name_ru,omitempty"`
	NameUA string `json:"name_ua,omitempty"`
	NameEN string `json:"name_en,omitempty"`
	Text   string `json:"text,omitempty"` // *
	TextRU string `json:"text_ru,omitempty"`
	TextUA string `json:"text_ua,omitempty"`
	TextEN string `json:"text_en,omitempty"`
	IsComp bool   `json:"is_comp,omitempty"` // *
	Logo   string `json:"logo,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

func (j *jsonMaker) getID() int64 {
	return j.ID
}

func (j *jsonMaker) getNameRU(_ string) string {
	return j.NameRU
}

func (j *jsonMaker) getNameUA(_ string) string {
	return j.NameUA
}

func (j *jsonMaker) getKey(p string) string {
	return genKey(p, j.ID)
}

func (j *jsonMaker) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id", j.ID,
		"id_node", j.IDNode,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
		"text_ru", j.TextRU,
		"text_ua", j.TextUA,
		"text_en", j.TextEN,
		"is_comp", j.IsComp,
		"logo", j.Logo,
		"slug", j.Slug,
	}
}

func (j *jsonMaker) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id",      // 0
		"id_node", // 1
		"name_ru", // 2
		"name_ua", // 3
		"name_en", // 4
		"text_ru", // 5
		"text_ua", // 6
		"text_en", // 7
		"is_comp", // 8
		"logo",    // 9
		"slug",    // 10
	}
}

func (j *jsonMaker) setValues(v ...interface{}) bool {
	for i := range v {
		switch i {
		case 0:
			j.ID, _ = redis.Int64(v[i], nil)
		case 1:
			j.IDNode, _ = redis.Int64(v[i], nil)
		case 2:
			j.NameRU, _ = redis.String(v[i], nil)
		case 3:
			j.NameUA, _ = redis.String(v[i], nil)
		case 4:
			j.NameEN, _ = redis.String(v[i], nil)
		case 5:
			j.TextRU, _ = redis.String(v[i], nil)
		case 6:
			j.TextUA, _ = redis.String(v[i], nil)
		case 7:
			j.TextEN, _ = redis.String(v[i], nil)
		case 8:
			j.IsComp, _ = redis.Bool(v[i], nil)
		case 9:
			j.Logo, _ = redis.String(v[i], nil)
		case 10:
			j.Slug, _ = redis.String(v[i], nil)
		}
	}
	return j.ID != 0
}

type jsonMakers []*jsonMaker

func (j jsonMakers) len() int {
	return len(j)
}

func (j jsonMakers) elem(i int) interface{} {
	return j[i]
}

func (j jsonMakers) nill(i int) {
	j[i] = nil
}

func jsonToMakers(data []byte) (jsonMakers, error) {
	var v []*jsonMaker
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonMakers(v), nil
}

func jsonToMakersFromIDs(data []byte) (jsonMakers, error) {
	v, err := jsonToIDs(data)
	if err != nil {
		return nil, err
	}
	return makeMakers(v...)
}

func makeMakers(x ...int64) (jsonMakers, error) {
	v := make([]*jsonMaker, len(x))
	for i := range v {
		v[i] = &jsonMaker{ID: x[i]}
	}
	return jsonMakers(v), nil
}

func getMakerXSync(h *dbxHelper, p string, d ...bool) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getMakerX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToMakersFromIDs(h.data)
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

func setMakerX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToMakers(h.data)
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

func delMakerX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToMakersFromIDs(h.data)
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

// MAKER

func getMakerSync(h *dbxHelper) (interface{}, error) {
	return getINNXSync(h, prefixMaker)
}

func getMakerSyncDel(h *dbxHelper) (interface{}, error) {
	return getINNXSync(h, prefixMaker)
}

func getMaker(h *dbxHelper) (interface{}, error) {
	return getINNX(h, prefixMaker)
}

func setMaker(h *dbxHelper) (interface{}, error) {
	return setINNX(h, prefixMaker)
}

func delMaker(h *dbxHelper) (interface{}, error) {
	return delINNX(h, prefixMaker)
}
