package _interface

import "backend-hacktober/services/user/model"

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	Get(searchedUsers ...*model.User) ([]*model.User, error)
	Add(user *model.User) (*model.User, error)
	AddTuition(userTuition *model.UserTuitionFee) (*model.UserTuitionFee, error)
	GetTuitions(searched ...*model.UserTuitionFee) ([]*model.UserTuitionFee, error)
	GetUserByTuition(tuition *model.UserTuitionFee) (*model.User, error)
	DeleteTuition(tuition *model.UserTuitionFee) error
}
