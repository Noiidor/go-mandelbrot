package messages

type WsParams struct {
	PointX           float64 `json:"pointx"`
	PointY           float64 `json:"pointY"`
	Zoom             uint64  `json:"zoom"`
	ResolutionWidth  uint32  `json:"resolutionWidth"`
	ResolutionHeight uint32  `json:"resolutionHeight"`
}
