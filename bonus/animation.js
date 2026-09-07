const track = document.querySelector(".animation-track");
const rectangle = document.querySelector(".rectangle");
const slowMotionButton = document.querySelector(".slow-motion-button");

let speedMultiplier = 1;

let position = 0;
let direction = 1;

const speed = 150; // pixels per second

let previousTime = null;

function animate(timestamp) {
    if (previousTime === null) {
        previousTime = timestamp;
    }

    const deltaTime = (timestamp - previousTime) / 1000;
    previousTime = timestamp;

    const maxPosition = track.clientWidth - rectangle.clientWidth;

    position += direction * speed * deltaTime * speedMultiplier;

    if (position >= maxPosition) {
        position = maxPosition;
        direction = -1;
    }

    if (position <= 0) {
        position = 0;
        direction = 1;
    }

    rectangle.style.transform = `translateX(${position}px)`;

    requestAnimationFrame(animate);
}

slowMotionButton.addEventListener("mousedown", () => {
    speedMultiplier = 1 / 3;
});

slowMotionButton.addEventListener("mouseup", () => {
    speedMultiplier = 1;
});

slowMotionButton.addEventListener("mouseleave", () => {
    speedMultiplier = 1;
});

requestAnimationFrame(animate);