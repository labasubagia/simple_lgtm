package repository

import (
	"context"
	"fmt"
	"sync"

	"simple_lgtm/internal/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Repository interface {
	CreateData(ctx context.Context, id string, value string) error
	GetData(ctx context.Context, id string) (string, error)
	UpdateData(ctx context.Context, id string, newValue string) error
	DeleteData(ctx context.Context, id string) error
	ListAllData(ctx context.Context) ([]model.DataItem, error)
}

type inMemoryRepository struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		data: make(map[string]string),
	}
}

func (r *inMemoryRepository) CreateData(ctx context.Context, id string, value string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "CreateDataInRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id), attribute.String("data.value", value))

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; exists {
		return fmt.Errorf("data with ID %s already exists", id)
	}

	r.data[id] = value
	return nil
}

func (r *inMemoryRepository) GetData(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "GetDataFromRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id))

	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.data[id]
	if !ok {
		return "", fmt.Errorf("data with ID %s not found", id)
	}

	return value, nil
}

func (r *inMemoryRepository) UpdateData(ctx context.Context, id string, newValue string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "UpdateDataInRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id), attribute.String("data.newValue", newValue))

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("data with ID %s not found", id)
	}

	r.data[id] = newValue
	return nil
}

func (r *inMemoryRepository) DeleteData(ctx context.Context, id string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "DeleteDataInRepo")
	defer span.End()

	span.SetAttributes(attribute.String("data.id", id))

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("data with ID %s not found", id)
	}

	delete(r.data, id)
	return nil
}

func (r *inMemoryRepository) ListAllData(ctx context.Context) ([]model.DataItem, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "ListAllDataInRepo")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a slice of DataItem to return the data
	dataList := make([]model.DataItem, 0, len(r.data))
	for key, value := range r.data {
		dataList = append(dataList, model.DataItem{ID: key, Value: value})
	}

	return dataList, nil
}
