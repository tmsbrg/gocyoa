#!/bin/bash -e
if [ $UID -ne 0 ]; then
	echo 'This script requires root' 1>&2
	exit 1
fi

# compile
go build

# create user
adduser --system --no-create-home --gecos '' --shell /bin/false --disabled-login --group gocyoa

# create directory
mkdir -m 775 -p /srv/gocyoa

# copy files
install -m 775 gocyoa /srv/gocyoa
install -m 664 *.html /srv/gocyoa

# install systemd service
install -m 664 gocyoa.service /etc/systemd/system
systemctl daemon-reload
