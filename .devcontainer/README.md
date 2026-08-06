# CloudLens Dev Container

This directory contains the development container configuration for CloudLens.

## Quick Start

### Using VS Code

1. Install the "Dev Containers" extension in VS Code
2. Open the CloudLens repository in VS Code
3. Press `F1` and select "Dev Containers: Reopen in Container"
4. VS Code will build and open the dev container environment

### Features

- **Go 1.19** - Pre-installed Go environment
- **Docker Support** - Access to host Docker daemon for running LocalStack and other containers
- **VS Code Extensions** - Pre-configured extensions for Go development, AWS, and Terraform
- **Go Formatting** - Automatic code formatting on save
- **Port Forwarding** - LocalStack port (4566) is automatically forwarded
- **Docker Socket Binding** - Seamlessly uses the host's Docker daemon without nested containers

## Local Development in Container

Once the dev container is open:

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Setup LocalStack:**
   ```bash
   make setup
   ```

3. **Build and run CloudLens:**
   ```bash
   make run
   ```

4. **Access LocalStack:**
   The LocalStack service will be available at `localhost:4566`

## Stopping LocalStack

To clean up LocalStack containers:
```bash
make setup-down
```

## Customization

Edit `devcontainer.json` to:
- Add additional VS Code extensions
- Modify Go version
- Add environment variables
- Configure additional port forwarding

For more information, refer to the [Dev Containers documentation](https://containers.dev/).
