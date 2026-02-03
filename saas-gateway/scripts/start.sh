#! /bin/bash
ulimit -n 100000
echo "Starting SaaS Gateway"
mkdir -p logs
PID=$(ps aux | grep saas-gateway | grep -v grep | awk '{print $2}')
if [ -z "$PID" ]
then
  echo "SaaS gateway is not running"
else
  echo "Bringing OOLB"
  curl -vXPUT http://localhost:$PORT/heartbeat/?beat=false
  sleep 15
  echo "Killing SaaS Gateway"
  kill -9 $PID
  echo "Killed SaaS Gateway"
fi

sleep 10
echo "Starting SaaS Gateway"
ENV=$CI_ENVIRONMENT_NAME ./saas-gateway >> "logs/out.log" 2>&1 &

echo "Started SaaS Gateway service"
rm -rf /bulbul/services/saas-gateway/logs
ln -s $(pwd)/logs /bulbul/services/saas-gateway/logs
echo "Sleeping before bringing in LB"
sleep 20
echo "Bringing in LB"
curl -vXPUT http://localhost:$PORT/heartbeat/?beat=true