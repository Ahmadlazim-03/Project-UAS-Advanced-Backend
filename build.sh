#!/bin/bash

# Build script for deployment
set -e

echo "🔨 Building Achievement System..."

# Build Frontend
echo "📦 Building frontend..."
cd frontend
npm install
npm run build
cd ..

echo "✅ Frontend build completed!"

# Build Backend
echo "🔧 Building backend..."
go build -o bin/achievement-server .

echo "✅ Backend build completed!"
echo "🎉 All builds completed successfully!"
