package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"internal/logger"

	"github.com/garyburd/redigo/redis"
)

var apiFunc = map[string]func(h *redisHelper) (interface{}, error){
	"get-inn":      getINN,
	"get-inn-sync": getINNSync,
	"set-inn":      setINN,
	"del-inn":      delINN,

	"get-drug":      getDrug,
	"get-drug-sync": getDrugSync,
	"set-drug":      setDrug,
	"set-drug-sale": setDrugSale,
	"del-drug":      delDrug,

	"get-maker":      getMaker,
	"get-maker-sync": getMakerSync,
	"set-maker":      setMaker,
	"del-maker":      delMaker,

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

func getINN(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
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
	return "FIXME", nil
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
	return "FIXME", nil
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
	return "FIXME", nil
}

func setDrugSale(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delDrug(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassATC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassATCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassATC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassATC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassNFC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassNFCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassNFC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassNFC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassFSC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassFSCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassFSC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassFSC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassBFC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassBFCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassBFC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassBFC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassCFC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassCFCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassCFC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassCFC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassMPC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassMPCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassMPC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassMPC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassCSC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassCSCSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassCSC(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassCSC(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassICD(h *redisHelper) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getClassICDSync(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func setClassICD(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delClassICD(h *redisHelper) (interface{}, error) {
	var v []int64
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
	return "FIXME", nil
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
	return "FIXME", nil
}

func setSpecINFSale(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
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
	return "FIXME", nil
}

func setSpecDECSale(h *redisHelper) (interface{}, error) {
	return "FIXME", nil
}

func delSpecDEC(h *redisHelper) (interface{}, error) {
	var v int64
	err := json.Unmarshal(h.data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

type jsonClass struct {
	ID     int64  `json:"id,omitempty"`
	IDNode int64  `json:"id_node,omitempty"`
	IDRoot int64  `json:"id_root,omitempty"`
	Code   string `json:"code,omitempty"`
	Name   string `json:"name,omitempty"` // *
	NameRU string `json:"name_ru,omitempty"`
	NameUA string `json:"name_ua,omitempty"`
	NameEN string `json:"name_en,omitempty"`
}

type jsonINN struct {
	ID        int64   `json:"id,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"` // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
	Name      string  `json:"name,omitempty"` // *
	NameRU    string  `json:"name_ru,omitempty"`
	NameUA    string  `json:"name_ua,omitempty"`
	NameEN    string  `json:"name_en,omitempty"`
	Slug      string  `json:"slug,omitempty"`
}

type jsonMake struct {
	ID        int64   `json:"id,omitempty"`
	IDSpec    []int64 `json:"id_spec,omitempty"` // *
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
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
	ID         int64   `json:"id,omitempty"`
	IDMake     int64   `json:"id_make,omitempty"`
	IDSpecDEC  int64   `json:"id_spec_dec,omitempty"`
	IDSpecINF  int64   `json:"id_spec_inf,omitempty"`
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
	Form       string  `json:"form,omitempty"` // *
	FormRU     string  `json:"form_ru,omitempty"`
	FormUA     string  `json:"form_ua,omitempty"`
	FormEN     string  `json:"form_en,omitempty"`
	Dose       string  `json:"dose,omitempty"` // *
	DoseRU     string  `json:"dose_ru,omitempty"`
	DoseUA     string  `json:"dose_ua,omitempty"`
	DoseEN     string  `json:"dose_en,omitempty"`
	Pack       string  `json:"pack,omitempty"` // *
	PackRU     string  `json:"pack_ru,omitempty"`
	PackUA     string  `json:"pack_ua,omitempty"`
	PackEN     string  `json:"pack_en,omitempty"`
	Note       string  `json:"note,omitempty"` // *
	NoteRU     string  `json:"note_ru,omitempty"`
	NoteUA     string  `json:"note_ua,omitempty"`
	NoteEN     string  `json:"note_en,omitempty"`
	Numb       string  `json:"numb,omitempty"`
	Make       string  `json:"make,omitempty"` // *
	MakeRU     string  `json:"make_ru,omitempty"`
	MakeUA     string  `json:"make_ua,omitempty"`
	MakeEN     string  `json:"make_en,omitempty"`
}

type jsonSpec struct {
	ID         int64   `json:"id,omitempty"`
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
