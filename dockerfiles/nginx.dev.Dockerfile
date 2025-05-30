FROM node:18-alpine AS build

WORKDIR /app

# Copy package.json and package-lock.json
COPY frontend/package*.json ./
RUN npm ci

# Copy the rest of the frontend code
COPY frontend/ ./

# Build the SvelteKit app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy the built files to nginx html directory
COPY --from=build /app/build /usr/share/nginx/html

# Copy nginx configuration
COPY dockerfiles/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]