package service

import (
	"errors"
	"file-sharing/jwt"
	errModel "file-sharing/model/error"
	"file-sharing/model/request"
	"file-sharing/model/response"
	"file-sharing/repository"
	"file-sharing/util"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type AuthService interface {
	Register(req request.RegisterRequest) (response.RegisterResponse, error)
	LoginUser(user request.LoginRequest) (response.LoginResponse, error)
	GetProfile(id int) (response.ProfileResponse, error)
}

type authService struct {
	repo repository.AuthRepo
}

func NewAuthService(repo repository.AuthRepo) AuthService {

	return authService{repo: repo}
}

func (a authService) Register(req request.RegisterRequest) (response.RegisterResponse, error) {
	log.Info("ActionLog.service.Register().start")
	hashedPass := util.HashPassword(req.Username, req.Password)
	message, err := a.repo.Register(request.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPass,
	})
	if err != nil {
		log.Info("Error while registering user")
		return message, &errModel.UnexpectedError
	}
	log.Info("ActionLog.service.Register().end")
	return response.RegisterResponse{Message: "User Registered Successfully!"}, nil
}
func (a authService) LoginUser(user request.LoginRequest) (response.LoginResponse, error) {
	var res response.LoginResponse
	resUser, err := a.repo.Login(user)
	if err != nil {
		switch {
		case errors.Is(err, &errModel.UsernameOrPasswordIsWrong):
			return response.LoginResponse{}, &errModel.UsernameOrPasswordIsWrong
		default:
			return response.LoginResponse{}, fmt.Errorf("getUser error: %w", err)
		}
	}
	userId := resUser.ID
	token, err := jwt.Create(userId)
	if err != nil {
		return response.LoginResponse{}, err
	}
	res.Token = token

	return res, nil

}

func (a authService) GetProfile(id int) (response.ProfileResponse, error) {
	profile, err := a.repo.GetProfile(id)
	if err != nil {
		log.Info("GetCategories.error: %v", err)

		switch {
		case errors.Is(err, &errModel.ProfileNotFound):
			return response.ProfileResponse{}, &errModel.ProfileNotFound
		case errors.Is(err, &errModel.UnexpectedError):
			return response.ProfileResponse{}, &errModel.UnexpectedError
		default:
			return response.ProfileResponse{}, &errModel.NotFoundError

		}
	}
	return profile, nil
}
