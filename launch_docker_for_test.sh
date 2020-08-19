#!/usr/bin/env bash
docker-compose -f docker_compo_test.yml up -d
MYSQL="mysql -h 127.0.0.1 -P 13306 -uroot -pMytest123! "
for i in $(seq 60 $END); do
  echo "SHOW STATUS like 'Connections';" | $MYSQL |& grep "ERROR"
  result=$?
  if [[ $result -ne 0 ]]; then
    echo "SHOW STATUS like 'Uptime';" | $MYSQL | grep 'Uptime'
    echo "Check mysql server working status successfully."
    break
  fi
  sleep 1
done
