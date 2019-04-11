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
	"net/http"
	"time"
)

func AddComments(db *database.MongoDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := util.IsLogin(r)

		if user == nil {
			fmt.Fprintf(w, util.ResponseErrorWithCode("你还未登录呢...", config.NotLogin))
			return
		}

		vars := util.ParseVars(r)

		id, _ := vars["id"]

		if !bson.IsObjectIdHex(id) {
			fmt.Fprintf(w, util.ResponseError("非法的Object Id"))
			return
		}

		comment, err := transferToComment(r, user)

		if err != nil {
			fmt.Fprintf(w, util.ResponseError(err.Error()))
			return
		}

		err = db.AddComment(bson.ObjectIdHex(id), comment)

		if err != nil {
			fmt.Fprintf(w, util.ResponseError(err.Error()))
			return
		}

		fmt.Fprintf(w, util.ResponseSuccess("添加成功"))
	}
}

func transferToComment(r *http.Request, user *config.GithubUser) (*model.Comment, error) {

	body, _ := ioutil.ReadAll(r.Body)

	var comment *model.Comment

	err := json.Unmarshal(body, &comment)

	if err != nil {
		return nil, err
	}
	if len(comment.Content) == 0 {
		return nil, fmt.Errorf("评论内容不能为空")
	}

	comment.Author = user.Name
	comment.Avator = user.AvatarUrl
	comment.Url = user.Home
	comment.Time = time.Now().Unix() * 1000
	return comment, nil
}