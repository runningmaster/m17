package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixSpecACT = "spec:act" // RU only
	prefixSpecINF = "spec:inf" // RU only
	prefixSpecDEC = "spec:dec" // UA only
)

type jsonSpec struct {
	ID         int64   `json:"id,omitempty"`
	IDINN      []int64 `json:"id_inn,omitempty"`
	IDDrug     []int64 `json:"id_drug,omitempty"`
	IDMake     []int64 `json:"id_make,omitempty"`
	IDMakeGP   int64   `json:"id_make_gp,omitempty"`
	IDSpecACT  []int64 `json:"id_spec_act,omitempty"`
	IDSpecDEC  []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF  []int64 `json:"id_spec_inf,omitempty"`
	IDClassATC []int64 `json:"id_class_atc,omitempty"`
	IDClassNFC []int64 `json:"id_class_nfc,omitempty"`
	IDClassFSC []int64 `json:"id_class_fsc,omitempty"`
	IDClassBFC []int64 `json:"id_class_bfc,omitempty"`
	IDClassCFC []int64 `json:"id_class_cfc,omitempty"`
	IDClassMPC []int64 `json:"id_class_mpc,omitempty"`
	IDClassCSC []int64 `json:"id_class_csc,omitempty"`
	IDClassICD []int64 `json:"id_class_icd,omitempty"`
	Name       string  `json:"name,omitempty"` // *
	NameRU     string  `json:"name_ru,omitempty"`
	NameRUSrc  string  `json:"name_ru_src,omitempty"`
	NameUA     string  `json:"name_ua,omitempty"`
	NameUASrc  string  `json:"name_ua_src,omitempty"`
	NameEN     string  `json:"name_en,omitempty"`
	NameENSrc  string  `json:"name_en_src,omitempty"`
	Head       string  `json:"head,omitempty"` // *
	HeadRU     string  `json:"head_ru,omitempty"`
	HeadUA     string  `json:"head_ua,omitempty"`
	HeadEN     string  `json:"head_en,omitempty"`
	Text       string  `json:"text,omitempty"` // *
	TextRU     string  `json:"text_ru,omitempty"`
	TextUA     string  `json:"text_ua,omitempty"`
	TextEN     string  `json:"text_en,omitempty"`
	Flag       int64   `json:"flag,omitempty"`
	Slug       string  `json:"slug,omitempty"`
	SlugGP     string  `json:"slug_gp,omitempty"` // [{"name":"foo", "slug": "bar"}]
	ImageOrg   string  `json:"image_org,omitempty"`
	ImageBox   string  `json:"image_box,omitempty"`
	CreatedAt  int64   `json:"created_at,omitempty"`
	UpdatedAt  int64   `json:"updated_at,omitempty"`

	Sale float64 `json:"sale,omitempty"`
}

func (j *jsonSpec) getID() int64 {
	return j.ID
}

func (j *jsonSpec) getNameRU(p string) string {
	if p != prefixSpecDEC {
		return j.NameRUSrc
	}
	return ""
}

func (j *jsonSpec) getNameUA(p string) string {
	if p == prefixSpecDEC {
		return j.NameUASrc
	}
	return ""
}

func (j *jsonSpec) getNameEN(p string) string {
	return j.NameENSrc
}

func (j *jsonSpec) lang(l, p string) {
	switch l {
	case "ru":
		if p != prefixSpecDEC {
			j.Name = j.NameRU
			j.Head = j.HeadRU
			j.Text = j.TextRU
		}
		if p == prefixSpecACT {
			j.Name = fmt.Sprintf("%s (%s)", j.NameRU, j.NameEN)
		}
	case "ua":
		if p == prefixSpecDEC {
			j.Name = j.NameUA
			j.Head = j.HeadUA
			j.Text = j.TextUA
		}
	}

	if l == "ru" || l == "ua" {
		j.NameRU = ""
		j.NameUA = ""
		j.NameEN = ""
		j.HeadRU = ""
		j.HeadUA = ""
		j.HeadEN = ""
		j.TextRU = ""
		j.TextUA = ""
		j.TextEN = ""
	}
}

func (j *jsonSpec) getFields(list bool) []interface{} {
	if list {
		return []interface{}{
			"id",      // 0
			"name_ru", // 1
			"name_ua", // 2
			"name_en", // 3
			"flag",    // 4
			"slug",    // 5
			"sale",    // 6
		}
	}
	return []interface{}{
		"id",         // 0
		"id_make_gp", // 1
		"name_ru",    // 2
		"name_ua",    // 3
		"name_en",    // 4
		"head_ru",    // 5
		"head_ua",    // 6
		"head_en",    // 7
		"text_ru",    // 8
		"text_ua",    // 9
		"text_en",    // 10
		"flag",       // 11
		"slug",       // 12
		"slug_gp",    // 13
		"image_org",  // 14
		"image_box",  // 15
		"created_at", // 16
		"updated_at", // 17
	}
}

