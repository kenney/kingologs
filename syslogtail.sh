#!/bin/sh
tail -F -q /var/log/system.log | \
while read -r line ; do
    # send to remote syslog daemon
    echo "<0>$line" | nc localhost 65140
done
