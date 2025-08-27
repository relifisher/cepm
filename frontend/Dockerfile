# Stage 1: Build the React application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application source code
COPY . .

# Build the project
RUN npm run build

# Stage 2: Serve the application using Nginx
FROM nginx:stable-alpine

# Copy the built assets from the builder stage
# The output directory for Create React App is "build"
COPY --from=builder /app/build /usr/share/nginx/html

# Copy the custom Nginx configuration to handle SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
