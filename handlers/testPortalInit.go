package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	// "log"
	"net/http"
	"time"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/models"

	"github.com/gin-gonic/gin"
)

// Format for testing portal
type PortalLoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	BirthDate string `json:"birthdate" binding:"required"` // YYYY-MM-DD
}

type PortalVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

type InitTestRequest struct {
	Email  string `json:"email" binding:"required,email"`
	Token  string `json:"token" binding:"required"`
	TestID uint32 `json:"test_id" binding:"required"`
}

type StartTestAttemptRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Token     string `json:"token" binding:"required"`
	TestID    uint32 `json:"test_id" binding:"required"`
	AttemptID uint32 `json:"attempt_id" binding:"required"`
}

type AnswerPattern struct {
	Sections []struct {
		SectionId int `json:"sectionId"`
		Answers   []struct {
			QuestionNumber int     `json:"questionNumber"`
			CorrectOption  *int    `json:"CorrectOption,omitempty"`
			CorrectOptions []int   `json:"CorrectOptions,omitempty"`
			Answer         *string `json:"answer,omitempty"`
		} `json:"answers"`
	} `json:"sections"`
}

type UpdateTestAttemptRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Token     string `json:"token" binding:"required"`
	AttemptId uint32 `json:"attempt_id" binding:"required"`
	Answer    struct {
		Sections []struct {
			SectionId int `json:"sectionId"`
			Answers   []struct {
				QuestionNumber int     `json:"questionNumber"`
				CorrectOption  *int    `json:"CorrectOption,omitempty"`
				CorrectOptions []int   `json:"CorrectOptions,omitempty"`
				Answer         *string `json:"answer,omitempty"`
			} `json:"answers"`
		} `json:"sections"`
	} `json:"answer" binding:"required"`
}

// APIs for testing portal
func Generate64AsciiToken() (string, error) {
	bytes := make([]byte, 48)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(bytes), nil
}

func TestPortalLogin(c *gin.Context) {
	var req PortalLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if birthdate matches
	if user.BirthDate.Format("2006-01-02") != req.BirthDate {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid birthdate"})
		return
	}

	// Add the email and the token to Redis with 5 hours expiry
	token, err := Generate64AsciiToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	if err := database.RedisClient.Set(context.Background(), "email:"+req.Email, token, 15*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token in Redis"})
		return
	}

	// Check for the tests assigned to the user and also send the test id and name in the response
	var assignedTests []models.TestAssignedToUser
	if err := database.DB.Where("candidate_email = ?", req.Email).Find(&assignedTests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No tests are assigned to this user, please contact the examiner"})
		return
	}

	// Get test names
	type TestInfo struct {
		TestID        uint32    `json:"test_id"`
		TestName      string    `json:"test_name"`
		TestStartTime time.Time `json:"test_start_time"`
		TestEndTime   time.Time `json:"test_end_time"`
	}

	// Prepare test info list
	var TestInfoList []TestInfo
	for _, assignedTest := range assignedTests {
		var test models.Test
		if err := database.DB.Where("test_id = ?", assignedTest.TestID).First(&test).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test information"})
			return
		}

		TestInfoList = append(TestInfoList, TestInfo{
			TestID:        test.TestID,
			TestName:      test.TestName,
			TestStartTime: test.TestStartTime,
			TestEndTime:   test.TestEndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"tests":   TestInfoList,
	})
}

func TestPortalVerify(c *gin.Context) {
	var req PortalVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if token is validon redis for the email, if yes then bump the expiry by next 15 minutes
	storedToken, err := database.RedisClient.Get(context.Background(), "email:"+req.Email).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token from Redis"})
		return
	}

	if storedToken != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Bump the expiry by next 15 minutes
	if err := database.RedisClient.Expire(context.Background(), "email:"+req.Email, 15*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend token expiry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token verified successfully"})
}

