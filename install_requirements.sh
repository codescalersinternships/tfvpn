#!/bin/bash

packages_to_install=()

if ! command -v sshuttle &> /dev/null; then
    echo "sshuttle could not be found"
    packages_to_install+=("sshuttle")
fi

if ! command -v python3 &> /dev/null; then
    echo "python3 could not be found"
    packages_to_install+=("python3")
fi

if [[ ${#packages_to_install[@]} -gt 0 ]]; then
    read -p "Do you want to install the missing packages? (Y/n) " -n 1 -r
    echo
    REPLY="${REPLY:-y}"

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        sudo apt-get install -y "${packages_to_install[@]}"
    else
        echo "Installation aborted. Exiting..."
        exit 1
    fi
fi
