#! /bin/bash
ulimit -n 100000
echo "Starting Coffee"
mkdir -p logs
PID=$(ps aux | grep coffee | grep -v grep | awk '{print $2}')
if [ -z "$PID" ]
then
  echo "Coffee is not running"
else
  echo "Bringing Out of LB"
  curl -vXPUT http://localhost:$PORT/heartbeat/?beat=false
  sleep 10
  echo "Killing Coffee"
  kill -9 $PID
  echo "Killed Coffee"
fi

sleep 10
echo "Starting go gateway"
cp ".env.$CI_ENVIRONMENT_NAME" .env
cp .env ./bin/.env
ENV=$CI_ENVIRONMENT_NAME ./bin/coffee >> "logs/out.log" 2>&1 &

echo "Started Coffee service"
rm -rf /bulbul/services/coffee/logs
ln -s $(pwd)/logs /bulbul/services/coffee/logs

rm -rf /bulbul/services/coffee/scripts
ln -s $(pwd)/scripts /bulbul/services/coffee/scripts

echo "Sleeping before bringing in LB"
sleep 10
echo "Bringing in LB"
curl -vXPUT http://localhost:$PORT/heartbeat/?beat=true
