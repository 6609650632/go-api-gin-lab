package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/student-api/models"
	"example.com/student-api/services"
)

type StudentHandler struct {
	Service *services.StudentService
}

// --------------------
// helper: validate
// --------------------
func validateStudent(s models.Student, requireID bool) (bool, string) {
	if requireID && s.Id == "" {
		return false, "id must not be empty"
	}
	if s.Name == "" {
		return false, "name must not be empty"
	}
	if s.GPA < 0 || s.GPA > 4 {
		return false, "gpa must be between 0.00 and 4.00"
	}
	return true, ""
}

// --------------------
// GET ALL
// --------------------
func (h *StudentHandler) GetStudents(c *gin.Context) {
	students, err := h.Service.GetStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
		return
	}
	c.JSON(http.StatusOK, students)
}

// --------------------
// GET BY ID
// --------------------
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")

	student, err := h.Service.GetStudentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// --------------------
// CREATE
// --------------------
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if ok, msg := validateStudent(student, true); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	if err := h.Service.CreateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// --------------------
// UPDATE
// --------------------
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")

	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if ok, msg := validateStudent(student, false); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	updated, err := h.Service.UpdateStudent(id, student)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// --------------------
// DELETE
// --------------------
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.DeleteStudent(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
