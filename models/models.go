package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID          uint32    `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	PublicEmail string    `json:"public_email" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	BirthDate   time.Time `json:"birth_date"`
	QSCoins     int64     `json:"qs_coins" gorm:"default:1500"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MagicLink struct {
	ID        uint32    `json:"id" gorm:"primaryKey"`
	UserID    uint32    `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	// Foreign keys
	User User `gorm:"foreignKey:UserID"`
}

// Test model
type Test struct {
	TestID                     uint32         `json:"test_id" gorm:"primaryKey"`
	ExaminerID                 uint32         `json:"examiner_id" gorm:"not null"`
	TestName                   string         `json:"test_name" gorm:"not null"`
	QSCoins                    int64          `json:"qs_coins" gorm:"default:0"`
	QuestionAnswerJSON         string         `json:"questions_json" gorm:"type:jsonb"`
	TestDuration               uint8          `json:"test_duration" gorm:"not null"`
	TotalMarks                 int16          `json:"total_marks" gorm:"not null"`
	NumberOfQuestionsPerTest   uint8          `json:"number_of_questions_per_test" gorm:"not null"`
	SizeOfQuestionPool         uint16         `json:"size_of_question_pool" gorm:"not null"`
	NumberOfTopics             uint8          `json:"number_of_topics" gorm:" not null"`
	Images                     pq.StringArray `json:"images" gorm:"type:text[]"`
	NumberOfStudents           uint32         `json:"number_of_students"`
	NumberOfAttemptsPerStudent uint8          `json:"number_of_attempts_per_student"`
	TestStartTime              time.Time      `json:"test_start_time" gorm:"not null"`
	TestEndTime                time.Time      `json:"test_end_time" gorm:"not null"`
	TestActive                 bool           `json:"test_active" gorm:"default:false"`
	CreatedAt                  time.Time      `json:"created_at"`
	// Foreign keys
	// Examiner User `gorm:"foreignKey:ExaminerID"`
}

// TestAssignedToUser model
type TestAssignedToUser struct {
	SomethingID      uint32 `json:"something_id" gorm:"primaryKey"`
	TestID           uint32 `json:"test_id" gorm:"not null"`
	CandidateID      uint32 `json:"candidate_id" gorm:"not null"`
	CandidateEmail   string `json:"candidate_email" gorm:"not null"`
	AttemptsAlloted  uint8  `json:"attempts_alloted" gorm:"not null"`
	AttemptRemaining uint8  `json:"attempt_remaining"`
	// Foreign keys
	// Candidate User `gorm:"foreignKey:CandidateID"`
}

// Answer model
type AnswerAttempt struct {
	AnswerID       uint64    `json:"answer_id" gorm:"primaryKey"`
	TestID         uint32    `json:"test_id" gorm:"not null"`
	CandidateID    uint32    `json:"candidate_id" gorm:"not null"`
	StartTime      time.Time `json:"start_time"`
	Duration       uint8     `json:"duration"`
	QuestionJSON   string    `json:"question_json" gorm:"type:jsonb"`
	AnswerJSON     string    `json:"answer_json" gorm:"type:jsonb"`
	EvaluationJSON string    `json:"evaluation_json" gorm:"type:jsonb"`
	AchievedMarks  uint8     `json:"achieved_marks"`
	// Foreign keys
	// Candidate User `gorm:"foreignKey:CandidateID"`
}

// PaymentTable model
type PaymentTable struct {
	OrderID           uint32    `json:"order_id" gorm:"primaryKey"`
	RazorpayOrderID   string    `json:"razorpay_order_id"`
	RazorpayPaymentID string    `json:"razorpay_payment_id"`
	RazorpaySignature string    `json:"razorpay_signature"`
	UserID            uint32    `json:"user_id" gorm:"not null"`
	Amount            int32     `json:"amount" gorm:"not null"`
	Currency          string    `json:"currency" gorm:"not null"`
	QSCoinsPurchased  int64     `json:"qs_coins_purchased" gorm:"not null"`
	PaymentStatus     string    `json:"payment_status" gorm:"default:'pending'"`
	DateTime          time.Time `json:"date_time"`
	// Foreign keys
	// User User `gorm:"foreignKey:UserID"`
}

// Certificate Table
