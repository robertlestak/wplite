#!/bin/bash

if [[ -f /var/www/html/wp-content/database/.ht.sqlite ]]; then
    exit 0
fi

wptitle=${WP_TITLE:-"My WordPress Site"}
wpuser=${WP_USER:-"admin"}
wppass=${WP_PASS:-"admin"}
wpemail=${WP_EMAIL:-"wplite@example.com"}
wp_port=${WP_PORT:-"80"}
siteurl="localhost"
curl -d "weblog_title=$wptitle&user_name=$wpuser&admin_password=$wppass&admin_password2=$wppass&admin_email=$wpemail" \
    http://$siteurl/wp-admin/install.php?step=2

if [[ $? -ne 0 ]]; then
    echo "WordPress installation failed"
    exit 1
fi

if [[ $wp_port -eq 80 ]]; then
    exit 0
fi

sleep 5

old_url="http://$siteurl"
new_url="http://$siteurl:$wp_port"

# this doesn't yet work with sqlite
# wp --allow-root search-replace "$old_url" "$new_url" --all-tables

# so we need to manually invoke sqlite3 to update the siteurl and home
sqlite3 /var/www/html/wp-content/database/.ht.sqlite \
    "UPDATE wp_options SET option_value = '$new_url' WHERE option_name = 'siteurl' OR option_name = 'home';"