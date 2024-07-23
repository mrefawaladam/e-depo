package user

import (
	"log"

	dto "e-depo/src/app/dto/user"

	repo "e-depo/src/infra/persistence/postgres/user"
)

type UserUCInterface interface {
	CreateUser(data *dto.CreateUserReqDTO) (*dto.RegisterRespDTO, error)
	Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error)
}

type userUseCase struct {
	Repo repo.UserRepository
}

func NewUserUseCase(repo repo.UserRepository) UserUCInterface {
	return &userUseCase{
		Repo: repo,
	}
}

func (uc *userUseCase) CreateUser(data *dto.CreateUserReqDTO) (*dto.RegisterRespDTO, error) {
	result, err := uc.Repo.StoreUser(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (uc *userUseCase) Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error) {

	result, err := uc.Repo.Login(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
