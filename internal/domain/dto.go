package domain

import (
	"fmt"
	"session_manager/internal/domain/response"
	"time"
)

type DataToJson interface {
	Key() string
	Marshal() []byte
}

type Campus struct {
	ID            string
	ComputerName  string
	IPAddress     string
	Login         string
	NextPing      time.Duration
	StartDateTime time.Time
	EndDateTime   time.Time
}

func (c Campus) Marshal() []byte {
	s := fmt.Sprintf(`{
		"id": "%s",
		"comp_name": "%s",
		"ip_addr": "%s",
		"login": "%s",
		"next_ping_sec": %v,
		"start_date_time": "%s",
		"end_date_time": "%s"
	}`,
		c.ID,
		c.ComputerName,
		c.IPAddress,
		c.Login,
		c.NextPing.Seconds(),
		c.StartDateTime.String(),
		c.EndDateTime.String(),
	)

	return []byte(s)
}

func (c Campus) Key() string {
	return fmt.Sprintf("%s / %s", c.Login, c.ComputerName)
}

type Platform struct {
	SessionID     string
	SessionType   string
	Login         string
	StartDateTime time.Time
	EndDateTime   time.Time
}

func (p Platform) Marshal() []byte {
	s := fmt.Sprintf(`{
		"session_id": "%s",
		"session_type": "%s",
		"login": "%s",
		"start_date_time": "%s",
		"end_date_time": "%s"
	}`,
		p.SessionID,
		p.SessionType,
		p.Login,
		p.StartDateTime.String(),
		p.EndDateTime.String(),
	)

	return []byte(s)
}

func (p Platform) Key() string {
	return p.Login
}

type UserActivity struct {
	SessionType string
	Login       string
	FromDate    time.Time
	ToDate      time.Time
	GroupBy     string
}

func (ua UserActivity) Marshal() []byte {
	s := fmt.Sprintf(`{
		"session_type": "%s",
		"login": "%s",
		"from_date": "%s",
		"to_date": "%s",
		"group_by": "%s"
	}`,
		ua.SessionType,
		ua.Login,
		ua.FromDate.String(),
		ua.ToDate.String(),
		ua.GroupBy,
	)

	return []byte(s)
}

func (ua UserActivity) Key() string {
	return ua.Login
}

type CadetTotalHoursRequest struct {
	ModuleID int
	FromDate time.Time
	ToDate   time.Time
}

type CadetTotalHoursResponse struct {
	ID         int
	Login      string
	TotalHours float32
	response.UserActivityByMonth
}
