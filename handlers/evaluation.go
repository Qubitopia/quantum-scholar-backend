package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/mail"
	"github.com/Qubitopia/quantum-scholar-backend/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

type evaluationQuestion struct {
	QuestionNumber int    `json:"questionNumber"`
	Type           string `json:"type"`
	SuccessMarks   int    `json:"successMarks"`
	FailureMarks   int    `json:"failureMarks"`
	QuestionText   string `json:"questionText,omitempty"`
	CorrectOption  *int   `json:"correctOption,omitempty"`
	CorrectOptions []int  `json:"correctOptions,omitempty"`
	ModelAnswer    string `json:"modelAnswer,omitempty"`
}

type evaluationSection struct {
	SectionID int                  `json:"sectionId"`
	Questions []evaluationQuestion `json:"questions"`
}

type evaluationTest struct {
	Title    string              `json:"title"`
	Sections []evaluationSection `json:"sections"`
}

type evaluationAnswer struct {
	QuestionNumber int     `json:"questionNumber"`
	CorrectOption  *int    `json:"CorrectOption,omitempty"`
	CorrectOptions []int   `json:"CorrectOptions,omitempty"`
	Answer         *string `json:"answer,omitempty"`
}

type evaluationAnswerSection struct {
	SectionID int                `json:"sectionId"`
	Answers   []evaluationAnswer `json:"answers"`
}

type evaluationAttemptAnswers struct {
	Sections []evaluationAnswerSection `json:"sections"`
}

type evaluationAnswerResult struct {
	QuestionNumber int    `json:"questionNumber"`
	Type           string `json:"type"`
	IsCorrect      bool   `json:"isCorrect"`
	AwardedMarks   int    `json:"awardedMarks"`
	SubmittedText  string `json:"submittedText,omitempty"`
}

type evaluationSectionResult struct {
	SectionID int                      `json:"sectionId"`
	Questions []evaluationAnswerResult `json:"questions"`
}

type evaluationResult struct {
	TestID          uint32                    `json:"testId"`
	TestTitle       string                    `json:"testTitle"`
	AnswerAttemptID uint64                    `json:"answerAttemptId"`
	CandidateID     uint32                    `json:"candidateId"`
	TotalMarks      int                       `json:"totalMarks"`
	ObtainedMarks   int                       `json:"obtainedMarks"`
	Sections        []evaluationSectionResult `json:"sections"`
}

