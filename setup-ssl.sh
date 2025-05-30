#!/bin/bash

echo "Setting up SSL certificates for serverpoli.go.ro..."

# Create directories
mkdir -p certbot/conf certbot/www

# Start nginx pentru initial setup
echo "Starting nginx for certificate verification..."
docker compose -f docker-compose.prod.yml up -d nginx

sleep 10

# Get SSL certificate cu certbot standalone (mai simplu)
echo "Getting SSL certificate..."
docker run -it --rm \
    -v $(pwd)/certbot/conf:/etc/letsencrypt \
    -v $(pwd)/certbot/www:/var/www/certbot \
    -p 80:80 \
    certbot/certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email andreirotaru366@gmail.com \
    --agree-tos \
    --no-eff-email \
    -d serverpoli.go.ro

# Restart cu SSL
echo "Restarting with SSL..."
docker compose -f docker-compose.prod.yml restart nginx

echo "Done! Access: https://serverpoli.go.ro:25565"