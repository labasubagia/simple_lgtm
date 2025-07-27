package service

import (
	"context"
	"fmt"
	"simple_lgtm/internal/model"
	"simple_lgtm/internal/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateData(ctx context.Context, id string, value string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "CreateDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id), attribute.String("service.value", value))

	err := s.repo.CreateData(ctx, id, value)
	if err != nil {
		return fmt.Errorf("failed to create data in repository: %w", err)
	}
	return nil
}

func (s *Service) GetData(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "GetDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id))

	data, err := s.repo.GetData(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get data from repository: %w", err)
	}
	return data, nil
}

func (s *Service) UpdateData(ctx context.Context, id string, newValue string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "UpdateDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id), attribute.String("service.newValue", newValue))

	err := s.repo.UpdateData(ctx, id, newValue)
	if err != nil {
		return fmt.Errorf("failed to update data in repository: %w", err)
	}
	return nil
}

func (s *Service) DeleteData(ctx context.Context, id string) error {
	_, span := otel.Tracer("app-tracer").Start(ctx, "DeleteDataService")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", id))

	err := s.repo.DeleteData(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete data from repository: %w", err)
	}
	return nil
}

func (s *Service) ListAllData(ctx context.Context) ([]model.DataItem, error) {
	_, span := otel.Tracer("app-tracer").Start(ctx, "ListAllDataService")
	defer span.End()

	data, err := s.repo.ListAllData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all data from repository: %w", err)
	}
	return data, nil
}
