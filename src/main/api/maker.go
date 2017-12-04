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
	MarkGP    bool    `json:"mark_gp,omitempty"`
	Logo      string  `json:"logo,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

func (j *jsonMaker) getID() int64 {
	return j.ID
}

func (j *jsonMaker) getSrchRU(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	if j.NameRU == "" {
		return s, r
	}
	s = append(s, normName(j.NameRU))
	if j.MarkGP {
		r = append(r, []rune(s[0])[0])
	}
	return s, r
}

func (j *jsonMaker) getSrchUA(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	if j.NameUA == "" {
		return s, r
	}
	s = append(s, normName(j.NameUA))
	if j.MarkGP {
		r = append(r, []rune(s[0])[0])
	}
	return s, r
}

func (j *jsonMaker) getSrchEN(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	if j.NameEN == "" {
		return s, r
	}
	s = append(s, normName(j.NameEN))
	if j.MarkGP {
		r = append(r, []rune(s[0])[0])
	}
	return s, r
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

func (j *jsonMaker) getFields(list bool) []interface{} {
	if list {
		return []interface{}{
			"id",      // 0
			"name_ru", // 1
			"name_ua", // 2
			"name_en", // 3
			"mark_gp", // 4
			"slug",    // 5
		}
	}
	return []interface{}{
		"id",      // 0
		"name_ru", // 1
		"name_ua", // 2
		"name_en", // 3
		"text_ru", // 4
		"text_ua", // 5
		"text_en", // 6
		"mark_gp", // 7
		"logo",    // 8
		"slug",    // 9
	}
}

func (j *jsonMaker) getValues() []interface{} {
	return []interface{}{
		j.ID,     // 0
		j.NameRU, // 1
		j.NameUA, // 2
		j.NameEN, // 3
		j.TextRU, // 4
		j.TextUA, // 5
		j.TextEN, // 6
		j.MarkGP, // 7
		j.Logo,   // 8
		j.Slug,   // 9
	}
}

func (j *jsonMaker) setValues(list bool, v ...interface{}) {
	for i := range v {
		if v[i] == nil {
			continue
		}
		if list {
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
				j.MarkGP, _ = redis.Bool(v[i], nil)
			case 5:
				j.Slug, _ = redis.String(v[i], nil)
			}
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
			j.TextRU, _ = redis.String(v[i], nil)
		case 5:
			j.TextUA, _ = redis.String(v[i], nil)
		case 6:
			j.TextEN, _ = redis.String(v[i], nil)
		case 7:
			j.MarkGP, _ = redis.Bool(v[i], nil)
		case 8:
			j.Logo, _ = redis.String(v[i], nil)
		case 9:
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

func (v jsonMakers) sort(lang string) {
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

func makeMakersFromJSON(data []byte) (jsonMakers, error) {
	var v []*jsonMaker
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonMakers(v), nil
}

func makeMakersFromIDs(v []int64, err error) (jsonMakers, error) {
	if err != nil {
		return nil, err
	}
	res := make([]*jsonMaker, len(v))
	for i := range res {
		res[i] = &jsonMaker{ID: v[i]}
	}
	return jsonMakers(res), nil
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

func getMakerXSync(h *ctxHelper, p string) ([]int64, error) {
	v, err := int64FromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func getMakerXAbcd(h *ctxHelper, p string) ([]string, error) {
	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcd(c, p, h.lang)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getMakerXAbcdLs(h *ctxHelper, p string) ([]int64, error) {
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

	return v, nil
}

func getMakerXList(h *ctxHelper, p string) (jsonMakers, error) {
	v, err := makeMakersFromIDs(int64sFromJSON(h.data))
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

func getMakerXListAZ(h *ctxHelper, p string) (jsonMakers, error) {
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
	return getMakerXList(h, p)
}

func getMakerX(h *ctxHelper, p string) (jsonMakers, error) {
	v, err := makeMakersFromIDs(int64sFromJSON(h.data))
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

	normLang(h.lang, p, v)

	return v, nil
}

func setMakerX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeMakersFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	x, err := makeMakersFromIDs(findExistsIDs(c, p, mineIDsFromHashers(v)...))
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

func delMakerX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeMakersFromIDs(int64sFromJSON(h.data))
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

func getMakerSync(h *ctxHelper) (interface{}, error) {
	return getMakerXSync(h, prefixMaker)
}

func getMakerAbcd(h *ctxHelper) (interface{}, error) {
	return getMakerXAbcd(h, prefixMaker)
}

func getMakerAbcdLs(h *ctxHelper) (interface{}, error) {
	return getMakerXAbcdLs(h, prefixMaker)
}

func getMakerList(h *ctxHelper) (interface{}, error) {
	return getMakerXList(h, prefixMaker)
}

func getMakerListAZ(h *ctxHelper) (interface{}, error) {
	return getMakerXListAZ(h, prefixMaker)
}

func getMaker(h *ctxHelper) (interface{}, error) {
	return getMakerX(h, prefixMaker)
}

func setMaker(h *ctxHelper) (interface{}, error) {
	return setMakerX(h, prefixMaker)
}

func delMaker(h *ctxHelper) (interface{}, error) {
	return delMakerX(h, prefixMaker)
}
