#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to configure the nostr relay
configure_nostr_rs_relay() {
  local domain_name="$1"
  local config_dir="/etc/nostr-rs-relay"
  local config_file="$config_dir/config.toml"

  # Check if the configuration file already exists
  if [ -f "$config_file" ]; then
    print_success "Configuration file already exists at $config_file."
    return
  fi

  # Prompt for configuration variables
  read -p "Enter the relay name: " relay_name
  read -p "Enter the relay description: " relay_description
  read -p "Enter the administrative contact pubkey (32-byte hex, not npub): " relay_pubkey

  # Create the configuration directory if it doesn't exist
  print_info "Creating configuration directory..."
  mkdir -p "$config_dir"

  # Create the config.toml file
  print_info "Creating config.toml file..."
  cat <<EOL >"$config_file"
[info]
# The advertised URL for the Nostr websocket.
relay_url = "wss://$domain_name"

# Relay information for clients. Put your unique server name here.
name = "$relay_name"

# Description
description = "$relay_description"

# Administrative contact pubkey (32-byte hex, not npub)
pubkey = "$relay_pubkey"

[database]
# Database engine to use.
engine = "sqlite"

# Directory to store the database in.
data_directory = "/var/lib/nostr-rs-relay/data"

[network]
# Bind to this network address
address = "127.0.0.1"

# Listen on this port
port = 8080
EOL

  print_success "Configuration of nostr-rs-relay completed."
}
