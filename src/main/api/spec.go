package api

import (
	"encoding/json"
	"net/http"

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
	IDSpec     []int64 `json:"id_spec,omitempty"` // *
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
	NameUA     string  `json:"name_ua,omitempty"`
	NameEN     string  `json:"name_en,omitempty"`
	Head       string  `json:"head,omitempty"` // *
	HeadRU     string  `json:"head_ru,omitempty"`
	HeadUA     string  `json:"head_ua,omitempty"`
	HeadEN     string  `json:"head_en,omitempty"`
	Text       string  `json:"text,omitempty"` // *
	TextRU     string  `json:"text_ru,omitempty"`
	TextUA     string  `json:"text_ua,omitempty"`
	TextEN     string  `json:"text_en,omitempty"`
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
		return j.NameRU
	}
	return ""
}

func (j *jsonSpec) getNameUA(p string) string {
	if p == prefixSpecDEC {
		return j.NameUA
	}
	return ""
}

func (j *jsonSpec) marshalToJSON(v interface{}) []byte {
	res, _ := json.Marshal(v)
	return res
}

func (j *jsonSpec) unmarshalFromJSON(b []byte, v interface{}) {
	_ = json.Unmarshal(b, v)
}

func (j *jsonSpec) getFields() []interface{} {
	return []interface{}{
		"id",         // 0
		"name_ru",    // 1
		"name_ua",    // 2
		"name_en",    // 3
		"head_ru",    // 4
		"head_ua",    // 5
		"head_en",    // 6
		"text_ru",    // 7
		"text_ua",    // 8
		"text_en",    // 9
		"slug",       // 10
		"slug_gp",    // 11
		"image_org",  // 12
		"image_box",  // 13
		"created_at", // 14
		"updated_at", // 15
	}
}

func (j *jsonSpec) getValues() []interface{} {
	return []interface{}{
		j.ID,        //0
		j.NameRU,    //1
		j.NameUA,    //2
		j.NameEN,    //3
		j.HeadRU,    //4
		j.HeadUA,    //5
		j.HeadEN,    //6
		j.TextRU,    //7
		j.TextUA,    //8
		j.TextEN,    //9
		j.Slug,      //10
		j.SlugGP,    //11
		j.ImageOrg,  //12
		j.ImageBox,  //13
		j.CreatedAt, //14
		j.UpdatedAt, //15
	}
}

func (j *jsonSpec) setValues(v ...interface{}) {
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
			j.HeadRU, _ = redis.String(v[i], nil)
		case 5:
			j.HeadUA, _ = redis.String(v[i], nil)
		case 6:
			j.HeadEN, _ = redis.String(v[i], nil)
		case 7:
			j.TextRU, _ = redis.String(v[i], nil)
		case 8:
			j.TextUA, _ = redis.String(v[i], nil)
		case 9:
			j.TextEN, _ = redis.String(v[i], nil)
		case 10:
			j.Slug, _ = redis.String(v[i], nil)
		case 11:
			j.SlugGP, _ = redis.String(v[i], nil)
		case 12:
			j.ImageOrg, _ = redis.String(v[i], nil)
		case 13:
			j.ImageBox, _ = redis.String(v[i], nil)
		case 14:
			j.CreatedAt, _ = redis.Int64(v[i], nil)
		case 15:
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

func (j jsonSpecs) nill(i int) {
	j[i] = nil
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

func getSpecXSync(h *dbxHelper, p string, d ...bool) ([]int64, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getSpecX(h *dbxHelper, p string) (jsonSpecs, error) {
	v, err := jsonToSpecsFromIDs(h.data)
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

	err = loadSpecLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func setSpecX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToSpecs(h.data)
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

func getSpecACTSyncDel(h *dbxHelper) (interface{}, error) {
	return getSpecXSync(h, prefixSpecACT)
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

func getSpecINFSyncDel(h *dbxHelper) (interface{}, error) {
	return getSpecXSync(h, prefixSpecINF)
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

func getSpecDECSyncDel(h *dbxHelper) (interface{}, error) {
	return getSpecXSync(h, prefixSpecDEC)
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
