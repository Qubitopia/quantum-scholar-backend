package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Qubitopia/QuantumScholar/server/database"
	"github.com/Qubitopia/QuantumScholar/server/models"
	"github.com/gin-gonic/gin"
)

type CreateNewTestRequest struct {
	TestName                   string `json:"test_name" binding:"required"`
	TestDuration               uint8  `json:"test_duration" binding:"required"`
	TotalMarks                 int16  `json:"total_marks" binding:"required"`
	NumberOfQuestions          uint8  `json:"number_of_questions" binding:"required"`
	NumberOfOpenEndedQuestions uint8  `json:"number_of_open_ended_questions" binding:"required"`
	NumberOfStudents           uint32 `json:"number_of_students" binding:"required"`
	NumberOfAttempts           uint8  `json:"number_of_attempts" binding:"required"`
	StudentsRemaining          uint32 `json:"students_remaining" binding:"required"`
}

type UpdateQuestionsAndAnswersInTestRequest struct {
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
		TestName:                   req.TestName,
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

func UpdateQuestionsAndAnswersInTest(c *gin.Context) {
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

	var req UpdateQuestionsAndAnswersInTestRequest
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

func GetAllTestsCreatedByUser(c *gin.Context) {
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

	type TestSummary struct {
		TestID                     uint32    `json:"test_id"`
		TestName                   string    `json:"test_name"`
		QSCoins                    int64     `json:"qs_coins"`
		TestActive                 bool      `json:"test_active"`
		TestDuration               uint8     `json:"test_duration"`
		TotalMarks                 int16     `json:"total_marks"`
		NumberOfQuestions          uint8     `json:"number_of_questions"`
		NumberOfOpenEndedQuestions uint8     `json:"number_of_open_ended_questions"`
		NumberOfStudents           uint32    `json:"number_of_students"`
		NumberOfAttempts           uint8     `json:"number_of_attempts"`
		StudentsRemaining          uint32    `json:"students_remaining"`
		DateTimeCreated            time.Time `json:"date_time_created"`
	}

	var tests []TestSummary
	if err := database.DB.Model(&models.Test{}).
		Where("examiner_id = ?", examiner.ID).
		Select("test_id, test_name, qs_coins, test_active, test_duration, total_marks, number_of_questions, number_of_open_ended_questions, number_of_students, number_of_attempts, students_remaining, date_time_created").
		Scan(&tests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tests": tests,
	})
}

func GetTestByID(c *gin.Context) {
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

	testIDParam := c.Param("id")
	var testID uint32
	_, err := fmt.Sscanf(testIDParam, "%d", &testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	var test models.Test
	if err := database.DB.First(&test, testID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}

	if test.ExaminerID != examiner.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this test"})
		return
	}

	c.JSON(http.StatusOK, test)
}
