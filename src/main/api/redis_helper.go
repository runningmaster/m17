package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"internal/logger"

	"github.com/garyburd/redigo/redis"
)

var statusOK = http.StatusText(http.StatusOK)

var apiFunc = map[string]func(h *dbxHelper) (interface{}, error){
	"get-class-atc":      getClassATC,
	"get-class-atc-sync": getClassATCSync,
	"set-class-atc":      setClassATC,
	"del-class-atc":      delClassATC,

	"get-class-nfc":      getClassNFC,
	"get-class-nfc-sync": getClassNFCSync,
	"set-class-nfc":      setClassNFC,
	"del-class-nfc":      delClassNFC,

	"get-class-fsc":      getClassFSC,
	"get-class-fsc-sync": getClassFSCSync,
	"set-class-fsc":      setClassFSC,
	"del-class-fsc":      delClassFSC,

	"get-class-bfc":      getClassBFC,
	"get-class-bfc-sync": getClassBFCSync,
	"set-class-bfc":      setClassBFC,
	"del-class-bfc":      delClassBFC,

	"get-class-cfc":      getClassCFC,
	"get-class-cfc-sync": getClassCFCSync,
	"set-class-cfc":      setClassCFC,
	"del-class-cfc":      delClassCFC,

	"get-class-mpc":      getClassMPC,
	"get-class-mpc-sync": getClassMPCSync,
	"set-class-mpc":      setClassMPC,
	"del-class-mpc":      delClassMPC,

	"get-class-csc":      getClassCSC,
	"get-class-csc-sync": getClassCSCSync,
	"set-class-csc":      setClassCSC,
	"del-class-csc":      delClassCSC,

	"get-class-icd":      getClassICD,
	"get-class-icd-sync": getClassICDSync,
	"set-class-icd":      setClassICD,
	"del-class-icd":      delClassICD,

	"get-inn":      getINN,
	"get-inn-sync": getINNSync,
	"set-inn":      setINN,
	"del-inn":      delINN,

	"get-maker":      getMaker,
	"get-maker-sync": getMakerSync,
	"set-maker":      setMaker,
	"del-maker":      delMaker,

	"get-drug":      getDrug,
	"get-drug-sync": getDrugSync,
	"set-drug":      setDrug,
	"del-drug":      delDrug,
	"set-drug-sale": setDrugSale,

	"get-spec-act":      getSpecACT,
	"get-spec-act-sync": getSpecACTSync,
	"set-spec-act":      setSpecACT,
	"del-spec-act":      delSpecACT,

	"get-spec-inf":      getSpecINF,
	"get-spec-inf-sync": getSpecINFSync,
	"set-spec-inf":      setSpecINF,
	"del-spec-inf":      delSpecINF,

	"get-spec-dec":      getSpecDEC,
	"get-spec-dec-sync": getSpecDECSync,
	"set-spec-dec":      setSpecDEC,
	"del-spec-dec":      delSpecDEC,

	"list-sugg":   listSugg,
	"find-sugg":   findSugg,
	"heat-search": heatSearch,
}

type rediser interface {
	Get() redis.Conn
}

type ruler interface {
	len() int
	elem(int) interface{}
}

type hasher interface {
	getKey(string) string
	getKeyAndUnixtimeID(string) []interface{}
	getKeyAndFieldValues(string) []interface{}
	getKeyAndFields(string) []interface{}
	setValues(...interface{}) bool
}

type niller interface {
	nill(int)
}

type searcher interface {
	getID() int64
	getNameRU(string) string
	getNameUA(string) string
}

type dbxHelper struct {
	ctx  context.Context
	rdb  rediser
	log  logger.Logger
	r    *http.Request
	w    http.ResponseWriter
	meta []byte
	data []byte
	//	lang string
}

func (h *dbxHelper) getConn() redis.Conn {
	return h.rdb.Get()
}

func (h *dbxHelper) delConn(c io.Closer) {
	_ = c.Close
}

func (h *dbxHelper) ping() (interface{}, error) {
	c := h.getConn()
	defer h.delConn(c)

	return redis.Bytes(c.Do("PING"))
}

func (h *dbxHelper) exec(s string) (interface{}, error) {
	if fn, ok := apiFunc[s]; ok {
		return fn(h)
	}

	return nil, fmt.Errorf("unknown func %q", s)
}

//type jsonSale struct {
//	ID int64   `json:"id,omitempty"`
//	Q  float64 `json:"q,omitempty"`
//	V  float64 `json:"v,omitempty"`
//}

