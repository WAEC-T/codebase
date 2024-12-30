# WAEC-T
This repository serves as the codebase for a study that examines the energy consumption of different web-frameworks implementing the same application called _Minitwit_. Measurements are done using the OTII-Arc-Pro hardware.

## Overview

This project provides tooling to conduct automated performance and energy consumption experiments across multiple web-application implementations. It includes capabilities for both sequential and "berries" (multi-client) testing scenarios.

### Folder structure

CODEBASE/
├── .github/
├── .pytest_cache/
├── aws-cdk/
├── c-sharp-razor/
├── database/
├── experiment-setup/
│   ├── bootstrap/
│   ├── execution/
│   └── loadtest_cloud_db/
├── requirements.txt
├── go-gin/
├── go-gorilla/
├── javascript-express/
├── python-flask/
├── ruby-sinatra/
├── rust-actix/
├── tests/

## Features

- Web-applications implementation support:
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
Configure the following addresses in the `execution.py`:
- Server URL (default: `http://10.7.7.144:5000`)
- Client URLs (default):
  - Client 1: `http://10.7.7.199:5001/trigger`
  - Client 2: `http://10.7.7.178:5001/trigger`
  - Client 3: `http://10.7.7.145:5001/trigger`

### Docker Compose Files
Place your Docker Compose files in the specified base location:
```python
BASE_COMPOSE_FILES_LOCATION = '/media/mmcblk0p2/setup/compose_files/'
