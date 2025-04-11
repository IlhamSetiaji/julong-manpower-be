package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IMajorUsecase interface {
	FindAll() (*[]response.FindAllMajorResponse, error)
	FindById(request request.FindByIdMajorRequest) (*response.FindAllMajorResponse, error)
	GetMajorsByEducationLevel(educationLevel string) (*[]response.FindAllMajorResponse, error)
	FindILikeMajor(major string) (*[]response.FindAllMajorResponse, error)
	FindILikeMajorAndEducationLevel(major string, educationLevel string) (*[]response.FindAllMajorResponse, error)
}

type MajorUsecase struct {
	Log  *logrus.Logger
	Repo repository.IMajorRepository
}

func NewMajorUsecase(log *logrus.Logger, repo repository.IMajorRepository) *MajorUsecase {
	return &MajorUsecase{Log: log, Repo: repo}
}

func (u *MajorUsecase) FindAll() (*[]response.FindAllMajorResponse, error) {
	majors, err := u.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	var majorResponses []response.FindAllMajorResponse
	for _, major := range *majors {
		majorResponses = append(majorResponses, response.FindAllMajorResponse{
			ID:             major.ID.String(),
			Major:          major.Major,
			EducationLevel: major.EducationLevel,
		})
	}

	return &majorResponses, nil
}

func (u *MajorUsecase) FindById(request request.FindByIdMajorRequest) (*response.FindAllMajorResponse, error) {
	id := request.ID
	major, err := u.Repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if major == nil {
		return nil, errors.New("Major not found")
	}

	majorResponse := response.FindAllMajorResponse{
		ID:             major.ID.String(),
		Major:          major.Major,
		EducationLevel: major.EducationLevel,
	}

	return &majorResponse, nil
}

func (u *MajorUsecase) GetMajorsByEducationLevel(educationLevel string) (*[]response.FindAllMajorResponse, error) {
	majors, err := u.Repo.GetMajorsByEducationLevel(entity.EducationLevelEnum(educationLevel))
	if err != nil {
		return nil, err
	}

	var majorResponses []response.FindAllMajorResponse
	for _, major := range *majors {
		majorResponses = append(majorResponses, response.FindAllMajorResponse{
			ID:             major.ID.String(),
			Major:          major.Major,
			EducationLevel: major.EducationLevel,
		})
	}

	return &majorResponses, nil
}

func (u *MajorUsecase) FindILikeMajor(major string) (*[]response.FindAllMajorResponse, error) {
	majors, err := u.Repo.FindILikeMajor(major)
	if err != nil {
		return nil, err
	}

	var majorResponses []response.FindAllMajorResponse
	for _, major := range *majors {
		majorResponses = append(majorResponses, response.FindAllMajorResponse{
			ID:             major.ID.String(),
			Major:          major.Major,
			EducationLevel: major.EducationLevel,
		})
	}

	return &majorResponses, nil
}

func (u *MajorUsecase) FindILikeMajorAndEducationLevel(major string, educationLevel string) (*[]response.FindAllMajorResponse, error) {
	majors, err := u.Repo.FindILikeMajorAndEducationLevel(major, educationLevel)
	if err != nil {
		return nil, err
	}

	var majorResponses []response.FindAllMajorResponse
	for _, major := range *majors {
		majorResponses = append(majorResponses, response.FindAllMajorResponse{
			ID:             major.ID.String(),
			Major:          major.Major,
			EducationLevel: major.EducationLevel,
		})
	}

	return &majorResponses, nil
}

func MajorUsecaseFactory(log *logrus.Logger) *MajorUsecase {
	repo := repository.MajorRepositoryFactory(log)
	return NewMajorUsecase(log, repo)
}
