package handlers

import (
	"encoding/json"
	"go-mandelbrot/pkg/messages"
	"go-mandelbrot/pkg/service"
	"image/png"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func GetMandelbrotImageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrate conn error: %w", err)
		return
	}
	defer conn.Close()

	for {
		_, mes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read message error: %w", err)
			return
		}

		var wsParams messages.WsParams
		err = json.Unmarshal(mes, &wsParams)
		if err != nil {
			log.Printf("unmarshal error: %w", err)
			return
		}

		img := service.GenerateMandelbrotImage(wsParams.PointX, wsParams.PointY, wsParams.Zoom, wsParams.MaxIters, wsParams.ResolutionWidth, wsParams.ResolutionHeight)

		wc, err := conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Printf("get next writer error: %w", err)
			return
		}
		err = png.Encode(wc, img)
		if err != nil {
			return
		}
		wc.Close()
	}

}

func RegenerateColorRegions(w http.ResponseWriter, r *http.Request) {
	service.RegenerateRegions()

	w.WriteHeader(http.StatusOK)
}
