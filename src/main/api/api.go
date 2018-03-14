package api

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"

	"internal/ctxutil"
	"internal/logger"
	"internal/mdware"
	"internal/router"

	"github.com/garyburd/redigo/redis"
	"github.com/nats-io/nuid"
	//"github.com/rogpeppe/fastuuid"
)

type handler struct {
	api    map[string]http.Handler
	rdb    rediser
	log    logger.Logger
	err404 http.Handler
	err405 http.Handler
}

func (h *handler) prepareAPI() *handler {
	pipe := mdware.NewPipe(8)

	pipe.BeforeJoin(
		mdware.Head(uuid),
		mdware.Auth(auth),
		mdware.Gzip,
		mdware.Body,
	)

	pipe.AfterJoin(
		mdware.Resp,
		mdware.Fail,
		mdware.Tail(h.log),
	)

	h.api = map[string]http.Handler{
		"GET /":            pipe.Join(mdware.Exec(home)),
		"GET /help":        pipe.Join(mdware.Exec(help())),
		"GET /:foo/bar":    pipe.Join(mdware.Exec(test)),
		"GET /test/:foo":   pipe.Join(mdware.Exec(test)),
		"GET /redis/ping":  pipe.Join(mdware.Exec(ping(h.rdb))),
		"POST /redis/ping": pipe.Join(mdware.Exec(ping(h.rdb))),

		// FIXME GET POST /
		"POST /get-class-atc-sync":       pipe.Join(mdware.Exec(exec(h, getClassATCSync))),
		"POST /get-class-atc-root":       pipe.Join(mdware.Exec(exec(h, getClassATCRoot))),
		"POST /get-class-atc-next":       pipe.Join(mdware.Exec(exec(h, getClassATCNext))),
		"POST /get-class-atc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassATCNextByID))),
		"POST /get-class-atc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassATCPathByID))),
		"POST /get-class-atc":            pipe.Join(mdware.Exec(exec(h, getClassATC))),
		"POST /set-class-atc":            pipe.Join(mdware.Exec(exec(h, setClassATC))),
		"POST /del-class-atc":            pipe.Join(mdware.Exec(exec(h, delClassATC))),

		"POST /get-class-nfc-sync":       pipe.Join(mdware.Exec(exec(h, getClassNFCSync))),
		"POST /get-class-nfc-root":       pipe.Join(mdware.Exec(exec(h, getClassNFCRoot))),
		"POST /get-class-nfc-next":       pipe.Join(mdware.Exec(exec(h, getClassNFCNext))),
		"POST /get-class-nfc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassNFCNextByID))),
		"POST /get-class-nfc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassNFCPathByID))),
		"POST /get-class-nfc":            pipe.Join(mdware.Exec(exec(h, getClassNFC))),
		"POST /set-class-nfc":            pipe.Join(mdware.Exec(exec(h, setClassNFC))),
		"POST /del-class-nfc":            pipe.Join(mdware.Exec(exec(h, delClassNFC))),

		"POST /get-class-fsc-sync":       pipe.Join(mdware.Exec(exec(h, getClassFSCSync))),
		"POST /get-class-fsc-root":       pipe.Join(mdware.Exec(exec(h, getClassFSCRoot))),
		"POST /get-class-fsc-next":       pipe.Join(mdware.Exec(exec(h, getClassFSCNext))),
		"POST /get-class-fsc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassFSCNextByID))),
		"POST /get-class-fsc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassFSCPathByID))),
		"POST /get-class-fsc":            pipe.Join(mdware.Exec(exec(h, getClassFSC))),
		"POST /set-class-fsc":            pipe.Join(mdware.Exec(exec(h, setClassFSC))),
		"POST /del-class-fsc":            pipe.Join(mdware.Exec(exec(h, delClassFSC))),

		"POST /get-class-bfc-sync":       pipe.Join(mdware.Exec(exec(h, getClassBFCSync))),
		"POST /get-class-bfc-root":       pipe.Join(mdware.Exec(exec(h, getClassBFCRoot))),
		"POST /get-class-bfc-next":       pipe.Join(mdware.Exec(exec(h, getClassBFCNext))),
		"POST /get-class-bfc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassBFCNextByID))),
		"POST /get-class-bfc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassBFCPathByID))),
		"POST /get-class-bfc":            pipe.Join(mdware.Exec(exec(h, getClassBFC))),
		"POST /set-class-bfc":            pipe.Join(mdware.Exec(exec(h, setClassBFC))),
		"POST /del-class-bfc":            pipe.Join(mdware.Exec(exec(h, delClassBFC))),

		"POST /get-class-cfc-sync":       pipe.Join(mdware.Exec(exec(h, getClassCFCSync))),
		"POST /get-class-cfc-root":       pipe.Join(mdware.Exec(exec(h, getClassCFCRoot))),
		"POST /get-class-cfc-next":       pipe.Join(mdware.Exec(exec(h, getClassCFCNext))),
		"POST /get-class-cfc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassCFCNextByID))),
		"POST /get-class-cfc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassCFCPathByID))),
		"POST /get-class-cfc":            pipe.Join(mdware.Exec(exec(h, getClassCFC))),
		"POST /set-class-cfc":            pipe.Join(mdware.Exec(exec(h, setClassCFC))),
		"POST /del-class-cfc":            pipe.Join(mdware.Exec(exec(h, delClassCFC))),

		"POST /get-class-mpc-sync":       pipe.Join(mdware.Exec(exec(h, getClassMPCSync))),
		"POST /get-class-mpc-root":       pipe.Join(mdware.Exec(exec(h, getClassMPCRoot))),
		"POST /get-class-mpc-next":       pipe.Join(mdware.Exec(exec(h, getClassMPCNext))),
		"POST /get-class-mpc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassMPCNextByID))),
		"POST /get-class-mpc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassMPCPathByID))),
		"POST /get-class-mpc":            pipe.Join(mdware.Exec(exec(h, getClassMPC))),
		"POST /set-class-mpc":            pipe.Join(mdware.Exec(exec(h, setClassMPC))),
		"POST /del-class-mpc":            pipe.Join(mdware.Exec(exec(h, delClassMPC))),

		"POST /get-class-csc-sync":       pipe.Join(mdware.Exec(exec(h, getClassCSCSync))),
		"POST /get-class-csc-root":       pipe.Join(mdware.Exec(exec(h, getClassCSCRoot))),
		"POST /get-class-csc-next":       pipe.Join(mdware.Exec(exec(h, getClassCSCNext))),
		"POST /get-class-csc-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassCSCNextByID))),
		"POST /get-class-csc-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassCSCPathByID))),
		"POST /get-class-csc":            pipe.Join(mdware.Exec(exec(h, getClassCSC))),
		"POST /set-class-csc":            pipe.Join(mdware.Exec(exec(h, setClassCSC))),
		"POST /del-class-csc":            pipe.Join(mdware.Exec(exec(h, delClassCSC))),

		"POST /get-class-icd-sync":       pipe.Join(mdware.Exec(exec(h, getClassICDSync))),
		"POST /get-class-icd-root":       pipe.Join(mdware.Exec(exec(h, getClassICDRoot))),
		"POST /get-class-icd-next":       pipe.Join(mdware.Exec(exec(h, getClassICDNext))),
		"POST /get-class-icd-next-by-id": pipe.Join(mdware.Exec(exec(h, getClassICDNextByID))),
		"POST /get-class-icd-path-by-id": pipe.Join(mdware.Exec(exec(h, getClassICDPathByID))),
		"POST /get-class-icd":            pipe.Join(mdware.Exec(exec(h, getClassICD))),
		"POST /set-class-icd":            pipe.Join(mdware.Exec(exec(h, setClassICD))),
		"POST /del-class-icd":            pipe.Join(mdware.Exec(exec(h, delClassICD))),

		"POST /get-inn-sync":    pipe.Join(mdware.Exec(exec(h, getINNSync))),
		"POST /get-inn-abcd":    pipe.Join(mdware.Exec(exec(h, getINNAbcd))),
		"POST /get-inn-abcd-ls": pipe.Join(mdware.Exec(exec(h, getINNAbcdLs))),
		"POST /get-inn-list":    pipe.Join(mdware.Exec(exec(h, getINNList))),
		"POST /get-inn-list-az": pipe.Join(mdware.Exec(exec(h, getINNListAZ))),
		"POST /get-inn":         pipe.Join(mdware.Exec(exec(h, getINN))),
		"POST /set-inn":         pipe.Join(mdware.Exec(exec(h, setINN))),
		"POST /del-inn":         pipe.Join(mdware.Exec(exec(h, delINN))),

		"POST /get-maker-sync":    pipe.Join(mdware.Exec(exec(h, getMakerSync))),
		"POST /get-maker-abcd":    pipe.Join(mdware.Exec(exec(h, getMakerAbcd))),
		"POST /get-maker-abcd-ls": pipe.Join(mdware.Exec(exec(h, getMakerAbcdLs))),
		"POST /get-maker-list":    pipe.Join(mdware.Exec(exec(h, getMakerList))),
		"POST /get-maker-list-az": pipe.Join(mdware.Exec(exec(h, getMakerListAZ))),
		"POST /get-maker":         pipe.Join(mdware.Exec(exec(h, getMaker))),
		"POST /set-maker":         pipe.Join(mdware.Exec(exec(h, setMaker))),
		"POST /del-maker":         pipe.Join(mdware.Exec(exec(h, delMaker))),

		"POST /get-drug-sync": pipe.Join(mdware.Exec(exec(h, getDrugSync))),
		"POST /get-drug":      pipe.Join(mdware.Exec(exec(h, getDrug))),
		"POST /get-drug-list": pipe.Join(mdware.Exec(exec(h, getDrugList))),
		"POST /set-drug":      pipe.Join(mdware.Exec(exec(h, setDrug))),
		"POST /set-drug-sale": pipe.Join(mdware.Exec(exec(h, setDrugSale))),
		"POST /del-drug":      pipe.Join(mdware.Exec(exec(h, delDrug))),

		"POST /get-spec-act-sync":      pipe.Join(mdware.Exec(exec(h, getSpecACTSync))),
		"POST /get-spec-act-abcd":      pipe.Join(mdware.Exec(exec(h, getSpecACTAbcd))),
		"POST /get-spec-act-abcd-ls":   pipe.Join(mdware.Exec(exec(h, getSpecACTAbcdLs))),
		"POST /get-spec-act-list":      pipe.Join(mdware.Exec(exec(h, getSpecACTList))),
		"POST /get-spec-act-list-az":   pipe.Join(mdware.Exec(exec(h, getSpecACTListAZ))),
		"POST /get-spec-act":           pipe.Join(mdware.Exec(exec(h, getSpecACT))),
		"POST /get-spec-act-with-deps": pipe.Join(mdware.Exec(exec(h, getSpecACTWithDeps))),
		"POST /set-spec-act":           pipe.Join(mdware.Exec(exec(h, setSpecACT))),
		"POST /del-spec-act":           pipe.Join(mdware.Exec(exec(h, delSpecACT))),

		"POST /get-spec-inf-sync":                      pipe.Join(mdware.Exec(exec(h, getSpecINFSync))),
		"POST /get-spec-inf-abcd":                      pipe.Join(mdware.Exec(exec(h, getSpecINFAbcd))),
		"POST /get-spec-inf-abcd-ls":                   pipe.Join(mdware.Exec(exec(h, getSpecINFAbcdLs))),
		"POST /get-spec-inf-list":                      pipe.Join(mdware.Exec(exec(h, getSpecINFList))),
		"POST /get-spec-inf-list-az":                   pipe.Join(mdware.Exec(exec(h, getSpecINFListAZ))),
		"POST /get-spec-inf-list-by-id-class-atc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassATC))),
		"POST /get-spec-inf-list-by-id-class-atc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassATCDeep))),
		"POST /get-spec-inf-list-by-id-class-nfc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassNFC))),
		"POST /get-spec-inf-list-by-id-class-nfc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassNFCDeep))),
		"POST /get-spec-inf-list-by-id-class-fsc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassFSC))),
		"POST /get-spec-inf-list-by-id-class-fsc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassFSCDeep))),
		"POST /get-spec-inf-list-by-id-class-bfc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassBFC))),
		"POST /get-spec-inf-list-by-id-class-bfc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassBFCDeep))),
		"POST /get-spec-inf-list-by-id-class-cfc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassCFC))),
		"POST /get-spec-inf-list-by-id-class-cfc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassCFCDeep))),
		"POST /get-spec-inf-list-by-id-class-mpc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassMPC))),
		"POST /get-spec-inf-list-by-id-class-mpc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassMPCDeep))),
		"POST /get-spec-inf-list-by-id-class-csc":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassCSC))),
		"POST /get-spec-inf-list-by-id-class-csc-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassCSCDeep))),
		"POST /get-spec-inf-list-by-id-class-icd":      pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassICD))),
		"POST /get-spec-inf-list-by-id-class-icd-deep": pipe.Join(mdware.Exec(exec(h, getSpecINFListByClassICDDeep))),
		"POST /get-spec-inf-list-by-id-inn":            pipe.Join(mdware.Exec(exec(h, getSpecINFListByINN))),
		"POST /get-spec-inf-list-by-id-maker":          pipe.Join(mdware.Exec(exec(h, getSpecINFListByMaker))),
		"POST /get-spec-inf-list-by-id-drug":           pipe.Join(mdware.Exec(exec(h, getSpecINFListByDrug))),
		"POST /get-spec-inf-list-by-id-spec-act":       pipe.Join(mdware.Exec(exec(h, getSpecINFListBySpecACT))),
		"POST /get-spec-inf-list-by-id-spec-dec":       pipe.Join(mdware.Exec(exec(h, getSpecINFListBySpecDEC))),
		"POST /get-spec-inf":                           pipe.Join(mdware.Exec(exec(h, getSpecINF))),
		"POST /get-spec-inf-with-deps":                 pipe.Join(mdware.Exec(exec(h, getSpecINFWithDeps))),
		"POST /set-spec-inf":                           pipe.Join(mdware.Exec(exec(h, setSpecINF))),
		"POST /set-spec-inf-sale":                      pipe.Join(mdware.Exec(exec(h, setSpecINFSale))),
		"POST /del-spec-inf":                           pipe.Join(mdware.Exec(exec(h, delSpecINF))),

		"POST /get-spec-dec-sync":                      pipe.Join(mdware.Exec(exec(h, getSpecDECSync))),
		"POST /get-spec-dec-abcd":                      pipe.Join(mdware.Exec(exec(h, getSpecDECAbcd))),
		"POST /get-spec-dec-abcd-ls":                   pipe.Join(mdware.Exec(exec(h, getSpecDECAbcdLs))),
		"POST /get-spec-dec-list":                      pipe.Join(mdware.Exec(exec(h, getSpecDECList))),
		"POST /get-spec-dec-list-az":                   pipe.Join(mdware.Exec(exec(h, getSpecDECListAZ))),
		"POST /get-spec-dec-list-by-id-class-atc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassATC))),
		"POST /get-spec-dec-list-by-id-class-atc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassATCDeep))),
		"POST /get-spec-dec-list-by-id-class-nfc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassNFC))),
		"POST /get-spec-dec-list-by-id-class-nfc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassNFCDeep))),
		"POST /get-spec-dec-list-by-id-class-fsc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassFSC))),
		"POST /get-spec-dec-list-by-id-class-fsc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassFSCDeep))),
		"POST /get-spec-dec-list-by-id-class-bfc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassBFC))),
		"POST /get-spec-dec-list-by-id-class-bfc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassBFCDeep))),
		"POST /get-spec-dec-list-by-id-class-cfc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassCFC))),
		"POST /get-spec-dec-list-by-id-class-cfc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassCFCDeep))),
		"POST /get-spec-dec-list-by-id-class-mpc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassMPC))),
		"POST /get-spec-dec-list-by-id-class-mpc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassMPCDeep))),
		"POST /get-spec-dec-list-by-id-class-csc":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassCSC))),
		"POST /get-spec-dec-list-by-id-class-csc-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassCSCDeep))),
		"POST /get-spec-dec-list-by-id-class-icd":      pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassICD))),
		"POST /get-spec-dec-list-by-id-class-icd-deep": pipe.Join(mdware.Exec(exec(h, getSpecDECListByClassICDDeep))),
		"POST /get-spec-dec-list-by-id-inn":            pipe.Join(mdware.Exec(exec(h, getSpecDECListByINN))),
		"POST /get-spec-dec-list-by-id-maker":          pipe.Join(mdware.Exec(exec(h, getSpecDECListByMaker))),
		"POST /get-spec-dec-list-by-id-drug":           pipe.Join(mdware.Exec(exec(h, getSpecDECListByDrug))),
		"POST /get-spec-dec-list-by-id-spec-act":       pipe.Join(mdware.Exec(exec(h, getSpecDECListBySpecACT))),
		"POST /get-spec-dec-list-by-id-spec-inf":       pipe.Join(mdware.Exec(exec(h, getSpecDECListBySpecINF))),
		"POST /get-spec-dec":                           pipe.Join(mdware.Exec(exec(h, getSpecDEC))),
		"POST /get-spec-dec-with-deps":                 pipe.Join(mdware.Exec(exec(h, getSpecDECWithDeps))),
		"POST /set-spec-dec":                           pipe.Join(mdware.Exec(exec(h, setSpecDEC))),
		"POST /set-spec-dec-sale":                      pipe.Join(mdware.Exec(exec(h, setSpecDECSale))),
		"POST /del-spec-dec":                           pipe.Join(mdware.Exec(exec(h, delSpecDEC))),

		"POST /get-sugg-by-text": pipe.Join(mdware.Exec(exec(h, listSugg))),
		"POST /get-list-by-sugg": pipe.Join(mdware.Exec(exec(h, findSugg))),

		//"POST /run-hotfix": pipe.Join(mdware.Exec(exec(h, runHotfix))),

		// => Debug mode only, when pref.Debug == true
		"GET /debug/vars":               pipe.Join(mdware.Exec(mdware.Stdh)), // expvar
		"GET /debug/pprof/":             pipe.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/cmdline":      pipe.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/profile":      pipe.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/symbol":       pipe.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/trace":        pipe.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/goroutine":    pipe.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
		"GET /debug/pprof/threadcreate": pipe.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
		"GET /debug/pprof/heap":         pipe.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
		"GET /debug/pprof/block":        pipe.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
	}

	h.err404 = mdware.Join(
		mdware.Head(uuid),
		mdware.Errc(http.StatusNotFound),
		mdware.Fail,
		mdware.Tail(h.log),
	)

	h.err405 = mdware.Join(
		mdware.Head(uuid),
		mdware.Errc(http.StatusMethodNotAllowed),
		mdware.Fail,
		mdware.Tail(h.log),
	)

	return h
}

