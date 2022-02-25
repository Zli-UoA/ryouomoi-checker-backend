package usecase

import "github.com/Zli-UoA/ryouomoi-checker-backend/repository"

type DeleteLovePointUseCase interface {
	Execute(userID, loverUserID int64) error
}

type deleteLovePointUseCaseImpl struct {
	ur repository.UserRepository
}

func (d *deleteLovePointUseCaseImpl) Execute(userID, loverUserID int64) error {
	_, err := d.ur.GetLovePoint(userID, loverUserID)
	if err != nil {
		return err
	}
	err = d.ur.DeleteLovePoint(userID, loverUserID)
	return err
}

func NewDeleteLovePointUseCase(ur repository.UserRepository) DeleteLovePointUseCase {
	return &deleteLovePointUseCaseImpl{ur: ur}
}
