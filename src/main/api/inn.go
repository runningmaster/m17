package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixINN = "inn"
)

type jsonINN struct {
	ID        int64   `json:"id,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?
	Name      string  `json:"name,omitempty"`        // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

func (j *jsonINN) getKey(p string) string {
	return p + ":" + strconv.Itoa(int(j.ID))
}

func (j *jsonINN) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		p + ":" + "sync",
		"CH",
		time.Now().Unix(),
		j.ID,
	}
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

func (j *jsonINN) setValues(v ...interface{}) {
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
}

type jsonINNs []*jsonINN

func (j jsonINNs) len() int {
	return len(j)
}

func (j jsonINNs) elem(i int) hasher {
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

func makeINNs(v ...int64) jsonINNs {
	out := make([]*jsonINN, len(v))
	for i := range out {
		out[i].ID = v[i]
	}

	return jsonINNs(out)
}

func getINN(h *dbxHelper) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	out := makeINNs(v...)
	err = loadHashers(c, prefixINN, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getINNSync(h *dbxHelper) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	s, err := loadSyncIDs(c, prefixINN, v)
	if err != nil {
		return nil, err
	}

	out := makeINNs(s...)
	err = loadHashers(c, prefixINN, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func setINN(h *dbxHelper) (interface{}, error) {
	v, err := jsonToINNs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = saveHashers(c, prefixINN, v)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func delINN(h *dbxHelper) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = freeHashers(c, prefixINN, makeINNs(v...))
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}
