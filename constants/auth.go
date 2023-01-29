package constants

import "time"

const (
	ACCESS_TOKEN_DURATION  time.Duration = time.Minute * time.Duration(15) // 15 minutes
	REFRESH_TOKEN_DURATION time.Duration = time.Minute * time.Duration(30) // 30 minutes
)
