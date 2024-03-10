<?php
/*
Plugin Name: wplite
Description: wplite
Version: 0.0.1
Author: shdw.tech
*/

function is_cli_running() {
    return defined( 'WP_CLI' ) && WP_CLI;
}

if ( is_cli_running() ) {
    require_once 'cli.php';
}

function wplite_export_complete() {
    exec('touch /wplite/sigs/.build_complete');
}

add_action('ss_completed', 'wplite_export_complete');
