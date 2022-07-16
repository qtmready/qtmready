package routes

type RegisterResponse struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (r *RegisterResponse) Validate() error {
	if r.Password != r.ConfirmPassword {
		return ErrorPasswordMismatch
	}
	return nil
}

func (r *RegisterResponse) Save() error {
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}
