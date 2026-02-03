#! /bin/bash
python3 -m pip  install -r requirements.txt
cp -r dags /gcc/airflow/dags
cp -r src/* /gcc/airflow/stir/
