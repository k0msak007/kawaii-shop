package middlewaresUsecases

import "github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecases interface {
	FindAccessToken(userId, access_token string) bool
}

type middlewaresUsecases struct {
	middlewaresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecases(middlewaresRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		middlewaresRepository: middlewaresRepository,
	}
}

func (u *middlewaresUsecases) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}
