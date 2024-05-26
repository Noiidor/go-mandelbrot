const ws = new WebSocket("ws://localhost:5050/v1/mandelbrot")

ws.binaryType = "arrayBuffer"

const canvas = document.getElementById("mandelbrotCanvas");
const ctx = canvas.getContext("2d");

ws.onopen = () => {
    console.log("connected to websocket");

    document.getElementById("requestImage").addEventListener("click", () => {
        const message = {
            pointX: 0,
            pointY: 0,
            zoom: 10,
            resolutionWidth: 500,
            resolutionHeight: 500
        };
        ws.send(JSON.stringify(message));
        console.log("sended message")
    });
};

ws.onmessage = (event) => {
    console.log("received message");
    if (event.data instanceof ArrayBuffer) {
        const blob = new Blob([event.data], {type: "image/png"});
        const url = URL.createObjectURL(blob);
        console.log("created url: " + url);

        const img = new Image();

        img.onload = function() {
            console.log("img loaded");
            ctx.drawImage(img, 0, 0, 500, 500);
        }

        img.src = url

    }else{
        console.log("expected binary data");
    }
}

ws.onerror = (error) => {
    console.log("websocket error: ", error);
}

ws.onclose = () => {
    console.log("websocket closed");
}