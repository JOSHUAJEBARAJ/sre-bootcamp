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
	ctx := c.Request.Context()
	if err := c.BindJSON(&newStudent); err != nil {
		log.WithError(err).Error("Error while Converting the data")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outputStudent, err := h.srv.AddStudent(ctx, newStudent)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}
	c.IndentedJSON(http.StatusCreated, outputStudent)

}

func (h *StudentHandler) Healthcheck(c *gin.Context) {
	ctx := c.Request.Context()
	err := h.srv.PingDB(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"dependencies": gin.H{
			"database": "connected",
		},
	})
}

func (h *StudentHandler) GetAllStudent(c *gin.Context) {
	// var students []models.StudentInput
	ctx := c.Request.Context()
	students, err := h.srv.GetAllStudent(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student"})
		return
	}
	c.IndentedJSON(http.StatusOK, students)

}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	// var students []models.StudentInput
	ctx := c.Request.Context()
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse student ID from path"})
		return
	}
	student, err := h.srv.GetStudent(ctx, int(idInt))
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
	ctx := c.Request.Context()
	idInt, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse student ID from path"})
		return
	}
	err = h.srv.DeleteStudent(ctx, int(idInt))
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
	ctx := c.Request.Context()
	updateStudent, err := h.srv.UpdateStudent(ctx, int(idInt), s)
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
