package util

import (
	"delonakc.com/api/config"
	"delonakc.com/api/redis"
	"encoding/json"
	"net/http"
)

func IsLogin(r *http.Request) *config.GithubUser {
	bid, err := r.Cookie("bid")

	if err != nil {
		return nil
	}

	v := redis.Get(bid.Value)

	if v == nil {
		return nil
	}

	var user *config.GithubUser

	json.Unmarshal(v, &user)

	return user
}

func ParseVars(r *http.Request) map[string]string {
	vars := r.Context().Value(config.RouteVarsKey)

	if vars == nil {
		return nil
	}

	v, ok := vars.(map[string]string)

	if !ok {
		return nil
	}

	return v
}