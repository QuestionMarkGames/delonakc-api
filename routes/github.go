package routes

import (
	"crypto/md5"
	"delonakc.com/api/config"
	"delonakc.com/api/redis"
	"delonakc.com/api/util"
	"fmt"
	"io"
	"net/http"
	"time"
)

func RedirectToGithub(conf *config.Conf) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		redirectUrl := r.Form.Get("redirectUrl")

		http.Redirect(w, r,
			fmt.Sprintf("%s?client_id=%s&login=%s&redirect_uri=%s",
				conf.Github.BaseUrl, conf.Github.ClientId, conf.Github.Login, redirectUrl),
			302,
		)
	}
}

func GetToken(conf *config.Conf) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		code := r.Form.Get("code")

		gitConf := conf.Github

		gitConf.Code = code

		_, err := gitConf.ExchangeToken()

		if err != nil {
			fmt.Fprintf(w, util.ResponseError(err.Error()))
		}

		user, err := gitConf.GetUserInfo()

		h := md5.New()

		io.WriteString(h, fmt.Sprintf("%d", user.Id))

		key := fmt.Sprintf("%x", h.Sum(nil))

		redis.Set(key, user)

		expireTime, _ := time.ParseDuration(conf.Redis.ExpireTime)
		cookie := fmt.Sprintf("bid=%s;Max-Age=%0.f;path=/;HttpOnly;", key, expireTime.Seconds())

		w.Header().Set("Set-Cookie", cookie)

		fmt.Fprintf(w, util.ResponseSuccess(user))

	}
}