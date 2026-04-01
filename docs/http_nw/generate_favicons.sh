#!/bin/bash
# Requires ImageMagick installed

# Convert is deprecated! replaced by magick PNG to multiple sizes
magick templates/img/favicon.png -resize 16x16 templates/img/favicon-16x16.png
magick templates/img/favicon.png -resize 32x32 templates/img/favicon-32x32.png
magick templates/img/favicon.png -resize 180x180 templates/img/apple-touch-icon.png
magick templates/img/favicon.png -resize 192x192 templates/img/android-chrome-192x192.png
magick templates/img/favicon.png -resize 512x512 templates/img/android-chrome-512x512.png

echo "Favicons generated successfully!"