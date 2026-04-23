package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Adopten123/go-todoapp-1/internal/core/errors"
	"github.com/google/uuid"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%s': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
