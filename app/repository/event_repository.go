package repository

import (
	"context"
	"errors"

	"eventhub/app/model"

	"github.com/jackc/pgx/v5"
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

func (r *eventRepository) GetAll(ctx context.Context) ([]model.Event, error) {
	query := `
		SELECT id, title, description, location, start_time, end_time, quota, is_active, created_at
		FROM events
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.Event, 0)
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Quota,
			&event.IsActive,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) GetByID(ctx context.Context, id int) (*model.Event, error) {
	query := `
		SELECT id, title, description, location, start_time, end_time, quota, is_active, created_at
		FROM events
		WHERE id = $1
	`
	var event model.Event
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.Quota,
		&event.IsActive,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &event, nil
}
