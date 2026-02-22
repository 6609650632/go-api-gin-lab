package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"example.com/student-api/handlers"
	"example.com/student-api/repositories"
	"example.com/student-api/services"
)

func main() {
	db, err := sql.Open("sqlite3", "./students.db")
	if err != nil {
		log.Fatal(err)
	}

	repo := &repositories.StudentRepository{DB: db}
	service := &services.StudentService{Repo: repo}
	handler := &handlers.StudentHandler{Service: service}

	r := gin.Default()

	r.GET("/students", handler.GetStudents)
	r.GET("/students/:id", handler.GetStudentByID)
	r.POST("/students", handler.CreateStudent)
	r.PUT("/students/:id", handler.UpdateStudent)
	r.DELETE("/students/:id", handler.DeleteStudent)

	r.Run(":8080")
}
