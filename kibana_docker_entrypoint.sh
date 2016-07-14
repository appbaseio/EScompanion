#!/bin/bash
set -e

# Add kibana as command if needed
if [[ "$1" == -* ]]; then
	set -- kibana "$@"
fi

# Run as user "kibana" if the command is "kibana"
if [ "$1" = 'kibana' ]; then
    
    args=("$@")
    length="$#"
    length=$((length-1))

    for i in $(seq 1 "$length");
    do
        kibana plugin --install ${args[$i]}
    done

	if [ "$ELASTICSEARCH_URL" ]; then
		sed -ri "s!^(\#\s*)?(elasticsearch\.url:).*!\2 '$ELASTICSEARCH_URL'!" /opt/kibana/config/kibana.yml
	fi

	set -- gosu kibana tini -- "$@"
fi

exec "$@"
