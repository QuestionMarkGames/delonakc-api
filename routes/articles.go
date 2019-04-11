package routes

import (
	"delonakc.com/api/database"
	"delonakc.com/api/util"
	"fmt"
	"net/http"
	"strconv"
)

type ArticlesHandler struct {
	DB *database.MongoDB
	Data string
}

func (a *ArticlesHandler) Get(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	page, err  := strconv.Atoi(r.Form.Get("page"))
	pageNum, err := strconv.Atoi(r.Form.Get("pageNum"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, util.ResponseError("非法的请求"))
		return
	}

	articles := a.DB.GetArticles(page, pageNum)

	fmt.Fprintf(w, util.ResponseSuccess(articles))

}

func (a *ArticlesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.Get(w, r)
		return
	}
}
