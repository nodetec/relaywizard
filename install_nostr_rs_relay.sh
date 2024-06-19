#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to install the nostr relay
install_nostr_rs_relay() {
  # Check if nostr-rs-relay is already installed
  if [ -f /usr/local/bin/nostr-rs-relay ]; then
    print_success "nostr-rs-relay is already installed."
    return
  fi

  # Clone the nostr-rs-relay repository and build the software
  print_info "Cloning the nostr-rs-relay repository..."
  cd /usr/local/src
  git clone https://git.sr.ht/~gheartsfield/nostr-rs-relay
  cd nostr-rs-relay

  print_info "Building the nostr-rs-relay software..."
  cargo build --release

  # Install the nostr-rs-relay executable
  print_info "Installing the nostr-rs-relay executable..."
  install target/release/nostr-rs-relay /usr/local/bin

  # Create the necessary directories for the nostr relay
  print_info "Creating necessary directories for nostr relay..."
  mkdir -p /var/lib/nostr-rs-relay/data

  print_success "nostr-rs-relay installation completed."
}
