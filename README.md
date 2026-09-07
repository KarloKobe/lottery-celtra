# Celtra Lottery

Full-stack lottery assignment built with Vanilla JavaScript, Go and SQLite.

## Run with Docker

Build the Docker image:

```bash
docker build -t celtra-lottery .
```

Run the application with a persistent SQLite volume:

```bash
docker run -p 8080:8080 \
  -v celtra-lottery-data:/app/data \
  celtra-lottery
```

Open:

`http://localhost:8080`

## Run without Docker

Requires Go.

```bash
go mod download
go run .
```

Open:

`http://localhost:8080`

## Testing

The solution can be tested manually by:

1. Submitting a name and a number between 1 and 30.
2. Waiting for the 30-second raffle to finish.
3. Verifying that the winning number and winners are displayed.
4. Verifying that the last 5 raffle results remain available after restarting the application.
5. Opening multiple browser sessions and verifying that they share the same raffle state and results.

## Bonus

The bonus animation is located in the `bonus` directory.

Open it with:

```bash
open bonus/index.html
```

Hold the **Slow-motion** button to slow the animation down 3x.