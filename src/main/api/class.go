package api

import (
	"encoding/json"
	"time"

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

func (j *jsonClass) getKey(p string) string {
	return genKey(p, j.ID)
}

func (j *jsonClass) getKeyNextAndIDNode(p string) []interface{} {
	return []interface{}{
		genKeyNext(p),
		j.IDNode,
	}
}

func (j *jsonClass) getKeyAndUnixtimeID(p string) []interface{} {
	return []interface{}{
		genKeySync(p),
		"CH",
		time.Now().Unix(),
		j.ID,
	}
}

func (j *jsonClass) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id", j.ID,
		"id_node", j.IDNode,
		"id_root", j.IDRoot,
		"code", j.Code,
		"name_ru", j.NameRU,
		"name_ua", j.NameUA,
		"name_en", j.NameEN,
	}
}

func (j *jsonClass) getKeyAndFields(p string) []interface{} {
	return []interface{}{
		j.getKey(p),
		"id",      // 0
		"id_node", // 1
		"id_root", // 2
		"code",    // 3
		"name_ru", // 4
		"name_ua", // 5
		"name_en", // 6
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
		}
	}
	return j.ID != 0
}

type jsonClasses []*jsonClass

func (j jsonClasses) len() int {
	return len(j)
}

func (j jsonClasses) elem(i int) hasher {
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

func makeClasses(v ...int64) jsonClasses {
	out := make([]*jsonClass, len(v))
	for i := range out {
		out[i] = &jsonClass{ID: v[i]}
	}

	return jsonClasses(out)
}

func cmdClassNext(c redis.Conn, cmd string, p string, v ...*jsonClass) error {
	var err error
	for i := range v {
		if v[i] == nil {
			continue
		}

		err = c.Send(cmd, v[i].getKeyNextAndIDNode(p)...)
		if err != nil {
			return err
		}
	}
	return c.Flush()
}

func setClassNext(c redis.Conn, p string, v ...*jsonClass) error {
	return cmdClassNext(c, "SADD", p, v...)
}

func remClassNext(c redis.Conn, p string, v ...*jsonClass) error {
	return cmdClassNext(c, "SREM", p, v...)
}

func getClass(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	out := makeClasses(v...)
	err = loadHashers(c, p, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getClassSync(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToID(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return loadSyncIDs(c, p, v)
}

func setClass(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToClasses(h.data)
	if err != nil {
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

	return statusOK, nil
}

func delClass(h *dbxHelper, p string) (interface{}, error) {
	v, err := jsonToIDs(h.data)
	if err != nil {
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	out := makeClasses(v...)
	err = freeHashers(c, p, out)
	if err != nil {
		return nil, err
	}

	err = remClassNext(c, p, out...)
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

func getClassATC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassATC)
}

func getClassATCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassATC)
}

func setClassATC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassATC)
}

func delClassATC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassATC)
}

func getClassNFC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassNFC)
}

func getClassNFCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassNFC)
}

func setClassNFC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassNFC)
}

func delClassNFC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassNFC)
}

func getClassFSC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassFSC)
}

func getClassFSCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassFSC)
}

func setClassFSC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassFSC)
}

func delClassFSC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassFSC)
}

func getClassBFC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassBFC)
}

func getClassBFCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassBFC)
}

func setClassBFC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassBFC)
}

func delClassBFC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassBFC)
}

func getClassCFC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassCFC)
}

func getClassCFCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassCFC)
}

func setClassCFC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassCFC)
}

func delClassCFC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassCFC)
}

func getClassMPC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassMPC)
}

func getClassMPCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassMPC)
}

func setClassMPC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassMPC)
}

func delClassMPC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassMPC)
}

func getClassCSC(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassCSC)
}

func getClassCSCSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassCSC)
}

func setClassCSC(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassCSC)
}

func delClassCSC(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassCSC)
}

func getClassICD(h *dbxHelper) (interface{}, error) {
	return getClass(h, prefixClassICD)
}

func getClassICDSync(h *dbxHelper) (interface{}, error) {
	return getClassSync(h, prefixClassICD)
}

func setClassICD(h *dbxHelper) (interface{}, error) {
	return setClass(h, prefixClassICD)
}

func delClassICD(h *dbxHelper) (interface{}, error) {
	return delClass(h, prefixClassICD)
}
