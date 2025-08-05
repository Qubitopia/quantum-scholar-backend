package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Name      string    `json:"Name" gorm:"unique;not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MagicLink struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	// Foreign keys
	User User `gorm:"foreignKey:UserID"`
}

// Test model
type Test struct {
	TestID          uint      `json:"test_id" gorm:"primaryKey"`
	OwnerID         uint      `json:"owner_id" gorm:"not null"`
	TestLength      uint16    `json:"test_length" gorm:"not null"`
	QuestionsJSON   string    `json:"questions_json" gorm:"type:jsonb"`
	AnswerJSON      string    `json:"answer_json" gorm:"type:jsonb"`
	TotalMarks      int       `json:"total_marks"`
	OrderID         uint      `json:"order_id" gorm:"not null"`
	DateTimeCreated time.Time `json:"date_time_created"`
	TestActive      bool      `json:"test_active" gorm:"default:true"`
	Paid            bool      `json:"paid" gorm:"default:false"`
	// Foreign keys
	Owner User `gorm:"foreignKey:OwnerID"`
}

// Answer model
type Answer struct {
	AnswerID       uint      `json:"answer_id" gorm:"primaryKey"`
	TestID         uint      `json:"test_id" gorm:"not null"`
	StudentID      uint      `json:"student_id" gorm:"not null"`
	DateTime       time.Time `json:"date_time"`
	AnswerJSON     string    `json:"answer_json" gorm:"type:jsonb"`
	EvaluationJSON string    `json:"evaluation_json" gorm:"type:jsonb"`
	AchievedMarks  uint8     `json:"achieved_marks"`
}

// PaymentTable model
type PaymentTable struct {
	OrderID           uint      `json:"order_id" gorm:"primaryKey"`
	RazorpayPaymentID string    `json:"razorpay_payment_id"`
	Amount            int       `json:"amount"`
	PaymentStatus     bool      `json:"payment_status" gorm:"default:false"`
	DateTime          time.Time `json:"date_time"`
	UserID            uint      `json:"user_id" gorm:"not null"`
	TestID            uint      `json:"test_id" gorm:"not null"`
	Creator           bool      `json:"creator" gorm:"default:false"`
	// Foreign keys
	User User `gorm:"foreignKey:UserID"`
}

// TestAssignedToUser model
type TestAssignedToUser struct {
	SomethingID      uint  `json:"something_id" gorm:"primaryKey"`
	TestID           uint  `json:"test_id" gorm:"not null"`
	UserID           uint  `json:"user_id" gorm:"not null"`
	OrderID          uint  `json:"order_id" gorm:"not null"`
	AttemptRemaining uint8 `json:"attempt_remaining"`
}
