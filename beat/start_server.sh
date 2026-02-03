#! /bin/bash
lsof -i :8000 | awk '{print $2}' | tail -n +2 | xargs kill -9
uvicorn server:app --reload