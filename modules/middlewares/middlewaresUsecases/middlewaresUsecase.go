package middlewaresUsecases

import "github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecases interface {
}

type middlewaresUsecases struct {
	middlewaresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecases(middlewaresRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		middlewaresRepository: middlewaresRepository,
	}
}
