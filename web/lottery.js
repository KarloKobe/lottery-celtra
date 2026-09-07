class LotteryWidget {
    constructor(rootElement, apiUrl) {
        this.root = rootElement;
        this.apiUrl = apiUrl;

        this.countdownInterval = null;

        this.countdownElement = this.root.querySelector(".countdown");
        this.resultsElement = this.root.querySelector(".results");
        this.form = this.root.querySelector(".entry-form");
        this.nameInput = this.root.querySelector(".name-input");
        this.guessInput = this.root.querySelector(".guess-input");
        this.messageElement = this.root.querySelector(".form-message");

        this.setupEntryForm();
        this.loadState();
        this.loadResults();
    }

    async loadState() {
        try {
            const response = await fetch(`${this.apiUrl}/state`);

            if (!response.ok) {
                throw new Error("Failed to load raffle state");
            }

            const state = await response.json();

            this.startCountdown(state.endsAt);
        } catch (error) {
            console.error(error);
            this.countdownElement.textContent = "Unable to load raffle";
        }
    }

    startCountdown(endsAt) {
        if (this.countdownInterval !== null) {
            clearInterval(this.countdownInterval);
        }

        const updateCountdown = () => {
            const now = new Date();
            const end = new Date(endsAt);

            const difference = end - now;

            if (difference <= 0) {
                this.countdownElement.textContent = "And the winner is...";

                clearInterval(this.countdownInterval);
                this.countdownInterval = null;

                setTimeout(() => {
                    this.loadState();
                    this.loadResults();
                }, 1000);

                return;
            }

            const seconds = Math.ceil(difference / 1000);

            this.countdownElement.textContent = `New winner in ${seconds}s...`;
        };

        updateCountdown();

        this.countdownInterval = setInterval(updateCountdown, 1000);
    }

    async loadResults() {
        try {
            const response = await fetch(`${this.apiUrl}/results`);

            if (!response.ok) {
                throw new Error("Failed to load results");
            }

            const results = await response.json();

            this.renderResults(results);
        } catch (error) {
            console.error(error);
        }
    }

renderResults(results) {
    this.resultsElement.innerHTML = "<h2>Previous results</h2>";

    results.forEach((result) => {
        const resultElement = document.createElement("div");
        resultElement.classList.add("result-row");

        let winnersText = "No lucky contestants";

        if (result.winners && result.winners.length > 0) {
            winnersText = result.winners.join(", ");
        }

        resultElement.innerHTML = `
            <span class="result-winners">${winnersText}</span>
            <span class="result-number">#${result.winningNumber}</span>
        `;

        this.resultsElement.appendChild(resultElement);
    });
}

    setupEntryForm() {
        this.form.addEventListener("submit", async (event) => {
            event.preventDefault();

            const name = this.nameInput.value.trim();
            const guess = Number(this.guessInput.value);

            this.messageElement.textContent = "";

            try {
                const response = await fetch(`${this.apiUrl}/entries`, {
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

                this.messageElement.textContent = "Entry submitted!";

                this.nameInput.value = "";
                this.guessInput.value = "";
            } catch (error) {
                this.messageElement.textContent = error.message;
                console.error(error);
            }
        });
    }
}

const widgetElements = document.querySelectorAll(".lottery-widget");

widgetElements.forEach((element) => {
    new LotteryWidget(element, "/api");
});