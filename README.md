# academy_users_session
manage user's pc session at domain.

### install migrate
```bash
curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash

sudo apt-get update

sudo apt-get install migrate
```

### Add new migrations
```bash
migrate create -ext sql -dir db/migrations -seq table_name
```

### Run service (REQ_LOG - optionally)
```bash
make run 'DATABASE_URL=postgres://user:password@host:port/db-name?search_path=session&sslmode=disable' 'REQ_LOG=true'
```

### Run service locally
for testing and be able to connect to the local database
```bash
make run_local 'DATABASE_URL=postgres://user:password@host:port/db-name?search_path=session&sslmode=disable' 'REQ_LOG=true'
```

### Add new users & computers
service [env-tracker](https://github.com/academie-one/env-tracker)

### APIs
#### Add new session on campus
The computer notifies the running script about the start of a session during user authorization
```http
POST http://localhost:8080/api/session-manager/session/on-campus
Content-Type: application/json
{
  "id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
  "comp_name": "academie-mac-pink0001",
  "ip_addr": "192.168.1.100",
  "login": "user_1",
  "next_ping_sec": 60,
  "date_time": "2023-09-06T12:30:00Z" // current time from pc
}

```
#### Add new session on any platform (after creating a main session in campus)
just calculate the session time
```http
POST http://localhost:8080/api/session-manager/session/on-platform
Content-Type: application/json
{
  "session_id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
  "session_type": "", // platform name empty
  "login": "user_1",
  "next_ping_sec": 60,
  "date_time": "2023-09-06T15:30:00Z" // current time from pc
}
```
or if you want to calculate the session time for individual platform
```http
POST http://localhost:8080/api/session-manager/session/on-platform
Content-Type: application/json
{
  "session_id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
  "session_type": "platform zero", // platform name
  "login": "user_1",
  "next_ping_sec": 60,
  "date_time": "2023-09-06T15:30:00Z" // current time from pc
}
```
the last notification (session or activity) sent from the computer will mean the end of the session ("date_time" + "next_ping_sec").
#### Get online sessions
```http
GET http://localhost:8080/api/session-manager/dashboard
```
response:
```json
// Content-Type: application/json
{
  "message": "Success",
  "data": [
    {
      "id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
      "comp_name": "academie-mac-pink0001",
      "ip_addr": "192.168.1.100",
      "login": "user_1",
      "start_date_time": "2023-09-06T08:00:00Z",
      "end_date_time": "2023-09-06T09:30:00Z" // ends at: last ping + "next_ping_sec"
    },
    {
      "id": "6b8a4f1a-1d09-4d08-8a03-86be3e3b9104",
      "comp_name": "academie-mac-blue0002",
      "ip_addr": "192.168.1.101",
      "login": "user_2",
      "start_date_time": "2023-09-06T10:15:00Z",
      "end_date_time": "2023-09-06T11:45:00Z"
    },
    // ...
  ]
}
```
#### Get user activity
query param
- `session_type` - ***"your platform"*** or empty
- `login` - ***"user_1"***
- `from_date` - ***"2022-09-01T00:00:00Z"*** or ***2006-01-02***
- `to_date` - ***"2022-12-31T00:00:00Z"*** or ***2006-01-02*** or empty
- `group_by` - ***"month"*** or ***"date"***
```
GET http://localhost:8080/api/session-manager/activity?session_type=xxx&login=xxx&from_date=xxx&to_date=xxx&group_by=xxx
```
response - group by month:
```json
// Content-Type: application/json
{
  "message": "Success",
  "data": {
    "login": "user_1",
    "total_hours": 436.63,
    "user_activity": [
      {
        "year": "2023",
        "month_num": "9",
        "hours": 90.616667
      },
      {
        "year": "2023",
        "month_num": "5",
        "hours": 79.016666
      },
      // ...
    ]
  }
}
```
response - group by date:
```json
{
  "message": "Success",
  "data": {
    "login": "user_1",
    "total_hours": 37.62,
    "user_activity": [
      {
        "date": "2023-01-01T00:00:00Z",
        "hours": 13.266666
      },
      {
        "date": "2023-01-02T00:00:00Z",
        "hours": 7.35
      },
      // ...
    ]
  }
}
```

### GRPC
#### Get user hours by module id
```
grpcb http://localhost:9191/CadetsTime/GetCadetsTimeByModuleID
```
```json
{
  "module_id": 66,
  "from_date": {
    "seconds": 1704102537,
    "nanos": 0
  },
  "to_date": {
    "seconds": 1712742522,
    "nanos": 0
  }
}
```
response
```json
{
  "message": "Success",
  "cadets": [
    {
      "login": "user_1",
      "total_hours": 271.9202880859375,
      "month": [
        {
          "year": "2024",
          "month": "april",
          "hours": 21.100276947021484
        },
        {
          "year": "2024",
          "month": "march",
          "hours": 73.82694244384766
        },
        {
          "year": "2024",
          "month": "february",
          "hours": 93.89555358886719
        },
        {
          "year": "2024",
          "month": "january",
          "hours": 83.09750366210938
        }
      ]
    },
    {
      "login": "user_2",
      "total_hours": 416.7369384765625,
      "month": [
        {
          "year": "2024",
          "month": "april",
          "hours": 19.4777774810791
        },
        {
          "year": "2024",
          "month": "march",
          "hours": 41.81444549560547
        },
        {
          "year": "2024",
          "month": "february",
          "hours": 88.83139038085938
        },
        {
          "year": "2024",
          "month": "january",
          "hours": 266.61334228515625
        }
      ]
    },
    // ...
  ]
}
```