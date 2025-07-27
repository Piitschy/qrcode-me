# qrcode-me

## Development

```bash
air
```

## Usage

```docker-compose
services:
  node:
    image: ghcr.com/piitschy/qrcode-me:latest
    volumes:
      - database:/pb/pb_data
    environment:
      - CORS_ENABLED=true
      - CORS_ORIGIN=*
      - SUPERUSER_EMAIL=${SUPERUSER_EMAIL}
      - SUPERUSER_PASSWORD=${SUPERUSER_PASSWORD}
      - CACHE_TTL=60s


volumes:
  database:
```
