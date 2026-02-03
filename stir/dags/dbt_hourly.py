import datetime as dt
from airflow import DAG
from airflow.utils.dates import days_ago
from airflow_dbt_python.operators.dbt import DbtRunOperator
from slack_connection import SlackNotifier
from airflow.operators.python import ShortCircuitOperator

def check_hour():
    hour = dt.datetime.now().hour
    if 19 <= hour < 22:
        return False
    return True

with DAG(
        dag_id="dbt_hourly",
        schedule_interval="0 1-23/2 * * *",
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        dagrun_timeout=dt.timedelta(minutes=180),
) as dag:
    check_hour_task = ShortCircuitOperator(
        task_id="check_hour",
        python_callable=check_hour,
    )
    
    dbt_run = DbtRunOperator(
        task_id="dbt_run_hourly",
        project_dir="/gcc/airflow/stir/gcc_social/",
        profiles_dir="/gcc/airflow/.dbt/",
        select=["tag:hourly"],
        exclude=["tag:deprecated", "tag:post_ranker", "tag:core"],
        target="production",
        profile="gcc_warehouse",
    )

    check_hour_task >> dbt_run
    