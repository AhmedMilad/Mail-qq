# Mail-qq

A simple mail queue built with Go, Gin, RabbitMQ and PostgreSQL.

You send an email request to the API. The API puts it on a RabbitMQ queue and answers right away. A background worker then sends it through your SMTP server and saves a log of it in PostgreSQL.

1. A client calls `POST /api/mails/` with the email details.
2. The API checks the request, puts it on the RabbitMQ queue and replies right away.
3. The worker takes the message from the queue and sends it through the SMTP server.
4. The worker saves a log entry in PostgreSQL.
5. A client calls `GET /api/mails/` to read the email log.

## Requirements

- Go
- RabbitMQ
- PostgreSQL
- An SMTP account (for example Gmail with an app password)

## Setup

1. Start RabbitMQ and PostgreSQL (Docker examples):

   ```bash
   docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

   docker run -d --name mailqq-db -p 5432:5432 \
     -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=mailqq \
     postgres:16
   ```

2. Clone the project and create your `.env` file:

   ```bash
   git clone https://github.com/AhmedMilad/Mail-qq.git
   cd Mail-qq
   cp .env.example .env
   ```

3. Fill in `.env`:

   ```env
   MAILER_EMAIL=me@example.com
   PASSWORD=your-app-password
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587

   AMQP_URL=amqp://guest:guest@localhost:5672/

   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=mailqq
   DB_PORT=5432
   DB_SSLMODE=disable
   ```

   `DB_SSLMODE` is read by the app but is not in `.env.example` yet, so add it yourself.

4. Run it:

   ```bash
   go mod download
   go run .
   ```

   The server starts on `http://localhost:8080`. The `email_logs` table is created automatically.

## API

### Queue an email

`POST /api/mails/`

```bash
curl -X POST http://localhost:8080/api/mails/ \
  -H "Content-Type: application/json" \
  -d '{"to":"someone@example.com","subject":"Hello","body":"Hi from Mail-qq"}'
```

All three fields are required, and `to` must be a valid email address.

| Status | Meaning |
|---|---|
| `202` | Email queued (it is sent in the background) |
| `400` | Missing or invalid field |
| `503` | RabbitMQ could not be reached |

### List sent emails

`GET /api/mails/`

```bash
curl http://localhost:8080/api/mails/
```

Returns the sent emails, newest first:

```json
{
  "message": "Fetched email log successfully.",
  "data": [
    { "id": 1, "sender": "me@example.com", "recipient": "someone@example.com", "subject": "Hello", "body": "Hi from Mail-qq" }
  ]
}
```

| Status | Meaning |
|---|---|
| `200` | Log fetched |
| `500` | Database error |

## How it works

- The exchange `app.events` (topic, durable) routes messages with keys matching `email.#` to the durable queue `email_queue`.
- The worker acknowledges a message only after the email is sent.
- If sending fails, the message is retried once. If it fails again, it is dropped. The worker waits 5 seconds after each failure.
- After a successful send, a row is saved in the `email_logs` table.
- Delivery is at-least-once, so in rare cases (for example a crash right after sending) an email can be sent twice.

## Project structure

```
main.go                  Entry point
controllers/mails/       HTTP handlers (send and list)
routes/api/              Route registration
queues/rabbitmq.go       Publisher and email worker
database/db.go           PostgreSQL connection and auto-migration
models/email_log.go      EmailLog table model
utils/request.go         Request body and validation (SendEmailRequest)
utils/response.go        JSON responses (Response, EmailLogResponse)
```