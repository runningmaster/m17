package api

import (
	"encoding/json"
	"net/http"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixDrug = "drug"
)

type jsonDrug struct {
	ID        int64   `json:"id,omitempty"`
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Form      string  `json:"form,omitempty"` // *
	FormRU    string  `json:"form_ru,omitempty"`
	FormUA    string  `json:"form_ua,omitempty"`
	FormEN    string  `json:"form_en,omitempty"`
	Dose      string  `json:"dose,omitempty"` // *
	DoseRU    string  `json:"dose_ru,omitempty"`
	DoseUA    string  `json:"dose_ua,omitempty"`
	DoseEN    string  `json:"dose_en,omitempty"`
	Pack      string  `json:"pack,omitempty"` // *
	PackRU    string  `json:"pack_ru,omitempty"`
	PackUA    string  `json:"pack_ua,omitempty"`
	PackEN    string  `json:"pack_en,omitempty"`
	Note      string  `json:"note,omitempty"` // *
	NoteRU    string  `json:"note_ru,omitempty"`
	NoteUA    string  `json:"note_ua,omitempty"`
	NoteEN    string  `json:"note_en,omitempty"`
	Numb      string  `json:"numb,omitempty"`
	Make      string  `json:"make,omitempty"` // *
	MakeRU    string  `json:"make_ru,omitempty"`
	MakeUA    string  `json:"make_ua,omitempty"`
	MakeEN    string  `json:"make_en,omitempty"`

	Sale float64 `json:"sale,omitempty"`
}

func (j *jsonDrug) getID() int64 {
	if j == nil {
		return 0
	}
	return j.ID
}

func (j *jsonDrug) getNameRU(_ string) string {
	if j == nil {
		return ""
	}
	return j.NameRU
}

func (j *jsonDrug) getNameUA(_ string) string {
	if j == nil {
		return ""
	}
	return j.NameUA
}

func (j *jsonDrug) getNameEN(_ string) string {
	if j == nil {
		return ""
	}
	return j.NameEN
}

func (j *jsonDrug) lang(l, _ string) {
	if j == nil {
		return
	}
	switch l {
	case "ru":
		j.Name = j.NameRU
		j.Form = j.FormRU
		j.Dose = j.DoseRU
		j.Pack = j.PackRU
		j.Note = j.NoteRU
		j.Make = j.MakeRU
		j.IDSpecDEC = nil
	case "ua":
		j.Name = j.NameUA
		j.Form = j.FormRU
		j.Dose = j.DoseRU
		j.Pack = j.PackRU
		j.Note = j.NoteRU
		j.Make = j.MakeRU
		j.IDSpecINF = nil
	}

	if l == "ru" || l == "ua" {
		j.NameRU = ""
		j.NameUA = ""
		j.NameEN = ""
		j.FormRU = ""
		j.FormUA = ""
		j.FormEN = ""
		j.DoseRU = ""
		j.DoseUA = ""
		j.DoseEN = ""
		j.PackRU = ""
		j.PackUA = ""
		j.PackEN = ""
		j.NoteRU = ""
		j.NoteUA = ""
		j.NoteEN = ""
		j.MakeRU = ""
		j.MakeUA = ""
		j.MakeEN = ""
	}
}

func (j *jsonDrug) getFields(_ bool) []interface{} {
	if j == nil {
		return nil
	}
	return []interface{}{
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

func (j *jsonDrug) getValues() []interface{} {
	if j == nil {
		return nil
	}
	return []interface{}{
		j.ID,     // 0
		j.NameRU, // 1
		j.NameUA, // 2
		j.NameEN, // 3
		j.FormRU, // 4
		j.FormUA, // 5
		j.FormEN, // 6
		j.DoseRU, // 7
		j.DoseUA, // 8
		j.DoseEN, // 9
		j.PackRU, // 10
		j.PackUA, // 11
		j.PackEN, // 12
		j.NoteRU, // 13
		j.NoteUA, // 14
		j.NoteEN, // 15
		j.Numb,   // 16
		j.MakeRU, // 17
		j.MakeUA, // 18
		j.MakeEN, // 19
	}
}

func (j *jsonDrug) setValues(_ bool, v ...interface{}) {
	if j == nil {
		return
	}
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
}

type jsonDrugs []*jsonDrug

func (j jsonDrugs) len() int {
	return len(j)
}

func (j jsonDrugs) elem(i int) interface{} {
	return j[i]
}

func (j jsonDrugs) null(i int) bool {
	return j[i] == nil
}

func (j jsonDrugs) nill(i int) {
	j[i] = nil
}

func (v jsonDrugs) sort(_ string) {
	//coll := newCollator(lang)
	//sort.Slice(v,
	//	func(i, j int) bool {
	//		if v[i] == nil || v[j] == nil {
	//			return true
	//		}
	//		return coll.CompareString(v[i].Slug, v[j].Slug) < 0
	//	},
	//)
}

func jsonToDrugs(data []byte) (jsonDrugs, error) {
	var v []*jsonDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonDrugs(v), nil
}

func jsonToDrugsFromIDs(data []byte) (jsonDrugs, error) {
	v, err := jsonToIDs(data)
	if err != nil {
		return nil, err
	}
	return makeDrugs(v...), nil
}

func makeDrugs(x ...int64) jsonDrugs {
	v := make([]*jsonDrug, len(x))
	for i := range v {
		v[i] = &jsonDrug{ID: x[i]}
	}
	return jsonDrugs(v)
}

func loadDrugLinks(c redis.Conn, p string, v []*jsonDrug) error {
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

func getDrugXSync(h *dbxHelper, p string) ([]int64, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func getDrugX(h *dbxHelper, p string) (jsonDrugs, error) {
	v, err := jsonToDrugsFromIDs(h.data)
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

	err = loadDrugLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	normLang(h.lang, p, v)

	return v, nil
}

func setDrugX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToDrugs(h.data)
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

	return statusOK, nil
}

func delDrugX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToDrugsFromIDs(h.data)
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

	return statusOK, nil
}

// DRUG

func getDrugSync(h *dbxHelper) (interface{}, error) {
	return getDrugXSync(h, prefixDrug)
}

func getDrug(h *dbxHelper) (interface{}, error) {
	return getDrugX(h, prefixDrug)
}

func setDrug(h *dbxHelper) (interface{}, error) {
	return setDrugX(h, prefixDrug)
}

func delDrug(h *dbxHelper) (interface{}, error) {
	return delDrugX(h, prefixDrug)
}

func setDrugSale(h *dbxHelper) (interface{}, error) {
	/*
		var v []struct{
			ID
		}
		err := json.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}

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
	*/
	return statusOK, nil
}
