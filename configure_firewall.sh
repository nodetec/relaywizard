#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to configure the firewall
configure_firewall() {
	print_info "Configuring firewall to allow HTTP (port 80) and HTTPS (port 443) traffic..."

	# Allow HTTP and HTTPS traffic
	ufw allow 'Nginx Full'

	# Reload the firewall to apply the changes
	ufw reload

	# Show the current firewall status
	ufw status verbose

	print_success "Firewall configuration completed successfully."
}

# Call the configure_firewall function
configure_firewall
