package handler

import (
	"errors"
	"file-sharing/config"
	"file-sharing/db"
	error "file-sharing/model/error"
	"file-sharing/model/request"
	"file-sharing/repository"
	"file-sharing/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler() *AuthHandler {
	log.Println("New AuthHandler generated")
	r := repository.NewAuthRepo(db.GetDb())
	s := service.NewAuthService(r)
	app := &AuthHandler{
		authService: s,
	}
	return app
}
func (a AuthHandler) Login(c *gin.Context) {

	var user request.LoginRequest

	if err := c.BindJSON(&user); err != nil {
		c.JSON(error.InvalidRequestError.Code, error.InvalidRequestError)

	}
	loginUser, err := a.authService.LoginUser(user)
	if err != nil {

		switch {
		case errors.Is(err, &error.UsernameOrPasswordIsWrong):
			c.JSON(error.UsernameOrPasswordIsWrong.Code, error.UsernameOrPasswordIsWrong)
			return
		}
	}

	c.JSON(http.StatusOK, loginUser)

}

func (a AuthHandler) GetProfile(c *gin.Context) {

	log.Println("ActionLog.application.GetProfile.start")

	// TODO: check error handling
	ctxUserId, ok := c.Get(config.Conf.UserId)

	if !ok {
		c.JSON(error.InvalidJWTToken.Code, error.InvalidJWTToken)
		return
	}

	userId, err := strconv.Atoi(ctxUserId.(string))

	if err != nil {
		c.JSON(error.InvalidJWTToken.Code, error.InvalidJWTToken)
	}

	user, err := a.authService.GetProfile(userId)

	if err != nil {
		c.JSON(error.InvalidJWTToken.Code, gin.H{"message": error.InvalidJWTToken.Message})
	}

	c.JSON(http.StatusOK, user)

	log.Info("ActionLog.application.GetProfile.end")

	return
}
func (a AuthHandler) Register(c *gin.Context) {
	log.Info("ActionLog.application.register.start", c.Request.Method, c.Request.URL)
	var user request.RegisterRequest
	//validation start
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(error.InvalidRequestError.Code, error.InvalidRequestError)
		return
	}
	if len(user.Password) < 6 {
		c.JSON(error.PasswordRequirementFailed.Code, error.PasswordRequirementFailed)
		return
	}
	//validation end
	message, err := a.authService.Register(user)
	if err != nil {
		c.JSON(error.InvalidRequestError.Code, error.InvalidRequestError)
		return
	}
	c.JSON(http.StatusOK, message)
	log.Info("ActionLog.application.RegisterHandler.end")
}
