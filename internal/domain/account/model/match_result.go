package model

type MatchResult string

const (
	LikeMatchResult         MatchResult = "like"
	DislikeMatchResult      MatchResult = "dislike"
	MatchAccountMatchResult MatchResult = "match"
	ErrorMatchResult        MatchResult = "error"
)
