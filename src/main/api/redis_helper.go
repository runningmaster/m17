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
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

var statusOK = http.StatusText(http.StatusOK)

var apiFunc = map[string]func(h *dbxHelper) (interface{}, error){
	"get-class-atc-sync": getClassATCSync,
	"get-class-atc":      getClassATC,
	"set-class-atc":      setClassATC,
	"del-class-atc":      delClassATC,

	"get-class-nfc-sync": getClassNFCSync,
	"get-class-nfc":      getClassNFC,
	"set-class-nfc":      setClassNFC,
	"del-class-nfc":      delClassNFC,

	"get-class-fsc-sync": getClassFSCSync,
	"get-class-fsc":      getClassFSC,
	"set-class-fsc":      setClassFSC,
	"del-class-fsc":      delClassFSC,

	"get-class-bfc-sync": getClassBFCSync,
	"get-class-bfc":      getClassBFC,
	"set-class-bfc":      setClassBFC,
	"del-class-bfc":      delClassBFC,

	"get-class-cfc-sync": getClassCFCSync,
	"get-class-cfc":      getClassCFC,
	"set-class-cfc":      setClassCFC,
	"del-class-cfc":      delClassCFC,

	"get-class-mpc-sync": getClassMPCSync,
	"get-class-mpc":      getClassMPC,
	"set-class-mpc":      setClassMPC,
	"del-class-mpc":      delClassMPC,

	"get-class-csc-sync": getClassCSCSync,
	"get-class-csc":      getClassCSC,
	"set-class-csc":      setClassCSC,
	"del-class-csc":      delClassCSC,

	"get-class-icd-sync": getClassICDSync,
	"get-class-icd":      getClassICD,
	"set-class-icd":      setClassICD,
	"del-class-icd":      delClassICD,

	"get-inn-sync":    getINNSync,
	"get-inn-abcd":    getINNAbcd,
	"get-inn-abcd-ls": getINNAbcdLs,
	"get-inn-list":    getINNList,
	"get-inn":         getINN,
	"set-inn":         setINN,
	"del-inn":         delINN,

	"get-maker-sync":    getMakerSync,
	"get-maker-abcd":    getMakerAbcd,
	"get-maker-abcd-ls": getMakerAbcdLs,
	"get-maker-list":    getMakerList,
	"get-maker":         getMaker,
	"set-maker":         setMaker,
	"del-maker":         delMaker,

	"get-drug-sync": getDrugSync,
	"get-drug":      getDrug,
	"set-drug":      setDrug,
	"set-drug-sale": setDrugSale,
	"del-drug":      delDrug,

	"get-spec-act-sync":    getSpecACTSync,
	"get-spec-act-abcd":    getSpecACTAbcd,
	"get-spec-act-abcd-ls": getSpecACTAbcdLs,
	"get-spec-act-list":    getSpecACTList,
	"get-spec-act":         getSpecACT,
	"set-spec-act":         setSpecACT,
	"del-spec-act":         delSpecACT,

	"get-spec-inf-sync":    getSpecINFSync,
	"get-spec-inf-abcd":    getSpecINFAbcd,
	"get-spec-inf-abcd-ls": getSpecINFAbcdLs,
	"get-spec-inf-list":    getSpecINFList,
	"get-spec-inf":         getSpecINF,
	"set-spec-inf":         setSpecINF,
	"del-spec-inf":         delSpecINF,

	"get-spec-dec-sync":    getSpecDECSync,
	"get-spec-dec-abcd":    getSpecDECAbcd,
	"get-spec-dec-abcd-ls": getSpecDECAbcdLs,
	"get-spec-dec-list":    getSpecDECList,
	"get-spec-dec":         getSpecDEC,
	"set-spec-dec":         setSpecDEC,
	"del-spec-dec":         delSpecDEC,

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
	getNameEN(string) string
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

func jsonToA(data []byte) (string, error) {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return "", err
	}
	return v, nil
}

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

