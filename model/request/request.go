package request

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ProfileRequest struct {
	Token string `json:"token"`
}
