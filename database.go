package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "foo:bar@(localhost)/ptt?loc=Local&charset=utf8mb4&parseTime=True")
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db.LogMode(false)
	return db, nil
}

func dbCheckExistArticle(url string, db *gorm.DB) (bool, error) {
	if result := db.Where("url = ?", url).First(&articles{}); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return false, nil
		} else {
			return false, result.Error
		}
	}
	return true, nil
}

func dbCreateArticle(article *articles, db *gorm.DB) error {
	article.Updated_at = time.Now().Local()
	if err := db.Create(&article).Error; err != nil {
		return err
	}
	log.Println("create	" + article.Url)
	return nil
}

func dbUpdateArticle(article *articles, db *gorm.DB) error {
	article.Updated_at = time.Now().Local()
	if err := db.Model(&article).Where("url = ?", article.Url).UpdateColumns(articles{Contains: article.Contains, Reply: article.Reply, Updated_at: article.Updated_at}).Error; err != nil {
		return err
	}
	log.Println("update	" + article.Url)
	return nil
}
