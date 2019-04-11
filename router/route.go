package router

import (
	"bytes"
	"context"
	"delonakc.com/api/config"
	"delonakc.com/api/database"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	DB *database.MongoDB
	handler http.Handler
	name string
	path string
	// conf
	//matchers []Matcher

	Options *RegexpMatcher
}

type RegexpMatcher struct {
	methods []string
	strict bool
	varsN []string
	varsM []*regexp.Regexp
	template string
	regexp *regexp.Regexp
	isMatchedRoute bool
}
//
//type Matcher interface {
//	Match(r *http.Request, match *router.RouteMatch) bool
//}

func (r *Route) Path(path string) *Route {
	err := r.AddRouteRegexp(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error In Path Method: %s \n%v", path, err))
	}
	return r
}

func (r *Route) Methods(methods... string) *Route {
	for _, method := range methods {
		r.Options.methods = append(r.Options.methods, strings.ToUpper(method))
	}
	return r
}

func (r *Route) AddRouteRegexp(path string) error {
	regexpMatcher, err := addNewRegexpMatcher(path)
	if err != nil {
		return err
	}
	r.Options = regexpMatcher
	return nil
}

func (r *Route) Match(req *http.Request) bool {
	path := req.URL.Path
	if !r.Options.isMatchedRoute {
		return r.Options.template == path && r.MatchMethod(req.Method)
	}
	//fmt.Println(r.Options.template, ":", r.Options.regexp.String(), ":", path)
	//fmt.Println(r.Options.regexp.MatchString(path))
	return r.Options.regexp.MatchString(path) && r.MatchMethod(req.Method)
}

func (r *Route) SetMatch(req *http.Request) *http.Request {
	path := req.URL.Path
	pattern := r.Options.regexp

	if pattern == nil {
		return req
	}

	ids := pattern.FindStringSubmatchIndex(path)

	if ids == nil {
		return req
	}

	var matcher = make(map[string]string)

	extractValues(path, r.Options.varsN, matcher, ids)
	req = setValues(req, config.RouteVarsKey, matcher)

	return req
}

func (r *Route) MatchMethod(method string) bool {
	for _, v := range r.Options.methods {
		if strings.ToUpper(method) == v {
			return true
		}
	}
	return false
}

func (r *Route) Handler() http.Handler {
	return r.handler
}

func (r *Route) HandleFunc(fn interface{}) *Route {

	handler, ok := fn.(http.Handler)

	if ok {
		r.handler = handler
		return r
	}
	HandlerFunc, _ := fn.(http.HandlerFunc)
	r.handler = http.HandlerFunc(HandlerFunc)
	return r
}

func addNewRegexpMatcher(path string) (*RegexpMatcher, error) {
	rr := &RegexpMatcher{ strict: true }

	// 非正则匹配路由
	if strings.Index(path, "{") == -1 {
		rr.isMatchedRoute = false
		rr.regexp = nil
		rr.template = path
		return rr, nil
	}

	// 正则匹配路由
	rr.isMatchedRoute = true
	rr.template = path

	ids, err := braceIndexs(path)
	defaultPattern := "[^/]+"
	if err != nil {
		return nil, errors.New("路径解析错误，请输入正确的匹配模式，如：/api/article/{id}")
	}

	tpl := path
	var start, end int
	var pattern = bytes.NewBufferString("")
	for i := 0; i < len(ids)/2; i += 2 {
		raw := tpl[end : ids[i]]
		start = ids[i]
		end = ids[i + 1]
		name := tpl[start + 1 : end]
		rr.varsN = append(rr.varsN, name)
		ptt, _ := regexp.Compile(fmt.Sprintf("^%s$", name))
		rr.varsM = append(rr.varsM, ptt)
		fmt.Fprintf(pattern, "%s(?P<%s>%s)", regexp.QuoteMeta(raw), varGroupName(i/2), defaultPattern)
	}

	if end + 1 < len(tpl) {
		pattern.WriteString(tpl[end+1:])
	}

	pat, _ := regexp.Compile(fmt.Sprintf("^%s$", pattern.String()))
	rr.regexp = pat

	return rr, nil
}

func braceIndexs(tpl string) ([]int, error) {
	var level int
	var ids []int

	for i := 0; i < len(tpl); i++ {
		switch tpl[i] {
			case '{':
				if level++; level == 1 {
					ids = append(ids, i)
				}
		case '}':
				if level--; level == 0 {
					ids = append(ids, i)
				}
		}
	}

	return ids, nil
}

func varGroupName(index int) string {
	return fmt.Sprintf("v%d", index)
}

func extractValues(str string, names []string, output map[string]string, ids []int) {
	for i, key := range names {
		output[key] = str[ids[i*2+2]:ids[i*2+3]]
	}
}

func setValues(r *http.Request, key string, value map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}