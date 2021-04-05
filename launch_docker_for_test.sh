#!/usr/bin/env bash
if [[ ! $(uname) = "Linux" ]]; then
	echo "Only support Linux right now"
	exit 1
fi

# Create ramdisk if it doesn't exist.
[[ ! -d /tmp/ramdisk ]] && sudo mkdir /tmp/ramdisk

# Mount ramdisk has been created.
[[ -d /tmp/ramdisk ]] && sudo mount -t tmpfs -o size=500M tmpfs /tmp/ramdisk/

docker-compose -f docker_compo_test.yml up -d
MYSQL="mysql -h 127.0.0.1 -P 13306 -uroot -pMytest123! "
for i in $(seq 60); do
  echo "SHOW STATUS like 'Connections';" | ${MYSQL} |& grep "ERROR"
  result=$?
  if [[ ${result} -ne 0 ]]; then
    echo "SHOW STATUS like 'Uptime';" | ${MYSQL} | grep 'Uptime'
    echo "Check mysql server working status successfully."
    break
  fi
  sleep 1
done
