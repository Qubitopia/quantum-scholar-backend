package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/models"

	"github.com/gin-gonic/gin"
)

type SendTestQuestionRequest struct {
	AnswerAttemptID uint64 `json:"answer_attempt_id" binding:"required"`
	SectionNumber   int    `json:"section_number" binding:"required"`
	QuestionNumber  int    `json:"question_number" binding:"required"`
}

type UpdateTestAttemptAnswerRequest struct {
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

func ListAssignedTestToUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	candidate, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	type TestInfo struct {
		TestID            uint32
		TestName          string
		TestDuration      uint8
		TestStartTime     time.Time
		TestEndTime       time.Time
		AttemptsAlloted   uint8
		AttemptsRemaining uint8
		TotalMarks        int16
	}

	var testInfoList []TestInfo

	err := database.DB.
		Table("test_assigned_to_users tau").
		Select(`
		t.test_id,
		t.test_name,
		t.test_duration,
		t.test_start_time,
		t.test_end_time,
		t.total_marks,
		tau.attempts_alloted,
		tau.attempt_remaining as attempts_remaining
		`).
		Joins("JOIN tests t ON t.test_id = tau.test_id").
		Where("tau.candidate_id = ?", candidate.ID).
		Scan(&testInfoList).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assigned tests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tests": testInfoList,
	})
}

func createQuestionAnswerJSON(test_id uint32, candidate_id uint32) (uint32, json.RawMessage, error) {
	// 1) Fetch test by id to get QuestionAnswerJSON
	var test models.Test
	if err := database.DB.Where("test_id = ?", test_id).First(&test).Error; err != nil {
		return 0, nil, err
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
		return 0, nil, err
	}

	// Build stored-test metadata JSON without the per-section questions payload.
	type storedSectionMeta struct {
		SectionID         int    `json:"section_id"`
		Title             string `json:"title"`
		NumberOfQuestions int    `json:"number_of_questions"`
	}
	type storedTestMeta struct {
		Title    string              `json:"title"`
		Sections []storedSectionMeta `json:"sections"`
	}

	metaTest := storedTestMeta{Title: sTest.Title}
	for _, sec := range sTest.Sections {
		metaTest.Sections = append(metaTest.Sections, storedSectionMeta{
			SectionID:         sec.SectionID,
			Title:             sec.Title,
			NumberOfQuestions: sec.QuestionsToDisplay,
		})
	}

	metaBytes, err := json.Marshal(metaTest)
	if err != nil {
		return 0, nil, err
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
		return 0, nil, err
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
		return 0, nil, err
	}

	return uint32(attempt.AnswerAttemptID), json.RawMessage(metaBytes), nil
}

func StartTestAttempt(c *gin.Context) {
	// 1) Get candidate from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	candidate, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	// 2) Get test_id from URL param and validate
	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	// 3) Create question set for candidate and store in AnswerAttempt
	answerAttemptID, testInfoJSON, err := createQuestionAnswerJSON(uint32(testID), candidate.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question set"})
		return
	}

	// 4) Update StartTime of attempt to now
	if err := database.DB.Model(&models.AnswerAttempt{}).Where("answer_attempt_id = ?", answerAttemptID).Update("start_time", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start test attempt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"answer_attempt_id": answerAttemptID,
		"test_id":           testID,
		"test":              testInfoJSON,
	})
}

func GetTestQuestion(c *gin.Context) {
	// This handler will get answer_attempt_id, section_number and question_number as input and return the question JSON for that question number which will consist of question text, question_type, options (if mcq/msq), success marks & failure marks.
	// 1) Get candidate from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	candidate, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	// 2) Get and validate input JSON
	var req SendTestQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3) Fetch AnswerAttempt by id  where candidate_id = candidate.ID and validate ownership
	var attempt models.AnswerAttempt
	if err := database.DB.Where("answer_attempt_id = ? AND candidate_id = ?", req.AnswerAttemptID, candidate.ID).First(&attempt).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}

	// 4) Check if attempt is still active based on start time and duration + 5 mins grace period
	if !attempt.StartTime.IsZero() {
		now := time.Now()
		endTime := attempt.StartTime.Add(time.Duration(attempt.Duration)*time.Minute + 5*time.Minute)
		if now.After(endTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Test attempt has expired"})
			return
		}
	}

	// 5) Unmarshal attempt.QuestionJSON and return the question matching section_number and question_number(as per index)
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

	var oTest outTest
	if err := json.Unmarshal([]byte(attempt.QuestionJSON), &oTest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse question JSON"})
		return
	}

	for _, sec := range oTest.Sections {
		if sec.SectionID == req.SectionNumber {
			for _, q := range sec.Questions {
				if q.QuestionNumber == req.QuestionNumber {
					c.JSON(http.StatusOK, gin.H{
						"question": q,
					})
					return
				}
			}
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Question not found"})

}

func UpdateTestAttemptAnswer(c *gin.Context) {
	// 1) Get candidate from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	candidate, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	// 2) Get and validate input JSON
	var req UpdateTestAttemptAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	

	// Check if the attempt exists and belongs to the user
	var attempt models.AnswerAttempt
	if err := database.DB.Where("answer_id = ? AND candidate_id = ?", req.AttemptId, candidate.ID).First(&attempt).Error; err != nil {
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


	// Validate request structure matches UpdateTestAttemptAnswerRequest
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
