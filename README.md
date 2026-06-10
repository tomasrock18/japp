# JAPP

An application for tracking daily macronutrients.

**Tech Stack:** Go (Backend), Python (Bot), PostgreSQL, Docker.

## Features are:

- **Global Product Catalog**: Add products by barcode. Once added, a product becomes available to all users.
- **Food Logging**: Log meals by providing a barcode and weight.
- **Daily Statistics**: Track consumed macros and display the remaining allowance until the daily target is reached.

## Requirements

- [Docker](https://www.docker.com/get-started) & [Docker Compose](https://docs.docker.com/compose/install/)
- A Telegram Bot Token (obtain it from [@BotFather](https://t.me/BotFather))

## Installation

1. **Clone the repository**:
   ```bash
   git clone git@github.com:tomasrock18/japp.git
   cd japp
   ```

2. **Configure environment variables**:
   ```bash
   cp .env.example .env
   ```
   **Note:** Make sure you fill `TELEGRAM_BOT_TOKEN`.

3. **Start the containers**:
   ```bash
   docker-compose up --build
   ```
