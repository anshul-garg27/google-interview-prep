from clickhouse_driver import Client
import concurrent.futures
import random
import requests
import time
from airflow import DAG
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier
from airflow.operators.python_operator import PythonOperator

eligible_users_with_campaigns_query = """
select cp_id,
       account_id,
       arraySlice(reward_ladder_campaign_ids, 1, 2) as challenge_campaign_ids,
       arraySlice(non_reward_ladder_campaign_ids, 1, 2) as collab_campaign_ids,
       webengage_user_id
from dbt.mart_user_eligible_campaigns
where cp_id > $cp_id_explored_value
order by cp_id asc
limit 100;
"""

campaign_wise_combined_data_query = """
select campaign_id,
       campaign_name,
       campaign_description,
       campaign_tentative_payout_text,
       campaign_image,
       campaign_start_time,
       campaign_end_time,
       # toString(campaign_start_time) as campaign_start_time,
       # toString(campaign_end_time) as campaign_end_time,
       campaign_available_slots_count,
       campaign_crm_invitation_link,
       campaign_type,
       reward_ladder_code,
       reward_ladder_level,
       campaign_invite_link,
       challenge_id,
       challenge_name,
       challenge_image,
       toString(challenge_end_time) as challenge_end_time,
       challenge_crm_invitation_link,
       challenge_link
from dbt.collab_challenge_slot_combined_data
where campaign_id in ($campaign_ids_for_combined_data)
"""


def send_webengage_event(url, data):
    r = random.random()
    # print(data)
    r >= 0.90 and print(data)
    response = requests.post(url, json=data)
    # print(response)
    r >= 0.90 and print(response)

