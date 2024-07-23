package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UserReqDTOInterface interface {
	Validate() error
}

type CreateUserReqDTO struct {
	UserName    string `json:"username"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

func (dto *CreateUserReqDTO) Validate() error {
	if err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.UserName, validation.Required),
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.PhoneNumber, validation.Required),
		validation.Field(&dto.Address, validation.Required),
		validation.Field(&dto.Password, validation.Required),
		validation.Field(&dto.Role, validation.Required),
	); err != nil {
		return err
	}
	return nil
}

type RegisterRespDTO struct {
	Token string `json:"token"`
}
type UserRespDTO struct {
	ID          int64  `json:"id"`
	UserName    string `json:"username"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}
type LoginReqDTO struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (dto *LoginReqDTO) Validate() error {
	if err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.UserName, validation.Required),
		validation.Field(&dto.Password, validation.Required),
	); err != nil {
		return err
	}
	return nil
}

type UserModel struct {
	ID          int64  `db:"id"`
	UserName    string `db:"username"`
	Name        string `db:"name"`
	PhoneNumber string `db:"phone_number"`
	Address     string `db:"address"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	Role        string `db:"role"`
}
