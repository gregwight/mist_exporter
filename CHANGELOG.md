# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

## [1.0.0] - 2025-08-12

This is the initial public release of the Mist Exporter.

### Added
- Initial implementation of the Prometheus exporter for Juniper Mist.
- Hybrid data collection model using both REST API and WebSockets.
- Real-time streaming of device (AP) and client statistics for low-latency metrics.
- On-demand scraping of organization-level metrics (alarms, tickets, sites).
- Automatic discovery of Organization ID from the API key.
- Dynamic site management: automatically starts/stops metric collection for added/removed sites.
- Configuration via YAML file and environment variables.
- Graceful shutdown handling for robust operation.
- Dockerfile for containerized deployment and a `systemd` service example.
- GitHub Actions workflow for building binaries and creating releases with GoReleaser.