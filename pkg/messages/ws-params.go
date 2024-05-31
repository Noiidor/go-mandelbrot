package messages

type WsParams struct {
	PointX           float64 `json:"pointx"`
	PointY           float64 `json:"pointY"`
	Zoom             uint64  `json:"zoom"`
	MaxIters         uint32  `json:"maxIters"`
	ResolutionWidth  uint32  `json:"resolutionWidth"`
	ResolutionHeight uint32  `json:"resolutionHeight"`
}
