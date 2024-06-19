#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to pull the GitHub repository and run the main.sh script
setup_environment() {
  local repo_url="https://github.com/nodetec/relaywizard.git"
  local repo_dir="/tmp/relaywizard_setup"

  # Clone the GitHub repository
  print_info "Cloning the repository from $repo_url..."
  git clone "$repo_url" "$repo_dir"

  if [ $? -ne 0 ]; then
    print_error "Failed to clone the repository."
    exit 1
  fi

  # Change to the repository directory
  cd "$repo_dir" || exit

  # Make all scripts executable
  print_info "Making all scripts executable..."
  chmod +x lib/print_colors.sh install_packages.sh configure_nginx_http.sh get_certificates.sh configure_nginx_https.sh install_nostr_rs_relay.sh configure_nostr_rs_relay.sh setup_nostr_rs_relay_service.sh main.sh

  # Prompt for domain name if not provided as an argument
  if [ -z "$1" ]; then
    read -p "Enter the domain name for the nostr relay site (e.g., example.com): " domain_name
  else
    domain_name="$1"
  fi

  # Prompt for email address if not provided as an argument
  if [ -z "$2" ]; then
    read -p "Enter the email address for SSL certificate registration: " email
  else
    email="$2"
  fi

  # Run the main.sh script
  print_info "Running the main.sh script..."
  ./main.sh "$domain_name" "$email"

  if [ $? -ne 0 ]; then
    print_error "Failed to execute the main.sh script."
    exit 1
  fi

  print_success "Environment setup completed successfully."
}

# Call the setup_environment function with the provided arguments
setup_environment "$1" "$2"
