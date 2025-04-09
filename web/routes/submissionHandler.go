package routes

import (
    "time"
)

// we should make it compatible and same as pooya's
type Submission struct {
    ID             string
    ProblemID      string
    ProblemTitle   string 
    UserID         string
    OwnerUsername  string 
    Code           string
    Status         string // "Pending", "OK", "Compile Error", "Wrong Answer", "Memory Limit", "Time Limit", "Runtime Error"
    TimeUsed       int  
    SubmittedAt    time.Time
}