# Salmon Ping
Online listing status checker by Salmon Fit

## Development
### Prerequisites
- [Go](go) 1.7+, but I use 1.21.3
- [Sqlx](https://docs.sqlc.dev/en/latest/overview/install.html) for development

### Setup
```sh
cp .env-template .env
# Then update variables accordingly
```

For the database I use https://neon.tech/

### Run
```sh
go run .

# Get History API used by salmonfit.com
curl -H "Origin: http://localhost:5173" http://localhost:8080/api/history

# Ping API
curl -H "X-API-Key: api_key" http://localhost:8080/api/ping
```

### Sqlc
1. Modify schema.sql or query.sql
2. Run `sqlc generate`
3. Files in `db` dir should be updated


### Telegram Alert (Optional)
1. Create a telegram bot: https://core.telegram.org/bots/tutorial to retrieve a bot token, store it as env "TELEGRAM_BOT_TOKEN"
2. Create a group chat
3. Invite your bot into the group chat
4. Run `curl https://api.telegram.org/bot[TELEGRAM_BOT_TOKEN]/getUpdates` to get the chat_id, store it as env "TELEGRAM_CHAT_ID". If it returns empty you can try removing the bot from group chat & add it again.

## Deployment

```sh
# Prepare secrets, this requires Secret Manager API to be enabled
echo -n "PUT_THE_API_KEY_HERE" | \
  gcloud secrets create "salmonping_API_KEY" --replication-policy "automatic" --data-file -
echo -n "PUT_THE_DATABASE_URL_HERE" | \
  gcloud secrets create "salmonping_DATABASE_URL" --replication-policy "automatic" --data-file -
# Optional
echo -n "PUT_TELEGRAM_BOT_TOKEN" | \
  gcloud secrets create "salmonping_TELEGRAM_BOT_TOKEN" --replication-policy "automatic" --data-file -
echo -n "PUT_TELEGRAM_CHAT_ID" | \
  gcloud secrets create "salmonping_TELEGRAM_CHAT_ID" --replication-policy "automatic" --data-file -
# https://developer.gobiz.com/docs/api/webhooks/receiving-notifications
echo -n "PUT_GOFOOD_NOTIFICATION_SECRET_KEY" | \
  gcloud secrets create "salmonping_GOFOOD_NOTIFICATION_SECRET_KEY" --replication-policy "automatic" --data-file -

# Create a service account and allow access to secret manager & cloud storage
gcloud iam service-accounts create SERVICE_ACCOUNT_NAME
gcloud projects add-iam-policy-binding SERVICE_ACCOUNT_NAME \
  --member='serviceAccount:SERVICE_ACCOUNT_NAME@PROJECT_NAME.iam.gserviceaccount.com' \
  --role='roles/secretmanager.secretAccessor' \
  --role='roles/storage.objectCreator'

# Deploy cloud run
gcloud run deploy salmonping --source . \
  --allow-unauthenticated \
  --concurrency 16 \
  --cpu 1 \
  --max-instances 5 \
  --memory 128Mi \
  --min-instances 0 \
  --region=asia-southeast2 \
  --service-account SERVICE_ACCOUNT_NAME@PROJECT_NAME.iam.gserviceaccount.com \
  --set-secrets API_KEY=salmonping_API_KEY:latest \
  --set-secrets DATABASE_URL=salmonping_DATABASE_URL:latest \
  # Uncomment if you choose to enable Gofood integration
  # --set-secrets GOFOOD_NOTIFICATION_SECRET_KEY=salmonping_GOFOOD_NOTIFICATION_SECRET_KEY:latest \
  # Uncomment if you choose to enable Telegram alert
  # --set-secrets TELEGRAM_BOT_TOKEN=salmonping_TELEGRAM_BOT_TOKEN:latest \
  # --set-secrets TELEGRAM_CHAT_ID=salmonping_TELEGRAM_CHAT_ID:latest \
  --tag=main \
  --timeout 15s \
  --update-labels service=salmonping

# Create a google cloud schedulers
gcloud scheduler jobs create http salmonping_ping --schedule="1,9,19,29,39,49,59 * * * *" --location="asia-southeast2" --time-zone="Asia/Jakarta" --uri=ENDPOINT_URL --http-method=GET --headers="X-API-Key=API_KEY"

# Create a cloud storage to debug html
gsutil mb -l asia-southeast2 gs://YOUR_BUCKET_NAME/
gsutil lifecycle set gcs-lifecycle-config.json gs://YOUR_BUCKET_NAME/
# Verify
gsutil lifecycle get gs://YOUR_BUCKET_NAME/
```
