package routes

import (
	"delonakc.com/api/config"
	"delonakc.com/api/util"
	"fmt"
	"net/http"
)

func GetUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := util.IsLogin(r)
		if user == nil {
			fmt.Fprintf(w, util.ResponseErrorWithCode("未登录", config.NotLogin))
			return
		}
		fmt.Fprintf(w, util.ResponseSuccess(user))
	}
}