#!/bin/bash
if [ ! $(uname) = "Linux" ]; then
	echo "Only support Linux right now"
	exit 1
fi

docker-compose -f docker_compo_test.yml down
sudo umount /tmp/ramdisk/