func (h *handler) withRouter(r router.Router) (router.Router, error) {
	var s []string
	var err error
	for k, v := range h.api {
		s = strings.Split(k, " ")
		if len(s) != 2 {
			panic("api: invalid pair method-path")
		}
		err = r.Add(s[0], s[1], v)
		if err != nil {
			return nil, err
		}
	}

	err = r.Set404(h.err404)
	if err != nil {
		return nil, err
	}

	err = r.Set405(h.err405)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewWithRouter returns http.Handler based on given router.
func NewWithRouter(r router.Router, options ...func(*handler) error) (http.Handler, error) {
	if r == nil {
		panic("nil router")
	}

	h := &handler{
		log: logger.NewDefault(),
		rdb: &redis.Pool{},
	}

	var err error
	for i := range options {
		err = options[i](h)
		if err != nil {
			return nil, err
		}
	}

	return h.prepareAPI().withRouter(r)
}

// Logger is option for passing logger interface.
func Logger(l logger.Logger) func(*handler) error {
	return func(h *handler) error {
		h.log = l
		return nil
	}
}

// Redis is interface for Redis Pool Connections.
func Redis(r rediser) func(*handler) error {
	return func(h *handler) error {
		h.rdb = r
		return nil
	}
}

func uuid() string {
	return nuid.Next()
}

func auth(r *http.Request) (string, int, error) {
	return "anonymous", http.StatusOK, nil
}

type ctxHelper struct {
	ctx  context.Context
	rdb  rediser
	log  logger.Logger
	r    *http.Request
	w    http.ResponseWriter
	meta []byte
	data []byte
	lang string
	atag string
	hack string
}

func (h *ctxHelper) getConn() redis.Conn {
	return h.rdb.Get()
}

func (h *ctxHelper) delConn(c io.Closer) {
	_ = c.Close
}

func (h *ctxHelper) clone() *ctxHelper {
	return &ctxHelper{
		h.ctx,
		h.rdb,
		h.log,
		h.r,
		h.w,
		h.meta,
		h.data,
		h.lang,
		h.atag,
		h.hack,
	}
}

func exec(h *handler, f func(*ctxHelper) (interface{}, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hlp := &ctxHelper{
			ctx,
			h.rdb,
			h.log,
			r,
			w,
			[]byte(r.Header.Get("Content-Meta")),
			ctxutil.BodyFrom(ctx),
			mineLang(r.Header.Get("Accept-Language")),
			mineATag(r.Header.Get("User-Agent-Tag")),
			"",
		}
		//FIXME temp workaround
		if (hlp.atag != "") && (strToSHA1(hlp.atag) == "fe5fca9e408b3f3c2346ae5dafa6d57e1856ac2a") {
			hlp.atag = ""
		}
		res, err := f(hlp)
		ctx = hlp.ctx // get ctx from func f
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func mineATag(s string) string {
	return strings.TrimSpace(s)
}

func mineLang(s string) string {
	s = strings.ToLower(s)
	if strings.Contains(s, "uk") || strings.Contains(s, "ua") {
		return "ua"
	}
	if strings.Contains(s, "ru") {
		return "ru"
	}
	if strings.Contains(s, "en") {
		return "en"
	}
	return ""
}

func btsToSHA1(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}

func strToSHA1(s string) string {
	return btsToSHA1([]byte(s))
}