// EvaluateTest evaluates all attempts for a test, persists marks, syncs attempt usage, and sends score emails.
func EvaluateTest(c *gin.Context) {
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

	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	var test models.Test
	if err := database.DB.Where("test_id = ? AND examiner_id = ?", testID, examiner.ID).First(&test).Error; err != nil {
		log.Printf("evaluate test %d: failed to fetch test: %v", testID, err)
		return
	}

	var canonical evaluationTest
	if err := json.Unmarshal([]byte(test.QuestionAnswerJSON), &canonical); err != nil {
		log.Printf("evaluate test %d: failed to parse test JSON: %v", testID, err)
		return
	}

	questionLookup := make(map[int]map[int]evaluationQuestion)
	for _, section := range canonical.Sections {
		if _, ok := questionLookup[section.SectionID]; !ok {
			questionLookup[section.SectionID] = make(map[int]evaluationQuestion)
		}
		for _, question := range section.Questions {
			questionLookup[section.SectionID][question.QuestionNumber] = question
		}
	}

	var attempts []models.AnswerAttempt
	if err := database.DB.Where("test_id = ?", testID).Find(&attempts).Error; err != nil {
		log.Printf("evaluate test %d: failed to fetch attempts: %v", testID, err)
		return
	}

	var assignments []models.TestAssignedToUser
	if err := database.DB.Where("test_id = ?", testID).Find(&assignments).Error; err != nil {
		log.Printf("evaluate test %d: failed to fetch assignments: %v", testID, err)
		return
	}

	assignmentByCandidate := make(map[uint32]*models.TestAssignedToUser, len(assignments))
	for i := range assignments {
		assignmentByCandidate[assignments[i].CandidateID] = &assignments[i]
	}

	attemptCountByCandidate := make(map[uint32]int)
	bestScoreByCandidate := make(map[uint32]int16)
	hasScoreByCandidate := make(map[uint32]bool)
	for i := range attempts {
		attempt := &attempts[i]
		attemptCountByCandidate[attempt.CandidateID]++

		score, resultJSON, err := evaluateAttemptScore(uint32(testID), canonical.Title, int(test.TotalMarks), attempt, questionLookup)
		if err != nil {
			log.Printf("evaluate test %d attempt %d: %v", testID, attempt.AnswerAttemptID, err)
			continue
		}

		attempt.AchievedMarks = int16(score)
		attempt.EvaluationJSON = resultJSON
		if err := database.DB.Save(attempt).Error; err != nil {
			log.Printf("evaluate test %d attempt %d: failed to save attempt: %v", testID, attempt.AnswerAttemptID, err)
			continue
		}

		if !hasScoreByCandidate[attempt.CandidateID] || attempt.AchievedMarks > bestScoreByCandidate[attempt.CandidateID] {
			bestScoreByCandidate[attempt.CandidateID] = attempt.AchievedMarks
			hasScoreByCandidate[attempt.CandidateID] = true
		}

		var candidate models.User
		if err := database.DB.Where("id = ?", attempt.CandidateID).First(&candidate).Error; err != nil {
			log.Printf("evaluate test %d attempt %d: failed to fetch candidate: %v", testID, attempt.AnswerAttemptID, err)
			continue
		}

		if err := mail.SendEmailTestEvaluation(candidate.Email, candidate.Name, canonical.Title, score, int(test.TotalMarks)); err != nil {
			log.Printf("evaluate test %d attempt %d: failed to send score email: %v", testID, attempt.AnswerAttemptID, err)
		}
	}

	for candidateID, assignment := range assignmentByCandidate {
		usedAttempts := attemptCountByCandidate[candidateID]
		remaining := int(assignment.AttemptsAlloted) - usedAttempts
		if remaining < 0 {
			remaining = 0
		}
		assignment.AttemptRemaining = uint8(remaining)
		if hasScoreByCandidate[candidateID] {
			assignment.BestScore = bestScoreByCandidate[candidateID]
		}
		if err := database.DB.Save(assignment).Error; err != nil {
			log.Printf("evaluate test %d candidate %d: failed to update attempts remaining: %v", testID, candidateID, err)
		}
	}
}

func evaluateAttemptScore(testID uint32, testTitle string, totalMarks int, attempt *models.AnswerAttempt, questionLookup map[int]map[int]evaluationQuestion) (int, string, error) {
	var attemptAnswers evaluationAttemptAnswers
	if err := json.Unmarshal([]byte(attempt.AnswerJSON), &attemptAnswers); err != nil {
		return 0, "", fmt.Errorf("parse answer json: %w", err)
	}

	result := evaluationResult{
		TestID:          testID,
		TestTitle:       testTitle,
		AnswerAttemptID: attempt.AnswerAttemptID,
		CandidateID:     attempt.CandidateID,
	}

	total := 0
	for _, section := range attemptAnswers.Sections {
		sectionResult := evaluationSectionResult{SectionID: section.SectionID}
		for _, answer := range section.Answers {
			question, ok := questionLookup[section.SectionID][answer.QuestionNumber]
			if !ok {
				sectionResult.Questions = append(sectionResult.Questions, evaluationAnswerResult{
					QuestionNumber: answer.QuestionNumber,
					Type:           "unknown",
					IsCorrect:      false,
					AwardedMarks:   0,
				})
				continue
			}

			marks, isCorrect, submittedText := scoreAnswer(question, answer)
			total += marks
			sectionResult.Questions = append(sectionResult.Questions, evaluationAnswerResult{
				QuestionNumber: answer.QuestionNumber,
				Type:           question.Type,
				IsCorrect:      isCorrect,
				AwardedMarks:   marks,
				SubmittedText:  submittedText,
			})
		}
		result.Sections = append(result.Sections, sectionResult)
	}

	result.TotalMarks = totalMarks
	result.ObtainedMarks = total
	encoded, err := json.Marshal(result)
	if err != nil {
		return 0, "", fmt.Errorf("marshal evaluation json: %w", err)
	}

	return total, string(encoded), nil
}

