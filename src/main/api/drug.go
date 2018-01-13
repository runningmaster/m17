package api

import (
	"encoding/json"
	"net/http"
	"sort"

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
	Quant     float64 `json:"q,omitempty"`
	Value     float64 `json:"v,omitempty"`
}

type jsonDrugSale struct {
	ID    int64   `json:"id,omitempty"`
	Quant float64 `json:"q,omitempty"`
	Value float64 `json:"v,omitempty"`
}

func (j *jsonDrug) getID() int64 {
	if j == nil {
		return 0
	}
	return j.ID
}

func (j *jsonDrug) lang(l, _ string) {
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
	case "en":
		j.Name = j.NameEN
		j.Form = j.FormEN
		j.Dose = j.DoseEN
		j.Pack = j.PackEN
		j.Note = j.NoteEN
		j.Make = j.MakeEN
	}

	if l != "" {
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
		"quant",   // 20
		"value",   // 21
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
		j.Quant,  // 20
		j.Value,  // 21
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
		case 20:
			j.Quant, _ = redis.Float64(v[i], nil)
		case 21:
			j.Value, _ = redis.Float64(v[i], nil)
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

func (v jsonDrugs) sort(lang string) {
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
			if v[i].Value > v[j].Value {
				return true
			} else if v[i].Value < v[j].Value {
				return false
			}
			return coll.CompareString(v[i].Name, v[j].Name) < 0
		},
	)
}

func makeDrugsFromJSON(data []byte) (jsonDrugs, error) {
	var v []*jsonDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonDrugs(v), nil
}

func makeDrugsFromIDs(v []int64, err error) (jsonDrugs, error) {
	if err != nil {
		return nil, err
	}
	res := make([]*jsonDrug, len(v))
	for i := range res {
		res[i] = &jsonDrug{ID: v[i]}
	}
	return jsonDrugs(res), nil
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

func getDrugXSync(h *ctxHelper, p string) ([]int64, error) {
	v, err := int64FromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func getDrugX(h *ctxHelper, p string) (jsonDrugs, error) {
	v, err := makeDrugsFromIDs(int64sFromJSON(h.data))
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

	err = loadDrugLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	normLang(h.lang, p, v)

	return v, nil
}

func getDrugXList(h *ctxHelper, p string) (jsonDrugs, error) {
	v, err := makeDrugsFromIDs(int64sFromJSON(h.data))
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

func setDrugX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeDrugsFromJSON(h.data)
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

func setDrugXSale(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeDrugsFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	err = saveHashers(c, p, v, true)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func delDrugX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeDrugsFromIDs(int64sFromJSON(h.data))
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

	return statusOK, nil
}

// DRUG

func getDrugSync(h *ctxHelper) (interface{}, error) {
	return getDrugXSync(h, prefixDrug)
}

func getDrug(h *ctxHelper) (interface{}, error) {
	return getDrugX(h, prefixDrug)
}

func getDrugList(h *ctxHelper) (interface{}, error) {
	return getDrugXList(h, prefixDrug)
}

func setDrug(h *ctxHelper) (interface{}, error) {
	return setDrugX(h, prefixDrug)
}

func setDrugSale(h *ctxHelper) (interface{}, error) {
	return setDrugXSale(h, prefixDrug)
}

func delDrug(h *ctxHelper) (interface{}, error) {
	return delDrugX(h, prefixDrug)
}
