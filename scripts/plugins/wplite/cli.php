<?php

class cli_wplite extends WP_CLI_Command {        
    /**
     * triggers a static site build.
     * 
     * ## OPTIONS
     * 
     * [--exit]
     * : If true, the main process will be stopped after the build is complete.
     * 
     */
    function build($args, $assoc_args) {
        WP_CLI::line( 'Building static site. Depending on the size of your site, this may take a while....' );
        do_action('simply_static_site_export_cron');
        // loop and wait for a file called /wplite/sigs/.build_complete to appear
        // once it does, we can exit
        while (!file_exists('/wplite/sigs/.build_complete')) {
            sleep(1);
        }
        unlink('/wplite/sigs/.build_complete');
        WP_CLI::line( 'Build complete!' );
        if ($assoc_args['exit'] == 'true') {
            // now create a file /wplite/sigs/.stop to signal the main process to stop
            touch('/wplite/sigs/.stop');
        }
    }
}
   
WP_CLI::add_command( 'wplite', 'cli_wplite' );