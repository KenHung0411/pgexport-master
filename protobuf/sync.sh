#!/usr/bin/env bash

TOKEN="Wa27egT5zMEGb32dCgud"
FILES=$(cat <<EOF
navyx/bluex-protocol booking_proto/booking_event.proto
navyx/bluex-protocol prebooking_proto/prebooking.proto
navyx/notification-hub protobuf/notifications.proto notifications/notifications.proto
EOF
)

echo "${FILES}" | while read repo file target
do
	echo -n "Downloading $repo:$file ... "
	target=${target:-$file}
	mkdir -p $(dirname $target)
	repo=${repo//\//%2F}
	file=${file//\//%2F}
	url="https://gitlab.com/api/v4/projects/$repo/repository/files/$file/raw?ref=master"
	curl -sSL --header "PRIVATE-TOKEN: ${TOKEN}" "$url" -o "${target}"
	echo "done"
done
