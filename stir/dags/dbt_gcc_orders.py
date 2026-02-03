import datetime as dt
from airflow import DAG
from airflow.utils.dates import days_ago
from airflow_dbt_python.operators.dbt import DbtRunOperator
from slack_connection import SlackNotifier

with DAG(
        dag_id="dbt_gcc_orders",
        schedule_interval="0 0 * * *",
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        dagrun_timeout=dt.timedelta(minutes=90),
) as dag:
    dbt_run = DbtRunOperator(
        task_id="dbt_run_hourly",
        project_dir="/gcc/airflow/stir/gcc_social/",
        profiles_dir="/gcc/airflow/.dbt/",
        select=["tag:gcc_orders"],
        exclude=["tag:deprecated", "tag:post_ranker", "tag:core"],
        target="production",
        profile="gcc_warehouse",
    )
