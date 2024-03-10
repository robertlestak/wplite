#!/bin/bash

PLUGINS_DIR=/var/www/html/wp-content/plugins

PLUGINS=(`cat /wplite/scripts/default-plugins.txt`)

function install_official_plugin() {
    local plugin=$1
    curl -o ${PLUGINS_DIR}/${plugin}.zip https://downloads.wordpress.org/plugin/${plugin}.zip && \
    unzip ${PLUGINS_DIR}/${plugin}.zip -d ${PLUGINS_DIR}/ && \
    rm ${PLUGINS_DIR}/${plugin}.zip
    wp --allow-root plugin activate $plugin
    echo "Plugin ${plugin} installed successfully"
}

function install_plugins() {
    for plugin in ${PLUGINS[@]}; do
        install_official_plugin $plugin
    done
    cp -r /wplite/scripts/plugins/wplite ${PLUGINS_DIR}/
    wp --allow-root plugin activate wplite
}

function main() {
    mkdir -p ${PLUGINS_DIR}
    install_plugins
}

main "$@"