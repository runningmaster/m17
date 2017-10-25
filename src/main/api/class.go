package api

import (
	"encoding/json"
	"net/http"

	"internal/ctxutil"

	"github.com/garyburd/redigo/redis"
)

const (
	prefixClassATC = "class:atc"
	prefixClassNFC = "class:nfc"
	prefixClassFSC = "class:fsc"
	prefixClassBFC = "class:bfc"
	prefixClassCFC = "class:cfc"
	prefixClassMPC = "class:mpc"
	prefixClassCSC = "class:csc"
	prefixClassICD = "class:icd"
)

type jsonClass struct {
	ID     int64 `json:"id,omitempty"`
	IDNode int64 `json:"id_node,omitempty"`
	IDRoot int64 `json:"id_root,omitempty"`

	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?

	Code   string `json:"code,omitempty"`
	Name   string `json:"name,omitempty"` // *
	NameRU string `json:"name_ru,omitempty"`
	NameUA string `json:"name_ua,omitempty"`
	NameEN string `json:"name_en,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

func (j *jsonClass) getID() int64 {
	return j.ID
}

func (j *jsonClass) getNameRU(_ string) string {
	return j.NameRU
}

func (j *jsonClass) getNameUA(_ string) string {
	return j.NameUA
}

func (j *jsonClass) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		genKey(p, j.ID),
		"id", j.ID,
		"id_node", j.IDNode,
		"id_root", j.IDRoot,
		"code", j.Code,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
		"slug", j.Slug,
	}
}

func (j *jsonClass) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		genKey(p, j.ID),
		"id",      // 0
		"id_node", // 1
		"id_root", // 2
		"code",    // 3
		"name_ru", // 4
		"name_ua", // 5
		"name_en", // 6
		"slug",    // 7
	}
}

func (j *jsonClass) setValues(v ...interface{}) bool {
	for i := range v {
		switch i {
		case 0:
			j.ID, _ = redis.Int64(v[i], nil)
		case 1:
			j.IDNode, _ = redis.Int64(v[i], nil)
		case 2:
			j.IDRoot, _ = redis.Int64(v[i], nil)
		case 3:
			j.Code, _ = redis.String(v[i], nil)
		case 4:
			j.NameRU, _ = redis.String(v[i], nil)
		case 5:
			j.NameUA, _ = redis.String(v[i], nil)
		case 6:
			j.NameEN, _ = redis.String(v[i], nil)
		case 7:
			j.Slug, _ = redis.String(v[i], nil)
		}
	}
	return j.ID != 0
}

type jsonClasses []*jsonClass

func (j jsonClasses) len() int {
	return len(j)
}

func (j jsonClasses) elem(i int) interface{} {
	return j[i]
}

func (j jsonClasses) nill(i int) {
	j[i] = nil
}

func jsonToClasses(data []byte) (jsonClasses, error) {
	var v []*jsonClass
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonClasses(v), nil
}

func jsonToClassesFromIDs(data []byte) (jsonClasses, error) {
	v, err := jsonToIDs(data)
	if err != nil {
		return nil, err
	}
	return makeClasses(v...)
}

func makeClasses(x ...int64) (jsonClasses, error) {
	v := make([]*jsonClass, len(x))
	for i := range v {
		v[i] = &jsonClass{ID: x[i]}
	}
	return jsonClasses(v), nil
}

func xxxClassNext(c redis.Conn, cmd string, p string, v ...*jsonClass) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}
		err = c.Send(cmd, genKey(p, v[i].IDNode, "next"), v[i].ID)
		if err != nil {
			return err
		}
	}
	return c.Flush()
}

func setClassNext(c redis.Conn, p string, v ...*jsonClass) error {
	return xxxClassNext(c, "SADD", p, v...)
}

func remClassNext(c redis.Conn, p string, v ...*jsonClass) error {
	return xxxClassNext(c, "SREM", p, v...)
}

func getClassXSync(h *dbxHelper, p string, d ...bool) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getClassX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToClassesFromIDs(h.data)
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

	return v, nil
}

func setClassX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToClasses(h.data)
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

	err = setClassNext(c, p, v...)
	if err != nil {
		return nil, err
	}

	if p == prefixClassATC {
		err = saveSearchers(c, p, v)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, nil
}

func delClassX(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToClassesFromIDs(h.data)
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

	err = remClassNext(c, p, v...)
	if err != nil {
		return nil, err
	}

	if p == prefixClassATC {
		err = freeSearchers(c, p, v)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, nil
}

// ATC

func getClassATCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassATC)
}

func getClassATCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassATC, true)
}

func getClassATC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassATC)
}

func setClassATC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassATC)
}

func delClassATC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassATC)
}

// NFC

func getClassNFCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassNFC)
}

func getClassNFCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassNFC, true)
}

func getClassNFC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassNFC)
}

func setClassNFC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassNFC)
}

func delClassNFC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassNFC)
}

// FSC

func getClassFSCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassFSC)
}

func getClassFSCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassFSC, true)
}

func getClassFSC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassFSC)
}

func setClassFSC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassFSC)
}

func delClassFSC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassFSC)
}

// BFC

func getClassBFCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassBFC)
}

func getClassBFCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassBFC, true)
}

func getClassBFC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassBFC)
}

func setClassBFC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassBFC)
}

func delClassBFC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassBFC)
}

// CFC

func getClassCFCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassCFC)
}

func getClassCFCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassCFC, true)
}

func getClassCFC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassCFC)
}

func setClassCFC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassCFC)
}

func delClassCFC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassCFC)
}

// MPC

func getClassMPCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassMPC)
}

func getClassMPCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassMPC, true)
}

func getClassMPC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassMPC)
}

func setClassMPC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassMPC)
}

func delClassMPC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassMPC)
}

// CSC

func getClassCSCSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassCSC)
}

func getClassCSCSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassCSC, true)
}

func getClassCSC(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassCSC)
}

func setClassCSC(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassCSC)
}

func delClassCSC(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassCSC)
}

// ICD

func getClassICDSync(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassICD)
}

func getClassICDSyncDel(h *dbxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassICD, true)
}

func getClassICD(h *dbxHelper) (interface{}, error) {
	return getClassX(h, prefixClassICD)
}

func setClassICD(h *dbxHelper) (interface{}, error) {
	return setClassX(h, prefixClassICD)
}

func delClassICD(h *dbxHelper) (interface{}, error) {
	return delClassX(h, prefixClassICD)
}