func CreateQuestionAnswerJSON(test_id uint32, candidate_id uint32) (uint32, error) {
	// 1) Fetch test by id to get QuestionAnswerJSON
	var test models.Test
	if err := database.DB.Where("test_id = ?", test_id).First(&test).Error; err != nil {
		return 0, err
	}

	// 2) Define structures matching stored Test.QuestionAnswerJSON (examiner view)
	type storedQuestion struct {
		QuestionNumber int      `json:"questionNumber"`
		Type           string   `json:"type"`
		SuccessMarks   int      `json:"successMarks"`
		FailureMarks   int      `json:"failureMarks"`
		QuestionText   string   `json:"questionText"`
		Options        []string `json:"options,omitempty"`
		// Fields to be omitted in candidate view
		CorrectOption  *int   `json:"correctOption,omitempty"`
		CorrectOptions []int  `json:"correctOptions,omitempty"`
		ModelAnswer    string `json:"modelAnswer,omitempty"`
	}
	type storedSection struct {
		SectionID          int              `json:"sectionId"`
		Title              string           `json:"title"`
		QuestionsToDisplay int              `json:"questionsToDisplay"`
		Questions          []storedQuestion `json:"questions"`
	}
	type storedTest struct {
		Title    string          `json:"title"`
		Sections []storedSection `json:"sections"`
	}

	// 3) Unmarshal stored JSON
	var sTest storedTest
	if err := json.Unmarshal([]byte(test.QuestionAnswerJSON), &sTest); err != nil {
		return 0, err
	}

	// 4) Define candidate-facing output structures (like tests/q1.json)
	type outQuestion struct {
		FailureMarks   int      `json:"failureMarks"`
		QuestionNumber int      `json:"questionNumber"`
		QuestionText   string   `json:"questionText"`
		SuccessMarks   int      `json:"successMarks"`
		Type           string   `json:"type"`
		Options        []string `json:"options,omitempty"`
	}
	type outSection struct {
		SectionID int           `json:"sectionId"`
		Title     string        `json:"title"`
		Questions []outQuestion `json:"questions"`
	}
	type outTest struct {
		Sections []outSection `json:"sections"`
		Title    string       `json:"title"`
	}

	// 5) Build candidate-facing question set, respecting questionsToDisplay and no repeats
	var oTest outTest
	oTest.Title = sTest.Title
	for _, sec := range sTest.Sections {
		// Determine how many questions to pick for this section
		nTotal := len(sec.Questions)
		k := sec.QuestionsToDisplay
		if k <= 0 || k > nTotal {
			k = nTotal
		}

		picked := map[int]bool{}
		outQs := make([]outQuestion, 0, k)
		for len(outQs) < k && nTotal > 0 {
			idx := randInt(nTotal, picked)
			picked[idx] = true

			q := sec.Questions[idx]
			oq := outQuestion{
				FailureMarks:   q.FailureMarks,
				QuestionNumber: q.QuestionNumber,
				QuestionText:   q.QuestionText,
				SuccessMarks:   q.SuccessMarks,
				Type:           q.Type,
			}
			// Include options for MCQ/MSQ, but never include correct answers/model answers
			if q.Type == "mcq" || q.Type == "msq" {
				oq.Options = append(oq.Options, q.Options...)
			}
			outQs = append(outQs, oq)
		}

		oTest.Sections = append(oTest.Sections, outSection{
			SectionID: sec.SectionID,
			Title:     sec.Title,
			Questions: outQs,
		})
	}

	// 6) Marshal output JSON for storing in AnswerAttempt.QuestionJSON
	qb, err := json.Marshal(oTest)
	if err != nil {
		return 0, err
	}

	// 7) Store in AnswerAttempt table
	attempt := models.AnswerAttempt{
		TestID:         test_id,
		CandidateID:    candidate_id,
		StartTime:      time.Time{},
		Duration:       test.TestDuration,
		QuestionJSON:   string(qb),
		AnswerJSON:     "{}",
		EvaluationJSON: "{}",
		AchievedMarks:  0,
	}

	if err := database.DB.Create(&attempt).Error; err != nil {
		return 0, err
	}

	return uint32(attempt.AnswerID), nil
}

// randInt returns a random int in [0, n) not in picked
func randInt(n int, picked map[int]bool) int {
	for {
		b := make([]byte, 1)
		rand.Read(b)
		idx := int(b[0]) % n
		if !picked[idx] {
			return idx
		}
	}
}