func jsonToID(data []byte) (int64, error) {
	var v int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func jsonToIDs(data []byte) ([]int64, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func freeLinkIDs(c redis.Conn, p1, p2 string, x int64, v ...int64) error {
	if len(v) == 0 {
		return nil
	}

	var key string
	var err error
	for i := range v {
		key = joinKey(genKey(p2, v[i]), p1)
		err = c.Send("SREM", key, x)
		if err != nil {
			return err
		}
	}

	key = joinKey(genKey(p1, x), p2)
	err = c.Send("DEL", key)
	if err != nil {
		return err
	}

	return c.Flush()
}

func saveLinkIDs(c redis.Conn, p1, p2 string, x int64, v ...int64) error {
	if len(v) == 0 {
		return nil
	}

	val := make([]interface{}, len(v)+1)
	var key string
	var err error
	for i := range v {
		key = joinKey(genKey(p2, v[i]), p1)
		err = c.Send("SADD", key, x)
		if err != nil {
			return err
		}
		val[i+1] = v[i]
	}

	key = joinKey(genKey(p1, x), p2)
	val[0] = key
	err = c.Send("SADD", val...)
	if err != nil {
		return err
	}

	return c.Flush()
}

func loadLinkIDs(c redis.Conn, p1, p2 string, x int64) ([]int64, error) {
	key := joinKey(genKey(p1, x), p2)
	res, err := redis.Values(c.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}

	out := make([]int64, len(res))
	for i := range res {
		out[i], _ = redis.Int64(res[i], nil)
	}

	return out, nil
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

func saveHashers(c redis.Conn, p string, v ruler) error {
	if v.len() == 0 {
		return nil
	}

	var err error
	for i := 0; i < v.len(); i++ {
		if h, ok := v.elem(i).(hasher); ok {
			err = c.Send("HMSET", h.getKeyAndFieldValues(p)...)
			if err != nil {
				return err
			}
			err = c.Send("ZADD", h.getKeyAndUnixtimeID(p)...)
			if err != nil {
				return err
			}
		}
	}

	return c.Flush()
}

func loadHashers(c redis.Conn, p string, v ruler) error {
	if v.len() == 0 {
		return nil
	}

	var err error
	for i := 0; i < v.len(); i++ {
		if h, ok := v.elem(i).(hasher); ok {
			err = c.Send("HMGET", h.getKeyAndFields(p)...)
			if err != nil {
				return err
			}
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
			return err
		}
		if h, ok := v.elem(i).(hasher); ok {
			if !h.setValues(r...) {
				if n, ok := v.elem(i).(niller); ok {
					n.nill(i)
				}
			}
		}
	}

	return nil
}

func freeHashers(c redis.Conn, p string, v ruler) error {
	if v.len() == 0 {
		return nil
	}

	var err error
	for i := 0; i < v.len(); i++ {
		if h, ok := v.elem(i).(hasher); ok {
			err = c.Send("DEL", h.getKey(p))
			if err != nil {
				return err
			}
		}
	}

	err = c.Flush()
	if err != nil {
		return err
	}

	var r bool
	for i := 0; i < v.len(); i++ {
		r, err = redis.Bool(c.Receive())
		if err != nil {
			return err
		}
		if !r {
			continue
		}
		if h, ok := v.elem(i).(hasher); ok {
			err = c.Send("ZADD", h.getKeyAndUnixtimeID(p)...)
			if err != nil {
				return err
			}
		}
	}

	return c.Flush()
}

func normName(s string) string {
	r := strings.NewReplacer(
		"®", "",
		"™", "",
		"*", "",
		"&", "",
		"†", "",
	)
	return strings.TrimSpace(strings.ToLower(r.Replace(s)))
}

func saveSearchers(c redis.Conn, p string, v ruler) error {
	var id int64
	var nameRU, nameUA string
	var err error
	for i := 0; i < v.len(); i++ {
		if s, ok := v.elem(i).(searcher); ok {
			id, nameRU, nameUA = s.getID(), s.getNameRU(p), s.getNameUA(p)
			if nameRU != "" {
				err = c.Send("ZADD", joinKey(p, "idx:ru"), id, normName(nameRU))
				if err != nil {
					return err
				}
			} else {
				if !strings.Contains(p, "spec") {
					println(joinKey(p, "idx:ru"), id, nameRU, nameUA)
				}
			}
			if nameUA != "" {
				err = c.Send("ZADD", joinKey(p, "idx:ua"), id, normName(nameUA))
				if err != nil {
					return err
				}
			} else {
				if !strings.Contains(p, "spec") {
					println(joinKey(p, "idx:ua"), id, nameRU, nameUA)
				}
			}
		}
	}

	return c.Flush()
}

func freeSearchers(c redis.Conn, p string, v ruler) error {
	var id int64
	var err error
	for i := 0; i < v.len(); i++ {
		if s, ok := v.elem(i).(searcher); ok {
			id = s.getID()
			err = c.Send("ZREMRANGEBYSCORE", joinKey(p, "idx:ru"), id, id)
			if err != nil {
				return err
			}
			err = c.Send("ZREMRANGEBYSCORE", joinKey(p, "idx:ua"), id, id)
			if err != nil {
				return err
			}
		}
	}
	return c.Flush()
}

func genKey(p string, v int64) string {
	return joinKey(p, strconv.Itoa(int(v)))
}

func genKeySync(p string) string {
	return joinKey(p, "sync")
}

func genKeyNext(p string) string {
	return joinKey(p, "next")
}

func joinKey(prefix, suffix string) string {
	return prefix + ":" + suffix
}
