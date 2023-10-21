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

## Deployment

```sh
# Prepare secrets, this requires Secret Manager API to be enabled
echo -n "PUT_THE_API_KEY_HERE" | \
  gcloud secrets create "salmonping_API_KEY" --replication-policy "automatic" --data-file -
echo -n "PUT_THE_DATABASE_URL_HERE" | \
  gcloud secrets create "salmonping_DATABASE_URL" --replication-policy "automatic" --data-file -

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
  --region=asia-southeast1 \
  --service-account SERVICE_ACCOUNT_NAME@PROJECT_NAME.iam.gserviceaccount.com \
  --set-secrets API_KEY=salmonping_API_KEY:latest \
  --set-secrets DATABASE_URL=salmonping_DATABASE_URL:latest \
  --tag=main \
  --timeout 15s \
  --update-labels service=salmonping

# Create a google cloud schedulers
API_KEY="PUT_THE_API_KEY_HERE"
ENDPOINT_URL="PUT_THE_ENDPOINT_URL_HERE"

# Times you specified
times=(05 30 55)

for hour in {8..21}; do
    for minute in "${times[@]}"; do
        job_name="salmonping_ping_${hour}${minute}"
        schedule="${minute} ${hour} * * 1-6"
        gcloud scheduler jobs create http $job_name --schedule="$schedule" --location="asia-southeast1" --time-zone="Asia/Jakarta" --uri=$ENDPOINT_URL --http-method=GET --headers=X-API-Key=$API_KEY
    done
done

# Create a cloud storage to debug html
gsutil mb -l ASIA-SOUTHEAST1 gs://YOUR_BUCKET_NAME/
gsutil lifecycle set gcs-lifecycle-config.json gs://YOUR_BUCKET_NAME/
# Verify
gsutil lifecycle get gs://YOUR_BUCKET_NAME/
```