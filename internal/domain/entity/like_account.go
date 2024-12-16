package entity

import "time"

type LikeAccount struct {
	ID        int64
	LikerID   int64
	LikedID   int64
	CreatedAt *time.Time
}
