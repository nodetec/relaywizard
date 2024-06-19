#!/bin/bash

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

echo "Updating package lists..."
apt update

# Check if nginx is installed, install if not
if command_exists nginx; then
  echo "nginx is already installed."
else
  echo "Installing nginx..."
  apt install -y nginx
fi

# Check if Certbot is installed, install if not
if command_exists certbot; then
  echo "Certbot is already installed."
else
  echo "Installing Certbot and dependencies..."
  apt install -y certbot python3-certbot-nginx
fi

# Check if the nostr relay dependencies are installed, install if not
if dpkg -l | grep -qw build-essential && dpkg -l | grep -qw cmake && dpkg -l | grep -qw protobuf-compiler && dpkg -l | grep -qw pkg-config && dpkg -l | grep -qw libssl-dev && command_exists git; then
  echo "All nostr relay dependencies are already installed."
else
  echo "Installing dependencies for nostr relay..."
  apt install -y build-essential cmake protobuf-compiler pkg-config libssl-dev git
fi

# Check if Rust is installed, install if not
if command_exists rustc; then
  echo "Rust is already installed."
else
  echo "Installing Rust..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

  echo "Sourcing Rust environment..."
  source "$HOME/.cargo/env"
fi

