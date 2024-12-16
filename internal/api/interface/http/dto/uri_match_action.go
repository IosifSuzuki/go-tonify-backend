package dto

type URIMatchAction struct {
	Action MatchAction `uri:"action" binding:"required,enum_validate" example:"like"`
}
