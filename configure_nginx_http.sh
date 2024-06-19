#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to extract the directory name from the domain
get_directory_name() {
  local domain_name="$1"
  IFS='.' read -r -a domain_parts <<<"$domain_name"
  if [[ ${#domain_parts[@]} -gt 2 ]]; then
    echo "${domain_parts[1]}"
  else
    echo "${domain_parts[0]}"
  fi
}

# Function to configure nginx for HTTP
configure_nginx_http() {
  local domain_name="$1"
  local dir_name
  dir_name=$(get_directory_name "$domain_name")

  print_info "Creating necessary directories..."
  mkdir -p /var/www/$dir_name
  mkdir -p /var/www/$dir_name/.well-known/acme-challenge/

  print_info "Removing existing nginx configuration if it exists..."
  rm -f /etc/nginx/conf.d/nostr_relay.conf

  print_info "Configuring nginx for HTTP..."
  cat <<EOL >/etc/nginx/conf.d/nostr_relay.conf
map \$http_upgrade \$connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket {
    server 127.0.0.1:8080;
}

# $domain_name
server {
    listen 80;
    listen [::]:80;
    server_name $domain_name;

    location /.well-known/acme-challenge/ {
        root /var/www/$dir_name;
        allow all;
    }

    location / {
        proxy_pass http://websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection \$connection_upgrade;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-For \$remote_addr;
    }
}
EOL

  print_info "Reloading nginx to apply the configuration..."
  systemctl reload nginx
}
