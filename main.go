package main

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// gorm.Modelの標準モデル
type Todo struct {
	gorm.Model
	// id
	// created_at
	// updated_at
	// deleted_at
	Text   string
	Status string
}

// マイグレーション
func dbInit() {
	// 第一引数：使用するDB　第二引数：ファイル名
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース接続不可(dbInit)")
	}

	// ファイルが無ければ生成、存在すれば処理スキップ
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

// DB追加
func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース接続不可(dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

// 全件取得
func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース接続不可(dbGetAll)")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

// DB1件取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース接続不可(dbGetOne)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

// DB更新
func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース接続不可(dbUpdate")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

// DB削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース接続不可(dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	dbInit()

	// Index
	r.GET("/", func(c *gin.Context) {
		todos := dbGetAll()
		c.HTML(http.StatusOK, "index.html", gin.H{"todos": todos})
	})

	// Create
	r.POST("/new", func(c *gin.Context) {
		text   := c.PostForm("text")
		status := c.PostForm("status")
		dbInsert(text, status)
		c.Redirect(http.StatusFound, "/")
	})

	// Detail
	r.GET("/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		c.HTML(http.StatusOK, "detail.html", gin.H{"todo": todo})
	})

	// Update
	r.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text   := c.PostForm("text")
		status := c.PostForm("status")
		dbUpdate(id, text, status)
		c.Redirect(http.StatusFound, "/")
	})

	// 削除確認
	r.GET("/delete_check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		c.HTML(http.StatusOK, "delete.html", gin.H{"todo": todo})
	})

	// Delete
	r.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		c.Redirect(http.StatusFound, "/")
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}