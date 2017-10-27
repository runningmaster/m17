package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"internal/ctxutil"
	"internal/logger"

	"github.com/garyburd/redigo/redis"
)

var statusOK = http.StatusText(http.StatusOK)

var apiFunc = map[string]func(h *dbxHelper) (interface{}, error){
	"get-class-atc-sync":     getClassATCSync,
	"get-class-atc-sync-del": getClassATCSyncDel,
	"get-class-atc":          getClassATC,
	"set-class-atc":          setClassATC,
	"del-class-atc":          delClassATC,

	"get-class-nfc-sync":     getClassNFCSync,
	"get-class-nfc-sync-del": getClassNFCSyncDel,
	"get-class-nfc":          getClassNFC,
	"set-class-nfc":          setClassNFC,
	"del-class-nfc":          delClassNFC,

	"get-class-fsc-sync":     getClassFSCSync,
	"get-class-fsc-sync-del": getClassFSCSyncDel,
	"get-class-fsc":          getClassFSC,
	"set-class-fsc":          setClassFSC,
	"del-class-fsc":          delClassFSC,

	"get-class-bfc-sync":     getClassBFCSync,
	"get-class-bfc-sync-del": getClassBFCSyncDel,
	"get-class-bfc":          getClassBFC,
	"set-class-bfc":          setClassBFC,
	"del-class-bfc":          delClassBFC,

	"get-class-cfc-sync":     getClassCFCSync,
	"get-class-cfc-sync-del": getClassCFCSyncDel,
	"get-class-cfc":          getClassCFC,
	"set-class-cfc":          setClassCFC,
	"del-class-cfc":          delClassCFC,

	"get-class-mpc-sync":     getClassMPCSync,
	"get-class-mpc-sync-del": getClassMPCSyncDel,
	"get-class-mpc":          getClassMPC,
	"set-class-mpc":          setClassMPC,
	"del-class-mpc":          delClassMPC,

	"get-class-csc-sync":     getClassCSCSync,
	"get-class-csc-sync-del": getClassCSCSyncDel,
	"get-class-csc":          getClassCSC,
	"set-class-csc":          setClassCSC,
	"del-class-csc":          delClassCSC,

	"get-class-icd-sync":     getClassICDSync,
	"get-class-icd-sync-del": getClassICDSyncDel,
	"get-class-icd":          getClassICD,
	"set-class-icd":          setClassICD,
	"del-class-icd":          delClassICD,

	"get-inn-sync":     getINNSync,
	"get-inn-sync-del": getINNSyncDel,
	"get-inn":          getINN,
	"set-inn":          setINN,
	"del-inn":          delINN,

	"get-maker-sync":     getMakerSync,
	"get-maker-sync-del": getMakerSyncDel,
	"get-maker":          getMaker,
	"set-maker":          setMaker,
	"del-maker":          delMaker,

	"get-drug-sync":     getDrugSync,
	"get-drug-sync-del": getDrugSyncDel,
	"get-drug":          getDrug,
	"set-drug":          setDrug,
	"set-drug-sale":     setDrugSale,
	"del-drug":          delDrug,

	"get-spec-act-sync":     getSpecACTSync,
	"get-spec-act-sync-del": getSpecACTSyncDel,
	"get-spec-act":          getSpecACT,
	"set-spec-act":          setSpecACT,
	"del-spec-act":          delSpecACT,

	"get-spec-inf-sync":     getSpecINFSync,
	"get-spec-inf-sync-del": getSpecINFSyncDel,
	"get-spec-inf":          getSpecINF,
	"set-spec-inf":          setSpecINF,
	"del-spec-inf":          delSpecINF,

	"get-spec-dec-sync":     getSpecDECSync,
	"get-spec-dec-sync-del": getSpecDECSyncDel,
	"get-spec-dec":          getSpecDEC,
	"set-spec-dec":          setSpecDEC,
	"del-spec-dec":          delSpecDEC,

	"list-sugg":   listSugg,
	"find-sugg":   findSugg,
	"heat-search": heatSearch,
}

type rediser interface {
	Get() redis.Conn
}

type ider interface {
	getID() int64
}

type ruler interface {
	len() int
	elem(int) interface{}
}

