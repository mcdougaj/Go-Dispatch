# API Key Setup Instructions

Since the Motive Interface directory is on your local machine, please follow these steps:

## Quick Setup

### Step 1: Locate Your Motive API Key

On your local computer, navigate to:
```
C:\Users\mcdou\Documents\Motive Interface\.env
```

Or on Mac:
```
/Users/mcdou/Documents/Motive Interface/.env
```

### Step 2: Open the .env File

Open the `.env` file in a text editor and look for your API key. It might look like:
```env
MOTIVE_API_KEY=sk_live_xxxxxxxxxxxxxxxxxxxx
```
or
```env
API_KEY=xxxxxxxxxxxxxxxxxxxx
```

### Step 3: Copy the Key

Copy just the key value (the part after the `=` sign)

### Step 4: Update This Project's .env File

In this Go-Dispatch project, create a `.env` file by copying from the example:
```bash
cp .env.example .env
```

Then edit `.env` and replace the placeholder with your actual keys:
```env
# Motive API Configuration
MOTIVE_API_KEY=paste_your_key_here
MOTIVE_BASE_URL=https://api.gomotive.com

# Google Maps API
GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here

# Claude API (Anthropic)
ANTHROPIC_API_KEY=your_anthropic_api_key_here

# Database
DATABASE_URL=postgres://dispatch_user:dispatch_pass@localhost:5432/motive_dispatch?sslmode=disable

# Redis Cache
REDIS_URL=redis://localhost:6379/0

# Server Configuration
PORT=8080
SYNC_INTERVAL_SECONDS=30
```

## Alternative: Paste Keys Here

If you prefer, you can paste the contents of your Motive Interface `.env` file here, and I'll extract the keys and set everything up for you.

Just copy and paste the file contents in the chat!
