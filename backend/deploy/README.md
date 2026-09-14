# VPS deployment

Single-host stack for sharing test builds: Caddy → API → Postgres. Needs Docker + Compose on the VPS, ports 12380/12443 open (or front it with Apache, see below).

```bash
# on the VPS
git clone <repo> && cd everyreview/backend/deploy
cp .env.example .env      # set DOMAIN and POSTGRES_PASSWORD
docker compose up -d --build
curl https://$DOMAIN/v1/products/0000000000000   # should return a placeholder product
```

The schema is applied automatically the first time Postgres starts (empty volume). Later migrations are not applied automatically — run them with `docker compose exec -T postgres psql -U everyreview -d everyreview < ../migrations/<file>`.

Redeploy after code changes: `git pull && docker compose up -d --build`. Uploaded product photos live in the `media-data` volume (`MEDIA_DIR=/data/media` inside the `api` container); the API builds their public URLs from the `Host`/`X-Forwarded-Proto` headers the proxy sends (the Caddyfile trusts those from private-range proxies such as Apache on the same host), or from `PUBLIC_BASE_URL` in `.env` if you set it. Check with `curl -s https://$DOMAIN/v1/products/<barcode>` that `image_url` starts with `https://`.

Migrations added so far, for an existing deployment (each is idempotent, safe to re-run):
```bash
docker compose exec -T postgres psql -U everyreview -d everyreview < ../migrations/0002_product_details.up.sql
```

Then build the APK against it:
```bash
cd android
./gradlew :app:assembleDebug -PAPI_BASE_URL=https://$DOMAIN/v1/
```

Without a domain (plain HTTP on the IP), the API is at `http://<vps-ip>:12380/v1/` — build the debug APK with `-PAPI_BASE_URL=http://<vps-ip>:12380/v1/`.

## Behind Apache on the same host

If Apache already owns ports 80/443, let **Apache terminate TLS** (with its existing certbot certificate) and reverse-proxy to Caddy on `127.0.0.1:12380`. Leave `DOMAIN=` empty in `.env` so Caddy serves plain HTTP internally and doesn't try to obtain its own certificate.

```bash
sudo a2enmod proxy proxy_http ssl headers
sudo certbot --apache -d reviews.example.com    # if you don't have a cert for this hostname yet
```

`/etc/apache2/sites-available/everyreview.conf`:
```apache
<VirtualHost *:80>
    ServerName reviews.example.com
    Redirect permanent / https://reviews.example.com/
</VirtualHost>

<VirtualHost *:443>
    ServerName reviews.example.com

    SSLEngine on
    SSLCertificateFile    /etc/letsencrypt/live/reviews.example.com/fullchain.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/reviews.example.com/privkey.pem

    ProxyPreserveHost On
    ProxyPass        / http://127.0.0.1:12380/
    ProxyPassReverse / http://127.0.0.1:12380/
    RequestHeader set X-Forwarded-Proto "https"

    ErrorLog  ${APACHE_LOG_DIR}/everyreview-error.log
    CustomLog ${APACHE_LOG_DIR}/everyreview-access.log combined
</VirtualHost>
```

```bash
sudo a2ensite everyreview && sudo apachectl configtest && sudo systemctl reload apache2
curl https://reviews.example.com/v1/products/0000000000000
```

Then build the APK with `-PAPI_BASE_URL=https://reviews.example.com/v1/` (HTTPS, so release builds work too). Since Apache fronts Caddy, you can also change the compose port mappings to `127.0.0.1:12380:80` so Caddy isn't reachable from outside directly.

### CentOS / RHEL differences

Service is `httpd`, config goes in `/etc/httpd/conf.d/everyreview.conf` (same vhost as above but log paths `/var/log/httpd/...`), mod_proxy is built in, and SELinux/firewalld must allow the proxy:

```bash
sudo dnf install -y mod_ssl epel-release && sudo dnf install -y certbot python3-certbot-apache
sudo setsebool -P httpd_can_network_connect 1
sudo firewall-cmd --permanent --add-service=http --add-service=https && sudo firewall-cmd --reload
sudo certbot --apache -d reviews.example.com
sudo apachectl configtest && sudo systemctl reload httpd
```
