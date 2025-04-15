package service

import (
	"github.com/CustomCloudStorage/repositories"
)

type file struct {
	repository *repositories.Repository
	storageDir string
}

type Service struct {
	File *file
}

func NewService(repository *repositories.Repository, storageDir string) *Service {
	return &Service{
		File: &file{
			repository: repository,
			storageDir: storageDir,
		},
	}
}
