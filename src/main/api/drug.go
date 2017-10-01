package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixDrug = "drug"
)

type jsonDrug struct {
	ID int64 `json:"id,omitempty"`
	//	IDMake     int64   `json:"id_make,omitempty"`
	//	IDSpecDEC  int64   `json:"id_spec_dec,omitempty"`
	//	IDSpecINF  int64   `json:"id_spec_inf,omitempty"`
	//	IDClassATC []int64 `json:"id_class_atc,omitempty"`
	IDClassNFC []int64 `json:"id_class_nfc,omitempty"`
	//	IDClassFSC []int64 `json:"id_class_fsc,omitempty"`
	//	IDClassBFC []int64 `json:"id_class_bfc,omitempty"`
	//	IDClassCFC []int64 `json:"id_class_cfc,omitempty"`
	//	IDClassMPC []int64 `json:"id_class_mpc,omitempty"`
	//	IDClassCSC []int64 `json:"id_class_csc,omitempty"`
	//	IDClassICD []int64 `json:"id_class_icd,omitempty"`
	Name   string `json:"name,omitempty"` // *
	NameRU string `json:"name_ru,omitempty"`
	NameUA string `json:"name_ua,omitempty"`
	NameEN string `json:"name_en,omitempty"`
	Form   string `json:"form,omitempty"` // *
	FormRU string `json:"form_ru,omitempty"`
	FormUA string `json:"form_ua,omitempty"`
	FormEN string `json:"form_en,omitempty"`
	Dose   string `json:"dose,omitempty"` // *
	DoseRU string `json:"dose_ru,omitempty"`
	DoseUA string `json:"dose_ua,omitempty"`
	DoseEN string `json:"dose_en,omitempty"`
	Pack   string `json:"pack,omitempty"` // *
	PackRU string `json:"pack_ru,omitempty"`
	PackUA string `json:"pack_ua,omitempty"`
	PackEN string `json:"pack_en,omitempty"`
	Note   string `json:"note,omitempty"` // *
	NoteRU string `json:"note_ru,omitempty"`
	NoteUA string `json:"note_ua,omitempty"`
	NoteEN string `json:"note_en,omitempty"`
	Numb   string `json:"numb,omitempty"`
	Make   string `json:"make,omitempty"` // *
	MakeRU string `json:"make_ru,omitempty"`
	MakeUA string `json:"make_ua,omitempty"`
	MakeEN string `json:"make_en,omitempty"`
}

func (j *jsonDrug) getKey(p string) string {
	return p + ":" + strconv.Itoa(int(j.ID))
}

func (j *jsonDrug) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		p + ":" + "sync",
		"CH",
		time.Now().Unix(),
		j.ID,
	}
}

func (j *jsonDrug) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id", j.ID,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
		"form_ru", j.FormRU,
		"form_ua", j.FormUA,
		"form_en", j.FormEN,
		"dose_ru", j.DoseRU,
		"dose_ua", j.DoseUA,
		"dose_en", j.DoseEN,
		"pack_ru", j.PackRU,
		"pack_ua", j.PackUA,
		"pack_en", j.PackEN,
		"note_ru", j.NoteRU,
		"note_ua", j.NoteUA,
		"note_en", j.NoteEN,
		"numb", j.Numb,
		"make_ru", j.MakeRU,
		"make_ua", j.MakeUA,
		"make_en", j.MakeEN,
	}
}

func (j *jsonDrug) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id",      // 0
		"name_ru", // 1
		"name_ua", // 2
		"name_en", // 3
		"form_ru", // 4
		"form_ua", // 5
		"form_en", // 6
		"dose_ru", // 7
		"dose_ua", // 8
		"dose_en", // 9
		"pack_ru", // 10
		"pack_ua", // 11
		"pack_en", // 12
		"note_ru", // 13
		"note_ua", // 14
		"note_en", // 15
		"numb",    // 16
		"make_ru", // 17
		"make_ua", // 18
		"make_en", // 19
	}
}

func (j *jsonDrug) setValues(v ...interface{}) bool {
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
			j.FormRU, _ = redis.String(v[i], nil)
		case 5:
			j.FormUA, _ = redis.String(v[i], nil)
		case 6:
			j.FormEN, _ = redis.String(v[i], nil)
		case 7:
			j.DoseRU, _ = redis.String(v[i], nil)
		case 8:
			j.DoseUA, _ = redis.String(v[i], nil)
		case 9:
			j.DoseEN, _ = redis.String(v[i], nil)
		case 10:
			j.PackRU, _ = redis.String(v[i], nil)
		case 11:
			j.PackUA, _ = redis.String(v[i], nil)
		case 12:
			j.PackEN, _ = redis.String(v[i], nil)
		case 13:
			j.NoteRU, _ = redis.String(v[i], nil)
		case 14:
			j.NoteUA, _ = redis.String(v[i], nil)
		case 15:
			j.NoteEN, _ = redis.String(v[i], nil)
		case 16:
			j.Numb, _ = redis.String(v[i], nil)
		case 17:
			j.MakeRU, _ = redis.String(v[i], nil)
		case 18:
			j.MakeUA, _ = redis.String(v[i], nil)
		case 19:
			j.MakeEN, _ = redis.String(v[i], nil)
		}
	}
	return j.ID != 0
}

type jsonDrugs []*jsonDrug

func (j jsonDrugs) len() int {
	return len(j)
}

func (j jsonDrugs) elem(i int) hasher {
	return j[i]
}

func (j jsonDrugs) nill(i int) {
	j[i] = nil
}

func jsonToDrugs(data []byte) (jsonDrugs, error) {
	var v []*jsonDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonDrugs(v), nil
}

func makeDrugs(v ...int64) jsonDrugs {
	out := make([]*jsonDrug, len(v))
	for i := range out {
		out[i] = &jsonDrug{ID: v[i]}
	}

	return jsonDrugs(out)
}

func getDrug(h *dbxHelper) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	out := makeDrugs(v...)
	err = loadHashers(c, prefixDrug, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getDrugSync(h *dbxHelper) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, prefixDrug, v)
}

func setDrug(h *dbxHelper) (interface{}, error) {
	v, err := jsonToDrugs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = saveHashers(c, prefixDrug, v)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func delDrug(h *dbxHelper) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = freeHashers(c, prefixDrug, makeDrugs(v...))
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}
