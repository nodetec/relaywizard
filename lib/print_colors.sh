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

# Function to print messages in yellow
print_warn() {
  echo -e "\033[0;33m☢ $1\033[0m"
}

# Function to print messages in yellow
print_zap() {
  echo -e "\033[0;33m⚡ $1\033[0m"
}

# Function to print messages in blue
print_blue() {
  echo -e "\033[0;34m$1\033[0m"
}

# Function to print messages in magenta
print_magenta() {
  echo -e "\033[0;35m$1\033[0m"
}

# Function to print messages in bold
print_bold() {
  echo -e "\033[1m$1\033[0m"
}

# Function to print messages in underline
print_underline() {
  echo -e "\033[4m$1\033[0m"
}

# Function to print messages with a green background
print_bg_green() {
  echo -e "\033[42m$1\033[0m"
}

# Function to print messages with a red background
print_bg_red() {
  echo -e "\033[41m$1\033[0m"
}

# Function to print messages with a yellow background
print_bg_yellow() {
  echo -e "\033[43m$1\033[0m"
}

# Function to print messages with a blue background
print_bg_blue() {
  echo -e "\033[44m$1\033[0m"
}

# Function to print messages with a magenta background
print_bg_magenta() {
  echo -e "\033[45m$1\033[0m"
}

# Function to print messages with a cyan background
print_bg_cyan() {
  echo -e "\033[46m$1\033[0m"
}

# Function to print blinking text
print_blink() {
  echo -e "\033[5m$1\033[0m"
}

# Function to print messages with inverted colors
print_inverted() {
  echo -e "\033[7m$1\033[0m"
}

# Function to print messages in high intensity green
print_high_intensity_green() {
  echo -e "\033[0;92m$1\033[0m"
}

# Function to print messages in dim text
print_dim() {
  echo -e "\033[2m$1\033[0m"
}
