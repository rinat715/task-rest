package rest

type UserForm struct {
	Email      string `json:"email" validate:"required,email"`
	Pass       string `json:"pass" validate:"required,eqcsfield=RepeatPass"`
	RepeatPass string `json:"repeat_pass" validate:"required"`
	IsAdmin    bool   `json:"is_admin" validate:"required"`
}

type TaskForm struct {
	Text string     `json:"text" validate:"required"`
	Date string     `json:"date" validate:"required"`
	Tags []*TagForm `json:"tags"`
	Done bool       `json:"done"`
}

type TagForm struct {
	Text string `json:"text" validate:"required"`
}
