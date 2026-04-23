package users_service

import (
	"context"
	"fmt"

	"github.com/Adopten123/go-todoapp-1/internal/core/domain"
	"github.com/google/uuid"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("update user in repository: %w", err)
	}

	return patchedUser, nil
}