def crm_recommendation_invitation():
    client = Client(host='172.31.28.68', user='old_crons', password='96qBac7D0U5')
    fetch_next = True
    eligible_users_query_param_cp_id = -1
    while fetch_next:
        # eligible_users_query_params = {'cp_id_explored_value': eligible_users_query_param_cp_id}  # initial
        formatted_eligible_users_query = eligible_users_with_campaigns_query.replace('$cp_id_explored_value', str(eligible_users_query_param_cp_id))
        eligible_users_data_resp = client.execute(formatted_eligible_users_query,
                                                  with_column_types=True)
        if (eligible_users_data_resp is None or eligible_users_data_resp[0] is None):
            break
        eligible_users_data_rows = eligible_users_data_resp[0]
        executor = concurrent.futures.ThreadPoolExecutor(max_workers=25)
        old_eligible_users_query_param_cp_id = eligible_users_query_param_cp_id
        if eligible_users_data_rows is not None and len(eligible_users_data_rows) > 0:
            campaign_ids_for_combined_data_query = []
            campaign_wise_combined_data = {}
            for eligible_influencer_data in eligible_users_data_rows:
                challenge_campaign_ids = eligible_influencer_data[2]
                collab_campaign_ids = eligible_influencer_data[3]
                if challenge_campaign_ids is not None and len(challenge_campaign_ids) > 0:
                    for challenge_campaign_id in challenge_campaign_ids:
                        if challenge_campaign_id is not None:
                            campaign_ids_for_combined_data_query.append(challenge_campaign_id)
                if collab_campaign_ids is not None and len(collab_campaign_ids) > 0:
                    for collab_campaign_id in collab_campaign_ids:
                        if collab_campaign_id is not None:
                            campaign_ids_for_combined_data_query.append(collab_campaign_id)
            if campaign_ids_for_combined_data_query is not None and len(campaign_ids_for_combined_data_query) > 0:
                # campaign_ids_for_combined_data_param = {'campaign_ids_for_combined_data': campaign_ids_for_combined_data_query}  # initial
                campaign_ids_for_combined_data_query_comma_separated_string = ', '.join(map(str, campaign_ids_for_combined_data_query))
                formatted_campaign_wise_combined_data_query = campaign_wise_combined_data_query.replace('$campaign_ids_for_combined_data',campaign_ids_for_combined_data_query_comma_separated_string)
                combined_campaign_data_resp = client.execute(formatted_campaign_wise_combined_data_query,
                                                             with_column_types=True)
                combined_campaign_data_rows = combined_campaign_data_resp[0]
                executor = concurrent.futures.ThreadPoolExecutor(max_workers=25)
                if combined_campaign_data_rows is not None and len(combined_campaign_data_rows) > 0:
                    for combined_campaign_data_row in combined_campaign_data_rows:
                        if combined_campaign_data_row is not None and combined_campaign_data_row[0] is not None:
                            campaign_wise_combined_data[combined_campaign_data_row[0]] = combined_campaign_data_row
            for eligible_influencer_data in eligible_users_data_rows:
                cp_id = eligible_influencer_data[0]
                account_id = eligible_influencer_data[1]
                challenge_campaign_ids = eligible_influencer_data[2]
                collab_campaign_ids = eligible_influencer_data[3]
                webengage_user_id = eligible_influencer_data[4]
                # if no account id present then jump to next records
                if account_id is None or webengage_user_id is None:
                    continue
                if cp_id is not None:
                    eligible_users_query_param_cp_id = cp_id  # extracting max cpId for next query
                collabCampaignId1 = None
                collabCampaignId2 = None
                challengeCampaignId1 = None
                challengeCampaignId2 = None
                if collab_campaign_ids is not None and len(collab_campaign_ids) > 0:
                    if (len(collab_campaign_ids) >= 2):
                        collabCampaignId1 = collab_campaign_ids[0]
                        collabCampaignId2 = collab_campaign_ids[1]
                    else:
                        collabCampaignId1 = collab_campaign_ids[0]
                if challenge_campaign_ids is not None and len(challenge_campaign_ids) > 0:
                    if (len(challenge_campaign_ids) >= 2):
                        challengeCampaignId1 = challenge_campaign_ids[0]
                        challengeCampaignId2 = challenge_campaign_ids[1]
                    else:
                        challengeCampaignId1 = challenge_campaign_ids[0]
                combinedCampaignData1 = None
                combinedCampaignData2 = None
                combinedChallengeData1 = None
                combinedChallengeData2 = None

                if campaign_wise_combined_data is not None:
                    if collabCampaignId1 is not None and campaign_wise_combined_data[collabCampaignId1] is not None:
                        combinedCampaignData1 = campaign_wise_combined_data[collabCampaignId1]
                    if collabCampaignId2 is not None and campaign_wise_combined_data[collabCampaignId2] is not None:
                        combinedCampaignData2 = campaign_wise_combined_data[collabCampaignId2]
                    if challengeCampaignId1 is not None and campaign_wise_combined_data[challengeCampaignId1] is not None:
                        # checking is challenge id present and>0 or not
                        if campaign_wise_combined_data[challengeCampaignId1][13] is not None and campaign_wise_combined_data[challengeCampaignId1][13] > 0:
                            combinedChallengeData1 = campaign_wise_combined_data[challengeCampaignId1]
                    if challengeCampaignId2 is not None and campaign_wise_combined_data[challengeCampaignId2] is not None:
                        # checking is challenge id present and>0 or not
                        if campaign_wise_combined_data[challengeCampaignId2][13] is not None and campaign_wise_combined_data[challengeCampaignId2][13] > 0:
                            combinedChallengeData2 = campaign_wise_combined_data[challengeCampaignId2]
                # if any of not null then only send event (i.e. atleast one should present)
                if (combinedCampaignData1 is not None
                        or combinedCampaignData2 is not None
                        or combinedChallengeData1 is not None
                        or combinedChallengeData2 is not None):
                    data_map = {}
                    data_map['appType'] = "USER"
                    data_map['userId'] = str(webengage_user_id)
                    data_map['eventName'] = "Final Campaign Recommendation"
                    data_map['eventTime'] = str(int(time.time() * 1000))
                    data_map['eventData'] = {}
                    # initializing with basic data

                    # data_map['eventData']['Collab 1 Id'] = ""
                    data_map['eventData']['User Id'] = str(webengage_user_id)
                    data_map['eventData']['Collab 1 Name'] = ""
                    data_map['eventData']['Collab 1 Description'] = ""
                    data_map['eventData']['Collab 1 Tentative Payout'] = ""
                    data_map['eventData']['Collab 1 Image URL'] = ""
                    # data_map['eventData']['Collab 1 Start Date & Time'] = ""
                    # data_map['eventData']['Collab 1 End Date & Time'] = ""
                    # data_map['eventData']['Collab 1 Slots Left'] = ""
                    data_map['eventData']['Collab 1 Deeplink'] = ""

                    # data_map['eventData']['Collab 2 Id'] = ""
                    data_map['eventData']['Collab 2 Name'] = ""
                    data_map['eventData']['Collab 2 Description'] = ""
                    data_map['eventData']['Collab 2 Tentative Payout'] = ""
                    data_map['eventData']['Collab 2 Image URL'] = ""
                    # data_map['eventData']['Collab 2 Start Date & Time'] = ""
                    # data_map['eventData']['Collab 2 End Date & Time'] = ""
                    # data_map['eventData']['Collab 2 Slots Left'] = ""
                    data_map['eventData']['Collab 2 Deeplink'] = ""

                    # data_map['eventData']['Challenge 1 Id'] = ""
                    data_map['eventData']['Challenge 1 Name'] = ""
                    data_map['eventData']['Challenge 1 Tentative Payout'] = ""
                    data_map['eventData']['Challenge 1 Image URL'] = ""
                    # data_map['eventData']['Challenge 1 Start Date & Time'] = ""
                    # data_map['eventData']['Challenge 1 End Date & Time'] = ""
                    data_map['eventData']['Challenge 1 Slots Left'] = ""
                    # data_map['eventData']['Challenge 1 Deeplink'] = ""

                    # data_map['eventData']['Challenge 2 Id'] = ""
                    data_map['eventData']['Challenge 2 Name'] = ""
                    data_map['eventData']['Challenge 2 Tentative Payout'] = ""
                    data_map['eventData']['Challenge 2 Image URL'] = ""
                    # data_map['eventData']['Challenge 2 Start Date & Time'] = ""
                    # data_map['eventData']['Challenge 2 End Date & Time'] = ""
                    data_map['eventData']['Challenge 2 Slots Left'] = ""
                    # data_map['eventData']['Challenge 2 Deeplink'] = ""

                    #### indexs of combined data table
                    # 0 campaign_id,
                    # 1 campaign_name,
                    # 2 campaign_description,
                    # 3 campaign_tentative_payout_text,
                    # 4 campaign_image,
                    # 5 campaign_start_time,
                    # 6 campaign_end_time,
                    # 7 campaign_available_slots_count,
                    # 8 campaign_crm_invitation_link,
                    # 9 campaign_type,
                    # 10 reward_ladder_code,
                    # 11 reward_ladder_level,
                    # 12 campaign_invite_link,
                    # 13 challenge_id,
                    # 14 challenge_name,
                    # 15 challenge_image,
                    # 16 challenge_end_time,
                    # 17 challenge_crm_invitation_link,
                    # 18 challenge_link

                    if combinedCampaignData1 is not None:
                        if combinedCampaignData1[0] is not None:
                            data_map['eventData']['Collab 1 Id'] = combinedCampaignData1[0]
                        if combinedCampaignData1[1] is not None:
                            data_map['eventData']['Collab 1 Name'] = combinedCampaignData1[1]
                        if combinedCampaignData1[2] is not None:
                            data_map['eventData']['Collab 1 Description'] = combinedCampaignData1[2]
                        if combinedCampaignData1[3] is not None:
                            data_map['eventData']['Collab 1 Tentative Payout'] = combinedCampaignData1[3]
                        if combinedCampaignData1[4] is not None:
                            data_map['eventData']['Collab 1 Image URL'] = combinedCampaignData1[4]
                        if combinedCampaignData1[5] is not None:
                            data_map['eventData']['Collab 1 Start Date & Time'] = str(combinedCampaignData1[5])
                        if combinedCampaignData1[6] is not None:
                            data_map['eventData']['Collab 1 End Date & Time'] = str(combinedCampaignData1[6])
                        if combinedCampaignData1[7] is not None:
                            data_map['eventData']['Collab 1 Slots Left'] = combinedCampaignData1[7]
                        if combinedCampaignData1[8] is not None:
                            data_map['eventData']['Collab 1 Deeplink'] = combinedCampaignData1[8]

                    if combinedCampaignData2 is not None:
                        if combinedCampaignData2[0] is not None:
                            data_map['eventData']['Collab 2 Id'] = combinedCampaignData2[0]
                        if combinedCampaignData2[1] is not None:
                            data_map['eventData']['Collab 2 Name'] = combinedCampaignData2[1]
                        if combinedCampaignData2[2] is not None:
                            data_map['eventData']['Collab 2 Description'] = combinedCampaignData2[2]
                        if combinedCampaignData2[3] is not None:
                            data_map['eventData']['Collab 2 Tentative Payout'] = combinedCampaignData2[3]
                        if combinedCampaignData2[4] is not None:
                            data_map['eventData']['Collab 2 Image URL'] = combinedCampaignData2[4]
                        if combinedCampaignData2[5] is not None:
                            data_map['eventData']['Collab 2 Start Date & Time'] = str(combinedCampaignData2[5])
                        if combinedCampaignData2[6] is not None:
                            data_map['eventData']['Collab 2 End Date & Time'] = str(combinedCampaignData2[6])
                        if combinedCampaignData2[7] is not None:
                            data_map['eventData']['Collab 2 Slots Left'] = combinedCampaignData2[7]
                        if combinedCampaignData2[8] is not None:
                            data_map['eventData']['Collab 2 Deeplink'] = combinedCampaignData2[8]

                    if combinedChallengeData1 is not None:
                        if combinedChallengeData1[13] is not None and combinedChallengeData1[13] > 0:
                            data_map['eventData']['Challenge 1 Id'] = combinedChallengeData1[13]
                        if combinedChallengeData1[14] is not None:
                            data_map['eventData']['Challenge 1 Name'] = combinedChallengeData1[14]
                        if combinedChallengeData1[3] is not None:
                            data_map['eventData']['Challenge 1 Tentative Payout'] = combinedChallengeData1[3]
                        if combinedChallengeData1[15] is not None:
                            data_map['eventData']['Challenge 1 Image URL'] = combinedChallengeData1[15]
                        if combinedChallengeData1[5] is not None:
                            data_map['eventData']['Challenge 1 Start Date & Time'] = str(combinedChallengeData1[5])
                        # if challenge end date then use it else fallback on campaign end date
                        if combinedChallengeData1[16] is not None:
                            data_map['eventData']['Challenge 1 End Date & Time'] = str(combinedChallengeData1[16])
                        else:
                            if combinedChallengeData1[6] is not None:
                                data_map['eventData']['Challenge 1 End Date & Time'] = str(combinedChallengeData1[6])
                        if combinedChallengeData1[7] is not None:
                            data_map['eventData']['Challenge 1 Slots Left'] = combinedChallengeData1[7]
                        if combinedChallengeData1[17] is not None:
                            data_map['eventData']['Challenge 1 Deeplink'] = combinedChallengeData1[17]

                    if combinedChallengeData2 is not None:
                        if combinedChallengeData2[13] is not None and combinedChallengeData2[13] > 0:
                            data_map['eventData']['Challenge 2 Id'] = combinedChallengeData2[13]
                        if combinedChallengeData2[14] is not None:
                            data_map['eventData']['Challenge 2 Name'] = combinedChallengeData2[14]
                        if combinedChallengeData2[3] is not None:
                            data_map['eventData']['Challenge 2 Tentative Payout'] = combinedChallengeData2[3]
                        if combinedChallengeData2[15] is not None:
                            data_map['eventData']['Challenge 2 Image URL'] = combinedChallengeData2[15]
                        if combinedChallengeData2[5] is not None:
                            data_map['eventData']['Challenge 2 Start Date & Time'] = str(combinedChallengeData2[5])
                        # if challenge end date then use it else fallback on campaign end date
                        if combinedChallengeData2[16] is not None:
                            data_map['eventData']['Challenge 2 End Date & Time'] = str(combinedChallengeData2[16])
                        else:
                            if combinedChallengeData2[6] is not None:
                                data_map['eventData']['Challenge 2 End Date & Time'] = str(combinedChallengeData2[6])
                        if combinedChallengeData2[7] is not None:
                            data_map['eventData']['Challenge 2 Slots Left'] = combinedChallengeData2[7]
                        if combinedChallengeData2[17] is not None:
                            data_map['eventData']['Challenge 2 Deeplink'] = combinedChallengeData2[17]

                    url = 'http://tpevents.bulbul.tv/webengage/event/'
                    executor.submit(send_webengage_event, url, data_map)
        else:
            fetch_next = False

        # add more correct check-> check if previous loop cpid and next loop cpid is same or not here  if same then break loop
        if eligible_users_query_param_cp_id == old_eligible_users_query_param_cp_id:
            fetch_next = False


with DAG(
        dag_id='crm_recommendation_invitation',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval=None, # Set to None to run the DAG only once
) as dag:
    task = PythonOperator(
        task_id='crm_recommendation_invitation',
        python_callable=crm_recommendation_invitation
    )
    task
