package model

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *ErrorResponse) Error() string {

	return fmt.Sprintf("Code :%d , Message:%s", err.Code, err.Message)
}

var (
	UnexpectedError           = ErrorResponse{Code: 500, Message: "Unexpected Error"}
	NotFoundError             = ErrorResponse{Code: 404, Message: "Resource Not Found"}
	InvalidRequestError       = ErrorResponse{Code: 422, Message: "Invalid Request"}
	UnauthorizedError         = ErrorResponse{Code: 401, Message: "Unauthorized"}
	EmptyListError            = ErrorResponse{Code: 403, Message: "Empty List Error"}
	UserAlreadyExists         = ErrorResponse{Code: 409, Message: "User already exists"}
	PasswordRequirementFailed = ErrorResponse{Code: 400, Message: "Password does not meet the required conditions"}
	UsernameOrPasswordIsWrong = ErrorResponse{Code: 401, Message: "Username or Password is wrong"}
	ProfileNotFound           = ErrorResponse{Code: 404, Message: "Profile Not Found"}
	InvalidJWTToken           = ErrorResponse{Code: 401, Message: "Invalid JWT Token"}
	ErrEmptyObjectName        = errors.New("object name cannot be empty")
	ErrRelativePath           = errors.New("file path must be absolute")
	ErrUploadFailed           = errors.New("file upload failed")
)
