package model

import "time"

type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type User struct {
	Model
	Username       string           `json:"username" gorm:"unique"`
	Password       string           `json:"password"`
	FullName       string           `json:"full_name"`
	Email          string           `json:"email"`
	Phone          string           `json:"phone"`
	TuitionFeeBase uint64           `json:"tuition_fee_base"`
	Tuition        []UserTuitionFee `json:"tuitions" gorm:"foreignkey:UserId;references:ID"`
}

type Semester uint

type UserTuitionFee struct {
	Model
	UserId        uint64   `json:"user_id" gorm:"index"`
	SemesterPay   Semester `json:"semester_pay"`
	InvoiceNumber string   `json:"invoice_number"`
}
