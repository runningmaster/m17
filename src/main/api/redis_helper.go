package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"internal/logger"

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

var apiFunc = map[string]func(h *redisHelper) (interface{}, error){
	"get-class-atc":      getClassATC,
	"get-class-atc-sync": getClassATCSync,
	"set-class-atc":      setClassATC,
	"del-class-atc":      delClassATC,

	"get-class-nfc":      getClassNFC,
	"get-class-nfc-sync": getClassNFCSync,
	"set-class-nfc":      setClassNFC,
	"del-class-nfc":      delClassNFC,

	"get-class-fsc":      getClassFSC,
	"get-class-fsc-sync": getClassFSCSync,
	"set-class-fsc":      setClassFSC,
	"del-class-fsc":      delClassFSC,

	"get-class-bfc":      getClassBFC,
	"get-class-bfc-sync": getClassBFCSync,
	"set-class-bfc":      setClassBFC,
	"del-class-bfc":      delClassBFC,

	"get-class-cfc":      getClassCFC,
	"get-class-cfc-sync": getClassCFCSync,
	"set-class-cfc":      setClassCFC,
	"del-class-cfc":      delClassCFC,

	"get-class-mpc":      getClassMPC,
	"get-class-mpc-sync": getClassMPCSync,
	"set-class-mpc":      setClassMPC,
	"del-class-mpc":      delClassMPC,

	"get-class-csc":      getClassCSC,
	"get-class-csc-sync": getClassCSCSync,
	"set-class-csc":      setClassCSC,
	"del-class-csc":      delClassCSC,

	"get-class-icd":      getClassICD,
	"get-class-icd-sync": getClassICDSync,
	"set-class-icd":      setClassICD,
	"del-class-icd":      delClassICD,

	"get-inn":      getINN,
	"get-inn-sync": getINNSync,
	"set-inn":      setINN,
	"del-inn":      delINN,

	"get-maker":      getMaker,
	"get-maker-sync": getMakerSync,
	"set-maker":      setMaker,
	"del-maker":      delMaker,

	"get-drug":      getDrug,
	"get-drug-sync": getDrugSync,
	"set-drug":      setDrug,
	"set-drug-sale": setDrugSale,
	"del-drug":      delDrug,

	"get-spec-act":      getSpecACT,
	"get-spec-act-sync": getSpecACTSync,
	"set-spec-act":      setSpecACT,
	"del-spec-act":      delSpecACT,

	"get-spec-inf":      getSpecINF,
	"get-spec-inf-sync": getSpecINFSync,
	"set-spec-inf":      setSpecINF,
	"set-spec-inf-sale": setSpecINFSale,
	"del-spec-inf":      delSpecINF,

	"get-spec-dec":      getSpecDEC,
	"get-spec-dec-sync": getSpecDECSync,
	"set-spec-dec":      setSpecDEC,
	"set-spec-dec-sale": setSpecDECSale,
	"del-spec-dec":      delSpecDEC,
}

type rediser interface {
	Get() redis.Conn
}

type redisHelper struct {
	ctx  context.Context
	rdb  rediser
	log  logger.Logger
	r    *http.Request
	w    http.ResponseWriter
	meta []byte
	data []byte
	//	lang string
}

func (h *redisHelper) getConn() redis.Conn {
	return h.rdb.Get()
}

func (h *redisHelper) delConn(c io.Closer) {
	_ = c.Close
}

func (h *redisHelper) ping() (interface{}, error) {
	c := h.getConn()
	defer h.delConn(c)

	return redis.Bytes(c.Do("PING"))
}

func (h *redisHelper) exec(s string) (interface{}, error) {
	if fn, ok := apiFunc[s]; ok {
		return fn(h)
	}

	return nil, fmt.Errorf("unknown func %q", s)
}

func (h *redisHelper) getSyncList(p string, v int64) ([]int64, error) {
	c := h.getConn()
	defer h.delConn(c)

	res, err := redis.Values(c.Do("ZRANGEBYSCORE", p+":"+"sync", v, "+inf"))
	if err != nil {
		return nil, err
	}

	out := make([]int64, len(res))
	for i := range res {
		out[i], _ = redis.Int64(res[i], err)
	}

	return out, nil
}

func (h *redisHelper) getHashers(p string, v ruleHasher) error {
	c := h.getConn()
	defer h.delConn(c)

	var err error
	for i := 0; i < v.len(); i++ {
		err = c.Send("HMGET", v.elem(i).getKeyAndFields(p)...)
		if err != nil {
			return err
		}
	}

	err = c.Flush()
	if err != nil {
		return err
	}

	var r []interface{}
	for i := 0; i < v.len(); i++ {
		r, err = redis.Values(c.Receive())
		if err != nil {
			if err == redis.ErrNil {
				v.nill(i)
				continue
			}
			return err
		}
		v.elem(i).setValues(r)
	}

	return nil
}

