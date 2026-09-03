let countdownInterval = null;

async function loadState() {
    try {
        const response = await fetch("/api/state");

        if (!response.ok) {
            throw new Error("Failed to load raffle state");
        }

        const state = await response.json();

        startCountdown(state.endsAt);
    } catch (error) {
        console.error(error);

        document.getElementById("countdown").textContent =
            "Unable to load raffle";
    }
}

async function loadResults() {
    try {
        const response = await fetch("/api/results");

        if (!response.ok) {
            throw new Error("Failed to load results");
        }

        const results = await response.json();

        renderResults(results);
    } catch (error) {
        console.error(error);
    }
}

function renderResults(results) {
    const resultsElement = document.getElementById("results");

    resultsElement.innerHTML = "<h2>Previous results</h2>";

    results.forEach(result => {
        const resultElement = document.createElement("div");

        let winnersText = "No winners";

        if (result.Winners && result.Winners.length > 0) {
            winnersText = result.Winners.join(", ");
        }

        resultElement.innerHTML = `
            <p>
                Winning number: <strong>${result.WinningNumber}</strong><br>
                Winners: ${winnersText}
            </p>
        `;

        resultsElement.appendChild(resultElement);
    });
}

function startCountdown(endsAt) {
    const countdownElement = document.getElementById("countdown");

    if (countdownInterval !== null) {
        clearInterval(countdownInterval);
    }

    function updateCountdown() {
        const now = new Date();
        const end = new Date(endsAt);

        const difference = end - now;

        if (difference <= 0) {
            countdownElement.textContent = "0";

            clearInterval(countdownInterval);
            countdownInterval = null;

            setTimeout(() => {
                loadState();
                loadResults();
            }, 1000);

            return;
        }

        const seconds = Math.ceil(difference / 1000);

        countdownElement.textContent = seconds;
    }

    updateCountdown();

    countdownInterval = setInterval(updateCountdown, 1000);
}

function setupEntryForm() {
    const form = document.getElementById("entry-form");
    const nameInput = document.getElementById("name");
    const guessInput = document.getElementById("guess");
    const messageElement = document.getElementById("form-message");

    form.addEventListener("submit", async (event) => {
        event.preventDefault();

        const name = nameInput.value.trim();
        const guess = Number(guessInput.value);

        messageElement.textContent = "";

        try {
            const response = await fetch("/api/entries", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    name: name,
                    guess: guess
                })
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                throw new Error(errorMessage);
            }

            messageElement.textContent = "Entry submitted!";

            nameInput.value = "";
            guessInput.value = "";
        } catch (error) {
            messageElement.textContent = error.message;
            console.error(error);
        }
    });
}

loadState();
loadResults();
setupEntryForm();

new LotteryWidget(
    document.getElementById("lottery-widget"),
    "/api"
);