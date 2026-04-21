package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Adopten123/go-todoapp-1/internal/core/errors"
	"github.com/google/uuid"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.UUID{}, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf(
			"path value='%s' by key='%s' not a valid uuid: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"no keys='%s' in path values",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path value '%s' by key='%s' not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}
