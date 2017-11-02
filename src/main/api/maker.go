package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixMaker = "maker"
)

type jsonMaker struct {
	ID        int64   `json:"id,omitempty"`
	IDNode    int64   `json:"id_node,omitempty"`
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Text      string  `json:"text,omitempty"` // *
	TextRU    string  `json:"text_ru,omitempty"`
	TextUA    string  `json:"text_ua,omitempty"`
	TextEN    string  `json:"text_en,omitempty"`
	IsGP      bool    `json:"is_gp,omitempty"`
	Logo      string  `json:"logo,omitempty"`
	Slug      string  `json:"slug,omitempty"`
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

func (j *jsonMaker) getNameEN(_ string) string {
	return j.NameEN
}

func (j *jsonMaker) lang(l, _ string) {
	switch l {
	case "ru":
		j.Name = j.NameRU
		j.Text = j.TextRU
		j.IDSpecDEC = nil
	case "ua":
		j.Name = j.NameUA
		j.Text = j.TextUA
		j.IDSpecINF = nil
	}

	if l == "ru" || l == "ua" {
		j.NameRU = ""
		j.NameUA = ""
		j.NameEN = ""
		j.TextRU = ""
		j.TextUA = ""
		j.TextEN = ""
	}
}

func (j *jsonMaker) getFields() []interface{} {
	return []interface{}{
		"id",      // 0
		"id_node", // 1
		"name_ru", // 2
		"name_ua", // 3
		"name_en", // 4
		"text_ru", // 5
		"text_ua", // 6
		"text_en", // 7
		"is_gp",   // 8
		"logo",    // 9
		"slug",    // 10
	}
}

func (j *jsonMaker) getValues() []interface{} {
	return []interface{}{
		j.ID,     // 0
		j.IDNode, // 1
		j.NameRU, // 2
		j.NameUA, // 3
		j.NameEN, // 4
		j.TextRU, // 5
		j.TextUA, // 6
		j.TextEN, // 7
		j.IsGP,   // 8
		j.Logo,   // 9
		j.Slug,   // 10
	}
}

func (j *jsonMaker) setValues(v ...interface{}) {
	for i := range v {
		if v[i] == nil {
			continue
		}
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
			j.IsGP, _ = redis.Bool(v[i], nil)
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

func (j jsonMakers) elem(i int) interface{} {
	return j[i]
}

func (j jsonMakers) null(i int) bool {
	return j[i] == nil
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
	return makeMakers(v...), nil
}

func makeMakers(x ...int64) jsonMakers {
	v := make([]*jsonMaker, len(x))
	for i := range v {
		v[i] = &jsonMaker{ID: x[i]}
	}
	return jsonMakers(v)
}

func loadMakerLinks(c redis.Conn, p string, v []*jsonMaker) error {
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

func getMakerXSync(h *dbxHelper, p string) ([]int64, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func getMakerXAbcd(h *dbxHelper, p string) ([]string, error) {
	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcd(c, p, h.lang)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getMakerXAbcdLs(h *dbxHelper, p string) ([]int64, error) {
	a, err := jsonToA(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcdLs(c, p, a, h.lang)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getMakerXList(h *dbxHelper, p string) (jsonMakers, error) {
	v, err := getMakerX(h, p)
	if err != nil {
		return nil, err
	}

	coll := newCollator(h.lang)
	sort.Slice(v,
		func(i, j int) bool {
			if v[i] == nil || v[j] == nil {
				return true
			}
			return coll.CompareString(v[i].Name, v[j].Name) < 0
		},
	)

	normLang(h.lang, p, v)

	return v, nil
}

func getMakerX(h *dbxHelper, p string) (jsonMakers, error) {
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

	err = loadMakerLinks(c, p, v)
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

// MAKER

func getMakerSync(h *dbxHelper) (interface{}, error) {
	return getMakerXSync(h, prefixMaker)
}

func getMakerAbcd(h *dbxHelper) (interface{}, error) {
	return getMakerXAbcd(h, prefixMaker)
}

func getMakerAbcdLs(h *dbxHelper) (interface{}, error) {
	return getMakerXAbcdLs(h, prefixMaker)
}

func getMakerList(h *dbxHelper) (interface{}, error) {
	return getMakerXList(h, prefixMaker)
}

func getMaker(h *dbxHelper) (interface{}, error) {
	return getMakerX(h, prefixMaker)
}

func setMaker(h *dbxHelper) (interface{}, error) {
	return setMakerX(h, prefixMaker)
}

func delMaker(h *dbxHelper) (interface{}, error) {
	return delMakerX(h, prefixMaker)
}
