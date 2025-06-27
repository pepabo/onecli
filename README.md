# onecli

A CLI tool for interacting with OneLogin API

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

### User Management

```bash
# List all users
onecli user list

# List users with filters
onecli user list --email user@example.com
onecli user list --username username
onecli user list --firstname John
onecli user list --lastname Doe
onecli user list --user-id 123

# Add a new user
onecli user add "John" "Doe" "john.doe@example.com"

# Modify user email
onecli user modify email "newemail@example.com" --email "oldemail@example.com"
```

### App Management

```bash
# List all apps
onecli app list

# List apps with user details
onecli app list --details
```

### Event Management

```bash
# List all events
onecli event list

# List events with filters
onecli event list --event-type-id 1
onecli event list --user-id 123
onecli event list --since 2023-01-01
onecli event list --until 2023-12-31

# List all event types
onecli event types

# List event types in JSON format
onecli event types --output json
```

## Output Formats

All list commands support multiple output formats:

- `yaml` (default)
- `json`

Example:
```bash
onecli user list --output json
```

## Configuration

Set the following environment variables:

- `ONELOGIN_CLIENT_ID`: Your OneLogin client ID
- `ONELOGIN_CLIENT_SECRET`: Your OneLogin client secret
- `ONELOGIN_SUBDOMAIN`: Your OneLogin subdomain

## Development

### Requirements

- Go 1.24 or later

### Building

```bash
go build
```

### Running Tests

```bash
# Run tests
go test ./...

# Run specific tests
go test ./onelogin -v -run TestGetEventTypes
```

## License

This project is licensed under the terms of the included LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