func (h *redisHelper) setHashers(p string, v ruleHasher) error {
	c := h.getConn()
	defer h.delConn(c)

	var err error
	for i := 0; i < v.len(); i++ {
		err = c.Send("HMSET", v.elem(i).getKeyAndFieldValues(p)...)
		if err != nil {
			return err
		}
		err = c.Send("ZADD", p+":"+"sync", "CH", time.Now().Unix(), v.elem(i).getID())
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

func (h *redisHelper) delHashers(p string, v ...int64) error {
	c := h.getConn()
	defer h.delConn(c)

	var err error
	for i := range v {
		err = c.Send("HDEL", p+":"+strconv.Itoa(int(v[i])))
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

type ruler interface {
	len() int
}

type hasher interface {
	getID() int64
	getKeyAndFieldValues(string) []interface{}
	getKeyAndFields(string) []interface{}
	setValues(...interface{})
}

type ruleHasher interface {
	ruler
	elem(int) hasher
	nill(int)
}

type jsonClassList []*jsonClass

func (j jsonClassList) len() int {
	return len(j)
}

func (j jsonClassList) elem(i int) hasher {
	return j[i]
}

func (j jsonClassList) nill(i int) {
	j[i] = nil
}

type jsonClass struct {
	ID        int64   `json:"id,omitempty"`
	IDNode    int64   `json:"id_node,omitempty"`
	IDRoot    int64   `json:"id_root,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?
	Code      string  `json:"code,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
}

func (j *jsonClass) getID() int64 {
	return j.ID
}

func (j *jsonClass) getKeyAndFieldValues(p string) []interface{} {
	return []interface{}{
		p + ":" + strconv.Itoa(int(j.ID)),
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
		p + ":" + strconv.Itoa(int(j.ID)),
		"id",
		"id_node",
		"id_root",
		"code",
		"name_ru",
		"name_ua",
		"name_en",
	}
}

func (j *jsonClass) setValues(v ...interface{}) {
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
}

func getClass(h *redisHelper, p string) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	return getClassByID(h, p, v...)
}

func getClassByID(h *redisHelper, p string, v ...int64) (interface{}, error) {
	out := make([]*jsonClass, len(v))
	for i := range v {
		out[i].ID = v[i]
	}

	err := h.getHashers(p, jsonClassList(out))
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getClassSync(h *redisHelper, p string) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	l, err := h.getSyncList(p, v)
	if err != nil {
		return nil, err
	}

	return getClassByID(h, p, l...)
}

func setClass(h *redisHelper, p string) (interface{}, error) {
	var v []*jsonClass
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	err = h.setHashers(p, jsonClassList(v))
	if err != nil {
		return nil, err
	}

	return http.StatusText(http.StatusOK), nil
}

func delClass(h *redisHelper, p string) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	err = h.delHashers(p, v...)
	if err != nil {
		return nil, err
	}

	return http.StatusText(http.StatusOK), nil
}

func getClassATC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassATC)
}

func getClassATCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassATC)
}

func setClassATC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassATC)
}

func delClassATC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassATC)
}

func getClassNFC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassNFC)
}

func getClassNFCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassNFC)
}

func setClassNFC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassNFC)
}

func delClassNFC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassNFC)
}

func getClassFSC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassFSC)
}

func getClassFSCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassFSC)
}

func setClassFSC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassFSC)
}

func delClassFSC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassFSC)
}

func getClassBFC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassBFC)
}

func getClassBFCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassBFC)
}

func setClassBFC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassBFC)
}

func delClassBFC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassBFC)
}

func getClassCFC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassCFC)
}

func getClassCFCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassCFC)
}

func setClassCFC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassCFC)
}

func delClassCFC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassCFC)
}

func getClassMPC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassMPC)
}

func getClassMPCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassMPC)
}

func setClassMPC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassMPC)
}

func delClassMPC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassMPC)
}

func getClassCSC(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassCSC)
}

func getClassCSCSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassCSC)
}

func setClassCSC(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassCSC)
}

func delClassCSC(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassCSC)
}

func getClassICD(h *redisHelper) (interface{}, error) {
	return getClass(h, prefixClassICD)
}

func getClassICDSync(h *redisHelper) (interface{}, error) {
	return getClassSync(h, prefixClassICD)
}

func setClassICD(h *redisHelper) (interface{}, error) {
	return setClass(h, prefixClassICD)
}

func delClassICD(h *redisHelper) (interface{}, error) {
	return delClass(h, prefixClassICD)
}

func getINN(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}

	res := make([]*jsonINN, len(v))
	// REDIS
	return res, nil
}

func getINNSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setINN(h *redisHelper) (interface{}, error) {
	var v []*jsonINN
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	for i := range v {
		_ = v[i].ID
	}
	return "OK", nil
}

func delINN(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getMaker(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getMakerSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setMaker(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func delMaker(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getDrug(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getDrugSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setDrug(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func setDrugSale(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func delDrug(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecACT(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecACTSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setSpecACT(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func delSpecACT(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecINF(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecINFSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setSpecINF(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func setSpecINFSale(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func delSpecINF(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecDEC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getSpecDECSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setSpecDEC(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func setSpecDECSale(h *redisHelper) (interface{}, error) {
	return "OK", nil
}

func delSpecDEC(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

type jsonINN struct {
	ID        int64   `json:"id,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?
	Name      string  `json:"name,omitempty"`        // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

type jsonMake struct {
	ID        int64   `json:"id,omitempty"`
	IDNode    int64   `json:"id_node,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"`     // ? // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"` // ?
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"` // ?
	ISComp    bool    `json:"is_comp,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Text      string  `json:"text,omitempty"` // *
	TextRU    string  `json:"text_ru,omitempty"`
	TextUA    string  `json:"text_ua,omitempty"`
	TextEN    string  `json:"text_en,omitempty"`
	Logo      string  `json:"logo,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

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

type jsonSale struct {
	ID int64   `json:"id,omitempty"`
	Q  float64 `json:"q,omitempty"`
	V  float64 `json:"v,omitempty"`
}
