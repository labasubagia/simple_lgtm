# simple_lgtm

This project is a Podman Compose setup designed for learning and experimenting with the LGTM (Loki, Grafana, Tempo, Mimir) stack, with a focus on using Alloy for collecting logs, metrics, and traces.

## Project Components:

*   **Loki:** A horizontally scalable, highly available, multi-tenant log aggregation system.
*   **Grafana:** The open-source analytics and monitoring solution for visualizing data.
*   **Tempo:** A high-volume, minimal dependency distributed tracing backend.
*   **Mimir:** A horizontally scalable, highly available, multi-tenant, long-term storage for Prometheus metrics.
*   **Alloy:** A vendor-agnostic, declarative, and extensible agent for collecting and processing observability data (logs, metrics, traces).

## Purpose:

The primary goal of this project is to provide a self-contained environment to:

*   Understand the integration of Loki, Grafana, Tempo, and Mimir.
*   Learn how to configure and use Alloy for collecting various types of observability data.
*   Experiment with different data sources and visualization techniques within Grafana.
*   Explore distributed tracing with Tempo and metric storage with Mimir.

This setup is ideal for local development, testing, and educational purposes related to the LGTM stack and observability practices.

## Prerequisites:

Before you begin, ensure you have the following installed on your system:

*   **Podman:** A daemonless container engine for developing, managing, and running OCI containers on your Linux system.
*   **Podman Compose:** A tool to define and run multi-container applications with Podman.

## Getting Started:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-repo/simple_lgtm.git
    cd simple_lgtm
    ```
    (Note: Replace `https://github.com/your-repo/simple_lgtm.git` with the actual repository URL if it's hosted elsewhere.)

2.  **Start the services:**
    ```bash
    podman compose up -d
    ```
    This command will build the necessary images (if not already present) and start all the services defined in `compose.yaml` in detached mode.

## Usage:

Once the services are up and running, you can access the following:

*   **Grafana:** Typically available at `http://localhost:3000`. You can log in with default credentials (admin/admin) and will be prompted to change the password.
*   **Loki:** Accessible internally by Grafana for logs.
*   **Tempo:** Accessible internally by Grafana for traces.
*   **Mimir:** Accessible internally by Grafana for metrics.
*   **Alloy:** Running as a service, configured to collect data from the `app` service and send it to the respective LGTM components.

## Configuration:

*   **`compose.yaml`**: Defines the services, networks, and volumes for the Podman Compose project.
*   **`app/`**: Contains the sample application that generates logs, metrics, and traces.
*   **`infra/alloy/config.alloy`**: The configuration file for the Alloy agent, defining how it collects and forwards observability data.

Feel free to modify these configuration files to experiment with different setups and data collection strategies.
