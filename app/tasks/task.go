package tasks

import (
    "time"
)

type TaskParams = map[string]string

type SyncTask struct {
    Id string
    Key string 
    Type string
    Params TaskParams
    LastRun time.Time
    ErrorMsg string
    State string
    Enabled bool
}
