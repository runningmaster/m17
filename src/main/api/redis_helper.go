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
