package runner

import (
	"delonakc.com/api/config"
	"delonakc.com/api/database"
	"delonakc.com/api/router"
	"delonakc.com/api/routes"
	"fmt"
	"log"
	"net/http"
)

func Run(r *router.Router, db *database.MongoDB, conf *config.Conf) {

	r.HandleFunc("/v1/article", &routes.ArticlesHandler{ DB: db }).Methods("GET")
	r.HandleFunc("/v1/article", &routes.ArticleHandler{ DB: db, Conf: conf }).Methods("POST")
	r.HandleFunc("/v1/article/{id}", &routes.ArticleHandler{ DB: db }).Methods("GET", "PUT")

	r.HandleFunc("/v1/article/{id}/comment", routes.AddComments(db)).Methods("POST")

	r.HandleFunc("/v1/auth/github", routes.RedirectToGithub(conf)).Methods("GET")
	r.HandleFunc("/v1/auth/github/token", routes.GetToken(conf)).Methods("GET")

	r.HandleFunc("/v1/user", routes.GetUserInfo()).Methods("GET")

	http.Handle("/", r)

	addr := fmt.Sprintf(":%d", conf.Server.Port)

	fmt.Println("Listening on Port:", conf.Server.Port)

	log.Fatal(http.ListenAndServe(addr, nil))
}