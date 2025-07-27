package service

import (
	"context"
	"fmt"
	"log/slog"
	"simple_lgtm/internal/model"
	"simple_lgtm/internal/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Service interface {
	CreateData(ctx context.Context, id string, value string) error
	GetData(ctx context.Context, id string) (string, error)
	UpdateData(ctx context.Context, id string, newValue string) error
	DeleteData(ctx context.Context, id string) error
	ListAllData(ctx context.Context) ([]model.DataItem, error)
}

type appService struct {
	repo repository.Repository
}

func NewAppService(repo repository.Repository) Service {
	return &appService{
		repo: repo,
	}
}

func (s *appService) CreateData(ctx context.Context, id string, value string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "CreateDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id), attribute.String("service.value", value))

	err := s.repo.CreateData(ctx, id, value)
	if err != nil {
		return fmt.Errorf("failed to create data in repository: %w", err)
	}
	return nil
}

func (s *appService) GetData(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "GetDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id))

	data, err := s.repo.GetData(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get data from repository: %w", err)
	}
	return data, nil
}

func (s *appService) UpdateData(ctx context.Context, id string, newValue string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "UpdateDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id), attribute.String("service.newValue", newValue))

	err := s.repo.UpdateData(ctx, id, newValue)
	if err != nil {
		return fmt.Errorf("failed to update data in repository: %w", err)
	}
	return nil
}

func (s *appService) DeleteData(ctx context.Context, id string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "DeleteDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id))

	err := s.repo.DeleteData(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete data from repository: %w", err)
	}
	return nil
}

func (s *appService) ListAllData(ctx context.Context) ([]model.DataItem, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "ListAllDataService")
	defer span.End()

	// TODO: maybe debug SQL query here
	slog.DebugContext(ctx, "Listing all data items")

	data, err := s.repo.ListAllData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all data from repository: %w", err)
	}
	return data, nil
}
