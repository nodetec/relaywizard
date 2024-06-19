#!/bin/bash

# Function to print messages in green
print_success() {
  echo -e "\033[0;32m✓ $1\033[0m"
}

# Function to print messages in red
print_error() {
  echo -e "\033[0;31m☠ $1\033[0m"
}

# Function to print messages in cyan
print_info() {
  echo -e "\033[0;36mℹ $1\033[0m"
}

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

  # Prompt for domain name if not provided as an argument
  read -p "Enter the domain name for the nostr relay site (e.g., example.com): " domain_name < /dev/tty

  read -p "Enter the email address for SSL certificate registration: " email < /dev/tty

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
setup_environment

