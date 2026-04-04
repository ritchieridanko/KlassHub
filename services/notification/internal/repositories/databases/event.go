package databases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/models"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
)

type EventDatabase interface {
	Create(ctx context.Context, data *models.CreateEventData) (err *ce.Error)
	GetByID(ctx context.Context, eventID uuid.UUID) (evt *models.Event, err *ce.Error)
	SetCompleted(ctx context.Context, eventID uuid.UUID) (err *ce.Error)
	SetLastProcessed(ctx context.Context, eventID uuid.UUID) (err *ce.Error)
}

type eventDatabase struct {
	database *database.Database
}

func NewEventDatabase(db *database.Database) EventDatabase {
	return &eventDatabase{database: db}
}

func (d *eventDatabase) Create(ctx context.Context, data *models.CreateEventData) *ce.Error {
	query := `
		INSERT INTO events (id, topic, payload)
		VALUES ($1, $2, $3)
	`

	err := d.database.Execute(
		ctx, query,
		data.ID, data.Topic, data.Payload,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeDBQueryExec,
			fmt.Errorf("failed to create event: %w", err),
			logger.NewField("event_id", data.ID.String()),
		)
	}

	return nil
}

func (d *eventDatabase) GetByID(ctx context.Context, eventID uuid.UUID) (*models.Event, *ce.Error) {
	query := `
		SELECT
			id, topic, payload, first_processed_at,
			last_processed_at, completed_at
		FROM
			events
		WHERE
			id = $1
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var evt models.Event
	err := d.database.Query(
		ctx, query,
		eventID,
	).Scan(
		&evt.ID, &evt.Topic, &evt.Payload,
		&evt.FirstProcessedAt, &evt.LastProcessedAt,
		&evt.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, nil
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			fmt.Errorf("failed to get event by id: %w", err),
			logger.NewField("event_id", eventID.String()),
		)
	}

	return &evt, nil
}

func (d *eventDatabase) SetCompleted(ctx context.Context, eventID uuid.UUID) *ce.Error {
	query := `
		UPDATE events
		SET completed_at = NOW()
		WHERE id = $1 AND completed_at IS NULL
	`

	err := d.database.Execute(
		ctx, query,
		eventID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to set event completed: %w", err)
		eventIDField := logger.NewField("event_id", eventID.String())

		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(ce.CodeEventNotFound, wrappedErr, eventIDField)
		}
		return ce.NewError(ce.CodeDBQueryExec, wrappedErr, eventIDField)
	}

	return nil
}

func (d *eventDatabase) SetLastProcessed(ctx context.Context, eventID uuid.UUID) *ce.Error {
	query := `
		UPDATE events
		SET last_processed_at = NOW()
		WHERE id = $1 AND completed_at IS NULL
	`

	err := d.database.Execute(
		ctx, query,
		eventID,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to set event last process time: %w", err)
		eventIDField := logger.NewField("event_id", eventID.String())

		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(ce.CodeEventNotFound, wrappedErr, eventIDField)
		}
		return ce.NewError(ce.CodeDBQueryExec, wrappedErr, eventIDField)
	}

	return nil
}
