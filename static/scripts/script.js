// I dont know JS

// Remove hard-coded urls
const WS_URL = "ws://localhost:5050/v1/mandelbrot";
const COLORS_URL = "http://localhost:5050/v1/mandelbrot/colors";

let ws = new WebSocket(WS_URL);
// ws.binaryType = "arraybuffer";

const canvas = document.getElementById("mandelbrotCanvas");
const ctx = canvas.getContext("2d");

const zoomSlider = document.getElementById("sliderZoom");
const zoomInput = document.getElementById("inputZoom");
const xInput = document.getElementById("inputX");
const yInput = document.getElementById("inputY");
const iterInput = document.getElementById("inputIters");

const decItersButton = document.getElementById("decreaseIters");
const incItersButton = document.getElementById("increaseIters");

const leftButton = document.getElementById("leftButton");
const rightButton = document.getElementById("rightButton");
const downButton = document.getElementById("downButton");
const upButton = document.getElementById("upButton");

zoomSlider.addEventListener("input", event => {
    zoomInput.value = event.target.value;
});

zoomInput.addEventListener("input", event => {
    zoomSlider.value = event.target.value;
});

function clamp(num, min, max) {
    return Math.min(Math.max(num, min), max);
}

decItersButton.addEventListener("click", () => {
    iterInput.value = clamp(Math.floor(parseInt(iterInput.value) - (iterInput.value / 4)), 1, Number.POSITIVE_INFINITY);
    sendGenerateRequest();
})

incItersButton.addEventListener("click", () => {
    iterInput.value = clamp(Math.floor(parseInt(iterInput.value) + (iterInput.value / 4)), 1, Number.POSITIVE_INFINITY);
    sendGenerateRequest();
})
// Handle canvas zoom with wheel
canvas.addEventListener("wheel", event => {
    const rect = canvas.getBoundingClientRect();
    const px = event.clientX - rect.left;
    const py = rect.bottom - event.clientY;
    const direction = Math.sign(event.deltaY);

    zoomInput.value = clamp((direction > 0 ? zoomInput.value / 2 : zoomInput.value * 2), 1, Number.POSITIVE_INFINITY);

    const x = transformPixelToCartesian(px, parseInt(canvas.width), -2, 2, parseInt(zoomInput.value), parseFloat(xInput.value));
    const y = transformPixelToCartesian(py, parseInt(canvas.height), -2, 2, parseInt(zoomInput.value), parseFloat(yInput.value));

    xInput.value = x;
    yInput.value = y;

    sendGenerateRequest();
});

// TODO: change hardcoded params to variable
leftButton.addEventListener("click", () => {
    const x = transformPixelToCartesian(0, parseInt(canvas.width), -2, 2, parseInt(zoomInput.value), parseFloat(xInput.value));
    xInput.value = x;

    sendGenerateRequest();
})

rightButton.addEventListener("click", () => {
    const x = transformPixelToCartesian(500, parseInt(canvas.width), -2, 2, parseInt(zoomInput.value), parseFloat(xInput.value));
    xInput.value = x;

    sendGenerateRequest();
})

upButton.addEventListener("click", () => {
    const y = transformPixelToCartesian(500, parseInt(canvas.height), -2, 2, parseInt(zoomInput.value), parseFloat(yInput.value));
    yInput.value = y;

    sendGenerateRequest();
})

downButton.addEventListener("click", () => {
    const y = transformPixelToCartesian(0, parseInt(canvas.height), -2, 2, parseInt(zoomInput.value), parseFloat(yInput.value));
    yInput.value = y;

    sendGenerateRequest();
})

function transformPixelToCartesian(point, pixelBounds, axisMin, axisMax, zoom, offset) {
    const adjustedAxisMin = (axisMin / zoom) + offset;
    const adjustedAxisMax = (axisMax / zoom) + offset;
    const transformed = adjustedAxisMin + (point / (pixelBounds - 1)) * (adjustedAxisMax - adjustedAxisMin);

    console.log(`Transformed: ${transformed}`);
    return transformed;
}

ws.onopen = () => {
    console.log("Connected to WebSocket");

    document.getElementById("requestImage").addEventListener("click", () => {
        sendGenerateRequest();
        zoomSlider.max = zoomInput.value;
        zoomSlider.step = zoomInput.value / 10; // not working properly
        zoomSlider.value = zoomInput.value;
    });

    document.getElementById("requestColors").addEventListener("click", async () => {
        try {
            const res = await fetch(COLORS_URL, { method: "PUT" });
            if (!res.ok) {
                console.error(`Colors request failed: ${res.status}`);
                return;
            }
            console.log("Colors request succeeded");
            sendGenerateRequest();
        } catch (error) {
            console.error("Error requesting colors:", error);
        }
    });

    sendGenerateRequest();
};

function sendGenerateRequest() {
    const message = {
        pointX: parseFloat(xInput.value),
        pointY: parseFloat(yInput.value),
        zoom: parseInt(zoomInput.value),
        maxIters: parseInt(iterInput.value),
        resolutionWidth: parseInt(canvas.width),
        resolutionHeight: parseInt(canvas.height)
    };
    ws.send(JSON.stringify(message));
    console.log("Sent message");
}

ws.onmessage = (event) => {
    console.log("Received message");
    if (event.data instanceof Blob) {
        const url = URL.createObjectURL(event.data);
        console.log(`Created URL: ${url}`);

        const img = new Image();
        img.onload = () => {
            console.log("Image loaded");
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
        };
        img.src = url;
    } else {
        console.log("Expected binary data");
    }
};


ws.onerror = error => {
    console.error("WebSocket error:", error);
};

ws.onclose = () => {
    console.log("WebSocket closed");
};
