package model

type MatchAction string

const (
	LikeMatchAction    MatchAction = "like"
	DislikeMatchAction MatchAction = "dislike"
	UnknownMatchAction MatchAction = "unknown"
)

func MathActionFromString(text string) MatchAction {
	switch text {
	case string(LikeMatchAction):
		return LikeMatchAction
	case string(DislikeMatchAction):
		return DislikeMatchAction
	default:
		return UnknownMatchAction
	}
}