func (j *jsonSpec) getValues() []interface{} {
	if len(j.TextRU) > 0 || len(j.TextUA) > 0 {
		j.Flag = 1
	}
	return []interface{}{
		j.ID,        // 0
		j.IDMakeGP,  // 1
		j.NameRU,    // 2
		j.NameUA,    // 3
		j.NameEN,    // 4
		j.HeadRU,    // 5
		j.HeadUA,    // 6
		j.HeadEN,    // 7
		j.TextRU,    // 8
		j.TextUA,    // 9
		j.TextEN,    // 10
		j.Flag,      // 11
		j.Slug,      // 12
		j.SlugGP,    // 13
		j.ImageOrg,  // 14
		j.ImageBox,  // 15
		j.CreatedAt, // 16
		j.UpdatedAt, // 17
	}
}

func (j *jsonSpec) setValues(list bool, v ...interface{}) {
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
				j.Flag, _ = redis.Int64(v[i], nil)
			case 5:
				j.Slug, _ = redis.String(v[i], nil)
			case 6:
				j.Sale, _ = redis.Float64(v[i], nil)
			}
			continue
		}
		switch i {
		case 0:
			j.ID, _ = redis.Int64(v[i], nil)
		case 1:
			j.IDMakeGP, _ = redis.Int64(v[i], nil)
		case 2:
			j.NameRU, _ = redis.String(v[i], nil)
		case 3:
			j.NameUA, _ = redis.String(v[i], nil)
		case 4:
			j.NameEN, _ = redis.String(v[i], nil)
		case 5:
			j.HeadRU, _ = redis.String(v[i], nil)
		case 6:
			j.HeadUA, _ = redis.String(v[i], nil)
		case 7:
			j.HeadEN, _ = redis.String(v[i], nil)
		case 8:
			j.TextRU, _ = redis.String(v[i], nil)
		case 9:
			j.TextUA, _ = redis.String(v[i], nil)
		case 10:
			j.TextEN, _ = redis.String(v[i], nil)
		case 11:
			j.Flag, _ = redis.Int64(v[i], nil)
		case 12:
			j.Slug, _ = redis.String(v[i], nil)
		case 13:
			j.SlugGP, _ = redis.String(v[i], nil)
		case 14:
			j.ImageOrg, _ = redis.String(v[i], nil)
		case 15:
			j.ImageBox, _ = redis.String(v[i], nil)
		case 16:
			j.CreatedAt, _ = redis.Int64(v[i], nil)
		case 17:
			j.UpdatedAt, _ = redis.Int64(v[i], nil)
		}
	}
}

type jsonSpecs []*jsonSpec

func (j jsonSpecs) len() int {
	return len(j)
}

func (j jsonSpecs) elem(i int) interface{} {
	return j[i]
}

func (j jsonSpecs) null(i int) bool {
	return j[i] == nil
}

func (j jsonSpecs) nill(i int) {
	j[i] = nil
}

func (v jsonSpecs) sort(lang string) {
	coll := newCollator(lang)
	sort.Slice(v,
		func(i, j int) bool {
			if v[i] == nil && v[j] == nil {
				return true
			}
			if v[i] != nil && v[j] == nil {
				return true
			}
			if v[i] == nil && v[j] != nil {
				return false
			}
			if v[i].Flag != 0 && v[j].Flag == 0 {
				return true
			} else if v[i].Flag == 0 && v[j].Flag != 0 {
				return false
			}
			if v[i].Sale > v[j].Sale {
				return true
			} else if v[i].Sale < v[j].Sale {
				return false
			}
			return coll.CompareString(v[i].Name, v[j].Name) < 0
		},
	)
}

func jsonToSpecs(data []byte) (jsonSpecs, error) {
	var v []*jsonSpec
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonSpecs(v), nil
}

func jsonToSpecsFromIDs(data []byte) (jsonSpecs, error) {
	v, err := jsonToIDs(data)
	if err != nil {
		return nil, err
	}
	return makeSpecs(v...), nil
}

func makeSpecs(x ...int64) jsonSpecs {
	v := make([]*jsonSpec, len(x))
	for i := range v {
		v[i] = &jsonSpec{ID: x[i]}
	}
	return jsonSpecs(v)
}

func makeSpecsFromSpecs(x ...*jsonSpec) jsonSpecs {
	v := make([]int64, 0, len(x))
	for i := range x {
		if x[i] != nil {
			v = append(v, x[i].ID)
		}
	}
	return makeSpecs(v...)
}

