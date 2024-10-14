package request

import (
	"errors"
	"fmt"
	"session_manager/internal/domain"
	"session_manager/internal/helper"
	"time"
)

type Campus struct {
	ID              string `json:"id"`
	ComputerName    string `json:"comp_name"`
	IPAddress       string `json:"ip_addr"`
	Login           string `json:"login"`
	NextPingSeconds int    `json:"next_ping_sec"`
	DateTime        string `json:"date_time"`
}

func (c *Campus) Validate() (*domain.Campus, error) {
	if c.ID == "" {
		return nil, errors.New("id is empty")
	}
	if c.ComputerName == "" {
		return nil, errors.New("comp_name is empty")
	}
	if c.Login == "" {
		return nil, errors.New("login is empty")
	}
	if c.NextPingSeconds <= 0 {
		return nil, errors.New("next ping duration less or eq 0")
	}
	dto := domain.Campus{
		ID:           c.ID,
		ComputerName: c.ComputerName,
		IPAddress:    c.IPAddress,
		Login:        c.Login,
		NextPing:     time.Duration(c.NextPingSeconds) * time.Second,
	}
	if c.DateTime == "" {
		dto.StartDateTime = time.Now()
	} else {
		t, err := helper.ParseDate(c.DateTime)
		if err != nil {
			return nil, err
		}
		dto.StartDateTime = t
	}
	dto.EndDateTime = dto.StartDateTime.Add(dto.NextPing)
	return &dto, nil
}

func (c *Campus) Print() string {
	return fmt.Sprintf(`Request data:
{
	'ID': %s
	'Computer': %s
	'IP': %s
	'Login': %s
	'NextPing': %d
	'Date': %s
}`, c.ID, c.ComputerName, c.IPAddress, c.Login, c.NextPingSeconds, c.DateTime)
}

type Platform struct {
	SessionID       string `json:"session_id"`
	SessionType     string `json:"session_type,omitempty"`
	Login           string `json:"login"`
	NextPingSeconds int    `json:"next_ping_sec"`
	DateTime        string `json:"date_time"`
}

func (p *Platform) Validate() (*domain.Platform, error) {
	if p.SessionID == "" {
		return nil, errors.New("session_id is empty")
	}
	if p.Login == "" {
		return nil, errors.New("login is empty")
	}
	if p.NextPingSeconds <= 0 {
		return nil, errors.New("next ping duration less or eq 0")
	}
	dto := domain.Platform{
		SessionID:   p.SessionID,
		SessionType: p.SessionType,
		Login:       p.Login,
	}
	if p.DateTime == "" {
		dto.StartDateTime = time.Now()
	} else {
		t, err := helper.ParseDate(p.DateTime)
		if err != nil {
			return nil, err
		}
		dto.StartDateTime = t
	}
	dto.EndDateTime = dto.StartDateTime.Add(time.Duration(p.NextPingSeconds) * time.Second)
	return &dto, nil
}

func (p *Platform) Print() string {
	return fmt.Sprintf(`Request data:
{
	'SessionID': %s
	'SessionType': %s
	'Login': %s
	'NextPing': %d
	'Date': %s
}`, p.SessionID, p.SessionType, p.Login, p.NextPingSeconds, p.DateTime)
}

type UserActivity struct {
	SessionType string `query:"session_type"` // parsing by link's queries ('omitempty' not working, do not add)
	Login       string `query:"login"`
	FromDate    string `query:"from_date"`
	ToDate      string `query:"to_date"`
	GroupBy     string `query:"group_by"`
}

const (
	GroupByMonth = "month"
	GroupByDate  = "date"
)

func (ua *UserActivity) Print() string {
	return fmt.Sprintf(`Request data:
{
	'SessionType': %s
	'Login': %s
	'FromDate': %s
	'ToDate': %s
	'GroupBy': %s
}`, ua.SessionType, ua.Login, ua.FromDate, ua.ToDate, ua.GroupBy)
}

func (ua *UserActivity) Validate() (*domain.UserActivity, error) {
	if ua.Login == "" {
		return nil, errors.New("login is empty")
	}
	if ua.GroupBy == "" {
		ua.GroupBy = GroupByDate
	}
	if ua.GroupBy != GroupByMonth && ua.GroupBy != GroupByDate {
		return nil, errors.New("group by must be 'month' or 'date'")
	}

	dto := domain.UserActivity{
		SessionType: ua.SessionType,
		Login:       ua.Login,
		GroupBy:     ua.GroupBy,
	}

	helper.ValidateFromTo(ua.FromDate, ua.ToDate)

	return &dto, nil
}

