package handlers

import (
	"net/http"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/models"
	"github.com/gin-gonic/gin"
)

type CreateNewTestRequest struct {
	TestDuration               uint8  `json:"test_duration" binding:"required"`
	TotalMarks                 int16  `json:"total_marks" binding:"required"`
	NumberOfQuestions          uint8  `json:"number_of_questions" binding:"required"`
	NumberOfOpenEndedQuestions uint8  `json:"number_of_open_ended_questions" binding:"required"`
	NumberOfStudents           uint32 `json:"number_of_students" binding:"required"`
	NumberOfAttempts           uint8  `json:"number_of_attempts" binding:"required"`
	StudentsRemaining          uint32 `json:"students_remaining" binding:"required"`
}

type UpdateTestRequest struct {
	TestID        uint32 `json:"test_id" binding:"required"`
	QuestionsJSON string `json:"questions_json"`
	AnswerJSON    string `json:"answer_json"`
}

func CreateNewTest(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	examiner, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	var req CreateNewTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the test in the database
	test := models.Test{
		ExaminerID:                 examiner.ID,
		TestDuration:               req.TestDuration,
		TotalMarks:                 req.TotalMarks,
		NumberOfQuestions:          req.NumberOfQuestions,
		NumberOfOpenEndedQuestions: req.NumberOfOpenEndedQuestions,
		NumberOfStudents:           req.NumberOfStudents,
		NumberOfAttempts:           req.NumberOfAttempts,
		StudentsRemaining:          req.StudentsRemaining,
		DateTimeCreated:            time.Now(),
		QuestionsJSON:              "{}",
		AnswerJSON:                 "{}",
	}
	if err := database.DB.Create(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Test created successfully", "test_id": test.TestID})
}

func UpdateTest(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	examiner, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	var req UpdateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the test from the database
	var test models.Test
	if err := database.DB.First(&test, req.TestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}

	// Verify the owner
	if test.ExaminerID != examiner.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this test"})
		return
	}

	// Update questions and answers if provided
	if req.QuestionsJSON != "" {
		test.QuestionsJSON = req.QuestionsJSON
	}
	if req.AnswerJSON != "" {
		test.AnswerJSON = req.AnswerJSON
	}

	if err := database.DB.Save(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test updated successfully"})
}
