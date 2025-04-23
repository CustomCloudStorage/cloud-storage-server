package service

import (
	"github.com/CustomCloudStorage/repositories"
)

type file struct {
	repository *repositories.Repository
	storageDir string
}

type multiPart struct {
	repository *repositories.Repository
	storageDir string
	tmpUpload  string
}

type Service struct {
	File      *file
	MultiPart *multiPart
}

func NewService(repository *repositories.Repository, storageDir, tmpUpload string) *Service {
	return &Service{
		File: &file{
			repository: repository,
			storageDir: storageDir,
		},
		MultiPart: &multiPart{
			repository: repository,
			storageDir: storageDir,
			tmpUpload:  tmpUpload,
		},
	}
}
