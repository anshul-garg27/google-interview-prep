from airflow.decorators import task
from airflow.models import Variable
from airflow.operators.empty import EmptyOperator
from airflow.providers.slack.operators.slack_webhook import SlackWebhookOperator
from airflow.hooks.base_hook import BaseHook
from airflow.notifications.basenotifier import BaseNotifier
import logging


class SlackNotifier(BaseNotifier):
    def slack_fail_alert(context):
        
        ti = context['ti'] 
        task_state = ti.state 
        SLACK_CONN_ID = None
        slack_webhook_token = None
        channel = None
        slack_msg = ''
        if task_state == 'success':
            return
            # SLACK_CONN_ID = Variable.get('slack_success_conn')
            # slack_webhook_token = BaseHook.get_connection(SLACK_CONN_ID).password
            # channel = BaseHook.get_connection(SLACK_CONN_ID).login
            # duration_minutes = context.get('task_instance').duration / 60

            # slack_msg = f"""
            # :white_check_mark: Task Succeeded.
            # *Task*: {context.get('task_instance').task_id}
            # *Dag*: {context.get('task_instance').dag_id}
            # *Execution Time*: {context.get('execution_date')}
            # *Duration*: {duration_minutes:.2f} minutes
            # """
            
        
        elif task_state =='failed':
            SLACK_CONN_ID = Variable.get('slack_failure_conn')
            slack_webhook_token = BaseHook.get_connection(SLACK_CONN_ID).password
            channel = BaseHook.get_connection(SLACK_CONN_ID).login
                        
            slack_msg = f"""
            :x: Task Failed.
            *Task*: {context.get('task_instance').task_id}
            *Dag*: {context.get('task_instance').dag_id}
            *Execution Time*: {context.get('execution_date')}
            <{context.get('task_instance').log_url}|*Logs*>
            """
        
        slack_alert = SlackWebhookOperator(
            task_id='slack_fail',
            webhook_token=slack_webhook_token,
            message=slack_msg,
            channel=channel,
            username='airflow',
            http_conn_id=SLACK_CONN_ID
        )

        return slack_alert.execute(context=context)
