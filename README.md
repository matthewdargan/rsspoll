# Rsspoll

Rsspoll polls RSS feeds for updates.

Usage:

    rsspoll [-d days] [file]

Rsspoll reads RSS feeds from a file and prints entries within the last
d days (default 1).

The `-d` flag specifies the number of days to recall.

The file containing RSS feeds of interest should either be passed as an
argument or exist at $XDG_CONFIG_HOME/rsspoll/config.txt. Each feed URL
should be on a separate line within the file.
