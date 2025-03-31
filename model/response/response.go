package response

type RegisterResponse struct {
	Message string `json:"register_response"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
type ProfileResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type FileResponse struct {
	Message string `json:"message"`
}
