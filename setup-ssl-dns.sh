#!/bin/bash

echo "Setting up SSL certificates for serverpoli.go.ro using DNS challenge..."

EMAIL="andreirotaru366@gmail.com"  # SCHIMBĂ ASTA!
DOMAIN="serverpoli.go.ro"

mkdir -p certbot/conf certbot/www

echo "Getting SSL certificate with DNS challenge..."
echo "You will need to add a TXT record to your DNS!"

docker run -it --rm \
    -v $(pwd)/certbot/conf:/etc/letsencrypt \
    certbot/certbot certonly \
    --manual \
    --preferred-challenges dns \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    -d $DOMAIN

if [ $? -eq 0 ]; then
    echo "Certificate obtained successfully!"
    echo "Building and starting services..."
    docker compose -f docker-compose.prod.yml build
    docker compose -f docker-compose.prod.yml up -d
    echo "Done! Your site should be accessible at https://$DOMAIN:25565"
else
    echo "Certificate request failed!"
    exit 1
fi