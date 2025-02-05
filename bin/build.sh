#!/bin/bash

# Build the project
echo "Building the project..."
docker build -t api-server:latest -f Dockerfile.api-server .
docker build -t data-aggregator -f Dockerfile.data-aggregator .
docker build -t store -f Dockerfile.store .
docker build -t frontend -f Dockerfile.frontend .

