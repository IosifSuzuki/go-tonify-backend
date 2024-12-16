package dto

type MatchAction string

const (
	LikeMatchAction    MatchAction = "like"
	DislikeMatchAction MatchAction = "dislike"
)

func (m MatchAction) Valid() bool {
	switch m {
	case LikeMatchAction, DislikeMatchAction:
		return true
	default:
		return false
	}
}
