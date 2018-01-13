package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

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
	ID        int64   `json:"id,omitempty"`
	IDNode    int64   `json:"id_node,omitempty"`
	IDRoot    int64   `json:"id_root,omitempty"`
	IDNext    []int64 `json:"id_next,omitempty"`
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
	Code      string  `json:"code,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Slug      string  `json:"slug,omitempty"`
	Full      bool    `json:"full,omitempty"`
}

func (j *jsonClass) getID() int64 {
	return j.ID
}

func (j *jsonClass) getSrchRU(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	s = append(s, normName(j.NameRU))
	return s, r
}

func (j *jsonClass) getSrchUA(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	s = append(s, normName(j.NameUA))
	return s, r
}

func (j *jsonClass) getSrchEN(_ string) ([]string, []rune) {
	var s []string
	var r []rune
	s = append(s, normName(j.NameEN))
	return s, r
}

func (j *jsonClass) lang(l, _ string) {
	switch l {
	case "ru":
		j.Name = j.NameRU
		j.IDSpecDEC = nil
	case "ua":
		j.Name = j.NameUA
		j.IDSpecINF = nil
	case "en":
		j.Name = j.NameEN
	}

	if l != "" {
		j.NameRU = ""
		j.NameUA = ""
		j.NameEN = ""
	}
}

func (j *jsonClass) getFields(_ bool) []interface{} {
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

func (j *jsonClass) setValues(_ bool, v ...interface{}) {
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

func (j jsonClasses) null(i int) bool {
	return j[i] == nil
}

func (j jsonClasses) nill(i int) {
	j[i] = nil
}

func (v jsonClasses) sort(lang string) {
	coll := newCollator(lang)
	sort.Slice(v,
		func(i, j int) bool {
			if v[i] == nil && v[j] == nil {
				return true
			}
			if v[i] == nil && v[j] != nil {
				return false
			}
			if v[i] != nil && v[j] == nil {
				return true
			}
			return coll.CompareString(v[i].Slug, v[j].Slug) < 0
		},
	)
}

func makeClassesFromJSON(data []byte) (jsonClasses, error) {
	var v []*jsonClass
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return jsonClasses(v), nil
}

func makeClassesFromIDs(v []int64, err error) (jsonClasses, error) {
	if err != nil {
		return nil, err
	}
	res := make([]*jsonClass, len(v))
	for i := range res {
		res[i] = &jsonClass{ID: v[i]}
	}
	return jsonClasses(res), nil
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

func getClassXSync(h *ctxHelper, p string) ([]int64, error) {
	v, err := int64FromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func mineClassRootIDs(c redis.Conn, p string, v []*jsonClass) ([]int64, error) {
	err := loadClassLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	// return if ICD
	if p == prefixClassICD {
		return v[0].IDNext, nil
	}

	if len(v[0].IDNext) == 0 {
		return nil, fmt.Errorf("something wrong with %s roots", p)
	}

	v[0].ID = v[0].IDNext[0]
	// fucking workaround for CFC
	if p == prefixClassCFC {
		v[0].ID = v[0].IDRoot
	}

	err = loadClassLinks(c, p, v)
	if err != nil {
		return nil, err
	}

	return v[0].IDNext, nil
}

func getClassXRoot(h *ctxHelper, p string) (jsonClasses, error) {
	h.data = []byte("[0]")
	v, err := makeClassesFromIDs(int64sFromJSON(h.data))
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	r, err := mineClassRootIDs(c, p, v)
	if err != nil {
		return nil, err
	}

	h.data = int64sToJSON(r)
	return getClassXNext(h, p)
}

func getClassXNext(h *ctxHelper, p string) (jsonClasses, error) {
	v, err := getClassX(h, p)
	if err != nil {
		return nil, err
	}

	v.sort(h.lang)

	return v, nil
}

func getClassXNextByID(h *ctxHelper, p string) (jsonClasses, error) {
	x, err := int64FromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	v, err := makeClassesFromIDs([]int64{x}, nil)
	if err != nil {
		return nil, err
	}

	var r []int64
	if v[0].ID != 0 {
		r, err = loadLinkIDs(c, p, "next", v[0].ID)
		if err != nil {
			return nil, err
		}
	} else {
		r, err = mineClassRootIDs(c, p, v)
		if err != nil {
			return nil, err
		}
	}

	h.data = int64sToJSON(r)
	v, err = getClassXNext(h, p)
	if err != nil {
		return nil, err
	}

	for i := range v {
		v[i].Full = len(v[i].IDSpecINF) > 0 || len(v[i].IDSpecDEC) > 0
		v[i].IDNode = 0
		v[i].IDNext = nil
		v[i].IDSpecINF = nil
		v[i].IDSpecDEC = nil
	}

	return v, nil
}

func getClassX(h *ctxHelper, p string) (jsonClasses, error) {
	v, err := makeClassesFromIDs(int64sFromJSON(h.data))
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

	normLang(h.lang, p, v)

	return v, nil
}

func setClassX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeClassesFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	x, err := makeClassesFromIDs(findExistsIDs(c, p, mineIDsFromHashers(v)...))
	if err != nil {
		return nil, err
	}

	if len(x) > 0 {
		err = loadHashers(c, p, x)
		if err != nil {
			return nil, err
		}
		err = loadClassLinks(c, p, x)
		if err != nil {
			return nil, err
		}
		if p == prefixClassATC {
			err = freeSearchers(c, p, x)
			if err != nil {
				return nil, err
			}
		}
		err = freeClassLinks(c, p, x...)
		if err != nil {
			return nil, err
		}
	}

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

func delClassX(h *ctxHelper, p string) (interface{}, error) {
	v, err := makeClassesFromIDs(int64sFromJSON(h.data))
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

func getClassATCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassATC)
}

func getClassATCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassATC)
}

func getClassATCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassATC)
}

func getClassATCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassATC)
}

func getClassATC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassATC)
}

func setClassATC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassATC)
}

func delClassATC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassATC)
}

// NFC

func getClassNFCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassNFC)
}

func getClassNFCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassNFC)
}

func getClassNFCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassNFC)
}

func getClassNFCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassNFC)
}

func getClassNFC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassNFC)
}

func setClassNFC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassNFC)
}

func delClassNFC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassNFC)
}

// FSC

func getClassFSCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassFSC)
}

func getClassFSCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassFSC)
}

func getClassFSCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassFSC)
}

func getClassFSCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassFSC)
}

func getClassFSC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassFSC)
}

func setClassFSC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassFSC)
}

func delClassFSC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassFSC)
}

// BFC

func getClassBFCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassBFC)
}

func getClassBFCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassBFC)
}

func getClassBFCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassBFC)
}

func getClassBFCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassBFC)
}

func getClassBFC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassBFC)
}

func setClassBFC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassBFC)
}

func delClassBFC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassBFC)
}

// CFC

func getClassCFCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassCFC)
}

func getClassCFCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassCFC)
}

func getClassCFCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassCFC)
}

func getClassCFCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassCFC)
}

func getClassCFC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassCFC)
}

func setClassCFC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassCFC)
}

func delClassCFC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassCFC)
}

// MPC

func getClassMPCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassMPC)
}

func getClassMPCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassMPC)
}

func getClassMPCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassMPC)
}

func getClassMPCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassMPC)
}

func getClassMPC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassMPC)
}

func setClassMPC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassMPC)
}

func delClassMPC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassMPC)
}

// CSC

func getClassCSCSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassCSC)
}

func getClassCSCRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassCSC)
}

func getClassCSCNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassCSC)
}

func getClassCSCNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassCSC)
}

func getClassCSC(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassCSC)
}

func setClassCSC(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassCSC)
}

func delClassCSC(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassCSC)
}

// ICD

func getClassICDSync(h *ctxHelper) (interface{}, error) {
	return getClassXSync(h, prefixClassICD)
}

func getClassICDRoot(h *ctxHelper) (interface{}, error) {
	return getClassXRoot(h, prefixClassICD)
}

func getClassICDNext(h *ctxHelper) (interface{}, error) {
	return getClassXNext(h, prefixClassICD)
}

func getClassICDNextByID(h *ctxHelper) (interface{}, error) {
	return getClassXNextByID(h, prefixClassICD)
}

func getClassICD(h *ctxHelper) (interface{}, error) {
	return getClassX(h, prefixClassICD)
}

func setClassICD(h *ctxHelper) (interface{}, error) {
	return setClassX(h, prefixClassICD)
}

func delClassICD(h *ctxHelper) (interface{}, error) {
	return delClassX(h, prefixClassICD)
}
