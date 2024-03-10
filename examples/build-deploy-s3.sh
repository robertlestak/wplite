#!/bin/bash

# This script generates a new site and deploys it to S3.

S3_BUCKET=my-bucket
CLOUDFRONT_DISTRIBUTION_ID=my-distribution-id

# create a temporary directory to work in and change to it
SITE_DIR=$(mktemp -d)
cd $SITE_DIR

# create a new site environment file
cat > .wplite-env <<EOF
WP_TITLE=my-site
WP_USER=admin
WP_PASS=password
WP_EMAIL=hello@example.com
WP_PORT=80
WP_THEME=twentytwentyfour
EOF

# build the site
wplite build

# deploy the site to S3
aws s3 sync ./wp-content/static/ s3://$S3_BUCKET --delete

# invalidate the CloudFront cache
if [ -n "$CLOUDFRONT_DISTRIBUTION_ID" ]; then
  aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DISTRIBUTION_ID --paths "/*"
fi

# clean up
cd ..
rm -rf $SITE_DIR