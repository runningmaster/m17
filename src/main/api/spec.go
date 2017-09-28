package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type jsonSpec struct {
	ID         int64   `json:"id,omitempty"`
	IDINN      []int64 `json:"id_inn,omitempty"`
	IDDrug     []int64 `json:"id_drug,omitempty"`
	IDMake     []int64 `json:"id_make,omitempty"`
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
	ImageOrg   string  `json:"image_org,omitempty"`
	ImageBox   string  `json:"image_box,omitempty"`
	CreatedAt  int64   `json:"created_at,omitempty"`
	UpdatedAt  int64   `json:"updated_at,omitempty"`
}

func (j *jsonSpec) getKey(p string) string {
	return p + ":" + strconv.Itoa(int(j.ID))
}

func (j *jsonSpec) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		p + ":" + "sync",
		"CH",
		time.Now().Unix(),
		j.ID,
	}
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
		"image_org",  // 11
		"image_box",  // 12
		"created_at", // 13
		"updated_at", // 14
	}
}

func (j *jsonSpec) setValues(v ...interface{}) {
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
			j.ImageOrg, _ = redis.String(v[i], nil)
		case 12:
			j.ImageBox, _ = redis.String(v[i], nil)
		case 13:
			j.CreatedAt, _ = redis.Int64(v[i], nil)
		case 14:
			j.UpdatedAt, _ = redis.Int64(v[i], nil)
		}
	}
}

func getSpecACT(h *dbxHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecACTSync(h *dbxHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setSpecACT(h *dbxHelper) (interface{}, error) {
	return "OK", nil
}

func delSpecACT(h *dbxHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecINF(h *dbxHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecINFSync(h *dbxHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setSpecINF(h *dbxHelper) (interface{}, error) {
	return "OK", nil
}

func setSpecINFSale(h *dbxHelper) (interface{}, error) {
	return "OK", nil
}

func delSpecINF(h *dbxHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecDEC(h *dbxHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecDECSync(h *dbxHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setSpecDEC(h *dbxHelper) (interface{}, error) {
	return "OK", nil
}

func setSpecDECSale(h *dbxHelper) (interface{}, error) {
	return "OK", nil
}

func delSpecDEC(h *dbxHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
