package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adopten123/go-todoapp-1/internal/core/domain"
	core_errors "github.com/Adopten123/go-todoapp-1/internal/core/errors"
	core_postgres_pool "github.com/Adopten123/go-todoapp-1/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	id uuid.UUID,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%s': %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return taskDomain, nil
}
