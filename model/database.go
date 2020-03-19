// mymusicalbum 데이터베이스, 테이블 생성 파일
package orm

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// 테이블 정의
type User struct {
	gorm.Model
	UserId string
	Pw     string
	Name   string
	Image  string
}

func MakeDatabase(c echo.Context) error {
	// 데이터베이스 생성
	db, _ := gorm.Open("mysql",
		"root:qwerz123@tcp(127.0.0.1:3306)/")
	defer db.Close()

	database := "mymusicalbum" // 생성할 데이터베이스 이름

	// Raw SQL로 데이터베이스 생성
	db.Exec("CREATE DATABASE IF NOT EXISTS " + database)
	db.Exec("commit;")

	return c.String(http.StatusOK, "makedatabase successed!")
}

// create table
func MakeTable(c echo.Context) error {
	db, err := gorm.Open("mysql",
		"root:qwerz123@tcp(127.0.0.1:3306)/mymusicalbum?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&User{})

	db.Create(&User{UserId: "jjmin321", Pw: "qwerz123", Name: "jejeongmin"})

	return c.String(http.StatusOK, "maketable successed!")
}
