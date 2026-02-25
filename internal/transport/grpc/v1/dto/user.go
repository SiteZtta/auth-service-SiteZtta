package dto

type SignUpInput struct {
	UserName string `validate: "required"`      // UserName of the user to create
	Email    string `validate:"required,email"` // Email of the user to create
	Phone    string `validate:"required,e164"`  // Phone of the user to create
	Password string `validate:"required,min=8"` // Password of the user to create
}

type SignInInput struct {
	Login    string `validate: "required"` // login of the user to sign in
	Password string `validate:"required"`  // Password of the user to sign in
}

type TokenInput struct {
	Token string `validate: "required"` // Token of the user
}

type AuthInfo struct {
	UserId int64 // User ID of the user
	Role   int32 // Role of the user
}
