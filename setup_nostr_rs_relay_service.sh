#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to set up the nostr relay service
setup_nostr_rs_relay_service() {
  local service_file="/etc/systemd/system/nostr-relay.service"

  # Check if the service file already exists
  if [ -f "$service_file" ]; then
    print_success "Service file already exists at $service_file."
    return
  fi

  # Create a user for the nostr relay service
  print_info "Creating user for nostr relay service..."
  adduser --disabled-login --gecos "" nostr

  # Set ownership of the data directory
  print_info "Setting ownership of the data directory..."
  chown -R nostr:nostr /var/lib/nostr-rs-relay

  # Create the systemd service file
  print_info "Creating systemd service file..."
  cat <<EOL >"$service_file"
[Unit]
Description=Nostr Relay
After=network.target

[Service]
Type=simple
User=nostr
WorkingDirectory=/home/nostr
Environment=RUST_LOG=info,nostr_rs_relay=info
ExecStart=/usr/local/bin/nostr-rs-relay --config /etc/nostr-rs-relay/config.toml
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOL

  # Reload systemd to apply the new service
  print_info "Reloading systemd daemon..."
  systemctl daemon-reload

  # Enable and start the nostr relay service
  print_info "Enabling and starting nostr relay service..."
  systemctl enable nostr-relay
  systemctl start nostr-relay

  print_success "Nostr relay service setup completed."
}
