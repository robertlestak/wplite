#!/bin/bash

WP_THEME=${WP_THEME:-"twentytwentyfour"}

wp --allow-root theme install ${WP_THEME} --activate