package usecase

import (
	"github.com/Zli-UoA/ryouomoi-checker-backend/service"
)

type GetTwitterLoginUrlUseCase interface {
	Execute() (string, error)
}

type getTwitterLoginUrlUseCaseImpl struct {
	ts service.TwitterService
}

func (g *getTwitterLoginUrlUseCaseImpl) Execute() (string, error) {
	url, err := g.ts.GetLoginUrl()
	return url, err
}

func NewGetTwitterLoginUrlUseCase(ts service.TwitterService) GetTwitterLoginUrlUseCase {
	return &getTwitterLoginUrlUseCaseImpl{
		ts: ts,
	}
}
