#!/bin/bash
chmod +x scripts/database/postgres/*
docker compose -f docker-compose.yml up --build --force-recreate --no-deps