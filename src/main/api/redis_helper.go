package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
}

type rediser interface {
	Get() redis.Conn
}

type ruler interface {
	len() int
}

type hasher interface {
	getKey(string) string
	getKeyAndUnixtimeID(string) []interface{}
	getKeyAndFieldValues(string) []interface{}
	getKeyAndFields(string) []interface{}
	setValues(...interface{})
}

type ruleHasher interface {
	ruler
	elem(int) hasher
	nill(int)
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

func (h *dbxHelper) getSyncList(p string, v int64) ([]int64, error) {
	c := h.getConn()
	defer h.delConn(c)

	res, err := redis.Values(c.Do("ZRANGEBYSCORE", p+":"+"sync", v, "+inf"))
	if err != nil {
		return nil, err
	}

	out := make([]int64, len(res))
	for i := range res {
		out[i], _ = redis.Int64(res[i], err)
	}

	return out, nil
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
