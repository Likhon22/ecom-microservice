#!/bin/sh
set -e

# Generate kong.yml from template in /tmp
perl -pe 'BEGIN{undef $/;open(F,q(/kong/descriptor.b64));$r=<F>;close F;chomp $r} s/\$\{PROTO_B64\}/$r/g' /kong/kong.tpl.yml > /tmp/kong.yml

# Start Kong with the generated config
exec /docker-entrypoint.sh kong docker-start
