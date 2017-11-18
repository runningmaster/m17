package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func genKey(v ...interface{}) string {
	a := make([]string, 0, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case []byte:
			a = append(a, string(x))
		case int:
			a = append(a, strconv.Itoa(x))
		case int64:
			a = append(a, strconv.Itoa(int(x)))
		case string:
			a = append(a, x)
		default:
			a = append(a, fmt.Sprintf("%v", x))
		}
	}
	return strings.Join(a, ":")
}

func normName(s string) string {
	r := strings.NewReplacer(
		"®", "",
		"™", "",
		"*", "",
		"&", "",
		"†", "",
	)
	return strings.TrimSpace(strings.ToLower(r.Replace(s)))
}

func stringFromJSON(data []byte) (string, error) {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return "", err
	}
	return v, nil
}

func int64FromJSON(data []byte) (int64, error) {
	var v int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func int64sFromJSON(data []byte) ([]int64, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func int64sToJSON(v []int64) []byte {
	r, _ := json.Marshal(v)
	return r
}

func uniqInt64(s []int64) []int64 {
	seen := make(map[int64]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func uniqString(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func notZeroValue(v interface{}) bool {
	return !(v == nil || reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()))
}
