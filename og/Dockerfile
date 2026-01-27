# Set the node version as a build argument
ARG NODE_VERSION

# Use the specified Node.js version as the base image
FROM node:${NODE_VERSION}-alpine

# Create and change to the app directory
WORKDIR /app

# Copy application dependency manifests to the container image
COPY package*.json ./

# Install dependencies
RUN npm ci --omit=dev

# Copy the source code
COPY src src

# Switch to a non-root user
USER node

# Run the web service on container startup
CMD ["npm", "run", "start"]

# Expose port for documentation, can be overriden if env variable PORT is set
EXPOSE 4000

# OCI-compliant labels
LABEL org.opencontainers.image.authors="ami@ccrsxx.com" \
    org.opencontainers.image.source="https://github.com/ccrsxx/api-go" \
    org.opencontainers.image.description="My personal API for my projects" \
    org.opencontainers.image.licenses="GPL-3.0"
