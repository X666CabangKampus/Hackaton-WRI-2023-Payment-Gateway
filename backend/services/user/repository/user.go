package repository

import (
	"backend-hacktober/services/user/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	req := u.DB.Where("username = ?", username).First(&user)
	if req.Error != nil {
		return nil, req.Error
	}
	return &user, nil
}

func (u *UserRepository) Get(searched ...*model.User) ([]*model.User, error) {
	var users []*model.User
	// If no parameter is given, return all users.
	if len(searched) == 0 {
		u.DB.Preload("Tuition").Find(&users)
		return users, nil
	}

	// Build the OR where clause.
	for _, searchedUser := range searched {
		var user []*model.User
		req := u.DB.Preload("Tuition")
		req = req.Find(&user, searchedUser)
		if req.Error != nil {
			return nil, req.Error
		}
		users = append(users, user...)
	}

	return users, nil
}

func (u UserRepository) GetUserByTuition(tuition *model.UserTuitionFee) (*model.User, error) {
	var user model.User
	req := u.DB.Preload("Tuition", tuition).First(&user)
	if req.Error != nil {
		return nil, req.Error
	}
	return &user, nil
}

func (u *UserRepository) Add(user *model.User) (*model.User, error) {
	req := u.DB.Create(user)
	if req.Error != nil {
		return nil, req.Error
	}
	return user, nil
}

func (u *UserRepository) AddTuition(userTuition *model.UserTuitionFee) (*model.UserTuitionFee, error) {
	req := u.DB.Create(userTuition)
	if req.Error != nil {
		return nil, req.Error
	}
	return userTuition, nil
}

func (u *UserRepository) GetTuitions(searched ...*model.UserTuitionFee) ([]*model.UserTuitionFee, error) {
	var tuitions []*model.UserTuitionFee
	// If no parameter is given, return all users.
	if len(searched) == 0 {
		u.DB.Find(&tuitions)
		return tuitions, nil
	}

	// Build the OR where clause.
	for _, searchedUserTuition := range searched {
		var tuition []*model.UserTuitionFee
		req := u.DB
		req = req.Find(&tuition, searchedUserTuition)
		if req.Error != nil {
			return nil, req.Error
		}
		tuitions = append(tuitions, tuition...)
	}

	return tuitions, nil
}

func (u *UserRepository) DeleteTuition(tuition *model.UserTuitionFee) error {
	req := u.DB.Delete(tuition)
	if req.Error != nil {
		return req.Error
	}
	return nil
}
