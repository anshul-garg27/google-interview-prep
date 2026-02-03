#! /bin/bash
ulimit -n 100000
echo "Starting Event Grpc"
mkdir -p logs
PID=$(ps aux | grep event-grpc | grep -v grep | awk '{print $2}')
if [ -z "$PID" ]
then
  echo "Event Grpc is not running"
else
  echo "Bringing OOLB"
  curl -vXPUT http://localhost:$GIN_PORT/heartbeat/?beat=false
  sleep 15
  echo "Killing Event Grpc"
  kill -9 $PID
  echo "Killed Event Grpc"
fi

sleep 10
echo "Starting Event Grpc"
ENV=$CI_ENVIRONMENT_NAME ./event-grpc >> "logs/out.log" 2>&1 &

echo "Started Event Grpc"
rm -rf /bulbul/services/event-grpc/logs
ln -s $(pwd)/logs /bulbul/services/event-grpc/logs

echo "Sleeping before bringing in LB"
sleep 20
echo "Bringing in LB"
curl -vXPUT http://localhost:$GIN_PORT/heartbeat/?beat=true