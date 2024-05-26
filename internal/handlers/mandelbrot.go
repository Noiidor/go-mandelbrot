package handlers

import (
	"encoding/json"
	"go-mandelbrot/pkg/messages"
	"go-mandelbrot/pkg/service"
	"image/png"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func GetMandelbrotImageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		_, mes, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var wsParams messages.WsParams
		err = json.Unmarshal(mes, &wsParams)
		if err != nil {
			return
		}

		img := service.GenerateMandelbrotImage(wsParams.PointX, wsParams.PointY, wsParams.Zoom, wsParams.ResolutionWidth, wsParams.ResolutionHeight)

		wc, err := conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			return
		}
		err = png.Encode(wc, img)
		if err != nil {
			return
		}
		wc.Close()
	}

}
