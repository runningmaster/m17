package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	ID     int64   `json:"id,omitempty"`
	IDNode int64   `json:"id_node,omitempty"`
	IDRoot int64   `json:"id_root,omitempty"`
	IDNext []int64 `json:"id_next,omitempty"`

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
	return j.NameRU + "|" + j.Code
}

func (j *jsonClass) getNameUA(_ string) string {
	return j.NameUA + "|" + j.Code
}

func (j *jsonClass) getFields() []interface{} {
	return []interface{}{
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

func (j *jsonClass) getValues() []interface{} {
	return []interface{}{
		j.ID,     // 0
		j.IDNode, // 1
		j.IDRoot, // 2
		j.Code,   // 3
		j.NameRU, // 4
		j.NameUA, // 5
		j.NameEN, // 6
		j.Slug,   // 7
	}
}

func (j *jsonClass) setValues(v ...interface{}) {
	for i := range v {
		if v[i] == nil {
			continue
		}
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
	return makeClasses(v...), nil
}

func makeClasses(x ...int64) jsonClasses {
	v := make([]*jsonClass, len(x))
	for i := range v {
		v[i] = &jsonClass{ID: x[i]}
	}
	return jsonClasses(v)
}

func loadClassLinks(c redis.Conn, p string, v []*jsonClass) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		v[i].IDNext, err = loadLinkIDs(c, p, "next", v[i].ID)
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
	}
	return nil
}

func saveClassLinks(c redis.Conn, p string, v ...*jsonClass) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		err = saveLinkIDs(c, "next", p, false, v[i].ID, v[i].IDNode)
		if err != nil {
			return err
		}
	}
	return nil
}

func freeClassLinks(c redis.Conn, p string, v ...*jsonClass) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		err = freeLinkIDs(c, "next", p, false, v[i].ID, v[i].IDNode)
		if err != nil {
			return err
		}
	}
	return nil
}

func getClassXSync(h *dbxHelper, p string, d ...bool) ([]int64, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v, d...)
}

func getClassXRoot(h *dbxHelper, p string) (jsonClasses, error) {
	c := h.getConn()
	defer h.delConn(c)

	r, err := loadLinkIDs(c, p, "next", 0)
	if err != nil {
		return nil, err
	}

	if len(r) == 0 {
		return nil, fmt.Errorf("something wrong with root %s", p)
	}

	h.data = []byte("[" + strings.Join(int64ToStrings(r...), ",") + "]")
	n, err := getClassX(h, p)

	if len(n) != 1 {
		return nil, fmt.Errorf("something wrong with root %s (%d)", p, len(n))
	}

	h.data = []byte("[" + strings.Join(int64ToStrings(n[0].IDNext...), ",") + "]")
	return getClassX(h, p)
}

func getClassX(h *dbxHelper, p string) (jsonClasses, error) {
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

	err = loadClassLinks(c, p, v)
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

	if p == prefixClassATC {
		err = saveSearchers(c, p, v)
		if err != nil {
			return nil, err
		}
	}

	err = saveClassLinks(c, p, v...)
	if err != nil {
		return nil, err
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

	if p == prefixClassATC {
		err = freeSearchers(c, p, v)
		if err != nil {
			return nil, err
		}
	}

	err = freeClassLinks(c, p, v...)
	if err != nil {
		return nil, err
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

func getClassATCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassATC)
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

func getClassNFCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassNFC)
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

func getClassFSCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassFSC)
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

func getClassBFCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassBFC)
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

func getClassCFCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassCFC)
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

func getClassMPCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassMPC)
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

func getClassCSCRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassCSC)
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

func getClassICDRoot(h *dbxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassICD)
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
