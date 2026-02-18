# API

Personal API for my projects and services. Built with Go and a Node.js for the Open Graph . It provides small utility endpoints and integrations used across my projects.

## Features

Currently available features/endpoints:

- Open Graph image generation (OG images).
- Spotify endpoints (currently-playing, top tracks, etc.).
- Jellyfin endpoints (media status).
- Tools endpoints (headers, IP info, etc.).
- Real-time updates for Spotify and Jellyfin via Server-Sent Events (SSE).

## Development

Steps to run the project locally:

1. Clone the repository

   ```bash
   git clone https://github.com/ccrsxx/api
   ```

1. Change directory to the project

   ```bash
   cd api
   ```

1. Install dependencies

   ```bash
   go mod download
   ```

1. Install Tools

   ```bash
   make setup-tools
   ```

1. Set up environment variables
   Create a copy of the `.env.example` file and name it `.env`. Fill in credentials as needed.

   ```bash
   cp .env.example .env
   ```

1. Run the app in development

   ```bash
   make dev
   ```
