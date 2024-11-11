#!/bin/bash

set -e

echo_info() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

echo_warning() {
    echo -e "\033[1;33m[WARNING]\033[0m $1"
}

echo_error() {
    echo -e "\033[1;31m[ERROR]\033[0m $1" >&2
}

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux*)     OS="linux";;
    darwin*)    OS="darwin";;
    *)          echo_error "Unsupported operating system: $OS"; exit 1;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH="amd64";;
    arm64|aarch64) ARCH="arm64";;
    *)             echo_error "Unsupported architecture: $ARCH"; exit 1;;
esac

REQUIRED_GO_VERSION="go1.23.1"

compare_go_version() {
    INSTALLED_VERSION=$(go version 2>/dev/null | awk '{print $3}')
    if [[ "$INSTALLED_VERSION" == "$REQUIRED_GO_VERSION" ]]; then
        return 0
    else
        return 1
    fi
}

if command -v go >/dev/null 2>&1; then
    if compare_go_version; then
        echo_info "Go $REQUIRED_GO_VERSION is already installed."
        USE_SYSTEM_GO=true
    else
        echo_warning "Installed Go version ($(go version | awk '{print $3}')) is not the required version ($REQUIRED_GO_VERSION)."
        USE_SYSTEM_GO=false
    fi
else
    echo_warning "Go is not installed on the system."
    USE_SYSTEM_GO=false
fi

if [ "$USE_SYSTEM_GO" = false ]; then
    if [[ "$OS" == "linux" ]]; then
        GO_URL="https://go.dev/dl/${REQUIRED_GO_VERSION}.linux-${ARCH}.tar.gz"
    elif [[ "$OS" == "darwin" ]]; then
        if [[ "$ARCH" == "amd64" ]]; then
            GO_URL="https://go.dev/dl/${REQUIRED_GO_VERSION}.darwin-amd64.pkg"
        else
            GO_URL="https://go.dev/dl/${REQUIRED_GO_VERSION}.darwin-arm64.pkg"
        fi
    else
        echo_error "Unsupported operating system for Go installation: $OS"
        exit 1
    fi

    GO_INSTALL_DIR="$HOME/go${REQUIRED_GO_VERSION}"
    mkdir -p "$GO_INSTALL_DIR"

    echo_info "Downloading Go $REQUIRED_GO_VERSION from $GO_URL..."
    cd "$HOME"

    if [[ "$OS" == "linux" ]]; then
        curl -L -o "go.tar.gz" "$GO_URL"
        tar -C "$GO_INSTALL_DIR" -xzf "go.tar.gz" --strip-components=1
        rm "go.tar.gz"
    elif [[ "$OS" == "darwin" ]]; then
        curl -L -o "go.pkg" "$GO_URL"
        pkgutil --expand-full "go.pkg" "go_pkg"
        mkdir -p "$GO_INSTALL_DIR"
        tar -C "$GO_INSTALL_DIR" -xzf "go_pkg/Payload~" --strip-components=3
        rm -rf "go.pkg" "go_pkg"
    fi

    export GOROOT="$GO_INSTALL_DIR"
    export PATH="$GOROOT/bin:$PATH"

    echo_info "Go $REQUIRED_GO_VERSION installed at $GOROOT"

    go version

else
    echo_info "Using Go $REQUIRED_GO_VERSION installed on the system."
fi

VENTANA_DIR="$HOME/ventana"
if [[ -d "$VENTANA_DIR" ]]; then
    echo_info "Directory $VENTANA_DIR already exists. Updating repository..."
    cd "$VENTANA_DIR"
    git pull
else
    echo_info "Cloning Ventana repository..."
    git clone https://github.com/pedrovian4/Ventana.git "$VENTANA_DIR"
fi

cd "$VENTANA_DIR"

echo_info "Building Ventana..."
go build -o ventana ./cmd/ventana

echo_info "Installing Ventana to /usr/local/bin..."
sudo mv ventana /usr/local/bin/

CONFIG_DIR="$HOME/.config/ventana"
mkdir -p "$CONFIG_DIR"

echo_info "Copying messages directory to $CONFIG_DIR..."
cp -r "$VENTANA_DIR/messages" "$CONFIG_DIR/"

if [ "$USE_SYSTEM_GO" = false ]; then
    echo_info "Cleaning up local Go installation..."
    rm -rf "$GO_INSTALL_DIR"
fi

echo_info "Installation complete! Check by running 'ventana'"

echo_info "Installation finished."
