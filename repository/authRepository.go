package repository

import (
	"errors"
	"file-sharing/model/entity"
	errModel "file-sharing/model/error"
	"file-sharing/model/request"
	"file-sharing/model/response"
	"file-sharing/util"
	"gorm.io/gorm"
)

type AuthRepo interface {
	Register(req request.RegisterRequest) (response.RegisterResponse, error)
	Login(req request.LoginRequest) (*entity.User, error)
	GetProfile(id int) (response.ProfileResponse, error)
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{db: db}
}

func (u authRepo) Register(req request.RegisterRequest) (response.RegisterResponse, error) {
	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	err := u.db.Create(&user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return response.RegisterResponse{}, &errModel.UserAlreadyExists
		default:
			return response.RegisterResponse{}, &errModel.UnexpectedError

		}

	}
	return response.RegisterResponse{Message: "User registered successfully"}, nil

}
func (u authRepo) Login(user request.LoginRequest) (*entity.User, error) {

	var res entity.User

	err := u.db.Where("email = ?", user.Email).First(&res).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, &errModel.UsernameOrPasswordIsWrong
		}
	}

	if !util.CheckPasswordHash(res.Username, user.Password, res.Password) {
		return nil, &errModel.UsernameOrPasswordIsWrong
	}

	return &res, nil

}
func (u authRepo) GetProfile(id int) (response.ProfileResponse, error) {

	user := entity.User{}
	err := u.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ProfileResponse{}, &errModel.NotFoundError
		}
	}
	profile := response.ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return profile, nil
}
