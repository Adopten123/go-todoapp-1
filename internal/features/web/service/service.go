package web_service

import "github.com/Adopten123/go-todoapp-1/internal/core/domain"

type WebService struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(filePath string) (domain.File, error)
}

func NewWebService(
	webRepository WebRepository,
) *WebService {
	return &WebService{
		webRepository: webRepository,
	}
}
