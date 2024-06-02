// I dont know JS

// Remove hard-coded urls
const WS_URL = "ws://localhost:5050/v1/mandelbrot";
const COLORS_URL = "http://localhost:5050/v1/mandelbrot/colors";

const ws = new WebSocket(WS_URL);
// ws.binaryType = "arraybuffer";

const canvas = document.getElementById("mandelbrotCanvas");
const ctx = canvas.getContext("2d");

const zoomSlider = document.getElementById("sliderZoom");
const zoomInput = document.getElementById("inputZoom");
const xInput = document.getElementById("inputX");
const yInput = document.getElementById("inputY");
const iterInput = document.getElementById("inputIters");

zoomSlider.addEventListener("input", event => {
    zoomInput.value = event.target.value;
});

zoomInput.addEventListener("input", event => {
    zoomSlider.value = event.target.value;
});

// Handle canvas zoom with wheel
canvas.addEventListener("wheel", event => {
    const rect = canvas.getBoundingClientRect();
    const px = event.clientX - rect.left;
    const py = rect.bottom - event.clientY;
    const direction = Math.sign(event.deltaY);

    zoomInput.value = direction > 0 ? zoomInput.value / 2 : zoomInput.value * 2;

    const x = transformPixelToCartesian(px, 500, -2, 2, parseInt(zoomInput.value), parseFloat(xInput.value));
    const y = transformPixelToCartesian(py, 500, -2, 2, parseInt(zoomInput.value), parseFloat(yInput.value));

    xInput.value = x;
    yInput.value = y;

    sendGenerateRequest();
});

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
        zoomSlider.step = zoomInput.value / 10;
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
            ctx.drawImage(img, 0, 0, 500, 500);
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
