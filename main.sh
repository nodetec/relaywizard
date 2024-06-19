#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Source the scripts containing the functions
source ./install_packages.sh
source ./configure_nginx_http.sh

# Prompt for domain name if not provided as an argument
if [ -z "$1" ]; then
  read -p -r "Enter the domain name for the nostr relay site (e.g., example.com): " domain_name
else
  domain_name="$1"
fi

# Call the function to install packages
print_info "Running install_packages function..."
install_packages

# Call the function to configure nginx for HTTP
print_info "Running configure_nginx_http function..."
configure_nginx_http "$domain_name"

print_success "Setup completed successfully."
