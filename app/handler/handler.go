package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type StudentHandler struct {
	srv *service.StudentService
}

func NewStudentHandler(s *service.StudentService) *StudentHandler {
	return &StudentHandler{s}
}
func (h *StudentHandler) AddStudent(c *gin.Context) {
	var newStudent models.StudentInput
	if err := c.BindJSON(&newStudent); err != nil {
		log.WithError(err).Error("Error while Converting the data")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outputStudent, err := h.srv.AddStudent(newStudent)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}
	c.IndentedJSON(http.StatusCreated, outputStudent)

}

func (h *StudentHandler) GetAllStudent(c *gin.Context) {
	// var students []models.StudentInput

	students, err := h.srv.GetAllStudent()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student"})
		return
	}
	c.IndentedJSON(http.StatusOK, students)

}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	// var students []models.StudentInput
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse student ID from path"})
		return
	}
	student, err := h.srv.GetStudent(int(idInt))
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Student"})
		return
	}

	c.IndentedJSON(http.StatusOK, student)

}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	// var students []models.StudentInput
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse student ID from path"})
		return
	}
	err = h.srv.DeleteStudent(int(idInt))
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.IndentedJSON(http.StatusOK, nil)

}

func (h *StudentHandler) PutStudent(c *gin.Context) {
	// var students []models.StudentInput
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse student ID from path"})
		return
	}
	var s models.StudentInput
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateStudent, err := h.srv.UpdateStudent(int(idInt), s)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.IndentedJSON(http.StatusOK, updateStudent)

}
