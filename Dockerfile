FROM scratch
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

# Add the binary and data files
ADD dist/informas /usr/local/bin/
ADD data /usr/local/share/informas/data

# Expose port 8000 by default
EXPOSE 8000

# The default command simply runs the binary
CMD [ \
    "informas", \
    "run", \
    "--data-dir", "/usr/local/share/informas/data" \
]
