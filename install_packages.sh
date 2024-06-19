#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

print_info "Updating package lists silently..."
apt update -qq >/dev/null 2>&1

# Check if nginx is installed, install if not
if command_exists nginx; then
  print_success "nginx is already installed."
else
  print_info "Installing nginx..."
  apt install -y -qq nginx >/dev/null 2>&1
fi

# Check if Certbot is installed, install if not
if command_exists certbot; then
  print_success "Certbot is already installed."
else
  print_info "Installing Certbot and dependencies..."
  apt install -y -qq certbot python3-certbot-nginx >/dev/null 2>&1
fi

# Check if the nostr relay dependencies are installed, install if not
if dpkg -l | grep -qw build-essential && dpkg -l | grep -qw cmake && dpkg -l | grep -qw protobuf-compiler && dpkg -l | grep -qw pkg-config && dpkg -l | grep -qw libssl-dev && command_exists git; then
  print_success "All nostr relay dependencies are already installed."
else
  print_info "Installing dependencies for nostr relay..."
  apt install -y -qq build-essential cmake protobuf-compiler pkg-config libssl-dev git >/dev/null 2>&1
fi

# Check if Rust is installed, install if not
if command_exists rustc; then
  print_success "Rust is already installed."
else
  print_info "Installing Rust..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

  print_info "Sourcing Rust environment..."
  source $HOME/.cargo/env
fi

# Example of other formatting functions
print_bg_blue "This is a message with a blue background"
print_blink "This is a blinking message"
print_inverted "This is a message with inverted colors"
print_high_intensity_green "This is a high intensity green message"
print_dim "This is a dim text message"

