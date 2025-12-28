#!/bin/sh
set -e

echo "Running migrations..."
./spotistats migrate

echo "Seeding database (skips if data exists)..."
./spotistats seed -f /data/tracks.csv

echo "Starting server..."
./spotistats serve
