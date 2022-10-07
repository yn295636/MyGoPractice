#!/usr/bin/env bash
[[ ! -d $HOME/mysql_data ]] && mkdir $HOME/mysql_data
[[ ! -d $HOME/mongo_data ]] && mkdir $HOME/mongo_data
[[ ! -d $HOME/nsq_data ]] && mkdir $HOME/nsq_data

docker-compose -f docker_compo_develop.yml up -d

MYSQL="mysql -h 127.0.0.1 -P 3306 -uroot -pMygopractice123! "
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
