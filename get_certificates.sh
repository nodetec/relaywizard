#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to get SSL certificates using Certbot
get_certificates() {
	local domain_name="$1"
	local email="$2"
	local dir_name
	dir_name=$(echo "$domain_name" | awk -F. '{print ($(NF-1) == "com" || $(NF-1) == "org" || $(NF-1) == "net" || $(NF-1) == "co") ? $(NF-2) : $(NF-1)}')

	# Check if certificates already exist
	if [ -f "/etc/letsencrypt/live/$domain_name/fullchain.pem" ] && [ -f "/etc/letsencrypt/live/$domain_name/privkey.pem" ]; then
		print_success "SSL certificates for $domain_name already exist."
	else
		print_info "Creating necessary directories for Certbot..."
		mkdir -p /var/www/$dir_name/.well-known/acme-challenge/

		print_info "Obtaining SSL certificates for $domain_name using Certbot..."
		certbot certonly --webroot -w /var/www/$dir_name -d $domain_name --email $email --agree-tos --no-eff-email -q

		if [ $? -ne 0 ]; then
			print_error "Certbot failed to obtain the certificate for $domain_name."
			exit 1
		fi
	fi
}
