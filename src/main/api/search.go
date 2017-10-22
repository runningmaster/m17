package api

import (
	"encoding/json"
)

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
	INF []*item `json:"inf,omitempty"`
	INN []*item `json:"inn,omitempty"`
	ATC []*item `json:"atc,omitempty"`
	ACT []*item `json:"act,omitempty"`
	ORG []*item `json:"org,omitempty"`
}

func jsonToText(data []byte) (int64, error) {
	var v int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func listSugg(h *dbxHelper) (interface{}, error) {
	return nil, nil
}

func findSugg(h *dbxHelper) (interface{}, error) {
	return nil, nil
}

/*
func compileFunctions() {
	if len(compilequeue) != 0 {
		sizeCalculationDisabled = true // not safe to calculate sizes concurrently
		if raceEnabled {
			// Randomize compilation order to try to shake out races.
			tmp := make([]*Node, len(compilequeue))
			perm := rand.Perm(len(compilequeue))
			for i, v := range perm {
				tmp[v] = compilequeue[i]
			}
			copy(compilequeue, tmp)
		} else {
			// Compile the longest functions first,
			// since they're most likely to be the slowest.
			// This helps avoid stragglers.
			obj.SortSlice(compilequeue, func(i, j int) bool {
				return compilequeue[i].Nbody.Len() > compilequeue[j].Nbody.Len()
			})
		}
		var wg sync.WaitGroup
		c := make(chan *Node, nBackendWorkers)
		for i := 0; i < nBackendWorkers; i++ {
			wg.Add(1)
			go func(worker int) {
				for fn := range c {
					compileSSA(fn, worker)
				}
				wg.Done()
			}(i)
		}
		for _, fn := range compilequeue {
			c <- fn
		}
		close(c)
		compilequeue = nil
		wg.Wait()
		sizeCalculationDisabled = false
	}
}
*/
