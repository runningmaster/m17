package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"internal/ctxutil"
)

var (
	mapPX = map[string][]string{
		"ru": []string{prefixINN, prefixMaker, prefixClassATC, prefixSpecINF, prefixSpecACT},
		"ua": []string{prefixINN, prefixMaker, prefixClassATC, prefixSpecDEC},
	}

	mapKB = map[string][]rune{
		"en": []rune("qwertyuiop[]\\asdfghjkl;'zxcvbnm,./`QWERTYUIOP{}|ASDFGHJKL:\"ZXCVBNM<>?~!@#$%^&*()_+"),
		"ru": []rune("йцукенгшщзхъ\\фывапролджэячсмитьбю.ёЙЦУКЕНГШЩЗХЪ/ФЫВАПРОЛДЖЭЯЧСМИТЬБЮ,Ё!\"№;%:?*()_+"),
		"uk": []rune("йцукенгшщзхї\\фівапролджєячсмитьбю.'ЙЦУКЕНГШЩЗХЇ/ФІВАПРОЛДЖЄЯЧСМИТЬБЮ,₴!\"№;%:?*()_+"),
	}
)

func convKB(s, from, to string) string {
	lang1 := mapKB[from]
	lang2 := mapKB[to]
	if lang1 == nil || lang2 == nil {
		return s
	}

	src := []rune(s)
	res := make([]rune, len(src))
	for i := range src {
		for j := range lang1 {
			if lang1[j] == src[i] {
				res[i] = lang2[j]
				break
			}
			res[i] = src[i]
		}
	}
	return string(res)
}

type spec struct {
	ID   int64   `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Sale float64 `json:"sale,omitempty"`
	Text bool    `json:"text,omitempty"`
}

type item struct {
	ID   int64   `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Spec []*spec `json:"spec,omitempty"`
}

type result struct {
	SPC []*item `json:"spc,omitempty"`
	INN []*item `json:"inn,omitempty"`
	ATC []*item `json:"atc,omitempty"`
	ACT []*item `json:"act,omitempty"`
	ORG []*item `json:"org,omitempty"`
}

/*
	convName := convKB(v.Name, "en", "ru")
	if langUA(r.Header) {
		convName = convKB(v.Name, "en", "uk")
	}
*/
func listSugg(h *dbxHelper) (interface{}, error) {
	s, err := jsonToA(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	spx := mapPX[h.lang]
	msg := make(chan string)
	end := make(chan bool)
	var wg sync.WaitGroup
	for i := range spx {
		p := spx[i]
		wg.Add(1)
		go func() {
			c := h.getConn()
			defer h.delConn(c)

			r, err := findIn(c, p, h.lang, "*"+s+"*")
			if err != nil {
				fmt.Println(err)
				return //nil, err
			}
			//fmt.Println(p, len(r))
			for i := range r {
				msg <- r[i].Name
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		end <- true
	}()

L:
	for {
		select {
		case m := <-msg:
			fmt.Println(m)
		case <-end:
			close(msg)
			close(end)
			break L
		}
	}

	return s, nil
}

func findSugg(h *dbxHelper) (interface{}, error) {
	s, err := jsonToA(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	c := h.getConn()
	defer h.delConn(c)

	return s, nil
}

func heatSearch(h *dbxHelper) (interface{}, error) {
	start := time.Now()

	c := h.getConn()
	defer h.delConn(c)

	var ids []int64
	var err error

	ids, err = loadSyncIDs(c, prefixClassATC, 0)
	if err != nil {
		return nil, err
	}
	atc := makeClasses(ids...)
	err = loadHashers(c, prefixClassATC, true, atc)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixClassATC, atc)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixINN, 0)
	if err != nil {
		return nil, err
	}
	inn := makeINNs(ids...)
	err = loadHashers(c, prefixINN, true, inn)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixINN, inn)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixMaker, 0)
	if err != nil {
		return nil, err
	}
	org := makeMakers(ids...)
	err = loadHashers(c, prefixMaker, true, org)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixMaker, org)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixSpecACT, 0)
	if err != nil {
		return nil, err
	}
	act := makeSpecs(ids...)
	err = loadHashers(c, prefixSpecACT, true, act)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecACT, act)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixSpecINF, 0)
	if err != nil {
		return nil, err
	}
	inf := makeSpecs(ids...)
	err = loadHashers(c, prefixSpecINF, true, inf)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecINF, inf)
	if err != nil {
		return nil, err
	}

	ids, err = loadSyncIDs(c, prefixSpecDEC, 0)
	if err != nil {
		return nil, err
	}
	dec := makeSpecs(ids...)
	err = loadHashers(c, prefixSpecDEC, true, dec)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecDEC, dec)
	if err != nil {
		return nil, err
	}

	return time.Since(start).String(), nil
}
