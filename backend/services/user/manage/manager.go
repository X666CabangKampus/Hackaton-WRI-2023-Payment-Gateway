package manage

import (
	"backend-hacktober/services/gateway"
	userInterface "backend-hacktober/services/user/manage/interface"
	"backend-hacktober/services/user/model"
	"backend-hacktober/util"
	"context"
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	paymentManage "github.com/imrenagi/go-payment/manage"
)

type PayTuitionsRequest struct {
	PaymentType      payment.PaymentType       `json:"payment_type"`
	Qty              int                       `json:"qty"`
	CreditCardDetail *invoice.CreditCardDetail `json:"credit_card,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

type Manager struct {
	userRepository userInterface.UserRepository
}

func NewManager(userRepository userInterface.UserRepository) *Manager {
	return &Manager{userRepository}
}

func (m Manager) Login(loginReq LoginRequest) (*model.User, error) {
	user, err := m.userRepository.FindByUsername(loginReq.Username)
	if err != nil {
		return nil, err
	}

	if user.Password != util.HashPassword(loginReq.Password) {
		return nil, nil
	}

	return user, nil
}

func (m Manager) GetUsers(searched ...*model.User) ([]*model.User, error) {
	return m.userRepository.Get(searched...)
}

func (m Manager) GetUser(searchedUser *model.User) (*model.User, error) {
	users, err := m.userRepository.Get(searchedUser)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (m Manager) CreateUser(user *model.User) (*model.User, error) {
	user.Password = util.HashPassword(user.Password)
	return m.userRepository.Add(user)
}

func (m Manager) PayTuition(paymentSrv *gateway.Server, ctx context.Context, username string, req *PayTuitionsRequest) (*invoice.Invoice, error) {
	user, err := m.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	invoiceReq := paymentManage.GenerateInvoiceRequest{
		Payment: struct {
			PaymentType      payment.PaymentType       `json:"payment_type"`
			CreditCardDetail *invoice.CreditCardDetail `json:"credit_card,omitempty"`
		}{
			PaymentType:      req.PaymentType,
			CreditCardDetail: req.CreditCardDetail,
		},
		Customer: struct {
			Name        string `json:"name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
		}{
			Name:        user.FullName,
			Email:       user.Email,
			PhoneNumber: user.Phone,
		},
		Items: []struct {
			Name         string  `json:"name"`
			Category     string  `json:"category"`
			MerchantName string  `json:"merchant"`
			Description  string  `json:"description"`
			Qty          int     `json:"qty"`
			Price        float64 `json:"price"`
			Currency     string  `json:"currency"`
		}{
			{
				Name:     "UKT",
				Category: "Tuition",
				Qty:      req.Qty,
				Price:    float64(user.TuitionFeeBase),
			},
		},
	}

	inv, err := paymentSrv.CreateInvoice(ctx, &invoiceReq)
	if err != nil {
		return nil, err
	}

	numberOfSemester := len(user.Tuition)
	for i := 0; i < req.Qty; i++ {
		numberOfSemester++
		_, err := m.userRepository.AddTuition(&model.UserTuitionFee{
			UserId:        user.ID,
			SemesterPay:   model.Semester(numberOfSemester),
			InvoiceNumber: inv.Number,
		})
		if err != nil {
			return nil, err
		}
	}
	return inv, nil
}

func (m Manager) GetTuitionByUsername(username string) ([]*model.UserTuitionFee, error) {
	user, err := m.userRepository.Get(&model.User{Username: username})
	if err != nil {
		return nil, err
	}
	return m.userRepository.GetTuitions(&model.UserTuitionFee{UserId: user[0].ID})
}

func (m Manager) GetUserFromInvoiceNumber(invoiceNumber string) (*model.User, error) {
	tuition, err := m.userRepository.GetTuitions(&model.UserTuitionFee{InvoiceNumber: invoiceNumber})
	if err != nil {
		return nil, err
	}
	return m.userRepository.GetUserByTuition(tuition[0])
}
