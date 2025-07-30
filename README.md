# Mist Prometheus Exporter

[![Build Status](https://github.com/gregwight/mist_exporter/actions/workflows/release.yml/badge.svg)](https://github.com/gregwight/mist_exporter/actions/workflows/release.yml)

A Prometheus exporter for Juniper Mist API metrics.

This exporter collects metrics from the Mist API, focusing on the operational state of wireless devices (APs) and their connected clients. It is designed for enterprise organizations using Prometheus and Grafana to monitor their Juniper Mist environment.

## Features

-   Collects metrics for Mist organizations, sites, devices (APs), and clients.
-   Automatic discovery of Organization ID from the API key if not specified.
-   Concurrent, rate-limited collection across sites for performance and to respect API limits.
-   Configurable via YAML file or environment variables.
-   Graceful shutdown and robust server lifecycle management.

## Getting Started

### Prerequisites

-   A Juniper Mist API token with at least `Read-Only` access to your organization. You can generate one from the Mist Dashboard under `My Account > API Tokens`.
-   Go 1.21+ (for building from source).
-   Docker (for running as a container).

### Configuration

The exporter is configured via a `config.yaml` file.

```yaml
---
# Mist Prometheus Exporter Configuration

# The ID of the Mist Organization to scrape.
# If omitted, the exporter will attempt to discover it from the API key.
# This is required if the API key has access to multiple organizations.
# Can be set with MIST_ORG_ID environment variable.
org_id: ${MIST_ORG_ID}

mist_api:
  # Mist API base URL.
  base_url: "https://api.mist.com"

  # Mist API token for authentication.
  # Can be set with MIST_API_KEY environment variable.
  api_key: ${MIST_API_KEY}
  
  # Mist API client request timeout.
  timeout: 10s

exporter:
  # Exporter bind address.
  address: 0.0.0.0
  
  # Exporter port.
  port: 9200

collector:
  # Timeout for a full scrape cycle. This should be less than your
  # Prometheus scrape_timeout setting.
  timeout: 30s
```

### Running with Docker

A Docker image can be used to run the exporter.

1.  Create your `config.yaml`.
2.  Run the container:

```sh
docker run -d \
  --name mist_exporter \
  -p 9200:9200 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  <your-docker-image-name>:latest
```

### Building from Source

1.  Clone the repository:
    ```sh
    git clone https://github.com/gregwight/mist_exporter.git
    cd mist_exporter
    ```
2.  Build the binary:
    ```sh
    go build -o mist_exporter ./cmd/main.go
    ```
3.  Run the exporter:
    ```sh
    ./mist_exporter --config config.yaml
    ```

### Running from a Pre-built Binary

Pre-built binaries for major platforms (Linux, Windows, macOS) are available on the project's GitHub Releases page. The release archives contain the exporter binary and an example `config.yaml`.

This is the recommended way to deploy the exporter without setting up a Go development environment.

**Example for Linux:**

```sh
# Replace <version> with the latest release version (e.g., v1.0.0)
VERSION="<latest_version>"
ARCH="amd64" # or arm64

# Download and extract the release archive
curl -sSL https://github.com/gregwight/mist_exporter/releases/download/${VERSION}/mist_exporter-${VERSION}-linux-${ARCH}.tar.gz | tar -xz

# Copy the example config and edit it with your API key and org ID
cp config.yaml.dist config.yaml
vim config.yaml # Or your favorite editor

# Run the exporter
./mist_exporter --config config.yaml
```

### Running as a systemd Service (Linux)

To run the Mist exporter as a long-running service on a Linux system with `systemd`, you can create a unit file. This ensures the exporter starts on boot and is restarted if it fails.

1.  **Install the exporter:**
    Follow the steps in "Running from a Pre-built Binary" to download and extract the exporter. Then, move the files to standard locations.

    ```sh
    sudo mv mist_exporter /usr/local/bin/
    sudo mkdir -p /etc/mist_exporter
    sudo mv config.yaml /etc/mist_exporter/
    sudo chmod +x /usr/local/bin/mist_exporter
    ```

2.  **Create a dedicated user (Recommended):**
    It is a security best practice to run services as a non-root user.
    ```sh
    sudo useradd --no-create-home --shell /bin/false prometheus
    sudo chown -R prometheus:prometheus /etc/mist_exporter
    ```

3.  **Create the systemd unit file:**
    Create a file at `/etc/systemd/system/mist_exporter.service` with the following content:

    ```ini
    [Unit]
    Description=Mist Prometheus Exporter
    Wants=network-online.target
    After=network-online.target

    [Service]
    User=prometheus
    Group=prometheus
    Type=simple
    ExecStart=/usr/local/bin/mist_exporter --config /etc/mist_exporter/config.yaml
    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
    ```

4.  **Enable and start the service:**
    ```sh
    # Reload the systemd daemon to recognize the new service
    sudo systemctl daemon-reload

    # Enable the service to start on boot
    sudo systemctl enable mist_exporter.service

    # Start the service now
    sudo systemctl start mist_exporter.service

    # Check the status
    sudo systemctl status mist_exporter.service
    ```

## Exposed Metrics

The exporter exposes the following metrics at the `/metrics` endpoint.

### Organization Metrics (`mist_org_*`)
| Metric | Description | Labels |
|---|---|---|
| `mist_org_alarms` | Number of alarms in the organization. | `alarm_type` |
| `mist_org_tickets` | Number of tickets in the organization. | `ticket_status` |
| `mist_org_sites` | Number of sites in the organization. | `site_id`, `site_name`, `country_code` |

### Device (AP) Metrics (`mist_device_*`)
| Metric | Description | Labels |
|---|---|---|
| `mist_device_last_seen` | Last time the device was seen (Unix timestamp). | `device_id`, `device_name`, `device_type`, `device_model`, `device_hw_rev`, `site_id` |
| `mist_device_uptime` | Device uptime in seconds. | (as above) |
| `mist_device_wlans` | Number of WLANs assigned to the device. | (as above) |
| `mist_device_tx_bps` | Transmit rate in bits per second. | (as above) |
| `mist_device_rx_bps` | Receive rate in bits per second. | (as above) |
| `mist_device_clients` | Number of clients connected to the device radio. | (as above), `radio` |
| `mist_device_tx_bytes` | Transmitted bytes on the radio. | (as above), `radio` |
| `mist_device_rx_bytes` | Received bytes on the radio. | (as above), `radio` |
| `mist_device_power` | Transmit power (in dBm) of the radio. | (as above), `radio` |
| `mist_device_channel` | Current channel of the radio. | (as above), `radio` |
| ...more radio and utilization metrics. |

### Client Metrics (`mist_client_*`)
| Metric | Description | Labels |
|---|---|---|
| `mist_client_last_seen` | Last time the client was seen (Unix timestamp). | `client_mac`, `client_hostname`, `client_os`, `site_id`, `ap_id`, `ssid` |
| `mist_client_uptime` | Client uptime in seconds. | (as above) |
| `mist_client_rssi` | Client's Received Signal Strength Indication. | (as above) |
| `mist_client_snr` | Client's Signal-to-Noise Ratio. | (as above) |
| `mist_client_tx_bytes` | Transmitted bytes from the client. | (as above) |
| `mist_client_rx_bytes` | Received bytes by the client. | (as above) |
| ...more client statistics. |

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for details.

## License

This project is licensed under the Apache 2.0 License. See the LICENSE file for details.
