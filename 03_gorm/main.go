package main

import (
	"base_go/03_gorm/lesson01"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Parent struct {
	ID   int `gorm:"primary_key"`
	Name string
}

type Child struct {
	Parent
	Age int
}

func InitDB(dst ...interface{}) *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:Conry@1238@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(dst...)
	return db
}

func main() {
	db, err := gorm.Open(mysql.Open("root:Conry@1238@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	lesson01.Run(db)

	//InitDB(&Parent{}, &Child{})
}
