#!/bin/bash

echo "Setting up SSL certificates using port 25565..."

EMAIL="andreirotaru366@gmail.com"
DOMAIN="serverpoli.go.ro"

mkdir -p certbot/conf certbot/www

echo "Step 1: Building containers..."
docker compose -f docker-compose.prod.yml build

echo "Step 2: Starting nginx on port 25565 for verification..."
docker compose -f docker-compose.prod.yml up -d

sleep 15

echo "Step 3: Getting certificate (port 25565 needs to be accessible as HTTP)..."
docker run -it --rm \
    -v $(pwd)/certbot/conf:/etc/letsencrypt \
    -v $(pwd)/certbot/www:/var/www/certbot \
    --network ourchat_ourchat-network \
    certbot/certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    -d $DOMAIN

if [ $? -eq 0 ]; then
    echo "Certificate obtained! Now switching to HTTPS..."

    # Oprește serviciile
    docker compose -f docker-compose.prod.yml down

    # Schimbă portul în docker-compose
    sed -i 's/25565:80/25565:443/' docker-compose.prod.yml

    # Pornește cu HTTPS
    docker compose -f docker-compose.prod.yml up -d

    echo "Done! Your site should be accessible at https://$DOMAIN:25565"
else
    echo "Certificate request failed!"
    exit 1
fi