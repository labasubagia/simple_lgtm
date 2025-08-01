package repository

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"simple_lgtm/internal/model"
	"simple_lgtm/pkg/errs"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Repository interface {
	CreateData(ctx context.Context, id string, value string) error
	GetData(ctx context.Context, id string) (string, error)
	UpdateData(ctx context.Context, id string, newValue string) error
	DeleteData(ctx context.Context, id string) error
	ListAllData(ctx context.Context) ([]model.DataItem, error)
}

type inMemoryRepository struct {
	data sync.Map
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		data: sync.Map{},
	}
}

func (r *inMemoryRepository) CreateData(ctx context.Context, id string, value string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "CreateDataInRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id), attribute.String("data.value", value))

	if _, exists := r.data.Load(id); exists {
		span.SetStatus(codes.Error, "data already exists")
		return errs.NewInvalidInput(fmt.Errorf("data with ID %s already exists", id))
	}

	r.data.Store(id, value)
	span.SetStatus(codes.Ok, "success")
	return nil
}

func (r *inMemoryRepository) GetData(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "GetDataFromRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id))

	value, ok := r.data.Load(id)
	if !ok {
		span.SetStatus(codes.Error, "data not found")
		return "", errs.NewNotFound(fmt.Errorf("data with ID %s not found", id))
	}

	span.SetStatus(codes.Ok, "success")
	return value.(string), nil
}

func (r *inMemoryRepository) UpdateData(ctx context.Context, id string, newValue string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "UpdateDataInRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id), attribute.String("data.newValue", newValue))

	if _, exists := r.data.Load(id); !exists {
		span.SetStatus(codes.Error, "data not found")
		return errs.NewInvalidInput(fmt.Errorf("data with ID %s not found", id))
	}

	r.data.Store(id, newValue)
	span.SetStatus(codes.Ok, "success")
	return nil
}

func (r *inMemoryRepository) DeleteData(ctx context.Context, id string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "DeleteDataInRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id))

	if _, exists := r.data.Load(id); !exists {
		span.SetStatus(codes.Error, "data not found")
		return errs.NewNotFound(fmt.Errorf("data with ID %s not found", id))
	}

	r.data.Delete(id)
	span.SetStatus(codes.Ok, "success")
	return nil
}

func (r *inMemoryRepository) ListAllData(ctx context.Context) ([]model.DataItem, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "ListAllDataInRepo")
	defer span.End()

	// TODO: maybe debug SQL query here
	slog.DebugContext(ctx, "Listing all data items")

	dataList := make([]model.DataItem, 0)
	r.data.Range(func(key, value any) bool {
		dataList = append(dataList, model.DataItem{ID: key.(string), Value: value.(string)})
		return true
	})

	span.SetStatus(codes.Ok, "success")
	return dataList, nil
}
