package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"main/version"

	"internal/ctxutil"
	"internal/router"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, version.WithBuildInfo(), runtime.Version())
}

func help() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		res, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), version.AppName()+".txt"))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
			*r = *r.WithContext(ctx)
			return
		}
		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Fprintf(w, "Hello, World! From test handler!\n")
	fmt.Fprintf(w, "Param foo: %s\n", router.ParamValueFrom(ctx, "foo"))
	fmt.Fprintf(w, "Query foo: %s\n", router.QueryValueFrom(ctx, "foo"))
	v, _ := ctx.Value("foo").(string)
	fmt.Fprintf(w, "Value foo: %s\n", v)

	*r = *r.WithContext(ctx)
}

func ping(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		dbx := &dbxHelper{
			ctx: ctx,
			rdb: rdb,
		}
		res, err := dbx.ping()
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func exec(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dbx := &dbxHelper{
			ctx,
			rdb,
			nil,
			r,
			w,
			[]byte(r.Header.Get("Content-Meta")),
			ctxutil.BodyFrom(ctx),
			mineRUorUA(r.Header.Get("Accept-Language")),
		}

		res, err := dbx.exec(router.ParamValueFrom(ctx, "func"))
		ctx = dbx.ctx // get ctx from func
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func mineRUorUA(s string) string {
	s = strings.ToLower(s)
	if strings.Contains(s, "uk") || strings.Contains(s, "ua") {
		return "ua"
	}
	if strings.Contains(s, "ru") {
		return "ru"
	}
	return ""
}

/*
func (j *jsonHead) uaInsteadRu(fn func(string, string) (string, string)) {
	if fn == nil {
		return
	}
	j.Name, j.NameUA = fn(j.Name, j.NameUA)
	j.NameFull, j.NameFullUA = fn(j.NameFull, j.NameFullUA)
	j.NameShort, j.NameShortUA = fn(j.NameShort, j.NameShortUA)
	j.Addr1Country, j.Addr1CountryUA = fn(j.Addr1Country, j.Addr1CountryUA)
	j.Addr1Area, j.Addr1AreaUA = fn(j.Addr1Area, j.Addr1AreaUA)
	j.Addr1Region, j.Addr1RegionUA = fn(j.Addr1Region, j.Addr1RegionUA)
	j.Addr1City, j.Addr1CityUA = fn(j.Addr1City, j.Addr1CityUA)
	j.Addr1Street, j.Addr1StreetUA = fn(j.Addr1Street, j.Addr1StreetUA)
	j.Addr2Country, j.Addr2CountryUA = fn(j.Addr2Country, j.Addr2CountryUA)
	j.Addr2Area, j.Addr2AreaUA = fn(j.Addr2Area, j.Addr2AreaUA)
	j.Addr2Area, j.Addr2AreaUA = fn(j.Addr2Area, j.Addr2AreaUA)
	j.Addr2City, j.Addr2CityUA = fn(j.Addr2City, j.Addr2CityUA)
	j.Addr2StreetUA, j.Addr2StreetUA = fn(j.Addr2StreetUA, j.Addr2StreetUA)

}
*/
