package user

import "time"

type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type User struct {
	Model
	Username       string           `json:"username" gorm:"unique_index"`
	Password       string           `json:"password"`
	FullName       string           `json:"-"`
	TuitionFeeBase uint64           `json:"-"`
	Tuition        []UserTuitionFee `json:"tuition"`
}

type Semester uint

type UserTuitionFee struct {
	Model
	UserId        uint64   `json:"user_id" gorm:"index:tuition_user_id_k"`
	SemesterPay   Semester `json:"semester_pay"`
	InvoiceNumber string   `json:"invoice_number"`
}
