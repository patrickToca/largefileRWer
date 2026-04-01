#!/bin/bash
# Requires ImageMagick installed

# Convert PNG to multiple sizes
convert templates/img/favicon.png -resize 16x16 templates/img/favicon-16x16.png
convert templates/img/favicon.png -resize 32x32 templates/img/favicon-32x32.png
convert templates/img/favicon.png -resize 180x180 templates/img/apple-touch-icon.png
convert templates/img/favicon.png -resize 192x192 templates/img/android-chrome-192x192.png
convert templates/img/favicon.png -resize 512x512 templates/img/android-chrome-512x512.png

echo "Favicons generated successfully!"