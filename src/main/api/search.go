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
	ID    int64   `json:"id,omitempty"`
	Code  string  `json:"code,omitempty"`
	Name  string  `json:"name,omitempty"`
	Full  bool    `json:"full,omitempty"`
	Slug  string  `json:"slug,omitempty"`
	Sale  float64 `json:"sale,omitempty"`
	Maker string  `json:"maker,omitempty"`
	UATag bool    `json:"uatag,omitempty"`
	List  []*item `json:"list,omitempty"`
}

type result struct {
	Kind string  `json:"kind,omitempty"`
	List []*item `json:"list,omitempty"`
	Sort int     `json:"sort,omitempty"`
}

func findSugg(h *ctxHelper) (interface{}, error) {
	s, err := stringFromJSON(h.data)
	if err != nil {
		h.ctx = ctxutil.WithCode(h.ctx, http.StatusBadRequest)
		return nil, err
	}

	spx := mapPX[h.lang]
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
				sugc <- &item{r[i].ID, p, "", false, "", 0, "", false, nil}
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
			pmap[s.Code] = append(pmap[s.Code], s.ID)
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

	return makeResult(h, pmap)
}

func makeResult(h *ctxHelper, m map[string][]int64) ([]*result, error) {
	errc := make(chan error)
	resc := make(chan *result)
	var wg sync.WaitGroup
	for k := range m {
		wg.Add(1)
		go func(p string, x []int64) {
			defer wg.Done()
			r := &result{
				Kind: p,
			}
			c := h.clone()
			c.data = int64sToJSON(uniqInt64(x))
			switch p {
			case prefixClassATC:
				v, err := getClassXNext(c, p)
				if err != nil {
					errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
					return
				}
				for i := range v {
					if v[i] == nil {
						continue
					}
					l, err := mineItemListByID(c, p, v[i].ID)
					if err != nil {
						errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
						return
					}
					r.List = append(r.List, &item{v[i].ID, v[i].Code, v[i].Name, false, v[i].Slug, 0, "", false, l})
					r.Sort = 1 // max's magic :)
				}
			case prefixINN:
				v, err := getINNXList(c, p)
				if err != nil {
					errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
					return
				}
				for i := range v {
					if v[i] == nil {
						continue
					}
					l, err := mineItemListByID(c, p, v[i].ID)
					if err != nil {
						errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
						return
					}
					r.List = append(r.List, &item{v[i].ID, "", v[i].Name, false, v[i].Slug, 0, "", false, l})
					r.Sort = 2
				}
			case prefixMaker:
				v, err := getMakerXList(c, p)
				if err != nil {
					errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
					return
				}
				for i := range v {
					if v[i] == nil {
						continue
					}
					l, err := mineItemListByID(c, p, v[i].ID)
					if err != nil {
						errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
						return
					}
					r.List = append(r.List, &item{v[i].ID, "", v[i].Name, false, v[i].Slug, 0, "", false, l})
					r.Sort = 4
				}
			default: // prefixSpecINF, prefixSpecDEC, prefixSpecACT
				v, err := getSpecXList(c, p)
				if err != nil {
					errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
					return
				}
				if c.atag != "" {
					v, err = crazyPermutation(c, p, v)
					if err != nil {
						errc <- fmt.Errorf("%s %s: %v", p, h.lang, err)
						return
					}
				}
				for i := range v {
					if v[i] == nil {
						continue
					}
					r.List = append(r.List, &item{v[i].ID, "", v[i].Name, v[i].Full, v[i].Slug, v[i].Sale, v[i].Maker, v[i].UATag, nil})
					if p == prefixSpecACT {
						r.Sort = 3
					}
				}
			}
			resc <- r
		}(k, m[k])
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	var err error
	res := make([]*result, 0, 5)
loop:
	for {
		select {
		case e := <-errc:
			err = fmt.Errorf("%v", e)
			if err != nil && strings.Compare(err.Error(), e.Error()) != 0 {
				err = fmt.Errorf("%v: %v", err, e)
			}
		case r := <-resc:
			res = append(res, r)
		case <-done:
			close(done)
			close(errc)
			close(resc)
			break loop
		}
	}

	if err != nil {
		return nil, err
	}

	sort.Slice(res,
		func(i, j int) bool {
			return res[i].Sort < res[j].Sort
		},
	)

	for i := range res {
		res[i].Sort = 0
	}

	// f2k
	r := &result{
		Kind: "x",
	}
	for i := range res {
		if res[i].Kind == prefixSpecINF || res[i].Kind == prefixSpecDEC {
			r.List = append(r.List, res[i].List[0])
			break
		}
	}

	if len(r.List) == 0 {
		for i := range res {
			if res[i].Kind == prefixINN {
				for j := range res[i].List {
					for k := range res[i].List[j].List {
						if res[i].List[j].List[k].Full {
							r.List = append(r.List, res[i].List[j].List[k])
							break
						}
					}
					if len(r.List) > 0 {
						break
					}
				}
				if len(r.List) > 0 {
					break
				}
			}
		}
	}

	if len(r.List) == 0 {
		for i := range res {
			if res[i].Kind == prefixClassATC {
				for j := range res[i].List {
					for k := range res[i].List[j].List {
						if res[i].List[j].List[k].Full {
							r.List = append(r.List, res[i].List[j].List[k])
							break
						}
					}
					if len(r.List) > 0 {
						break
					}
				}
				if len(r.List) > 0 {
					break
				} else {
					r.List = append(r.List, res[i].List[0].List[0])
					break
				}
			}
		}
	}

	if len(r.List) == 0 {
		for i := range res {
			if res[i].Kind == prefixMaker {
				for j := range res[i].List {
					for k := range res[i].List[j].List {
						if res[i].List[j].List[k].Full {
							r.List = append(r.List, res[i].List[j].List[k])
							break
						}
					}
					if len(r.List) > 0 {
						break
					}
				}
				if len(r.List) > 0 {
					break
				}
			}
		}
	}

	if len(r.List) > 0 {
		res = append(res, r)
	}

	return res, nil
}

func mineItemListByID(h *ctxHelper, p string, x int64) ([]*item, error) {
	h.data = int64ToJSON(x)
	s := prefixSpecINF
	if h.lang == "ua" {
		s = prefixSpecDEC
	}
	if p == prefixMaker {
		h.atag = ""
	}

	v, err := getSpecXListByWithCrazyPermutation(h, s, p, p == prefixClassATC)
	if err != nil {
		return nil, err
	}

	res := make([]*item, 0, len(v))
	for i := range v {
		if v[i] == nil {
			continue
		}
		res = append(res, &item{v[i].ID, "", v[i].Name, v[i].Full, v[i].Slug, v[i].Sale, v[i].Maker, v[i].UATag, nil})
	}

	return res, nil
}

/*


	if len(r.List) == 0 {
		for i := range res {
			if res[i].Kind != prefixINN {
				continue
			}
			for j := range res[i].List {
				h.data = int64sToJSON([]int64{res[i].List[j].ID})
				v, err := getSpecXListByWithCrazyPermutation(h, prefixSpecINF, prefixINN)
				if err != nil {
					return nil, err
				}
				for k := range v {
					if v[k] == nil {
						continue
					}
					if v[k].Full {
						r.List = append(r.List, &item{v[k].ID, "", v[k].Name, v[k].Full, v[k].Slug, v[k].Sale, v[k].Maker, v[k].UATag, nil})
						break
					}
				}
				if len(res[i].List) > 0 {
					break
				}
			}
		}
	}

	if len(r.List) == 0 {
		for i := range res {
			if res[i].Kind == prefixClassATC {
				for j := range res[i].List {
					h.data = int64sToJSON([]int64{res[i].List[j].ID})
					v, err := getSpecXListByWithCrazyPermutation(h, prefixSpecINF, prefixClassATC)
					if err != nil {
						return nil, err
					}
				}
				if len(res[i].List) > 0 {
					break
				}
			}
		}
	}


*/
