#!/bin/bash

# Deploy script
set -e

echo "🚀 Deploying Achievement System..."

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI not found. Installing..."
    npm install -g @railway/cli
fi

# Build the project
echo "📦 Building project..."
./build.sh

# Deploy to Railway
echo "🚂 Deploying to Railway..."
railway up

echo "✅ Deployment completed!"
echo "🌐 Your app should be live shortly!"
