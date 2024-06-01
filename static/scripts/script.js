// I dont know JS

const ws = new WebSocket("ws://localhost:5050/v1/mandelbrot")

ws.binaryType = "arrayBuffer"

const canvas = document.getElementById("mandelbrotCanvas");
const ctx = canvas.getContext("2d");

const zoomSlider = document.getElementById("sliderZoom");
const zoomInput = document.getElementById("inputZoom");
const xInput = document.getElementById("inputX")
const yInput = document.getElementById("inputY")
const iterInput = document.getElementById("inputIters")

zoomSlider.addEventListener("input", (event) => {
    zoomInput.value = event.target.value;
})

zoomInput.addEventListener("input", (event) => {
    zoomSlider.value = event.target.value;
})

canvas.addEventListener("click", function(e) {
    const rect = canvas.getBoundingClientRect();
    // 0,0 = bottom left corner
    const px = event.clientX - rect.left;
    const py = rect.bottom - event.clientY;
    console.log("X: "+px + " Y: "+py);

    const x = transformPixelToCartesian(px, 500, -2, 2, parseInt(zoomInput.value), parseFloat(xInput.value));
    const y = transformPixelToCartesian(py, 500, -2, 2, parseInt(zoomInput.value), parseFloat(yInput.value));

    xInput.value = x;
    yInput.value = y;

    sendGenerateRequest();
})

function getClickPosition(canvas, event) {
}

function transformPixelToCartesian(point, pixelBounds, axisMin, axisMax, zoom, offset) {
    axisMin = (axisMin / zoom) + offset;
    axisMax = (axisMax / zoom) + offset;

    const transformed = axisMin + (point/(pixelBounds-1))*(axisMax-axisMin);

    console.log("Transformed: " + transformed);

    return transformed;
}

ws.onopen = () => {
    console.log("connected to websocket");

    document.getElementById("requestImage").addEventListener("click", function() {
        sendGenerateRequest();
    });
};

function sendGenerateRequest() {
    const message = {
        pointX: parseFloat(xInput.value),
        pointY: parseFloat(yInput.value),
        zoom: parseInt(zoomInput.value),
        maxIters: parseInt(iterInput.value),
        resolutionWidth: 500,
        resolutionHeight: 500
    };
    ws.send(JSON.stringify(message));
    console.log("sended message");
}

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

        img.src = url;

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