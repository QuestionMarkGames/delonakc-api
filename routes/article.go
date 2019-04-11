package routes

import (
	"delonakc.com/api/config"
	"delonakc.com/api/database"
	"delonakc.com/api/model"
	"delonakc.com/api/util"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
)

type ArticleHandler struct {
	DB *database.MongoDB
	Vars map[string]string
	Conf *config.Conf
}

func (a *ArticleHandler) Get(w http.ResponseWriter, r *http.Request) {
	val, _ := a.Vars["id"]

	if !bson.IsObjectIdHex(val) {
		fmt.Fprintf(w, util.ResponseError("非法的object id"))
		return
	}

	articles := a.DB.GetArticleById(bson.ObjectIdHex(val))

	if articles == nil {
		fmt.Fprintf(w, util.ResponseError("不存在该id"))
		return
	}
	fmt.Fprintf(w, util.ResponseSuccess(articles))
}

func (a *ArticleHandler) Post(w http.ResponseWriter, r *http.Request) {

	user := util.IsLogin(r)

	if user == nil {
		fmt.Fprintf(w, util.ResponseErrorWithCode("你还未登录", config.NotLogin))
		return
	}

	ar, err := TransferFormToArticle(r)

	if err != nil {
		fmt.Fprintf(w, util.ResponseError(err.Error()))
		return
	}

	err = a.DB.AddArticle(ar)

	if err != nil {
		log.Printf("添加文章错误 %v", err)
		fmt.Fprintf(w, util.ResponseError("添加失败"))
		return
	}

	fmt.Fprintf(w, util.ResponseSuccess("添加成功"))
}

func (a *ArticleHandler) Put(w http.ResponseWriter, r *http.Request) {

	user := util.IsLogin(r)

	if user == nil {
		fmt.Fprintf(w, util.ResponseErrorWithCode("你还未登录", config.NotLogin))
		return
	}

	val, _ := a.Vars["id"]

	if !bson.IsObjectIdHex(val) {
		fmt.Fprintf(w, util.ResponseError("非法的object id"))
		return
	}

	ar, err := TransferFormToArticle(r)

	if err != nil {
		log.Println(err)
		util.ResponseError("参数解析错误")
		return
	}

	err = valiedFormKey(ar)

	if err != nil {
		fmt.Fprintf(w, util.ResponseError(err.Error()))
		return
	}

	err = a.DB.UpdateArticle(bson.ObjectIdHex(val), ar)

	if err != nil {
		log.Printf("更新文章错误: id: %s \n,%v", val, err)
		fmt.Fprintf(w, util.ResponseError("更新文章错误"))
		return
	}

	fmt.Fprintf(w, util.ResponseSuccess("更新成功"))
}

func (a *ArticleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.Vars = util.ParseVars(r)

	if r.Method == http.MethodGet {
		a.Get(w, r)
		return
	}
	if r.Method == http.MethodPut {
		a.Put(w, r)
		return
	}
	if r.Method == http.MethodPost {
		a.Post(w, r)
		return
	}
}

func TransferFormToArticle(r *http.Request) (*model.Article, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var article *model.Article

	err := json.Unmarshal(body, &article)

	if err != nil {
		return nil, err
	}
	return article, nil
}

func valiedFormKey(a *model.Article) error {
	if a.Title == "" {
		return fmt.Errorf("Title 字段不能为空")
	}
	if a.Content == "" {
		return fmt.Errorf("Content 字段不能为空")
	}
	if a.Type == "" {
		return fmt.Errorf("Type 字段不能为空")
	}
	return nil
}