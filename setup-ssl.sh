#!/bin/bash

echo "Setting up SSL certificates for serverpoli.go.ro..."

# Creează directoarele necesare
mkdir -p certbot/conf certbot/www

# Înlocuiește email-ul cu al tău
EMAIL="your-email@example.com"  # SCHIMBĂ ASTA!
DOMAIN="serverpoli.go.ro"

echo "Step 1: Building containers..."
docker compose -f docker-compose.prod.yml build

echo "Step 2: Starting nginx for certificate verification..."
docker compose -f docker-compose.prod.yml up -d

echo "Step 3: Waiting for services to start..."
sleep 20

echo "Step 4: Testing if nginx is accessible..."
curl -f http://localhost:80 || echo "Warning: nginx might not be ready"

echo "Step 5: Requesting SSL certificate for $DOMAIN..."
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
    --dry-run \
    -d $DOMAIN

if [ $? -eq 0 ]; then
    echo "Dry run successful! Running actual certificate request..."
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
        echo "Certificate obtained successfully!"
        echo "Restarting nginx with SSL..."
        docker compose -f docker-compose.prod.yml restart nginx
        echo "Done! Your site should be accessible at https://$DOMAIN:25565"
    else
        echo "Certificate request failed!"
        exit 1
    fi
else
    echo "Dry run failed! Check your configuration."
    exit 1
fi