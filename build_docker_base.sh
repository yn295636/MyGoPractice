#!/usr/bin/env bash
docker build -f Dockerfile.gomod -t gomod:latest .
docker build -f Dockerfile.build -t build:latest .