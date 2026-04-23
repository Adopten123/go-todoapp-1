package web_fs_repository

import (
	"errors"
	"fmt"
	"os"

	"github.com/Adopten123/go-todoapp-1/internal/core/domain"
	core_errors "github.com/Adopten123/go-todoapp-1/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) (domain.File, error) {
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.File{}, fmt.Errorf(
				"file: %s: %w",
				filePath,
				core_errors.ErrNotFound,
			)
		}

		return domain.File{}, fmt.Errorf(
			"get file: %s: %w",
			filePath,
			err,
		)
	}

	file := domain.NewFile(buffer)

	return file, nil
}
