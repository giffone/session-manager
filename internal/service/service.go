package service

import (
	"context"
	"fmt"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/repository/postgres"
)

type Service interface {
	CreateSessionOnCampus(ctx context.Context, dto *domain.Campus) ([]response.Session, error)
	CreateSessionOnPlatform(ctx context.Context, dto *domain.Platform) error
	GetOnlineDashboard(ctx context.Context) ([]response.Session, error)
	GetUserActivity(ctx context.Context, dto *domain.UserActivity) (activity *response.UserActivity, err error)
}

func New(storage postgres.Storage) Service {
	return &service{storage: storage}
}

type service struct {
	storage postgres.Storage
}

func (s *service) CreateSessionOnCampus(ctx context.Context, dto *domain.Campus) ([]response.Session, error) {
	// first check if session already exists
	if sessions, err := s.storage.IsSessionExists(ctx, dto.Login); err != nil {
		return nil, fmt.Errorf("IsSessionExists: %w", err)
	} else if len(sessions) != 0 {
		// multi sessions
		if sessions[0].Multi {
			// create session
			return nil, s.storage.CreateSessionOnCampus(ctx, dto)
		}
		return sessions, response.ErrAccessDenied
	}

	// create session
	return nil, s.storage.CreateSessionOnCampus(ctx, dto)
}

func (s *service) CreateSessionOnPlatform(ctx context.Context, dto *domain.Platform) error {
	return s.storage.CreateSessionOnPlatform(ctx, dto)
}

func (s *service) GetOnlineDashboard(ctx context.Context) ([]response.Session, error) {
	return s.storage.GetOnlineDashboard(ctx)
}

func (s *service) GetUserActivity(ctx context.Context, dto *domain.UserActivity) (activity *response.UserActivity, err error) {
	if dto.GroupBy == request.GroupByMonth {
		if dto.SessionType == "" {
			// no need sort - get from main table on_campus
			return s.storage.GetUserActivityByMonthInCampus(ctx, dto)
		}
		// need sort
		return s.storage.GetUserActivityByMonth(ctx, dto)
	}

	if dto.SessionType == "" {
		// no need sort - get from main table on_campus
		return s.storage.GetUserActivityByDateInCampus(ctx, dto)
	}
	// need sort
	return s.storage.GetUserActivityByDate(ctx, dto)
}
