import os
from datetime import date, datetime


def keyword_collection_logs(message, job_id):
    current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S.%f")[:-3]
    level_name = "DEBUG"
    log_message = f"{current_time} - {level_name} || {message}"
    current_date = date.today()
    directory_path = f"{os.path.expanduser('~')}/keyword_logs/{current_date}"
    if not os.path.exists(directory_path):
        os.makedirs(directory_path)
    keywords_log_path = f"{directory_path}/{job_id}.log"
    with open(keywords_log_path, 'a+') as file:
        file.write(log_message + "\n")
