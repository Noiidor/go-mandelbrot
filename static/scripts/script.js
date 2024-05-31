// I dont know JS

const ws = new WebSocket("ws://localhost:5050/v1/mandelbrot")

ws.binaryType = "arrayBuffer"

const canvas = document.getElementById("mandelbrotCanvas");
const ctx = canvas.getContext("2d");

const zoomSlider = document.getElementById("sliderZoom");
const zoomInput = document.getElementById("inputZoom");

zoomSlider.addEventListener("input", (event) => {
    zoomInput.value = event.target.value;
})

zoomInput.addEventListener("input", (event) => {
    zoomSlider.value = event.target.value;
})

ws.onopen = () => {
    console.log("connected to websocket");

    document.getElementById("requestImage").addEventListener("click", () => {
        const message = {
            pointX: parseFloat(document.getElementById("inputX").value),
            pointY: parseFloat(document.getElementById("inputY").value),
            zoom: parseInt(document.getElementById("inputZoom").value),
            maxIters: parseInt(document.getElementById("inputIters").value),
            resolutionWidth: 500,
            resolutionHeight: 500
        };
        ws.send(JSON.stringify(message));
        console.log("sended message")
    });
};

ws.onmessage = (event) => {
    console.log("received message");
    if (event.data instanceof Blob) {
        // const blob = new Blob([event.data], {type: "image/png"});
        const url = URL.createObjectURL(event.data);
        console.log("created url: " + url);

        const img = new Image();

        img.onload = function() {
            console.log("img loaded");
            ctx.clearRect(0, 0, canvas.width, canvas.height);
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