func loadSpecLinks(c redis.Conn, p string, v []*jsonSpec) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		v[i].IDINN, err = loadLinkIDs(c, p, prefixINN, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDDrug, err = loadLinkIDs(c, p, prefixDrug, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDMake, err = loadLinkIDs(c, p, prefixMaker, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDSpecACT, err = loadLinkIDs(c, p, prefixSpecACT, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDSpecDEC, err = loadLinkIDs(c, p, prefixSpecDEC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDSpecINF, err = loadLinkIDs(c, p, prefixSpecINF, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassATC, err = loadLinkIDs(c, p, prefixClassATC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassNFC, err = loadLinkIDs(c, p, prefixClassNFC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassFSC, err = loadLinkIDs(c, p, prefixClassFSC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassBFC, err = loadLinkIDs(c, p, prefixClassBFC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassCFC, err = loadLinkIDs(c, p, prefixClassCFC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassMPC, err = loadLinkIDs(c, p, prefixClassMPC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassCSC, err = loadLinkIDs(c, p, prefixClassCSC, v[i].ID)
		if err != nil {
			return err
		}
		v[i].IDClassICD, err = loadLinkIDs(c, p, prefixClassICD, v[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveSpecLinks(c redis.Conn, p string, v ...*jsonSpec) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		err = saveLinkIDs(c, p, prefixINN, true, v[i].ID, v[i].IDINN...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixDrug, true, v[i].ID, v[i].IDDrug...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixMaker, true, v[i].ID, v[i].IDMake...)
		if err != nil {
			return err
		}
		if v[i].IDMakeGP != 0 {
			err = saveLinkIDs(c, p, prefixMaker, false, v[i].ID, v[i].IDMakeGP)
			if err != nil {
				return err
			}
			err = saveMakerGPLinks(c, v[i].IDMakeGP, v[i].IDMake...)
			if err != nil {
				return err
			}
		}

		err = saveLinkIDs(c, p, prefixSpecACT, true, v[i].ID, v[i].IDSpecACT...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixSpecDEC, true, v[i].ID, v[i].IDSpecDEC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixSpecINF, true, v[i].ID, v[i].IDSpecINF...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassATC, true, v[i].ID, v[i].IDClassATC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassNFC, true, v[i].ID, v[i].IDClassNFC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassFSC, true, v[i].ID, v[i].IDClassFSC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassBFC, true, v[i].ID, v[i].IDClassBFC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassCFC, true, v[i].ID, v[i].IDClassCFC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassMPC, true, v[i].ID, v[i].IDClassMPC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassCSC, true, v[i].ID, v[i].IDClassCSC...)
		if err != nil {
			return err
		}

		err = saveLinkIDs(c, p, prefixClassICD, true, v[i].ID, v[i].IDClassICD...)
		if err != nil {
			return err
		}
	}

	return nil
}

func freeSpecLinks(c redis.Conn, p string, v ...*jsonSpec) error {
	var val []int64
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		val, _ = loadLinkIDs(c, p, prefixINN, v[i].ID)
		err = freeLinkIDs(c, p, prefixINN, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixDrug, v[i].ID)
		err = freeLinkIDs(c, p, prefixDrug, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		if v[i].IDMakeGP != 0 {
			err = freeMakerGPLinks(c, v[i].IDMakeGP, v[i].IDMake...)
			if err != nil {
				return err
			}
			err = freeLinkIDs(c, p, prefixMaker, false, v[i].ID, v[i].IDMakeGP)
			if err != nil {
				return err
			}
		}
		val, _ = loadLinkIDs(c, p, prefixMaker, v[i].ID)
		err = freeLinkIDs(c, p, prefixMaker, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixSpecDEC, v[i].ID)
		err = freeLinkIDs(c, p, prefixSpecDEC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixSpecINF, v[i].ID)
		err = freeLinkIDs(c, p, prefixSpecINF, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassATC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassATC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassNFC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassNFC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassFSC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassFSC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassBFC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassBFC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassCFC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassCFC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassMPC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassMPC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassCSC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassCSC, true, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassICD, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassICD, true, v[i].ID, val...)
		if err != nil {
			return err
		}
	}

	return nil
}

func getSpecXSync(h *dbxHelper, p string) ([]int64, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func getSpecXAbcd(h *dbxHelper, p string) ([]string, error) {
	c := h.getConn()
	defer h.delConn(c)

	v, err := loadAbcd(c, p, h.lang)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getSpecXAbcdLs(h *dbxHelper, p string) ([]int64, error) {
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

func getSpecXList(h *dbxHelper, p string) (jsonSpecs, error) {
	v, err := jsonToSpecsFromIDs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = loadHashers(c, p, true, v)
	if err != nil {
		return nil, err
	}

	normLang(h.lang, p, v)

	v.sort(h.lang)

	return v, nil
}

func getSpecXListAZ(h *dbxHelper, p string) (jsonSpecs, error) {
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

	h.data = []byte("[" + strings.Join(int64ToStrings(v...), ",") + "]")
	return getSpecXList(h, p)
}

func getSpecX(h *dbxHelper, p string) (jsonSpecs, error) {
	v, err := jsonToSpecsFromIDs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = loadHashers(c, p, false, v)
	if err != nil {
		return nil, err
	}

	err = loadSpecLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	normLang(h.lang, p, v)

	return v, nil
}

func setSpecX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToSpecs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}
	x := makeSpecsFromSpecs(v...)

	c := h.getConn()
	defer h.delConn(c)

	err = loadHashers(c, p, false, x)
	if err != nil {
		return nil, err
	}
	err = loadSpecLinks(c, p, x)
	if err != nil {
		return nil, err
	}
	err = freeSearchers(c, p, x)
	if err != nil {
		return nil, err
	}
	err = freeSpecLinks(c, p, x...)
	if err != nil {
		return nil, err
	}

	err = saveHashers(c, p, v)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, p, v)
	if err != nil {
		return nil, err
	}
	err = saveSpecLinks(c, p, v...)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func delSpecX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToSpecsFromIDs(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = loadHashers(c, p, false, v)
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

	err = freeSpecLinks(c, p, v...)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

// ACT

func getSpecACTSync(h *dbxHelper) (interface{}, error) {
	return getSpecXSync(h, prefixSpecACT)
}

func getSpecACTAbcd(h *dbxHelper) (interface{}, error) {
	return getSpecXAbcd(h, prefixSpecACT)
}

func getSpecACTAbcdLs(h *dbxHelper) (interface{}, error) {
	return getSpecXAbcdLs(h, prefixSpecACT)
}

func getSpecACTList(h *dbxHelper) (interface{}, error) {
	return getSpecXList(h, prefixSpecACT)
}

func getSpecACTListAZ(h *dbxHelper) (interface{}, error) {
	return getSpecXListAZ(h, prefixSpecACT)
}

func getSpecACT(h *dbxHelper) (interface{}, error) {
	return getSpecX(h, prefixSpecACT)
}

func setSpecACT(h *dbxHelper) (interface{}, error) {
	return setSpecX(h, prefixSpecACT)
}

func delSpecACT(h *dbxHelper) (interface{}, error) {
	return delSpecX(h, prefixSpecACT)
}

// INF

func getSpecINFSync(h *dbxHelper) (interface{}, error) {
	return getSpecXSync(h, prefixSpecINF)
}

func getSpecINFAbcd(h *dbxHelper) (interface{}, error) {
	return getSpecXAbcd(h, prefixSpecINF)
}

func getSpecINFAbcdLs(h *dbxHelper) (interface{}, error) {
	return getSpecXAbcdLs(h, prefixSpecINF)
}

func getSpecINFList(h *dbxHelper) (interface{}, error) {
	return getSpecXList(h, prefixSpecINF)
}

func getSpecINFListAZ(h *dbxHelper) (interface{}, error) {
	return getSpecXListAZ(h, prefixSpecINF)
}

func getSpecINF(h *dbxHelper) (interface{}, error) {
	return getSpecX(h, prefixSpecINF)
}

func setSpecINF(h *dbxHelper) (interface{}, error) {
	return setSpecX(h, prefixSpecINF)
}

func delSpecINF(h *dbxHelper) (interface{}, error) {
	return delSpecX(h, prefixSpecINF)
}

// DEC

func getSpecDECSync(h *dbxHelper) (interface{}, error) {
	return getSpecXSync(h, prefixSpecDEC)
}

func getSpecDECAbcd(h *dbxHelper) (interface{}, error) {
	return getSpecXAbcd(h, prefixSpecDEC)
}

func getSpecDECAbcdLs(h *dbxHelper) (interface{}, error) {
	return getSpecXAbcdLs(h, prefixSpecDEC)
}

func getSpecDECList(h *dbxHelper) (interface{}, error) {
	return getSpecXList(h, prefixSpecDEC)
}

func getSpecDECListAZ(h *dbxHelper) (interface{}, error) {
	return getSpecXListAZ(h, prefixSpecDEC)
}

func getSpecDEC(h *dbxHelper) (interface{}, error) {
	return getSpecX(h, prefixSpecDEC)
}

func setSpecDEC(h *dbxHelper) (interface{}, error) {
	return setSpecX(h, prefixSpecDEC)
}

func delSpecDEC(h *dbxHelper) (interface{}, error) {
	return delSpecX(h, prefixSpecDEC)
}
