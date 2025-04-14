# onecli

CLI tool for OneLogin

## Description

onecli is a command-line interface tool for interacting with OneLogin services. It provides a convenient way to manage and automate OneLogin-related tasks from the terminal.

## Prerequisites

Before using onecli, you need to set the following environment variables:

```bash
# Required environment variables
export ONELOGIN_CLIENT_ID="your_client_id"
export ONELOGIN_CLIENT_SECRET="your_client_secret"
export ONELOGIN_SUBDOMAIN="your_subdomain"
export ONELOGIN_TIMEOUT="timeout_in_seconds"  # e.g., "30"
```

## Features

- OneLogin API integration
- Command-line interface for OneLogin operations
- Version management
- YAML configuration support

## Installation

### From Source

```bash
git clone https://github.com/pepabo/onecli.git
cd onecli
go build
```

### Binary Release

Download the latest release from the [releases page](https://github.com/pepabo/onecli/releases).

## Usage

```bash
# Check version
onecli --version

# Run commands
onecli [command] [options]
```

## Development

### Requirements

- Go 1.24 or later

### Building

```bash
go build
```

## License

This project is licensed under the terms of the included LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
