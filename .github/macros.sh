#!/bin/sh

scriptdir=$(dirname $(readlink -f $0))

echo "# Macros"
echo

$scriptdir/../cmd/carapace/carapace --macros \
  | sed 's_^\([^ ]*\.\)\([^. ]\+\) \+\(.*\)_- [\1\2]\(https://pkg.go.dev/github.com/rsteube/carapace-bin/pkg/actions/\1#Action\2) \3_' \
  | sed -e ':loop' -e 's_\(carapace-bin/pkg/actions/[^#]*\)[.]_\1/_' -e 't loop' 
