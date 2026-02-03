#! /bin/bash
PORT=$1
ENV=$2
cd /home/gitlab/beat
python3.11 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
mkdir -p logs
cp ".env.$ENV" .env
pid=$(cat beat.pid)
if [ -z "$pid" ]
then
  echo "Beat main is not running"
else
  echo "killing beat main"
  ps aux | grep python | grep beat | grep main.py | awk '{print $2}' | xargs kill -9
  echo "killed main"
fi

pid=$(cat beat_server.pid)
if [ -z "$pid" ]
then
  echo "Beat server is not running"
else
  echo "Set heartbeat to false"
  curl -vXPUT http://localhost:$PORT/heartbeat?beat=false
  sleep 10
  echo "Killing beat server"
  ps aux | grep python | grep beat | grep server.py | awk '{print $2}' | xargs kill -9
  echo "Killed beat server"
fi

sleep 10
# ps aux | grep python | grep beat | awk '{print $2}' | xargs kill -9
python main.py --parent beat  &>> logs/main_out.log &
echo "started main"
python server.py --parent beat &>> logs/server_out.log &
echo "started server"
sleep 10
echo "Set heartbeat to true"
curl -vXPUT http://localhost:$PORT/heartbeat?beat=true