#!/usr/bin/env bash
[ ! -d mysql_data ] && mkdir mysql_data
[ ! -d mongo_data ] && mkdir mongo_data

docker-compose -f docker_compo.yml up -d

MYSQL="mysql -h 127.0.0.1 -P 3306 -uroot -pMygopractice123! "
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
