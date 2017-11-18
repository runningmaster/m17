package api

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
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
		"ua": []rune("йцукенгшщзхї\\фівапролджєячсмитьбю.'ЙЦУКЕНГШЩЗХЇ/ФІВАПРОЛДЖЄЯЧСМИТЬБЮ,₴!\"№;%:?*()_+"),
	}
)

func convLayout(s, from, to string) string {
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

func listSugg(h *dbxHelper) (interface{}, error) {
	s, err := stringFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	// FIXME: check len(rune(s))

	res := make([]string, 0, 100)
	spx := mapPX[h.lang]
	if len(spx) == 0 {
		return res, nil
	}

	errc := make(chan error)
	sugc := make(chan string)
	var wg sync.WaitGroup
	for i := range spx {
		p := spx[i]
		wg.Add(1)
		go func(s string) {
			defer wg.Done()

			c := h.getConn()
			defer h.delConn(c)

			r, err := findIn(c, p, h.lang, s, true)
			if err != nil {
				errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
				return
			}

			// workaround for en layout
			if len(r) == 0 {
				s = convLayout(s, "en", h.lang)
				r, err = findIn(c, p, h.lang, s, true)
				if err != nil {
					errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
					return
				}
			}

			for i := range r {
				sugc <- r[i].Name
			}
		}(s)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

loop:
	for {
		select {
		case e := <-errc:
			err = fmt.Errorf("%v", e)
			if err != nil && strings.Compare(err.Error(), e.Error()) != 0 {
				err = fmt.Errorf("%v: %v", err, e)
			}
		case s := <-sugc:
			res = append(res, s)
		case <-done:
			close(done)
			close(errc)
			close(sugc)
			break loop
		}
	}

	if err != nil {
		return nil, err
	}

	//var n int // FIXME: remove n
	for i := range res {
		//n = strings.Index(res[i], " (")
		//if n > 0 {
		//	res[i] = res[i][:n]
		//}
		//n = strings.Index(res[i], "|")
		//if n > 0 {
		//	res[i] = res[i][:n]
		//}
		res[i] = strings.ToUpper(res[i])
	}
	res = uniqString(res)

	s = strings.ToUpper(s)
	coll := newCollator(h.lang)
	sort.Slice(res,
		func(i, j int) bool {
			if strings.HasPrefix(res[i], s) && !strings.HasPrefix(res[j], s) {
				return true
			} else if !strings.HasPrefix(res[i], s) && strings.HasPrefix(res[j], s) {
				return false
			}

			return coll.CompareString(res[i], res[j]) < 0
		},
	)
	return res, nil
}

//type spec struct {
//	ID   int64   `json:"id,omitempty"`
//	Name string  `json:"name,omitempty"`
//	Sale float64 `json:"sale,omitempty"`
//	Text bool    `json:"text,omitempty"`
//}

type item struct {
	ID        int64   `json:"id,omitempty"`
	IDSpecDEC []int64 `json:"id_spec_dec,omitempty"`
	IDSpecINF []int64 `json:"id_spec_inf,omitempty"`
	Name      string  `json:"name,omitempty"`
}

type result struct {
	Kind string  `json:"kind,omitempty"`
	Item []*item `json:"item,omitempty"`
}

func findSugg(h *dbxHelper) (interface{}, error) {
	s, err := stringFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	var res []*result
	spx := mapPX[h.lang]
	if len(spx) == 0 {
		return res, nil
	}

	errc := make(chan error)
	sugc := make(chan *item)
	var wg sync.WaitGroup
	for i := range spx {
		p := spx[i]
		wg.Add(1)
		go func(s string) {
			defer wg.Done()

			c := h.getConn()
			defer h.delConn(c)

			r, err := findIn(c, p, h.lang, s, false)
			if err != nil {
				errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
				return
			}

			for i := range r {
				sugc <- &item{r[i].ID, nil, nil, p}
			}
		}(s)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	pmap := make(map[string][]int64, 5)
loop:
	for {
		select {
		case e := <-errc:
			err = fmt.Errorf("%v", e)
			if err != nil && strings.Compare(err.Error(), e.Error()) != 0 {
				err = fmt.Errorf("%v: %v", err, e)
			}
		case s := <-sugc:
			pmap[s.Name] = append(pmap[s.Name], s.ID)
		case <-done:
			close(done)
			close(errc)
			close(sugc)
			break loop
		}
	}

	if err != nil {
		return nil, err
	}

	res = makeResult(pmap)
	//res = mineResult(res)
	// 			res = append(res, s)
	return res, nil
}

func makeResult(m map[string][]int64) []*result {
	res := make([]*result, 0, 5)
	var r *result
	for k := range m {
		r = &result{
			Kind: k,
		}
		switch k {
		case prefixSpecINF:
			r.Item = []*item{
				&item{
					ID:        0,
					IDSpecINF: m[k],
					IDSpecDEC: nil,
					Name:      "",
				},
			}
		case prefixSpecDEC:
			r.Item = []*item{
				&item{
					ID:        0,
					IDSpecINF: nil,
					IDSpecDEC: m[k],
					Name:      "",
				},
			}
		default:
			for _, v := range m[k] {
				r.Item = append(r.Item, &item{v, nil, nil, ""})
			}
		}

		res = append(res, r)
	}

	return res
}

func heatSearch(h *dbxHelper) (interface{}, error) {
	start := time.Now()

	c := h.getConn()
	defer h.delConn(c)

	var err error

	atc, err := makeClassesFromIDs(loadSyncIDs(c, prefixClassATC, 0))
	if err != nil {
		return nil, err
	}
	err = loadHashers(c, prefixClassATC, true, atc)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixClassATC, atc)
	if err != nil {
		return nil, err
	}

	inn, err := makeINNsFromIDs(loadSyncIDs(c, prefixINN, 0))
	if err != nil {
		return nil, err
	}
	err = loadHashers(c, prefixINN, true, inn)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixINN, inn)
	if err != nil {
		return nil, err
	}

	org, err := makeMakersFromIDs(loadSyncIDs(c, prefixMaker, 0))
	if err != nil {
		return nil, err
	}
	err = loadHashers(c, prefixMaker, true, org)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixMaker, org)
	if err != nil {
		return nil, err
	}

	act, err := makeSpecsFromIDs(loadSyncIDs(c, prefixSpecACT, 0))
	if err != nil {
		return nil, err
	}
	err = loadHashers(c, prefixSpecACT, true, act)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecACT, act)
	if err != nil {
		return nil, err
	}

	inf, err := makeSpecsFromIDs(loadSyncIDs(c, prefixSpecINF, 0))
	if err != nil {
		return nil, err
	}
	err = loadHashers(c, prefixSpecINF, true, inf)
	if err != nil {
		return nil, err
	}
	err = saveSearchers(c, prefixSpecINF, inf)
	if err != nil {
		return nil, err
	}

	dec, err := makeSpecsFromIDs(loadSyncIDs(c, prefixSpecDEC, 0))
	if err != nil {
		return nil, err
	}
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
