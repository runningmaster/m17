package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixClassATC = "class:atc"
	prefixClassNFC = "class:nfc"
	prefixClassFSC = "class:fsc"
	prefixClassBFC = "class:bfc"
	prefixClassCFC = "class:cfc"
	prefixClassMPC = "class:mpc"
	prefixClassCSC = "class:csc"
	prefixClassICD = "class:icd"
)

type jsonClass struct {
	ID        int64   `json:"id,omitempty"`
	IDNode    int64   `json:"id_node,omitempty"`
	IDRoot    int64   `json:"id_root,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?
	Code      string  `json:"code,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
}

func (j *jsonClass) getKey(p string) string {
	return p + ":" + strconv.Itoa(int(j.ID))
}

func (j *jsonClass) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		p + ":" + "sync",
		"CH",
		time.Now().Unix(),
		j.ID,
	}
}

func (j *jsonClass) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id", j.ID,
		"id_node", j.IDNode,
		"id_root", j.IDRoot,
		"code", j.Code,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
	}
}

func (j *jsonClass) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id",      // 0
		"id_node", // 1
		"id_root", // 2
		"code",    // 3
		"name_ru", // 4
		"name_ua", // 5
		"name_en", // 6
	}
}

func (j *jsonClass) setValues(v ...interface{}) {
	for i := range v {
		switch i {
		case 0:
			j.ID, _ = redis.Int64(v[i], nil)
		case 1:
			j.IDNode, _ = redis.Int64(v[i], nil)
		case 2:
			j.IDRoot, _ = redis.Int64(v[i], nil)
		case 3:
			j.Code, _ = redis.String(v[i], nil)
		case 4:
			j.NameRU, _ = redis.String(v[i], nil)
		case 5:
			j.NameUA, _ = redis.String(v[i], nil)
		case 6:
			j.NameEN, _ = redis.String(v[i], nil)
		}
	}
}

func getclass(c redis.Conn, p string, v ...int64) ([]*jsonClass, error) {
	out := make([]*jsonClass, len(v))
	for i := range out {
		out[i].ID = v[i]
	}

	var err error
	for i := range out {
		err = c.Send("HMGET", out[i].getKeyAndFields(p)...)
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	var r []interface{}
	for i := range out {
		r, err = redis.Values(c.Receive())
		if err != nil {
			if err == redis.ErrNil {
				out[i] = nil
				continue
			}
			return nil, err
		}
		out[i].setValues(r)
	}

	return out, nil
}

func setclass(c redis.Conn, p string, v ...*jsonClass) error {
	var err error
	for i := range v {
		err = c.Send("HMSET", v[i].getKeyAndFieldValues(p)...)
		if err != nil {
			return err
		}
		err = c.Send("ZADD", v[i].getKeyAndUnixtimeID(p)...)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func delclass(c redis.Conn, p string, v ...int64) error {
	out := make([]*jsonClass, len(v))
	for i := range out {
		out[i].ID = v[i]
	}

	var err error
	for i := range out {
		err = c.Send("DEL", out[i].getKey(p))
		if err != nil {
			return err
		}
		err = c.Send("ZADD", out[i].getKeyAndUnixtimeID(p)...)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func getClass(h *dbxHelper, p string) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return getclass(c, p, v...)
}

func getClassSync(h *dbxHelper, p string) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	l, err := h.getSyncList(p, v)
	if err != nil {
		return nil, err
	}

	return getclass(c, p, l...)
}

func setClass(h *dbxHelper, p string) (interface{}, error) {
	var v []*jsonClass
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = setclass(c, p, v...)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}

func delClass(h *dbxHelper, p string) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = delclass(c, p, v...)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}

func getClassATC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassATC)
}

func getClassATCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassATC)
}

func setClassATC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassATC)
}

func delClassATC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassATC)
}

func getClassNFC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassNFC)
}

func getClassNFCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassNFC)
}

func setClassNFC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassNFC)
}

func delClassNFC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassNFC)
}

func getClassFSC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassFSC)
}

func getClassFSCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassFSC)
}

func setClassFSC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassFSC)
}

func delClassFSC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassFSC)
}

func getClassBFC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassBFC)
}

func getClassBFCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassBFC)
}

func setClassBFC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassBFC)
}

func delClassBFC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassBFC)
}

func getClassCFC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassCFC)
}

func getClassCFCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassCFC)
}

func setClassCFC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassCFC)
}

func delClassCFC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassCFC)
}

func getClassMPC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassMPC)
}

func getClassMPCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassMPC)
}

func setClassMPC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassMPC)
}

func delClassMPC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassMPC)
}

func getClassCSC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassCSC)
}

func getClassCSCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassCSC)
}

func setClassCSC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassCSC)
}

func delClassCSC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassCSC)
}

func getClassICD(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassICD)
}

func getClassICDSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassICD)
}

func setClassICD(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassICD)
}

func delClassICD(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassICD)
}
