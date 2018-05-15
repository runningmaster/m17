package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func getCheck(h *ctxHelper) (interface{}, error) {
	c := h.getConn()
	defer h.delConn(c)

	//
	p := prefixSpecDEC
	//v, err := makeSpecsFromIDs(loadSyncIDs(c, p, 0))
	v, err := makeSpecsFromIDs(scanIDs(c, p))
	if err != nil {
		return nil, err
	}
	err = loadSpecLinks(c, p, v)
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{}, len(v)*3)
	for i := range v {
		if len(v[i].IDSpecINF) == 0 {
			fmt.Println(v[i].ID)
		}
		for j := range v[i].IDSpecINF {
			_, ok := m[fmt.Sprintf("%d+%d", v[i].ID, v[i].IDSpecINF[j])]
			if ok {
				fmt.Println(fmt.Sprintf("%d+%d", v[i].ID, v[i].IDSpecINF[j]))
			}
			m[fmt.Sprintf("%d+%d", v[i].ID, v[i].IDSpecINF[j])] = struct{}{}
		}
	}

	fmt.Println(len(v), len(m))
	//
	p = prefixSpecINF
	//v, err = makeSpecsFromIDs(loadSyncIDs(c, p, 0))
	v, err = makeSpecsFromIDs(scanIDs(c, p))
	if err != nil {
		return nil, err
	}
	err = loadSpecLinks(c, p, v)
	if err != nil {
		return nil, err
	}
	for i := range v {
		for j := range v[i].IDSpecDEC {
			delete(m, fmt.Sprintf("%d+%d", v[i].IDSpecDEC[j], v[i].ID))
		}
	}

	fmt.Println(len(v), len(m))
	for k := range m {
		fmt.Println(k)
	}
	return statusOK, nil
}

func scanIDs(c redis.Conn, p string) ([]int64, error) {
	res := make([]int64, 0, 25000)
	var next int
	var vals []interface{}
	for done := false; !done; {
		v, err := redis.Values(c.Do("SCAN", next, "MATCH", p+"*", "COUNT", 100))
		if err != nil {
			return nil, err
		}

		next, _ = redis.Int(v[0], err)
		vals, _ = redis.Values(v[1], err)

		var s string
		var r int64
		for i := range vals {
			s, _ = redis.String(vals[i], err)
			if strings.Count(s, ":") != 2 {
				continue
			}

			r, err = strconv.ParseInt(strings.Split(s, ":")[2], 10, 64)
			if err != nil {
				continue
			}
			res = append(res, r)
		}
		done = next == 0
	}
	return res, nil
}
