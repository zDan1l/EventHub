package repository

import (
	"context"

	"eventhub/app/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	GetAll(ctx context.Context) ([]model.Event, error)
	GetByID(ctx context.Context, id int) (*model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	Delete(ctx context.Context, id int) error
}

type eventRepository struct {
	db *pgxpool.Pool
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
	query := `
		INSERT INTO events (title, description, location, start_time, end_time, quota, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		event.Title,
		event.Description,
		event.Location,
		event.StartTime,
		event.EndTime,
		event.Quota,
		event.IsActive,
	).Scan(&event.ID, &event.CreatedAt)
}
