#!/bin/bash
# Copy API Keys from Motive Interface to Go-Dispatch

echo "=========================================="
echo "Motive Dispatch App - API Key Setup"
echo "=========================================="
echo ""
echo "This script will help you copy API keys from your Motive Interface folder"
echo ""

# Check if running on local machine
if [ -d "$HOME/Documents" ]; then
    echo "Searching for Motive Interface in Documents..."
    MOTIVE_DIR=$(find "$HOME/Documents" -type d -iname "*motive*interface*" 2>/dev/null | head -1)

    if [ -n "$MOTIVE_DIR" ]; then
        echo "Found Motive Interface directory: $MOTIVE_DIR"

        if [ -f "$MOTIVE_DIR/.env" ]; then
            echo "Found .env file!"
            echo ""
            echo "Extracting Motive API key..."

            # Extract Motive API key
            MOTIVE_KEY=$(grep -i "MOTIVE.*API.*KEY\|API.*KEY" "$MOTIVE_DIR/.env" | cut -d'=' -f2 | tr -d ' "'"'"'')

            if [ -n "$MOTIVE_KEY" ]; then
                echo "✓ Motive API key found"

                # Update .env file
                cp .env.example .env
                sed -i "s/your_motive_api_key_here/$MOTIVE_KEY/" .env

                echo "✓ Updated .env file with Motive API key"
                echo ""
                echo "Next steps:"
                echo "1. Add your Google Maps API key to .env"
                echo "2. Add your Anthropic API key to .env"
                echo "3. Run: make setup"
            else
                echo "✗ Could not extract API key from .env file"
            fi
        else
            echo "✗ No .env file found in $MOTIVE_DIR"
        fi
    else
        echo "✗ Could not find Motive Interface directory"
    fi
else
    echo "This script is meant to run on your local machine"
fi

echo ""
echo "=========================================="
echo "Manual Setup Instructions"
echo "=========================================="
echo ""
echo "1. Open: ~/Documents/[Motive Interface Folder]/.env"
echo "2. Copy the API key value"
echo "3. Edit: .env in this project"
echo "4. Replace 'your_motive_api_key_here' with your actual key"
echo ""
