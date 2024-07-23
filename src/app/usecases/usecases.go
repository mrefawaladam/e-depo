package usecases

import (
	userUC "e-depo/src/app/usecases/user"
)

type AllUseCases struct {
	UserUC userUC.UserUCInterface
}
