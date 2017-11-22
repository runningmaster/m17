package api

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

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

func listSugg(h *ctxHelper) (interface{}, error) {
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

	for i := range res {
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
	ID   int64  `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

type result struct {
	Kind string  `json:"kind,omitempty"`
	List []*item `json:"list,omitempty"`
}

func findSugg(h *ctxHelper) (interface{}, error) {
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

	res, err = mineResult(h, makeResult(pmap))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func makeResult(h *ctxHelper, m map[string][]int64) []*result {
	res := make([]*result, 0, 5)
	var r *result
	for k := range m {
		r = &result{
			Kind: k,
		}
		switch k {
		case prefixSpecINF:
			r.List = []*item{
				&item{
					ID:        0,
					IDSpecINF: uniqInt64(m[k]),
					IDSpecDEC: nil,
					Name:      "",
				},
			}
		case prefixSpecDEC:
			r.List = []*item{
				&item{
					ID:        0,
					IDSpecINF: nil,
					IDSpecDEC: uniqInt64(m[k]),
					Name:      "",
				},
			}
		default:
			for _, v := range m[k] {
				r.List = append(r.List, &item{v, nil, nil, ""})
			}
		}
		res = append(res, r)
	}

	return res
}

/*
h.data = int64sToJSON(x)
	return getSpecXList(h, p1)
*/
func mineResult(h *ctxHelper, v []*result) ([]*result, error) {
	errc := make(chan error)
	var wg sync.WaitGroup
	for i := range v {
		//// FIXME: load list ?
		//if strings.HasPrefix(v[i].Kind, prefixSpecINF) || strings.HasPrefix(v[i].Kind, prefixSpecDEC) {
		//	continue
		//}
		p := v[i].Kind
		for j := range v[i].List {
			x := v[i].List[j]
			wg.Add(1)
			go func(p string, v *item) {
				defer wg.Done()

				c := h.getConn()
				defer h.delConn(c)

				switch p {
				case prefixClassATC:
				case prefixClassATC:
				}

				//var err error
				//v.Name, err = loadHashFieldAsString(c, p, iifString(h.lang == "ru", "name_ru", "name_ua"), v.ID)
				//if err != nil {
				//	errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
				//	return
				//}
				//if h.lang == "ru" {
				//	v.IDSpecINF, err = loadLinkIDs(c, p, prefixSpecINF, v.ID)
				//	if err != nil {
				//		errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
				//		return
				//	}
				//} else {
				//	v.IDSpecDEC, err = loadLinkIDs(c, p, prefixSpecDEC, v.ID)
				//	if err != nil {
				//		errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
				//		return
				//	}
				//}
			}(p, x)
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	var err error
loop:
	for {
		select {
		case e := <-errc:
			err = fmt.Errorf("%v", e)
			if err != nil && strings.Compare(err.Error(), e.Error()) != 0 {
				err = fmt.Errorf("%v: %v", err, e)
			}
		case <-done:
			close(done)
			close(errc)
			break loop
		}
	}

	if err != nil {
		return nil, err
	}

	return v, nil
}
