package database

import (
	"delonakc.com/api/model"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (m *MongoDB) GetArticles(page, pageNum int) []*model.Article {
	m.SetCollection("articles")
	var articles []*model.Article
	m.C.Find(nil).Skip((page - 1)*pageNum).Limit(pageNum).Sort("-updateTime").All(&articles)
	return articles
}

func (m *MongoDB) GetArticleById(id bson.ObjectId) []*model.Article {
	m.SetCollection("articles")
	var article []*model.Article
	m.C.FindId(id).Limit(1).All(&article)
	return article
}

func (m *MongoDB) UpdateArticle(id bson.ObjectId, article *model.Article) error {
	article.UpdateTime = time.Now().Unix() * 1000
	article.Id = id
	m.SetCollection("articles")
	err := m.C.UpdateId(id, article)
	return err
}

func (m *MongoDB) AddArticle(article *model.Article) error {
	article.SaveTime = time.Now().Unix() * 1000
	article.UpdateTime = time.Now().Unix() * 1000
	article.Id = bson.NewObjectId()
	m.SetCollection("articles")
	return m.C.Insert(article)
}

func (m *MongoDB) AddComment(id bson.ObjectId, comment *model.Comment) error {
	m.SetCollection("articles")
	var article *model.Article

	m.C.FindId(id).Limit(1).One(&article)

	if article == nil {
		return fmt.Errorf("没有这个id:%s的文章", id.String())
	}

	article.Comments = append([]*model.Comment{ comment}, article.Comments...)

	err := m.UpdateArticle(id, article)

	if err != nil {
		return err
	}

	return nil
}