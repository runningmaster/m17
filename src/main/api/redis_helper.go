package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

var statusOK = http.StatusText(http.StatusOK)

type rediser interface {
	Get() redis.Conn
}

type ider interface {
	getID() int64
}

type ruler interface {
	len() int
	elem(int) interface{}
	null(int) bool
}

type hasher interface {
	ider
	getFields(bool) []interface{}
	getValues() []interface{}
	setValues(bool, ...interface{})
}

type niller interface {
	nill(int)
}

type langer interface {
	lang(string, string)
}

type searcher interface {
	ider
	getSrchRU(string) ([]string, []rune)
	getSrchUA(string) ([]string, []rune)
	getSrchEN(string) ([]string, []rune)
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

func normLang(s, p string, v ruler) {
	if s == "" {
		return
	}

	for i := 0; i < v.len(); i++ {
		if v.null(i) {
			continue
		}

		if l, ok := v.elem(i).(langer); ok {
			l.lang(s, p)
		}
	}
}

func mixKeyAndFields(p string, list bool, h hasher) []interface{} {
	f := h.getFields(list)
	r := make([]interface{}, 0, 1+len(f))
	return append(append(r, genKey(p, h.getID())), f...)
}

func mixKeyAndFieldsAndValues(p string, h hasher) []interface{} {
	f := h.getFields(false)
	v := h.getValues()
	r := make([]interface{}, 0, 1+len(f)*2)
	r = append(r, genKey(p, h.getID()))
	for i := range v {
		if !zeroValue(v[i]) {
			r = append(r, f[i], v[i])
		}
	}
	return r
}

func saveHashers(c redis.Conn, p string, v ruler, onlyUpdate ...bool) error {
	if v.len() == 0 {
		return nil
	}

	var err error
	for i := 0; i < v.len(); i++ {
		if v.null(i) {
			continue
		}

		if h, ok := v.elem(i).(hasher); ok {
			if h.getID() == 0 {
				return fmt.Errorf("ID must have value (%s)", p)
			}
			if len(onlyUpdate) == 0 {
				err = c.Send("DEL", genKey(p, h.getID()))
				if err != nil {
					return err
				}
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

func loadHashFieldAsString(c redis.Conn, p, s string, x int64) (string, error) {
	key := genKey(p, x)
	res, err := redis.Strings(c.Do("HMGET", key, s))
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", fmt.Errorf("got %v", len(res))
	}
	return res[0], nil
}

func loadHashers(c redis.Conn, p string, v ruler, mustBeList ...bool) error {
	if v.len() == 0 {
		return nil
	}

	l := len(mustBeList) > 0
	var err error
	for i := 0; i < v.len(); i++ {
		if v.null(i) {
			continue
		}

		if h, ok := v.elem(i).(hasher); ok {
			err = c.Send("HMGET", mixKeyAndFields(p, l, h)...)
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
		if v.null(i) {
			continue
		}

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
			h.setValues(l, r...)
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
		if v.null(i) {
			continue
		}

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
		if v.null(i) {
			continue
		}

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

func mineIDsFromHashers(v ruler) []int64 {
	x := make([]int64, 0, v.len())
	for i := 0; i < v.len(); i++ {
		if v.null(i) {
			continue
		}
		if h, ok := v.elem(i).(hasher); ok {
			x = append(x, h.getID())
		}
	}
	return x
}

func findExistsIDs(c redis.Conn, p string, v ...int64) ([]int64, error) {
	if len(v) == 0 {
		return nil, nil
	}

	var err error
	for i := range v {
		err = c.Send("EXISTS", genKey(p, v[i]))
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	res := make([]int64, 0, len(v))
	var r bool
	for i := range v {
		r, err = redis.Bool(c.Receive())
		if err != nil {
			return nil, err
		}
		if r {
			res = append(res, v[i])
		}
	}

	return res, nil
}

func saveSearchers(c redis.Conn, p string, v ruler) error {
	if v.len() == 0 {
		return nil
	}

	var id int64
	var sx string
	var nameRU, nameUA, nameEN []string
	var abcdRU, abcdUA, abcdEN []rune
	var err error
	for i := 0; i < v.len(); i++ {
		if v.null(i) {
			continue
		}

		if s, ok := v.elem(i).(searcher); ok {
			id = s.getID()
			sx = "|" + strconv.Itoa(int(id))
			nameRU, abcdRU = s.getSrchRU(p)
			nameUA, abcdUA = s.getSrchUA(p)
			nameEN, abcdEN = s.getSrchEN(p)

			for _, v := range nameRU {
				err = c.Send("ZADD", genKey(p, "srch", "ru"), id, v+sx)
				if err != nil {
					return err
				}
			}
			for _, v := range abcdRU {
				err = c.Send("ZADD", genKey(p, "abcd", "ru"), v, id)
				if err != nil {
					return err
				}
				err = c.Send("ZINCRBY", genKey(p, "rune", "ru"), 1, v)
				if err != nil {
					return err
				}
			}

			for _, v := range nameUA {
				err = c.Send("ZADD", genKey(p, "srch", "ua"), id, v+sx)
				if err != nil {
					return err
				}
			}
			for _, v := range abcdUA {
				err = c.Send("ZADD", genKey(p, "abcd", "ua"), v, id)
				if err != nil {
					return err
				}
				err = c.Send("ZINCRBY", genKey(p, "rune", "ua"), 1, v)
				if err != nil {
					return err
				}
			}

			for _, v := range nameEN {
				err = c.Send("ZADD", genKey(p, "srch", "en"), id, v+sx)
				if err != nil {
					return err
				}
			}
			for _, v := range abcdEN {
				err = c.Send("ZADD", genKey(p, "abcd", "en"), v, id)
				if err != nil {
					return err
				}
				err = c.Send("ZINCRBY", genKey(p, "rune", "en"), 1, v)
				if err != nil {
					return err
				}
			}
		}
	}

	return c.Flush()
}

func freeSearchers(c redis.Conn, p string, v ruler) error {
	if v.len() == 0 {
		return nil
	}

	var id int64
	var nameRU, nameUA, nameEN []string
	var abcdRU, abcdUA, abcdEN []rune
	var err error
	for i := 0; i < v.len(); i++ {
		if v.null(i) {
			continue
		}

		if s, ok := v.elem(i).(searcher); ok {
			id = s.getID()
			nameRU, abcdRU = s.getSrchRU(p)
			nameUA, abcdUA = s.getSrchUA(p)
			nameEN, abcdEN = s.getSrchEN(p)

			for range nameRU {
				err = c.Send("ZREMRANGEBYSCORE", genKey(p, "srch", "ru"), id, id)
				if err != nil {
					return err
				}
			}
			for _, v := range abcdRU {
				err = c.Send("ZREM", genKey(p, "abcd", "ru"), id)
				if err != nil {
					return err
				}
				err = c.Send("ZINCRBY", genKey(p, "rune", "ru"), -1, v)
				if err != nil {
					return err
				}
			}

			for range nameUA {
				err = c.Send("ZREMRANGEBYSCORE", genKey(p, "srch", "ua"), id, id)
				if err != nil {
					return err
				}
			}
			for _, v := range abcdUA {
				err = c.Send("ZREM", genKey(p, "abcd", "ua"), id)
				if err != nil {
					return err
				}
				err = c.Send("ZINCRBY", genKey(p, "rune", "ua"), -1, v)
				if err != nil {
					return err
				}
			}

			for range nameEN {
				err = c.Send("ZREMRANGEBYSCORE", genKey(p, "srch", "en"), id, id)
				if err != nil {
					return err
				}
			}
			for _, v := range abcdEN {
				err = c.Send("ZREM", genKey(p, "abcd", "en"), id)
				if err != nil {
					return err
				}
				err = c.Send("ZINCRBY", genKey(p, "rune", "en"), -1, v)
				if err != nil {
					return err
				}
			}
		}
	}
	return c.Flush()
}

func newCollator(lang string) *collate.Collator {
	switch lang {
	case "ru":
		return collate.New(language.Russian)
	case "ua":
		return collate.New(language.Ukrainian)
	default:
		return collate.New(language.English)
	}
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

	coll := newCollator(lang)
	coll.SortStrings(out)

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

	return out, nil
}

type findRes struct {
	ID   int64
	Name string
}

//
func findIn(c redis.Conn, p, lang, text string, conj bool) ([]*findRes, error) {
	text = strings.ToLower(text)
	var flds []string
	if conj {
		flds = strings.Fields(text)
		text = flds[0]
	}

	res := make([]*findRes, 0, 100)
	var next int
	var vals []interface{}
	for done := false; !done; {
		v, err := redis.Values(c.Do("ZSCAN", genKey(p, "srch", lang), next, "MATCH", "*"+text+"*", "COUNT", 100))
		if err != nil {
			return nil, err
		}

		next, _ = redis.Int(v[0], err)
		vals, _ = redis.Values(v[1], err)

		var r *findRes
		var s string
		var y bool
		for i := range vals {
			if i&1 != 1 {
				continue // ignore even
			}
			s, _ = redis.String(vals[i-1], err)
			y = true
			if conj {
				for j := 1; j < len(flds); j++ {
					y = strings.Contains(s, flds[j])
					if !y {
						break
					}
				}
			}
			if y {
				r = &findRes{}
				r.ID, _ = redis.Int64(vals[i], err)
				r.Name = strings.Split(s, "|")[0]
				res = append(res, r)
			}
		}
		done = next == 0
	}

	return res, nil
}
