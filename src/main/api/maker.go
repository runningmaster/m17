package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type jsonMaker struct {
	ID        int64   `json:"id,omitempty"`
	IDNode    int64   `json:"id_node,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?
	Name      string  `json:"name,omitempty"`        // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Text      string  `json:"text,omitempty"` // *
	TextRU    string  `json:"text_ru,omitempty"`
	TextUA    string  `json:"text_ua,omitempty"`
	TextEN    string  `json:"text_en,omitempty"`
	IsComp    bool    `json:"is_comp,omitempty"`
	Logo      string  `json:"logo,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

func (j *jsonMaker) getKey(p string) string {
	return p + ":" + strconv.Itoa(int(j.ID))
}

func (j *jsonMaker) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		p + ":" + "sync",
		"CH",
		time.Now().Unix(),
		j.ID,
	}
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

func (j *jsonMaker) setValues(v ...interface{}) {
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
}

type jsonMakers []*jsonMaker

func (j jsonMakers) len() int {
	return len(j)
}

func (j jsonMakers) elem(i int) hasher {
	return j[i]
}

func (j jsonMakers) nill(i int) {
	j[i] = nil
}

func jsonToMakers(data []byte) ([]*jsonMaker, error) {
	var v []*jsonMaker
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func makeMakers(v ...int64) []*jsonMaker {
	out := make([]*jsonMaker, len(v))
	for i := range out {
		out[i].ID = v[i]
	}

	return out
}

func loadSyncIDs(c redis.Conn, p string, v int64) ([]int64, error) {
	res, err := redis.Values(c.Do("ZRANGEBYSCORE", p+":"+"sync", v, "+inf"))
	if err != nil {
		return nil, err
	}

	out := make([]int64, len(res))
	for i := range res {
		out[i], _ = redis.Int64(res[i], nil)
	}

	return out, nil
}

func saveHashers(c redis.Conn, p string, v ruleHasher) error {
	var err error
	for i := 0; i < v.len(); i++ {
		err = c.Send("HMSET", v.elem(i).getKeyAndFieldValues(p)...)
		if err != nil {
			return err
		}
		err = c.Send("ZADD", v.elem(i).getKeyAndUnixtimeID(p)...)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func loadHashers(c redis.Conn, p string, v ruleHasher) error {
	var err error
	for i := 0; i < v.len(); i++ {
		err = c.Send("HMGET", v.elem(i).getKeyAndFields(p)...)
		if err != nil {
			return err
		}
	}

	err = c.Flush()
	if err != nil {
		return err
	}

	var r []interface{}
	for i := 0; i < v.len(); i++ {
		r, err = redis.Values(c.Receive())
		if err != nil {
			if err == redis.ErrNil {
				v.nill(i)
				continue
			}
			return err
		}
		v.elem(i).setValues(r)
	}

	return nil
}

func freeHashers(c redis.Conn, p string, v ruleHasher) error {
	var err error
	for i := 0; i < v.len(); i++ {
		err = c.Send("DEL", v.elem(i).getKey(p))
		if err != nil {
			return err
		}
		err = c.Send("ZADD", v.elem(i).getKeyAndUnixtimeID(p)...)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func getMaker(h *dbxHelper) (interface{}, error) {
	x, err := jsonToInt64s(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	v := makeMakers(x...)
	err = loadHashers(c, "maker", jsonMakers(v))
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getMakerSync(h *dbxHelper) (interface{}, error) {
	x, err := jsonToInt64(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	s, err := loadSyncIDs(c, "maker", x)
	if err != nil {
		return nil, err
	}

	v := makeMakers(s...)
	err = loadHashers(c, "maker", jsonMakers(v))
	if err != nil {
		return nil, err
	}

	return v, nil
}

func setMaker(h *dbxHelper) (interface{}, error) {
	v, err := jsonToMakers(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = saveHashers(c, "maker", jsonMakers(v))
	if err != nil {
		return nil, err
	}

	return "OK", nil
}

func delMaker(h *dbxHelper) (interface{}, error) {
	p, err := jsonToInt64s(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	v := makeMakers(p...)
	err = freeHashers(c, "maker", jsonMakers(v))
	if err != nil {
		return nil, err
	}

	return "OK", nil
}
