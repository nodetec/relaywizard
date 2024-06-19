#!/bin/bash

# Source the print_colors.sh script to use color printing functions
source ./lib/print_colors.sh

# Function to configure nginx for HTTPS
configure_nginx_https() {
	local domain_name="$1"
	local dir_name
	dir_name=$(echo "$domain_name" | awk -F. '{print ($(NF-1) == "com" || $(NF-1) == "org" || $(NF-1) == "net" || $(NF-1) == "co") ? $(NF-2) : $(NF-1)}')

	print_info "Removing existing nginx configuration if it exists..."
	rm -f /etc/nginx/conf.d/nostr_relay.conf

	print_info "Configuring nginx for HTTPS..."
	cat <<EOL >/etc/nginx/conf.d/nostr_relay.conf
map \$http_upgrade \$connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket {
    server 127.0.0.1:8080;
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name $domain_name;

    location / {
        proxy_pass http://websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection \$connection_upgrade;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-For \$remote_addr;
    }

    #### SSL Configuration ####
    ssl_certificate /etc/letsencrypt/live/$domain_name/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$domain_name/privkey.pem;

    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_protocols TLSv1.2 TLSv1.1 TLSv1;
    ssl_prefer_server_ciphers on;
    ssl_ciphers "EECDH+ECDSA+AESGCM EECDH+aRSA+AESGCM EECDH+ECDSA+SHA384 EECDH+ECDSA+SHA256 EECDH+aRSA+SHA384 EECDH+aRSA+SHA256 EECDH+aRSA+RC4 EECDH EDH+aRSA RC4 !aNULL !eNULL !LOW !3DES !MD5 !EXP !PSK !SRP !DSS";
    ssl_stapling on;
    ssl_stapling_verify on;
    ssl_ecdh_curve secp384r1;

    add_header Strict-Transport-Security "max-age=31536000; includeSubdomains";
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Referrer-Policy same-origin;
    add_header Feature-Policy "geolocation none;midi none;notifications none;push none;sync-xhr none;microphone none;camera none;magnetometer none;gyroscope none;speaker self;vibrate none;fullscreen self;payment none;";
}

server {
    listen 80;
    listen [::]:80;
    server_name $domain_name;

    location /.well-known/acme-challenge/ {
        root /var/www/$dir_name;
        allow all;
    }

    location / {
        return 301 https://$domain_name\$request_uri;
    }
}
EOL

	print_info "Reloading nginx to apply the configuration..."
	systemctl reload nginx
}
