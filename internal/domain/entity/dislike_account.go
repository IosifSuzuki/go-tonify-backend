package entity

import "time"

type DislikeAccount struct {
	ID         int64
	DislikerID int64
	DislikedID int64
	CreatedAt  *time.Time
}
