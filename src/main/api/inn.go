package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixINN = "inn"
)

type jsonINN struct {
	ID        int64   `json:"id,omitempty"`
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

func (j *jsonINN) getID() int64 {
	return j.ID
}

func (j *jsonINN) getSrchRU(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	if j.NameRU == "" {
		return s, r
	}
	s = append(s, normName(j.NameRU))
	r = append(r, []rune(s[0])[0])
	return s, r
}

func (j *jsonINN) getSrchUA(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	if j.NameUA == "" {
		return s, r
	}
	s = append(s, normName(j.NameUA))
	r = append(r, []rune(s[0])[0])
	return s, r
}

func (j *jsonINN) getSrchEN(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	if j.NameEN == "" {
		return s, r
	}
	s = append(s, normName(j.NameEN))
	r = append(r, []rune(s[0])[0])
	return s, r
}

func (j *jsonINN) lang(l, _ string) {
	switch l {
	case "ru":
		j.Name = fmt.Sprintf("%s (%s)", j.NameEN, j.NameRU)
		j.IDSpecDEC = nil
	case "ua":
		j.Name = fmt.Sprintf("%s (%s)", j.NameEN, j.NameUA)
		j.IDSpecINF = nil
	}

	if l == "ru" || l == "ua" {
		j.NameRU = ""
		j.NameUA = ""
		j.NameEN = ""
	}
}

func (j *jsonINN) getFields(_ bool) []interface{} {
	return []interface{}{
		"id",      // 0
		"name_ru", // 1
		"name_ua", // 2
		"name_en", // 3
		"slug",    // 4
	}
}

func (j *jsonINN) getValues() []interface{} {
	return []interface{}{
		j.ID,     // 0
		j.NameRU, // 1
		j.NameUA, // 2
		j.NameEN, // 3
		j.Slug,   // 4
	}
}

func (j *jsonINN) setValues(_ bool, v ...interface{}) {
	for i := range v {
		if v[i] == nil {
			continue
		}
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

func (j jsonINNs) elem(i int) interface{} {
	return j[i]
}

func (j jsonINNs) null(i int) bool {
	return j[i] == nil
}

func (j jsonINNs) nill(i int) {
	j[i] = nil
}

func (v jsonINNs) sort(lang string) {
	coll := newCollator(lang)
	sort.Slice(v,
		func(i, j int) bool {
			if v[i] == nil && v[j] == nil {
				return true
			}
			if v[i] == nil && v[j] != nil {
				return false
			}
			if v[i] != nil && v[j] == nil {
				return true
			}
			return coll.CompareString(v[i].Name, v[j].Name) < 0
		},
	)
}

func makeINNsFromJSON(data []byte) (jsonINNs, error) {
	var v []*jsonINN
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonINNs(v), nil
}

func makeINNsFromIDs(v []int64, err error) (jsonINNs, error) {
	if err != nil {
		return nil, err
	}
	res := make([]*jsonINN, len(v))
	for i := range res {
		res[i] = &jsonINN{ID: v[i]}
	}
	return jsonINNs(res), nil
}

func loadINNLinks(c redis.Conn, p string, v []*jsonINN) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		v[i].IDSpecDEC, err = loadLinkIDs(c, p, prefixSpecDEC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDSpecINF, err = loadLinkIDs(c, p, prefixSpecINF, v[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getINNXSync(h *ctxHelper, p string) ([]int64, error) {
	v, err := int64FromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func getINNXAbcd(h *ctxHelper, p string) ([]string, error) {
	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcd(c, p, h.lang)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getINNXAbcdLs(h *ctxHelper, p string) ([]int64, error) {
	s, err := stringFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcdLs(c, p, s, h.lang)
	if err != nil {
		return nil, err
	}

	//	h.data = int64sToJSON(v)
	//	x, err := getINNXList

	return v, nil
}

func getINNXList(h *ctxHelper, p string) (jsonINNs, error) {
	v, err := makeINNsFromIDs(int64sFromJSON(h.data))
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = loadHashers(c, p, v, true)
	if err != nil {
		return nil, err
	}

	normLang(h.lang, p, v)

	v.sort(h.lang)

	return v, nil
}

func getINNXListAZ(h *ctxHelper, p string) (jsonINNs, error) {
	s, err := stringFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcdLs(c, p, s, h.lang)
	if err != nil {
		return nil, err
	}

	h.data = int64sToJSON(v)
	return getINNXList(h, p)
}

func getINNX(h *ctxHelper, p string) (jsonINNs, error) {
	v, err := makeINNsFromIDs(int64sFromJSON(h.data))
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

	err = loadINNLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	normLang(h.lang, p, v)

	return v, nil
}

func setINNX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeINNsFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	x, err := makeINNsFromIDs(findExistsIDs(c, p, mineIDsFromHashers(v)...))
	if err != nil {
		return nil, err
	}

	if len(x) > 0 {
		err = loadHashers(c, p, x)
		if err != nil {
			return nil, err
		}
		err = freeSearchers(c, p, x)
		if err != nil {
			return nil, err
		}
	}

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

func delINNX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeINNsFromIDs(int64sFromJSON(h.data))
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

func getINNSync(h *ctxHelper) (interface{}, error) {
	return getINNXSync(h, prefixINN)
}

func getINNAbcd(h *ctxHelper) (interface{}, error) {
	return getINNXAbcd(h, prefixINN)
}

func getINNAbcdLs(h *ctxHelper) (interface{}, error) {
	return getINNXAbcdLs(h, prefixINN)
}

func getINNList(h *ctxHelper) (interface{}, error) {
	return getINNXList(h, prefixINN)
}

func getINNListAZ(h *ctxHelper) (interface{}, error) {
	return getINNXListAZ(h, prefixINN)
}

func getINN(h *ctxHelper) (interface{}, error) {
	return getINNX(h, prefixINN)
}

func setINN(h *ctxHelper) (interface{}, error) {
	return setINNX(h, prefixINN)
}

func delINN(h *ctxHelper) (interface{}, error) {
	return delINNX(h, prefixINN)
}
