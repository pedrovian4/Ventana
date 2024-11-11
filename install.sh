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

MIN_GO_VERSION="1.22"

version_greater_equal() {
   local IFS=.
    local i ver1=($1) ver2=($2)
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++)); do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++)); do
        if [[ -z ${ver2[i]} ]]; then
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]})); then
            return 0
        elif ((10#${ver1[i]} < 10#${ver2[i]})); then
            return 1
        fi
    done
    return 0
}

compare_go_version() {
    INSTALLED_VERSION=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    if version_greater_equal "$INSTALLED_VERSION" "$MIN_GO_VERSION"; then
        return 0
    else
        return 1
    fi
}

if command -v go >/dev/null 2>&1; then
    if compare_go_version; then
        echo_info "A versão do Go ($(go version | awk '{print $3}')) é suficiente."
    else
        echo_error "A versão do Go instalada ($(go version | awk '{print $3}')) é mais antiga que a versão mínima requerida (go$MIN_GO_VERSION)."
        exit 1
    fi
else
    echo_error "O Go não está instalado no sistema."
    exit 1
fi

VENTANA_DIR="$HOME/ventana"
if [[ -d "$VENTANA_DIR" ]]; then
    echo_info "O diretório $VENTANA_DIR já existe. Atualizando o repositório..."
    cd "$VENTANA_DIR"
    git pull
else
    echo_info "Clonando o repositório Ventana..."
    git clone https://github.com/pedrovian4/Ventana.git "$VENTANA_DIR"
fi

cd "$VENTANA_DIR"

echo_info "Compilando o Ventana..."
go build -o ventana ./cmd/ventana

echo_info "Instalando o Ventana em /usr/local/bin..."
sudo mv ventana /usr/local/bin/

CONFIG_DIR="$HOME/.config/ventana"
mkdir -p "$CONFIG_DIR"

echo_info "Copiando o diretório de mensagens para $CONFIG_DIR..."
cp -r "$VENTANA_DIR/messages" "$CONFIG_DIR/"

echo_info "Instalação concluída! Verifique executando 'ventana'"

echo_info "Instalação finalizada."