type hasher interface {
	ider
	getFields() []interface{}
	getValues() []interface{}
	setValues(...interface{})
}

type niller interface {
	nill(int)
}

type searcher interface {
	ider
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

	h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
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
		key = genKey(p2, v[i], p1)
		err = c.Send("SREM", key, x)
		if err != nil {
			return err
		}
	}

	key = genKey(p1, x, p2)
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
		key = genKey(p2, v[i], p1)
		err = c.Send("SADD", key, x)
		if err != nil {
			return err
		}
		val[i+1] = v[i]
	}

	key = genKey(p1, x, p2)
	val[0] = key
	err = c.Send("SADD", val...)
	if err != nil {
		return err
	}

	return c.Flush()
}

func loadLinkIDs(c redis.Conn, p1, p2 string, x int64) ([]int64, error) {
	key := genKey(p1, x, p2)
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

func loadSyncIDs(c redis.Conn, p string, v int64, deleted ...bool) ([]int64, error) {
	key := genKey(p, "sync")
	if len(deleted) > 0 {
		key = genKey(p, "sync", "del")
	}
	res, err := redis.Values(c.Do("ZRANGEBYSCORE", key, v, "+inf"))
	if err != nil {
		return nil, err
	}

	out := make([]int64, len(res))
	for i := range res {
		out[i], _ = redis.Int64(res[i], nil)
	}

	return out, nil
}

func mixKeyAndFields(p string, h hasher) []interface{} {
	f := h.getFields()
	r := make([]interface{}, 0, 1+len(f))
	return append(append(r, genKey(p, h.getID())), f...)
}

func mixKeyAndFieldsAndValues(p string, h hasher) []interface{} {
	f := h.getFields()
	v := h.getValues()
	r := make([]interface{}, 0, 1+len(f)*2)
	r = append(r, genKey(p, h.getID()))
	for i := range v {
		// ? check zero values here
		r = append(r, f[i], v[i])
	}
	return r
}

func saveHashers(c redis.Conn, p string, v ruler) error {
	if v.len() == 0 {
		return nil
	}

	var err error
	for i := 0; i < v.len(); i++ {
		if h, ok := v.elem(i).(hasher); ok {
			if h.getID() == 0 {
				return fmt.Errorf("ID must have value (%s)", p)
			}
			err = c.Send("HMSET", mixKeyAndFieldsAndValues(p, h)...)
			if err != nil {
				return err
			}
			err = c.Send("ZADD", genKey(p, "sync"), "CH", time.Now().Unix(), h.getID())
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
			err = c.Send("HMGET", mixKeyAndFields(p, h)...)
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
		if len(r) > 0 && r[0] == nil {
			if n, ok := v.(niller); ok {
				n.nill(i)
			}
			continue
		}
		if h, ok := v.elem(i).(hasher); ok {
			h.setValues(r...)
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
			err = c.Send("DEL", genKey(p, h.getID()))
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
			err = c.Send("ZREM", genKey(p, "sync"), h.getID())
			if err != nil {
				return err
			}
			err = c.Send("ZADD", genKey(p, "sync", "del"), "CH", time.Now().Unix(), h.getID())
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
				err = c.Send("ZADD", genKey(p, "idx", "ru"), id, normName(nameRU))
				if err != nil {
					return err
				}
			}
			if nameUA != "" {
				err = c.Send("ZADD", genKey(p, "idx", "ua"), id, normName(nameUA))
				if err != nil {
					return err
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
			err = c.Send("ZREMRANGEBYSCORE", genKey(p, "idx", "ru"), id, id)
			if err != nil {
				return err
			}
			err = c.Send("ZREMRANGEBYSCORE", genKey(p, "idx", "ua"), id, id)
			if err != nil {
				return err
			}
		}
	}
	return c.Flush()
}

func genKey(v ...interface{}) string {
	a := make([]string, 0, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case []byte:
			a = append(a, string(x))
		case int:
			a = append(a, strconv.Itoa(x))
		case int64:
			a = append(a, strconv.Itoa(int(x)))
		case string:
			a = append(a, x)
		default:
			a = append(a, fmt.Sprintf("%v", x))
		}
	}
	return strings.Join(a, ":")
}