func InitTestForCandidate(c *gin.Context) {
	var req InitTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if token is validon redis for the email, if yes then proceed
	storedToken, err := database.RedisClient.Get(context.Background(), "email:"+req.Email).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token from Redis"})
		return
	}

	if storedToken != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Bump the expiry by next 15 minutes
	if err := database.RedisClient.Expire(context.Background(), "email:"+req.Email, 15*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend token expiry"})
		return
	}

	// Check if the user is assigned to the test and has attempts remaining
	var assignedTest models.TestAssignedToUser
	if err := database.DB.Where("candidate_email = ? AND test_id = ?", req.Email, req.TestID).First(&assignedTest).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not assigned to this test"})
		return
	}
	if assignedTest.AttemptRemaining <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "No attempts remaining for this test"})
		return
	}

	// Create question set for this candidate and store in AnswerAttempt table
	answerID, err := CreateQuestionAnswerJSON(req.TestID, assignedTest.CandidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AnswerAttempt: " + err.Error()})
		return
	}

	// Decrement attempt remaining by 1
	assignedTest.AttemptRemaining -= 1
	if err := database.DB.Save(&assignedTest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attempts remaining"})
		return
	}

	// Send success response with attempt id
	c.JSON(http.StatusOK, gin.H{"message": "Test initialized successfully", "attempt_id": answerID})
}

func StartTestAttempt(c *gin.Context) {
	var req StartTestAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if token is validon redis for the email, if yes then proceed
	storedToken, err := database.RedisClient.Get(context.Background(), "email:"+req.Email).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token from Redis"})
		return
	}

	if storedToken != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Bump the expiry by next 15 minutes
	if err := database.RedisClient.Expire(context.Background(), "email:"+req.Email, 15*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend token expiry"})
		return
	}

	// Check if the start time is null, if not then return error
	var attempt models.AnswerAttempt
	if err := database.DB.Where("answer_id = ? AND test_id = ?", req.AttemptID, req.TestID).First(&attempt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Your test attempt has not been initialized"})
		return
	}
	if !attempt.StartTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Test has already been started"})
		return
	}
	// Update the start time to current time
	// Start the test
	attempt.StartTime = time.Now()
	if err := database.DB.Save(&attempt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start test"})
		return
	}

	// Fetch test duration (in minutes) from the Test table
	var test models.Test
	if err := database.DB.Select("test_duration").Where("test_id = ?", req.TestID).First(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test duration"})
		return
	}

	// send success response with QuestionJSON
	c.JSON(http.StatusOK, gin.H{
		"message":          "Test started successfully",
		"question_json":    attempt.QuestionJSON,
		"duration_minutes": test.TestDuration,
	})
}

func UpdateTestAttempt(c *gin.Context) {
	var req UpdateTestAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate request structure matches UpdateTestAttemptRequest
	if len(req.Answer.Sections) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Answer must contain at least one section"})
		return
	}
	lastSectionId := 0
	for i, section := range req.Answer.Sections {
		if section.SectionId <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Each section must have a valid sectionId (>0)"})
			return
		}
		if section.Answers == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Each section must contain answers array"})
			return
		}
		lastSectionId = section.SectionId
		for _, ans := range section.Answers {
			if ans.QuestionNumber <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Each answer must have a valid questionNumber (>0)"})
				return
			}
			if ans.CorrectOption == nil && len(ans.CorrectOptions) == 0 && (ans.Answer == nil || *ans.Answer == "") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Each answer must have one of: CorrectOption, CorrectOptions, or answer"})
				return
			}
		}
		// Ensure section ids are sequential and match index+1
		if section.SectionId != i+1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SectionId must be sequential starting from 1"})
			return
		}
	}
	if len(req.Answer.Sections) != lastSectionId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number of sections mismatch"})
		return
	}

	// Check if token is validon redis for the email, if yes then proceed
	storedToken, err := database.RedisClient.Get(context.Background(), "email:"+req.Email).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token from Redis"})
		return
	}

	if storedToken != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Bump the expiry by next 15 minutes
	if err := database.RedisClient.Expire(context.Background(), "email:"+req.Email, 15*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend token expiry"})
		return
	}

	// Check if the attempt exists and belongs to the test
	var attempt models.AnswerAttempt
	if err := database.DB.Where("answer_id = ?", req.AttemptId).First(&attempt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attempt not found"})
		return
	}

	// Ensure the attempt has been started and is within the allowed duration (+5 min grace)
	if attempt.StartTime.IsZero() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Test attempt has not been started"})
		return
	}
	allowedEnd := attempt.StartTime.Add(time.Duration(attempt.Duration) * time.Minute).Add(5 * time.Minute)
	if time.Now().After(allowedEnd) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Time window for this attempt has expired"})
		return
	}

	// Update the answers and reset marks
	answerJSONBytes, err := json.Marshal(req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal answer"})
		return
	}
	attempt.AnswerJSON = string(answerJSONBytes)
	attempt.AchievedMarks = 0

	if err := database.DB.Save(&attempt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test attempt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test attempt updated successfully"})
}