func freeLinkIDs(c redis.Conn, p1, p2 string, s bool, x int64, v ...int64) error {
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

	if s { // symmetrically
		key = genKey(p1, x, p2)
		err = c.Send("DEL", key)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func saveLinkIDs(c redis.Conn, p1, p2 string, s bool, x int64, v ...int64) error {
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

	if s { // symmetrically
		key = genKey(p1, x, p2)
		val[0] = key
		err = c.Send("SADD", val...)
		if err != nil {
			return err
		}
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

func loadSyncIDs(c redis.Conn, p string, v int64) ([]int64, error) {
	val := make([]interface{}, 0, 3)
	val = append(val, genKey(p, "sync"))
	if v >= 0 {
		val = append(val, v, "+inf")
	} else {
		val = append(val, "-inf", v)
	}
	res, err := redis.Values(c.Do("ZRANGEBYSCORE", val...))
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
			err = c.Send("ZADD", genKey(p, "sync"), "CH", -1*time.Now().Unix(), h.getID())
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
	var nameRU, nameUA, nameEN string
	var abcdRU, abcdUA, abcdEN rune
	var err error
	for i := 0; i < v.len(); i++ {
		if s, ok := v.elem(i).(searcher); ok {
			id = s.getID()
			nameRU = normName(s.getNameRU(p))
			nameUA = normName(s.getNameUA(p))
			nameEN = normName(s.getNameEN(p))

			if nameRU != "" {
				err = c.Send("ZADD", genKey(p, "srch", "ru"), id, nameRU)
				if err != nil {
					return err
				}

				abcdRU = []rune(nameRU)[0]
				err = c.Send("ZADD", genKey(p, "abcd", "ru"), abcdRU, id)
				if err != nil {
					return err
				}

				err = c.Send("ZINCRBY", genKey(p, "rune", "ru"), 1, abcdRU)
				if err != nil {
					return err
				}
			}

			if nameUA != "" {
				err = c.Send("ZADD", genKey(p, "srch", "ua"), id, nameUA)
				if err != nil {
					return err
				}

				abcdUA = []rune(nameUA)[0]
				err = c.Send("ZADD", genKey(p, "abcd", "ua"), abcdUA, id)
				if err != nil {
					return err
				}

				err = c.Send("ZINCRBY", genKey(p, "rune", "ua"), 1, abcdUA)
				if err != nil {
					return err
				}
			}

			if nameEN != "" {
				err = c.Send("ZADD", genKey(p, "srch", "en"), id, nameEN)
				if err != nil {
					return err
				}

				abcdEN = []rune(nameEN)[0]
				err = c.Send("ZADD", genKey(p, "abcd", "en"), abcdEN, id)
				if err != nil {
					return err
				}

				err = c.Send("ZINCRBY", genKey(p, "rune", "en"), 1, abcdEN)
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
	var nameRU, nameUA, nameEN string
	var abcdRU, abcdUA, abcdEN rune
	var err error
	for i := 0; i < v.len(); i++ {
		if s, ok := v.elem(i).(searcher); ok {
			id = s.getID()
			nameRU = normName(s.getNameRU(p))
			nameUA = normName(s.getNameUA(p))
			nameEN = normName(s.getNameEN(p))

			if nameRU != "" {
				err = c.Send("ZREMRANGEBYSCORE", genKey(p, "srch", "ru"), id, id)
				if err != nil {
					return err
				}

				err = c.Send("ZREM", genKey(p, "abcd", "ru"), id)
				if err != nil {
					return err
				}

				abcdRU = []rune(nameRU)[0]
				err = c.Send("ZINCRBY", genKey(p, "rune", "ru"), -1, abcdRU)
				if err != nil {
					return err
				}
			}

			if nameUA != "" {
				err = c.Send("ZREMRANGEBYSCORE", genKey(p, "srch", "ua"), id, id)
				if err != nil {
					return err
				}

				err = c.Send("ZREM", genKey(p, "abcd", "ua"), id)
				if err != nil {
					return err
				}

				abcdUA = []rune(nameUA)[0]
				err = c.Send("ZINCRBY", genKey(p, "rune", "ua"), -1, abcdUA)
				if err != nil {
					return err
				}
			}

			if nameEN != "" {
				err = c.Send("ZREMRANGEBYSCORE", genKey(p, "srch", "en"), id, id)
				if err != nil {
					return err
				}

				err = c.Send("ZREM", genKey(p, "abcd", "en"), id)
				if err != nil {
					return err
				}

				abcdEN = []rune(nameEN)[0]
				err = c.Send("ZINCRBY", genKey(p, "rune", "en"), -1, abcdEN)
				if err != nil {
					return err
				}
			}
		}
	}
	return c.Flush()
}

func loadAbcd(c redis.Conn, p, lang string) ([]string, error) {
	res, err := redis.Ints(c.Do("ZRANGEBYSCORE", genKey(p, "rune", lang), "-inf", "+inf"))
	if err != nil {
		return nil, err
	}

	out := make([]string, len(res))
	for i := range res {
		out[i] = strings.ToUpper(string(rune(res[i])))
	}

	// Sorting
	var col *collate.Collator
	switch lang {
	case "ua":
		col = collate.New(language.Ukrainian)
	case "en":
		col = collate.New(language.English)
	default:
		col = collate.New(language.Russian)
	}
	col.SortStrings(out)

	return out, nil
}

func loadAbcdLs(c redis.Conn, p, a, lang string) ([]int64, error) {
	r := []rune(normName(a))
	if len(r) == 0 {
		return nil, fmt.Errorf("someting wrong with abcd %s", p)
	}

	res, err := redis.Ints(c.Do("ZRANGEBYSCORE", genKey(p, "abcd", lang), r[0], r[0]))
	if err != nil {
		return nil, err
	}

	out := make([]int64, len(res))
	for i := range res {
		out[i] = int64(res[i])
	}

	// FIXME sort magic (info first)

	return out, nil
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

func intsToInt64s(v ...int) []int64 {
	r := make([]int64, len(v))
	for i := range v {
		r[i] = int64(v[i])
	}
	return r
}

func int64ToStrings(v ...int64) []string {
	r := make([]string, len(v))
	for i := range v {
		r[i] = strconv.Itoa(int(v[i]))
	}
	return r
}

/*
func SliceUniqMap(s []int) []int {
    seen := make(map[int]struct{}, len(s))
    j := 0
    for _, v := range s {
        if _, ok := seen[v]; ok {
            continue
        }
        seen[v] = struct{}{}
        s[j] = v
        j++
    }
    return s[:j]
}
*/
