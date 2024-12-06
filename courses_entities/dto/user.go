package dto

type UserInputDto struct {
	ID       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GetJWTInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type UserOutputDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserListOutputDto struct {
	Users []UserOutputDto `json:"users"`
}
