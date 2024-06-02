// I dont know JS

// Remove hard-coded url
const ws = new WebSocket("ws://localhost:5050/v1/mandelbrot")
const colorsUrl = "http://localhost:5050/v1/mandelbrot/colors"

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

canvas.addEventListener("wheel", function(e) {
    const rect = canvas.getBoundingClientRect();
    // 0,0 = bottom left corner
    const px = e.clientX - rect.left;
    const py = rect.bottom - e.clientY;

    var dir = Math.sign(e.deltaY);

    if (dir > 0) {
        zoomInput.value /= 2
    } else{
        zoomInput.value *= 2
    }

    const x = transformPixelToCartesian(px, 500, -2, 2, parseInt(zoomInput.value), parseFloat(xInput.value));
    const y = transformPixelToCartesian(py, 500, -2, 2, parseInt(zoomInput.value), parseFloat(yInput.value));

    xInput.value = x;
    yInput.value = y;

    sendGenerateRequest();
})

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

        zoomSlider.max = zoomInput.value;
        zoomSlider.step = zoomInput.value/10
        zoomSlider.value = zoomInput.value
    });

    document.getElementById("requestColors").addEventListener("click", function() {
        const res = fetch(colorsUrl, {
            method: "PUT",
        });

        if (!res.ok) {
            console.log("Colors not OK: " + res.status);
        }
        sendGenerateRequest();
    })
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