#!/bin/bash

# Function to pull the GitHub repository and run the main.sh script
setup_environment() {
  local repo_url="https://github.com/nodetec/relaywizard.git"
  local repo_dir="/tmp/relaywizard_setup"

  # Clone the GitHub repository
  echo "Cloning the repository from $repo_url..."
  git clone "$repo_url" "$repo_dir"

  if [ $? -ne 0 ]; then
    echo "Failed to clone the repository."
    exit 1
  fi

  # Change to the repository directory
  cd "$repo_dir" || exit

  # Make all scripts executable
  echo "Making all scripts executable..."
  chmod +x lib/print_colors.sh install_packages.sh configure_nginx_http.sh get_certificates.sh configure_nginx_https.sh install_nostr_rs_relay.sh configure_nostr_rs_relay.sh setup_nostr_rs_relay_service.sh main.sh

  # Run the main.sh script
  echo "Running the main.sh script..."
  ./main.sh

  if [ $? -ne 0 ]; then
    echo "Failed to execute the main.sh script."
    exit 1
  fi

  echo "Environment setup completed successfully."
}

# Call the setup_environment function with the provided arguments
setup_environment
