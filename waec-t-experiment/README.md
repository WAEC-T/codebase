# WAEC-T
This repository serves as the codebase for a study that examines the energy consumption of different web-frameworks implementing the same application called _Minitwit_. Measurements are done using the OTII-Arc-Pro hardware.

## Overview

This project provides tooling to conduct automated performance and energy consumption experiments across multiple web-application implementations. It includes capabilities for both sequential and "berries" (multi-client) testing scenarios.

## Features

- Web-applications:
  - Rust-Actix
  - Python-Flask
  - Go-Gorilla
  - Go-Gin
  - Ruby-Sinatra
  - C-Sharp-Razor
- Automated power consumption measurements using OTII hardware
- Support for different testing scenarios:
  - Sequential API testing
  - Sequential UI testing
- Automated data collection and analysis from an AWS S3-bucket
- AWS CDK for RDS (Postgres) and S3-bucket
- Docker-based deployment
- Multi-client orchestration

## Prerequisites

- Python 3.x
- Docker and Docker Compose
- OTII hardware setup
- Network configuration for distributed testing
- Pandas library
- AsyncIO support

## Installation

1. Clone the repository
2. Install required Python dependencies
3. Configure the OTII hardware
4. Set up the network configuration for distributed testing

## Configuration

### Network Setup
Configure you own addresses with your own IP in the `execution.py`:

### Docker Compose Files
Place your Docker Compose files in the specified base location:
```python
BASE_COMPOSE_FILES_LOCATION = '/media/mmcblk0p2/setup/compose_files/'
```

# How to run on Raspberry Pis

Each MiniTwit Implementation will contain a readme specifying how to run the project on a Raspberry Pi with a specific operating system

## Index
- [Ruby Sinatra](./ruby-sinatra/README.md#how-to-run-on-raspberry-pi)
- [Rust Actix](./rust-actix/README.md#how-to-run-on-raspberry-pi)