func scoreAnswer(question evaluationQuestion, answer evaluationAnswer) (int, bool, string) {
	questionType := strings.ToLower(strings.TrimSpace(question.Type))
	switch questionType {
	case "mcq":
		if question.CorrectOption != nil && answer.CorrectOption != nil && *question.CorrectOption == *answer.CorrectOption {
			return question.SuccessMarks, true, fmt.Sprintf("%d", *answer.CorrectOption)
		}
		if answer.CorrectOption != nil {
			return question.FailureMarks, false, fmt.Sprintf("%d", *answer.CorrectOption)
		}
	case "msq":
		if len(answer.CorrectOptions) > 0 && equalIntSets(question.CorrectOptions, answer.CorrectOptions) {
			return question.SuccessMarks, true, fmt.Sprintf("%v", answer.CorrectOptions)
		}
		if len(answer.CorrectOptions) > 0 {
			return question.FailureMarks, false, fmt.Sprintf("%v", answer.CorrectOptions)
		}
	default:
		if answer.Answer != nil {
			submitted := normalizeText(*answer.Answer)
			expected := normalizeText(question.ModelAnswer)

			if submitted != "" {
				aiMarks, aiCorrect, aiReason, aiErr := scoreOpenEndedWithGemini(question, *answer.Answer)
				if aiErr == nil {
					submittedWithReason := strings.TrimSpace(*answer.Answer)
					if aiReason != "" {
						submittedWithReason = fmt.Sprintf("%s | AI: %s", submittedWithReason, aiReason)
					}
					return aiMarks, aiCorrect, submittedWithReason
				}

				log.Printf("gemini scoring failed for question %d: %v", question.QuestionNumber, aiErr)
			}

			if submitted != "" && submitted == expected {
				return question.SuccessMarks, true, *answer.Answer
			}
			if submitted != "" {
				return question.FailureMarks, false, *answer.Answer
			}
		}
	}

	return question.FailureMarks, false, ""
}

type geminiOpenEndedGrade struct {
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

func scoreOpenEndedWithGemini(question evaluationQuestion, submittedAnswer string) (int, bool, string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return 0, false, "", fmt.Errorf("create genai client: %w", err)
	}

	prompt := fmt.Sprintf(
		"You are grading one open-ended exam answer. Return ONLY compact JSON with fields: score (number from 0 to 1) and reason (short string under 120 chars).\\nQuestion: %s\\nExpected model answer: %s\\nCandidate answer: %s",
		strings.TrimSpace(question.QuestionText),
		strings.TrimSpace(question.ModelAnswer),
		strings.TrimSpace(submittedAnswer),
	)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return 0, false, "", fmt.Errorf("generate gemini content: %w", err)
	}

	rawText := strings.TrimSpace(result.Text())
	if rawText == "" {
		return 0, false, "", fmt.Errorf("empty gemini response")
	}
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)
	if strings.Contains(rawText, "{") && strings.Contains(rawText, "}") {
		start := strings.Index(rawText, "{")
		end := strings.LastIndex(rawText, "}")
		if start >= 0 && end > start {
			rawText = rawText[start : end+1]
		}
	}

	var grade geminiOpenEndedGrade
	if err := json.Unmarshal([]byte(rawText), &grade); err != nil {
		return 0, false, "", fmt.Errorf("parse gemini grade json: %w", err)
	}

	if grade.Score < 0 {
		grade.Score = 0
	}
	if grade.Score > 1 {
		grade.Score = 1
	}

	rawMarks := float64(question.FailureMarks) + grade.Score*float64(question.SuccessMarks-question.FailureMarks)
	awarded := int(math.Round(rawMarks))
	isCorrect := awarded >= question.SuccessMarks

	return awarded, isCorrect, strings.TrimSpace(grade.Reason), nil
}

func equalIntSets(expected []int, actual []int) bool {
	if len(expected) != len(actual) {
		return false
	}

	expectedCopy := append([]int(nil), expected...)
	actualCopy := append([]int(nil), actual...)
	sort.Ints(expectedCopy)
	sort.Ints(actualCopy)
	for i := range expectedCopy {
		if expectedCopy[i] != actualCopy[i] {
			return false
		}
	}

	return true
}

func normalizeText(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(value))), " ")
}

func clampToUint8(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

// Candidate facing endpoint to get scores of all attempts for a test
func GetScoreOfAllAttempts(c *gin.Context) {
	// fetch all the attempts of the candidate for the test and return their scores
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

	testIDStr := c.Param("test_id")
	testID, err := strconv.ParseUint(testIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	var attempts []models.AnswerAttempt
	if err := database.DB.Where("test_id = ? AND candidate_id = ?", testID, candidate.ID).Find(&attempts).Error; err != nil {
		log.Printf("get scores for test %d candidate %d: failed to fetch attempts: %v", testID, candidate.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attempts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"attempts": attempts})

}
