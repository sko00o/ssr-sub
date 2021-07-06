#!/usr/bin/env bash
###
# File: launch-ssrs.sh
# Author: Ming Cheng<mingcheng@outlook.com>
#
# Created Date: Thursday, August 8th 2019, 3:32:38 pm
# Last Modified: Tuesday, April 27th 2021, 7:40:19 pm
#
# http://www.opensource.org/licenses/MIT
###

if [ -z $SSR_SUBSCRIBER ]; then
  SSR_SUBSCRIBER=http://ssr-subscriber.default.svc.cluster.local/random/json
fi

CHECK_SOCK5_URL=https://www.google.com
SSR_CONF_FILE=/tmp/ss-local.json
SSR_PID_FILE=/var/run/ss-local.pid

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] ${1}"
}

start_ssr() {
  if [ -f $SSR_PID_FILE ]; then
    log "ss-local maybe is running, so stop first"
    stop_ssr
  fi

  curl -sSkL -o $SSR_CONF_FILE $SSR_SUBSCRIBER
  cat $SSR_CONF_FILE
  ss-local $SSR_OPT -v -c $SSR_CONF_FILE -l 1086 -b 0.0.0.0 -f $SSR_PID_FILE &

  while true; do
    sleep 600
    check_ssr
  done
}

check_ssr() {
  PROXY_ADDR=127.0.0.1:1086
  curl_command="curl -sSkL -w %{http_code} \
	  -o /dev/null -X HEAD \
	  -x socks5://${PROXY_ADDR} \
	  --connect-timeout 10 --max-time 30 \
	  ${CHECK_SOCK5_URL}"

  echo $curl_command
  if [ $($curl_command) == "200" ]; then
    log "ss-local connection check is ok"
    return 0
  else
    log "ss-local connection check is failed"
    return 255
  fi
}

stop_ssr() {
  kill -INT $(cat $SSR_PID_FILE)
}

case $1 in
start)
  start_ssr
  ;;
stop)
  stop_ssr
  ;;
check)
  check_ssr
  ;;
*)
  echo "Usage $0 start|stop|check"
  ;;
esac
