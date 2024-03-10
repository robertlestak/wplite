#!/bin/bash

pid=0

function install_wp() {
    if [[ -f /var/www/html/wp-content/database/.ht.sqlite ]]; then
        echo "WordPress already installed"
        return 0
    fi
    bash /wplite/scripts/install-wp.sh
    bash /wplite/scripts/install-theme.sh
    bash /wplite/scripts/install-plugins.sh
    bash /wplite/scripts/configure-simply-static.sh
}

function run_entrypoint() {
    docker-entrypoint.sh apache2-foreground
}

function stop_entrypoint() {
    kill -TERM $pid
    exit 0
}

function stop_when_file_exists() {
    mkdir -p /wplite/sigs
    chmod -R 777 /wplite/sigs
    while true; do
        if [[ -f /wplite/sigs/.stop ]]; then
            stop_entrypoint
        fi
        sleep 5
    done
}

function handle_signals() {
    kill -TERM $pid
    wait $pid
    exit 0
}

function main() {
    trap 'handle_signals' SIGINT SIGTERM
    chown -R www-data:www-data /var/www/html
    run_entrypoint &
    pid=$!
    sleep 4
    stop_when_file_exists &
    install_wp
    if [[ $? -eq 0 ]]; then
        echo "=============================================="
        echo "WordPress up and running on http://localhost:${WP_PORT}"
        echo "=============================================="
    else
        echo "WordPress installation failed"
        exit 1
    fi
    wait $pid
    exit 0
}

main "$@"