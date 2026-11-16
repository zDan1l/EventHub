package repository

import (
	"context"

	"eventhub/app/model"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	GetAll(ctx context.Context) ([]model.Event, error)
	GetByID(ctx context.Context, id int) (*model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	Delete(ctx context.Context, id int) error
}
