package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/models"
	"github.com/gin-gonic/gin"
)

type CreateNewTestRequest struct {
	TestName                 string `json:"test_name" binding:"required"`
	TestDuration             uint8  `json:"test_duration" binding:"required"`
	TotalMarks               int16  `json:"total_marks" binding:"required"`
	NumberOfQuestionsPerTest uint8  `json:"number_of_questions_per_test" binding:"required"`
	TestStartTime            string `json:"test_start_time" binding:"required"`
	TestEndTime              string `json:"test_end_time" binding:"required"`
}

type Question struct {
	QuestionNumber int      `json:"questionNumber" binding:"required"`
	Type           string   `json:"type" binding:"required"`
	SuccessMarks   int      `json:"successMarks" binding:"required"`
	FailureMarks   int      `json:"failureMarks" binding:"required"`
	QuestionText   string   `json:"questionText" binding:"required"`
	Options        []string `json:"options,omitempty"`
	CorrectOption  int      `json:"correctOption,omitempty"`
	CorrectOptions []int    `json:"correctOptions,omitempty"`
	ModelAnswer    string   `json:"modelAnswer,omitempty"`
}
type Section struct {
	SectionID          int        `json:"sectionId" binding:"required"`
	Title              string     `json:"title" binding:"required"`
	QuestionsToDisplay int        `json:"questionsToDisplay" binding:"required"`
	Questions          []Question `json:"questions" binding:"required"`
}
type TestFormat struct {
	Title    string    `json:"title" binding:"required"`
	Sections []Section `json:"sections" binding:"required"`
}
type UpdateQuestionsAndAnswersInTestRequest struct {
	TestID             uint32     `json:"test_id" binding:"required"`
	QuestionAnswerJSON TestFormat `json:"test" binding:"required"`
}

type AddCandidatesToTestRequest struct {
	TestID           uint32   `json:"test_id" binding:"required"`
	NumberOfAttempts uint8    `json:"number_of_attempts" binding:"required"`
	CandidateEmails  []string `json:"candidate_emails" binding:"required"`
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

	// Parse start and end times from string to time.Time
	testStartTime, err := time.Parse(time.RFC3339, req.TestStartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test_start_time format. Use RFC3339 format."})
		return
	}
	testEndTime, err := time.Parse(time.RFC3339, req.TestEndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test_end_time format. Use RFC3339 format."})
		return
	}

	// Create the test in the database
	test := models.Test{
		ExaminerID:               examiner.ID,
		TestName:                 req.TestName,
		TestDuration:             req.TestDuration,
		TotalMarks:               req.TotalMarks,
		NumberOfQuestionsPerTest: req.NumberOfQuestionsPerTest,
		SizeOfQuestionPool:       uint16(req.NumberOfQuestionsPerTest),
		NumberOfTopics:           0,
		TestStartTime:            testStartTime,
		TestEndTime:              testEndTime,
		CreatedAt:                time.Now(),
		QuestionAnswerJSON:       "{}",
		TopicJson:                "{}",
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

	// Recursive validation of the JSON structure (Go struct fields)
	if err := validateTestFormat(req.QuestionAnswerJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionAnswerJSONBytes, err := json.Marshal(req.QuestionAnswerJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question and answer format"})
		return
	}

	log.Println(req.QuestionAnswerJSON)

	// Update questions and answers in the test
	test.QuestionAnswerJSON = string(questionAnswerJSONBytes)

	// Save the updated test to the database
	if err := database.DB.Save(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test updated successfully"})
}

// validateTestFormat recursively validates required fields in TestFormat, Section, and Question
func validateTestFormat(tf TestFormat) error {
	if tf.Title == "" {
		return fmt.Errorf("test title is required")
	}
	if len(tf.Sections) == 0 {
		return fmt.Errorf("at least one section is required")
	}
	for i, section := range tf.Sections {
		if section.SectionID == 0 {
			return fmt.Errorf("section %d: sectionId is required", i+1)
		}
		if section.Title == "" {
			return fmt.Errorf("section %d: title is required", i+1)
		}
		if section.QuestionsToDisplay == 0 {
			return fmt.Errorf("section %d: questionsToDisplay is required", i+1)
		}
		if len(section.Questions) == 0 {
			return fmt.Errorf("section %d: at least one question is required", i+1)
		}
		for j, q := range section.Questions {
			if q.QuestionNumber == 0 {
				return fmt.Errorf("section %d, question %d: questionNumber is required", i+1, j+1)
			}
			if q.Type == "" {
				return fmt.Errorf("section %d, question %d: type is required", i+1, j+1)
			}
			if q.QuestionText == "" {
				return fmt.Errorf("section %d, question %d: questionText is required", i+1, j+1)
			}
			if q.SuccessMarks <= 0 {
				return fmt.Errorf("section %d, question %d: successMarks is required", i+1, j+1)
			}
			if q.FailureMarks > 0 {
				return fmt.Errorf("section %d, question %d: failureMarks should be zero or negative", i+1, j+1)
			}
			log.Println(len(q.Options))
			if q.Type == "mcq" && (len(q.Options) < 2 || (q.CorrectOption < 1 || q.CorrectOption > len(q.Options))) {
				return fmt.Errorf("section %d, question %d: mcq type requires at least 2 options and a correct option", i+1, j+1)
			}
			if q.Type == "msq" && (len(q.Options) < 2 || len(q.CorrectOptions) == 0) {
				return fmt.Errorf("section %d, question %d: msq type requires at least 2 options and at least one correct option", i+1, j+1)
			}
			if q.Type == "open-ended" && q.ModelAnswer == "" {
				return fmt.Errorf("section %d, question %d: open-ended type requires a model answer", i+1, j+1)
			}
		}
	}
	return nil
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

	// Fetch all tests created by the user
	var tests []models.Test
	if err := database.DB.Where("examiner_id = ?", examiner.ID).Find(&tests).Error; err != nil {
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

func AddCandidatesToTest(c *gin.Context) {
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

	var req AddCandidatesToTestRequest
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

	// Add candidates to the TestAssignedToUser from the request, if user not found, cretate a new user with default values
	for _, email := range req.CandidateEmails {
		var candidate models.User
		if err := database.DB.Where("email = ?", email).First(&candidate).Error; err != nil {
			// Create a new user with default values
			candidate = models.User{
				Email:       email,
				PublicEmail: email,
				Name:        email,
				QSCoins:     1500,
				IsActive:    true,
			}
			if err := database.DB.Create(&candidate).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		}

		// Add the candidate to the test
		testAssigned := models.TestAssignedToUser{
			TestID:           test.TestID,
			CandidateID:      candidate.ID,
			AttemptRemaining: req.NumberOfAttempts,
		}
		if err := database.DB.Create(&testAssigned).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add candidate to test"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Candidates added to test successfully"})
}
