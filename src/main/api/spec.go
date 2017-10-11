package api

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixSpecACT = "spec:act" // RU only
	prefixSpecINF = "spec:inf" // RU only
	prefixSpecDEC = "spec:dec" // UA only
)

type slug struct {
	Name   string `json:"name,omitempty"` // *
	NameRU string `json:"name_ru,omitempty"`
	NameUA string `json:"name_ua,omitempty"`
	NameEN string `json:"name_en,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

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
	Slug       string  `json:"slug,omitempty"` // FIXME
	Slugs      []*slug `json:"slugs,omitempty"`
	ImageOrg   string  `json:"image_org,omitempty"`
	ImageBox   string  `json:"image_box,omitempty"`
	CreatedAt  int64   `json:"created_at,omitempty"`
	UpdatedAt  int64   `json:"updated_at,omitempty"`

	Sale float64 `json:"sale,omitempty"`
}

func (j *jsonSpec) getKey(p string) string {
	return genKey(p, j.ID)
}

func (j *jsonSpec) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		genKeySync(p),
		"CH",
		time.Now().Unix(),
		j.ID,
	}
}

func (j *jsonSpec) marshalToJSON(v interface{}) []byte {
	res, _ := json.Marshal(v)
	return res
}

func (j *jsonSpec) unmarshalFromJSON(b []byte, v interface{}) {
	_ = json.Unmarshal(b, v)
}

func (j *jsonSpec) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id", j.ID,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
		"head_ru", j.HeadRU,
		"head_ua", j.HeadUA,
		"head_en", j.HeadEN,
		"text_ru", j.TextRU,
		"text_ua", j.TextUA,
		"text_en", j.TextEN,
		"slug", j.Slug,
		"slugs", j.marshalToJSON(j.Slugs),
		"image_org", j.ImageOrg,
		"image_box", j.ImageBox,
		"created_at", j.CreatedAt,
		"updated_at", j.UpdatedAt,
	}
}

func (j *jsonSpec) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
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
		"slugs",      // 11
		"image_org",  // 12
		"image_box",  // 13
		"created_at", // 14
		"updated_at", // 15
	}
}

func (j *jsonSpec) setValues(v ...interface{}) bool {
	var b []byte
	for i := range v {
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
			b, _ = redis.Bytes(v[i], nil)
			j.unmarshalFromJSON(b, &j.Slugs)
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
	return j.ID != 0
}

type jsonSpecs []*jsonSpec

func (j jsonSpecs) len() int {
	return len(j)
}

func (j jsonSpecs) elem(i int) hasher {
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

func makeSpecs(v ...int64) jsonSpecs {
	out := make([]*jsonSpec, len(v))
	for i := range out {
		out[i] = &jsonSpec{ID: v[i]}
	}

	return jsonSpecs(out)
}

func cmdSpecLinkSendOnly(c redis.Conn, cmd, ps, pl string, x int64, v ...int64) error {
	key := joinKey(genKey(ps, x), pl)
	var err error
	for i := range v {
		err = c.Send(cmd, key, v[i])
		if err != nil {
			return err
		}
		err = c.Send(cmd, joinKey(genKey(pl, v[i]), ps), x)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmdSpecLink(c redis.Conn, cmd string, p string, v ...*jsonSpec) error {
	var err error
	for i := range v {
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixINN, v[i].ID, v[i].IDINN...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixDrug, v[i].ID, v[i].IDDrug...)
		if err != nil {
			return err
		}

		err = cmdSpecLinkSendOnly(c, cmd, p, prefixMaker, v[i].ID, v[i].IDMake...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixSpecDEC, v[i].ID, v[i].IDSpecDEC...)
		if err != nil {
			return err
		}

		err = cmdSpecLinkSendOnly(c, cmd, p, prefixSpecINF, v[i].ID, v[i].IDSpecINF...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassATC, v[i].ID, v[i].IDClassATC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassNFC, v[i].ID, v[i].IDClassNFC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassFSC, v[i].ID, v[i].IDClassFSC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassBFC, v[i].ID, v[i].IDClassBFC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassCFC, v[i].ID, v[i].IDClassCFC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassMPC, v[i].ID, v[i].IDClassMPC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassCSC, v[i].ID, v[i].IDClassCSC...)
		if err != nil {
			return err
		}
		err = cmdSpecLinkSendOnly(c, cmd, p, prefixClassICD, v[i].ID, v[i].IDClassICD...)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func setSpecLink(c redis.Conn, p string, v ...*jsonSpec) error {
	return cmdSpecLink(c, "SADD", p, v...)
}

func remSpecLink(c redis.Conn, p string, v ...*jsonSpec) error {
	return cmdSpecLink(c, "SREM", p, v...)
}

func getSpec(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	out := makeSpecs(v...)
	err = loadHashers(c, p, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getSpecSync(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func setSpec(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToSpecs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = saveHashers(c, p, v)
	if err != nil {
		return nil, err
	}

	err = setSpecLink(c, p, v...)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func delSpec(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = freeHashers(c, p, makeSpecs(v...))
	if err != nil {
		return nil, err
	}

	err = remSpecLink(c, p, makeSpecs(v...)...)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func getSpecACT(h *dbxHelper) (interface{}, error) {
	return getSpec(h, prefixSpecACT)
}

func getSpecACTSync(h *dbxHelper) (interface{}, error) {
	return getSpecSync(h, prefixSpecACT)
}

func setSpecACT(h *dbxHelper) (interface{}, error) {
	return setSpec(h, prefixSpecACT)
}

func delSpecACT(h *dbxHelper) (interface{}, error) {
	return delSpec(h, prefixSpecACT)
}

func getSpecINF(h *dbxHelper) (interface{}, error) {
	return getSpec(h, prefixSpecINF)
}

func getSpecINFSync(h *dbxHelper) (interface{}, error) {
	return getSpecSync(h, prefixSpecINF)
}

func setSpecINF(h *dbxHelper) (interface{}, error) {
	return setSpec(h, prefixSpecINF)
}

func delSpecINF(h *dbxHelper) (interface{}, error) {
	return delSpec(h, prefixSpecINF)
}

func getSpecDEC(h *dbxHelper) (interface{}, error) {
	return getSpec(h, prefixSpecDEC)
}

func getSpecDECSync(h *dbxHelper) (interface{}, error) {
	return getSpecSync(h, prefixSpecDEC)
}

func setSpecDEC(h *dbxHelper) (interface{}, error) {
	return setSpec(h, prefixSpecDEC)
}

func delSpecDEC(h *dbxHelper) (interface{}, error) {
	return delSpec(h, prefixSpecDEC)
}
