package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/notification/internal/models"
	"github.com/ritchieridanko/klasshub/services/notification/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
)

type EventRepository interface {
	Create(ctx context.Context, data *models.CreateEventData) (err *ce.Error)
	GetByID(ctx context.Context, eventID uuid.UUID) (evt *models.Event, err *ce.Error)
	SetCompleted(ctx context.Context, eventID uuid.UUID) (err *ce.Error)
	SetLastProcessed(ctx context.Context, eventID uuid.UUID) (err *ce.Error)
}

type eventRepository struct {
	database databases.EventDatabase
}

func NewEventRepository(db databases.EventDatabase) EventRepository {
	return &eventRepository{database: db}
}

func (r *eventRepository) Create(ctx context.Context, data *models.CreateEventData) *ce.Error {
	return r.database.Create(ctx, data)
}

func (r *eventRepository) GetByID(ctx context.Context, eventID uuid.UUID) (*models.Event, *ce.Error) {
	return r.database.GetByID(ctx, eventID)
}

func (r *eventRepository) SetCompleted(ctx context.Context, eventID uuid.UUID) *ce.Error {
	return r.database.SetCompleted(ctx, eventID)
}

func (r *eventRepository) SetLastProcessed(ctx context.Context, eventID uuid.UUID) *ce.Error {
	return r.database.SetLastProcessed(ctx, eventID)
}
