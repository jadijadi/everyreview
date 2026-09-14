# VPS deployment

Single-host stack for sharing test builds: Caddy → API → Postgres. Needs Docker + Compose on the VPS, ports 80/443 open.

```bash
# on the VPS
git clone <repo> && cd everyreview/backend/deploy
cp .env.example .env      # set DOMAIN and POSTGRES_PASSWORD
docker compose up -d --build
curl https://$DOMAIN/v1/products/0000000000000   # should return a placeholder product
```

The schema is applied automatically the first time Postgres starts (empty volume). Later migrations are not applied automatically — run them with `docker compose exec -T postgres psql -U everyreview -d everyreview < ../migrations/<file>`.

Redeploy after code changes: `git pull && docker compose up -d --build`.

Then build the APK against it:
```bash
cd android
./gradlew :app:assembleDebug -PAPI_BASE_URL=https://$DOMAIN/v1/
```
