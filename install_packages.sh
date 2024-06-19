#!/bin/bash

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to print messages in green with a checkmark
print_green() {
  echo -e "\033[0;32m✓ $1\033[0m"
}

echo "Updating package lists..."
apt update -qq >/dev/null 2>&1

# Check if nginx is installed, install if not
if command_exists nginx; then
  print_green "nginx is already installed."
else
  echo "Installing nginx..."
  apt install -y -qq nginx >/dev/null 2>&1
fi

# Check if Certbot is installed, install if not
if command_exists certbot; then
  print_green "Certbot is already installed."
else
  echo "Installing Certbot and dependencies..."
  apt install -y -qq certbot python3-certbot-nginx >/dev/null 2>&1
fi

# Check if the nostr relay dependencies are installed, install if not
if dpkg -l | grep -qw build-essential && dpkg -l | grep -qw cmake && dpkg -l | grep -qw protobuf-compiler && dpkg -l | grep -qw pkg-config && dpkg -l | grep -qw libssl-dev && command_exists git; then
  print_green "All nostr relay dependencies are already installed."
else
  echo "Installing dependencies for nostr relay..."
  apt install -y -qq build-essential cmake protobuf-compiler pkg-config libssl-dev git >/dev/null 2>&1
fi

# Check if Rust is installed, install if not
if command_exists rustc; then
  print_green "Rust is already installed."
else
  echo "Installing Rust..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

  echo "Sourcing Rust environment..."
  source $HOME/.cargo/env
fi

