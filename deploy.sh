#!/bin/bash -e
if [ $UID -ne 0 ]; then
	echo 'This script requires root' 1>&2
	exit 1
fi

# compile
echo compiling... 1>&2
go build

# create user
echo creating user... 1>&2
adduser --system --no-create-home --gecos '' --shell /bin/false --disabled-login --group gocyoa

echo installing... 1>&2

# create directory
mkdir -m 775 -p /srv/gocyoa

# copy files
install -m 775 gocyoa /srv/gocyoa
install -m 664 *.html /srv/gocyoa

# install systemd service
install -m 664 gocyoa.service /etc/systemd/system
systemctl daemon-reload

echo starting... 1>&2
systemctl restart gocyoa

echo deployment successful 1>&2
