#!/bin/bash

current_config=`wp --allow-root option get simply-static --format=json`

# set the following fields:
# delivery_method: local
# local_dir: /var/www/html/wp-content/static
# clear_directory_before_export: "1"

updated_config=`echo $current_config | jq '.temp_files_dir = "/tmp/" | .delivery_method = "local" | .local_dir = "/var/www/html/wp-content/static" | .clear_directory_before_export = "1"'`


wp --allow-root option update \
    simply-static \
    "$updated_config" \
    --format=json