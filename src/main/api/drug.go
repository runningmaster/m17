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
	ID     int64 `json:"id,omitempty"`
	IDMake int64 `json:"id_make,omitempty"`

	IDClassATC []int64 `json:"id_class_atc,omitempty"`
	IDClassNFC []int64 `json:"id_class_nfc,omitempty"`
	IDClassFSC []int64 `json:"id_class_fsc,omitempty"`
	IDClassBFC []int64 `json:"id_class_bfc,omitempty"`
	IDClassCFC []int64 `json:"id_class_cfc,omitempty"`
	IDClassMPC []int64 `json:"id_class_mpc,omitempty"`
	IDClassCSC []int64 `json:"id_class_csc,omitempty"`
	IDClassICD []int64 `json:"id_class_icd,omitempty"`

	IDSpecACT []int64 `json:"id_spec_act,omitempty"`
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`

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

	Sale float64 `json:"sale,omitempty"`
}

func (j *jsonDrug) getID() int64 {
	return j.ID
}

func (j *jsonDrug) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		genKey(p, j.ID),
		"id", j.ID,
		"id_make", j.Make,
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
		genKey(p, j.ID),
		"id",      // 0
		"id_make", // 1
		"name_ru", // 2
		"name_ua", // 3
		"name_en", // 4
		"form_ru", // 5
		"form_ua", // 6
		"form_en", // 7
		"dose_ru", // 8
		"dose_ua", // 9
		"dose_en", // 10
		"pack_ru", // 11
		"pack_ua", // 12
		"pack_en", // 13
		"note_ru", // 14
		"note_ua", // 15
		"note_en", // 16
		"numb",    // 17
		"make_ru", // 18
		"make_ua", // 19
		"make_en", // 20
	}
}

func (j *jsonDrug) setValues(v ...interface{}) bool {
	for i := range v {
		switch i {
		case 0:
			j.ID, _ = redis.Int64(v[i], nil)
		case 1:
			j.IDMake, _ = redis.Int64(v[i], nil)
		case 2:
			j.NameRU, _ = redis.String(v[i], nil)
		case 3:
			j.NameUA, _ = redis.String(v[i], nil)
		case 4:
			j.NameEN, _ = redis.String(v[i], nil)
		case 5:
			j.FormRU, _ = redis.String(v[i], nil)
		case 6:
			j.FormUA, _ = redis.String(v[i], nil)
		case 7:
			j.FormEN, _ = redis.String(v[i], nil)
		case 8:
			j.DoseRU, _ = redis.String(v[i], nil)
		case 9:
			j.DoseUA, _ = redis.String(v[i], nil)
		case 10:
			j.DoseEN, _ = redis.String(v[i], nil)
		case 11:
			j.PackRU, _ = redis.String(v[i], nil)
		case 12:
			j.PackUA, _ = redis.String(v[i], nil)
		case 13:
			j.PackEN, _ = redis.String(v[i], nil)
		case 14:
			j.NoteRU, _ = redis.String(v[i], nil)
		case 15:
			j.NoteUA, _ = redis.String(v[i], nil)
		case 16:
			j.NoteEN, _ = redis.String(v[i], nil)
		case 17:
			j.Numb, _ = redis.String(v[i], nil)
		case 18:
			j.MakeRU, _ = redis.String(v[i], nil)
		case 19:
			j.MakeUA, _ = redis.String(v[i], nil)
		case 20:
			j.MakeEN, _ = redis.String(v[i], nil)
		}
	}
	return j.ID != 0
}

type jsonDrugs []*jsonDrug

func (j jsonDrugs) len() int {
	return len(j)
}

func (j jsonDrugs) elem(i int) interface{} {
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

func jsonToDrugsFromIDs(data []byte) (jsonDrugs, error) {
	v, err := jsonToIDs(data)
	if err != nil {
		return nil, err
	}
	return makeDrugs(v...)
}

func makeDrugs(x ...int64) (jsonDrugs, error) {
	v := make([]*jsonDrug, len(x))
	for i := range v {
		v[i] = &jsonDrug{ID: x[i]}
	}
	return jsonDrugs(v), nil
}

func loadDrugLinks(c redis.Conn, p string, v []*jsonDrug) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
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

func saveDrugLinks(c redis.Conn, p string, v ...*jsonDrug) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		err = saveLinkIDs(c, p, prefixClassATC, v[i].ID, v[i].IDClassATC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassNFC, v[i].ID, v[i].IDClassNFC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassFSC, v[i].ID, v[i].IDClassFSC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassBFC, v[i].ID, v[i].IDClassBFC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassCFC, v[i].ID, v[i].IDClassCFC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassMPC, v[i].ID, v[i].IDClassMPC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassCSC, v[i].ID, v[i].IDClassCSC...)
		if err != nil {
			return err
		}
		err = saveLinkIDs(c, p, prefixClassICD, v[i].ID, v[i].IDClassICD...)
		if err != nil {
			return err
		}
	}
	return nil
}

func freeDrugLinks(c redis.Conn, p string, v ...*jsonDrug) error {
	var val []int64
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		val, _ = loadLinkIDs(c, p, prefixClassATC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassATC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassNFC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassNFC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassFSC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassFSC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassBFC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassBFC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassCFC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassCFC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassMPC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassMPC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassCSC, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassCSC, v[i].ID, val...)
		if err != nil {
			return err
		}

		val, _ = loadLinkIDs(c, p, prefixClassICD, v[i].ID)
		err = freeLinkIDs(c, p, prefixClassICD, v[i].ID, val...)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDrugXSync(h *dbxHelper, p string, d ...bool) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getDrugX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToDrugsFromIDs(h.data)
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

	err = saveDrugLinks(c, p, v...)
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

	err = freeHashers(c, p, v)
	if err != nil {
		return nil, err
	}

	err = freeDrugLinks(c, p, v...)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

// DRUG

func getDrugSync(h *dbxHelper) (interface{}, error) {
	return getDrugXSync(h, prefixDrug)
}

func getDrugSyncDel(h *dbxHelper) (interface{}, error) {
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
