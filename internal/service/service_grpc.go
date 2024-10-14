package service

import (
	"context"
	"errors"
	"session_manager/internal/domain"
	"time"

	"session_manager/internal/repository/pb/session_manager"
	"session_manager/internal/repository/postgres"
)

func NewServiceGrpc(storage postgres.Storage) *ServiceGrpc {
	return &ServiceGrpc{storage: storage}
}

type ServiceGrpc struct {
	storage postgres.Storage
	session_manager.UnimplementedSessionManagerServer
}

func (s *ServiceGrpc) GetCadetsTimeByModuleID(ctx context.Context, in *session_manager.CadetsTimeRequest) (*session_manager.CadetsTimeResponse, error) {
	// validate date
	if !in.FromDate.IsValid() {
		return nil, errors.New("field \"from_date\" is not valid")
	}

	toDate := time.Now().Add(24 * time.Hour)

	if in.ToDate.IsValid() {
		toDate = in.ToDate.AsTime()
	}

	req := domain.CadetTotalHoursRequest{
		ModuleID: int(in.ModuleId),
		FromDate: in.FromDate.AsTime(),
		ToDate:   toDate,
	}

	cadetsHours, err := s.storage.GetTotalHours(ctx, &req)
	if err != nil {
		return nil, err
	}

	if len(cadetsHours) == 0 {
		return &session_manager.CadetsTimeResponse{
			Message: "no data",
		}, nil
	}

	res := session_manager.CadetsTimeResponse{}

	res.Cadets = make([]*session_manager.Cadet, 0, len(cadetsHours))

	cadet := make(map[string]int)

	for i := 0; i < len(cadetsHours); i++ {
		index, ok := cadet[cadetsHours[i].Login]
		if !ok {
			// create cadet
			res.Cadets = append(res.Cadets, &session_manager.Cadet{
				Id:    int32(cadetsHours[i].ID),
				Total: cadetsHours[i].TotalHours,
				Month: make([]*session_manager.MonthNum, 0, 12),
			})
			cadet[cadetsHours[i].Login] = len(res.Cadets) - 1
			// update changes
			index = cadet[cadetsHours[i].Login]
		}
		// update cadet
		res.Cadets[index].Month = append(res.Cadets[index].Month, &session_manager.MonthNum{
			Year:  cadetsHours[i].Year,
			Month: session_manager.MonthNumber(cadetsHours[i].MonthNumber),
			Hours: cadetsHours[i].Hours,
		})

	}

	res.Message = "Success"

	return &res, nil
}
