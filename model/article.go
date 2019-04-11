package model

import (
	"gopkg.in/mgo.v2/bson"
)

type ArticleDatabase interface {
	//获取文章列表
	GetArticles(page, pageNum int) []*Article

	//通过文章id获取文章详情
	//GetArticlesById(id string) *Article

	//添加新的文章
	//AddNewArticle(article *Article) error

	//更新文章
	//UpdateArticle(id string, article *Article) error
}

//评论
type Comment struct {
	Author string `json:"author"`
	Avator string `json:"avator"`
	Content string `json:"content"`
	Time int64 `json:"time"`
	Url string `json:"url"`
}

//文章
type Article struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	SaveTime int64 `json:"saveTime"`
	UpdateTime int64 `json:"updateTime"`
	Type string `json:"type"`
	Visit int `json:"visit"`
  Comments []*Comment
}