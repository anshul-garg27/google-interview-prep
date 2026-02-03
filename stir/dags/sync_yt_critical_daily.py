from airflow import DAG
from airflow.hooks.base import BaseHook
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
from slack_connection import SlackNotifier


def batchify(lst, max_len=50):
    result = []
    for i in range(0, len(lst), max_len):
        result.append(lst[i:i + max_len])
    return result


def create_scrape_request_log():
    import os
    import json
    import requests

    os.environ["no_proxy"] = "*"
    handles = ['UC36yfZqAOteL4irbOSIAWHA', 'UC5qhqSIcahBKKMEdLPFqeVA', 'UC6mZXPjbDrjQFzXbiBHoaOA',
               'UC6tQtwAuADVKuxIRKfwMZfw', 'UC6wFD2VyMe43gAWN6176q5w', 'UC84q_FHnZBcZJH9G7ZjonxA',
               'UC9kbfpWblPuvjBvg5VeJ8yA', 'UCak8KQ1SJXPx7bHIINsWFEA', 'UCBepWoxn8E1TQTdxsPgLX7Q',
               'UCCpuWQs7haqrmpedX9s3xcA', 'UCCY9YBQCGx0lVu9voSvN9lw', 'UCD7U0EUBYg9YlL-zVZjkjfA',
               'UCdGNGxYoN0rPBM8CeB3bD0g', 'UCfy6FXs2068R9zo8Cfdo2mA', 'UCgnVhWJlQVP8t53Osq8b91g',
               'UCHpBLdUcA3Fuwx-aTQQOrAg', 'UCJd0jl8aKCqxF94by-R_3ZQ', 'UCKgVZn8W_DiJUtXI_40e8sw',
               'UCKVqnev2idvmUNKc2b91B8g', 'UCl4qNDVAboXyBqPRvHRtcRg', 'UClRx-1O56Z_VQgxBPQQkW3w',
               'UCMv7rbhyOwFpAGdpAfTJM7A', 'UCMxDYxxFFJkB2D-PNdeo0kg', 'UCn_YmO_n1NkHBGiSvTCnI2w',
               'UCNhMdV20V-AyjzxqD45T9vQ', 'UCNkBtAfx0b2Gs4_eMXE6qcw', 'UCO6msHL6wSfYUNX1Ug7bx-g',
               'UCOujVRoOa29CTcgtxDXHPgw', 'UCQ3g0glwFaeIBftFy6ZyHtg', 'UCR57rFpJ4rABT6voCOv954g',
               'UCrgcBTLLgV-ZbaSBPkqjvow', 'UCRiVehxaB5Q1jHmf3cLmHFg', 'UCrUxV9rsE-ivYxrgwkQhrjQ',
               'UCrYcHTAo4lCo5yPdXosWedQ', 'UCTvrRxv4ng3SzTK6yy5gUYg', 'UCUI30rMByYDIEAEUFtG_Mtg',
               'UCuJXacHnH_o6xzSqse6VZ7A', 'UCu-luPcM9OzpFRzOVAG7-cQ', 'UCUTrMET_8V4kR8xPGtmP8Aw',
               'UCV0wtCows-JFIGk_xk2trag', 'UCXNHba51JL2Jt705hhlSQJA', 'UC-xTA_bW6pyEqK3byfbwWww',
               'UCz3yLP_63OUzsBcT0n3xcxA', 'UCZKdxkPSW9eaVFxdXed7cBA', 'UCPP3etACgdUWvizcES1dJ8Q',
               'UCt4t-jeY85JegMlZ-E5UWtA', 'UCRWFSbif-RFENbBrSiez1DA', 'UCOutOIcn_oho8pyVN3Ng-Pg',
               'UC7wXt18f2iA3EDXeqAVuKng', 'UCIvaYmXn910QMdemBG3v1pQ', 'UC9CYT9gSNLevX5ey2_6CK0Q',
               'UCttspZesZIDEwwpVIgoZtWQ', 'UCl8wUKzoUVzg7U6ky1ZR8hQ', 'UCuzS3rPQAYqHcLWqOFuY0pw',
               'UCbf0XHULBkTfv2hBjaaDw9Q', 'UCv3rFzn-GHGtqzXiaq3sWNg', 'UCdF5Q5QVbYstYrTfpgUl0ZA',
               'UCHCR4UFsGwd_VcDa0-a4haw', 'UCajVjEHDoVn_AHsunUZz_EQ', 'UCESY9RzEDPsb9lRLMxAirwQ',
               'UCNvCQpcafnbW4KQ8X7oQ9kg', 'UCJ3I6MHOz5exARlTW_meOGQ', 'UCef1-8eOpJgud7szVPlZQAQ',
               'UCYPvAwZP8pZhSMW8qs7cVCw', 'UCZFMm1mMw0F81Z37aaEzTUA', 'UC6RJ7-PaXg6TIH2BzZfTV7w',
               'UCckHqySbfy5FcPP6MD_S-Yg', 'UCm7lHFkt2yB_WzL67aruVBQ', 'UCwqusr8YDwM-3mEYTDeJHzw',
               'UC3njZ48-FDxLleBYaP0SZIg', 'UCJEDFSxHHOW1PpBccdSxOTA', 'UCAjBd-r8JWfnRjfhg23nqLQ',
               'UC1JkYCOb4wKvdWreGQ0CR9g', 'UCrQHRYuJG8jmpUVALIC9Gkw', 'UCBMgVLD4TaNFiZLNB60WMgg',
               'UCUK49UvmYWYLiB7_bZFuFZQ', 'UCCgLMMp4lv7fSD2sBz1Ai6Q', 'UC90RW5ZmBBqp4r2QIQxfACA',
               'UC1vOCYEu3DfoI4fYcfWKmWA', 'UC33tp6m-j1pEmYnkYIzSdoQ', 'UCat88i6_rELqI_prwvjspRA',
               'UC8Z-VjXBtDJTvq6aqkIskPg', 'UCYlh4lH762HvHt6mmiecyWQ', 'UCi-MPuZaplEjdCTwquSMETQ',
               'UC-JFyL0zDFOsPMpuWu39rPA', 'UC2f4w_ppqHplvjiNaoTAK9w', 'UCmyKnNRH0wH-r8I-ceP-dsg',
               'UC-mMi78WJST4N5o8_i1FoXw', 'UCP0uG-mcMImgKnJz-VjJZmQ', 'UCup3etEdjyF1L3sRbU-rKLw',
               'UCwXrBBZnIh2ER4lal6WbAHw', 'UCf8w5m0YsRa8MHQ5bwSGmbw', 'UCNVkxRPqsBNejO6B9thG9Xw',
               'UCa-vioGhe2btBcZneaPonKA', 'UCU_Dzyqehuz3pVbZrOqOT3A', 'UCl-OodciBGZ0k8K8rBZGe4w',
               'UC8dnBi4WUErqYQHZ4PfsLTg', 'UCjElJyiXmQXnWmceQ1JyKrA', 'UC531MlZA5LUbeGwEN_zcppw',
               'UCZUjHLJivN0OZPC_Wj8JvkA', 'UCBdxSSyIlnwo3Aj5O0c6RWQ', 'UCnAp2J0bR9b8pM-Avp1GFOQ',
               'UC9c0UZND0vLUsMAWJAlkSqw', 'UCafYgzpyw7aIUYOLjjADu7w', 'UCUGwnDFBHY52YhgVjn-Tvww',
               'UC8Xb-hkR53QVN4GlUHi1iZg', 'UCWmwc3u2ikBG2Xfitc_26gg', 'UCskG03x3CoEW9W7s2c3IgbA',
               'UCeNP3dAy9_A32IWYZq_PZMQ', 'UCBwc2cbPpvxNCNEI2-8YrqQ', 'UCQ79XyTiFbb28NvhmFPZOiw',
               'UCBc13XYipnBIBE3Ff8QaaGg', 'UCL_3_9PhDyZXu1oexm8gtFQ', 'UCV8uDD-i5zBnwdyHYRHjlEg',
               'UC7FlLbNo66YsCEAPuJITiNg', 'UC_cZYFovNQTsKw1ZLmO_AMA', 'UCfnvGSUXKIuaCv-9E9a2UJg',
               'UCWk-7Yosyvzln9ZzJg8BvVg', 'UCgJ02VYByDODEIkOlxjg4xg', 'UC-crZTQNRzZgzyighTKF0nQ',
               'UC3_TTermAnf2fcLe-QX7BVQ', 'UCQLEbraENUGWh6p1Rv664rQ', 'UCYGZ0qW3w_dWExE3QzMkwZA',
               'UCmhe0p2-m_KURZYQNchK-QQ', 'UC68JhJpCRYYhOcO054zw_Ew', 'UCIb1wDIMrOwK3zEzi3ROuXg',
               'UCrcpw88HvKJ0skdsHniCJtQ', 'UCdOSeEq9Cs2Pco7OCn2_i5w', 'UCH7nv1A9xIrAifZJNvt7cgA',
               'UC6cxTsUnfSZrj96KNHhRTHQ', 'UCVbsFo8aCgvIRIO9RYwsQMA', 'UCQIycDaLsBpMKjOCeaKUYVg',
               'UCkXopQ3ubd-rnXnStZqCl2w', 'UCmRbHAgG2k2vDUvb3xsEunQ', 'UCI_mwTKUhicNzFrhm33MzBQ',
               'UCJEtD6JDgNKxcPJPBOAf-Tg', 'UCqfX_f4rBbPhq_KESE4saCA', 'UCeJWZgSMlzqYEDytDnvzHnw',
               'UCMX41X1am8oYxT336dqk4sA', 'UC3C6_1ETXfE807LltDbKYxg', 'UCkNL_TQio--h85-14lUVY3A',
               'UCiAH2s_M6nPfGZk-PpfyPkg', 'UC9IQhNgS43lmUdU-KdGyb0A', 'UCkSllUpc8cedRe85JkqLiBQ',
               'UCVpZYVnk7OOnz_BlHSwRKWA', 'UC3prwMn9aU2z5Y158ZdGyyA', 'UCSaf-7p3J_N-02p7jHzm5tA',
               'UCuyRsHZILrU7ZDIAbGASHdA', 'UChWtJey46brNr7qHQpN6KLQ', 'UChftTVI0QJmyXkajQYt2tiQ',
               'UC3uJIdRFTGgLWrUziaHbzrg', 'UCSWSOS6YXUbNMzTH-tV7Pfw', 'UCw5TLrz3qADabwezTEcOmgQ',
               'UCLnuTrdp3MBzo9ZPig5szbQ', 'UCYL4XTsXXffnnM_RveJPilQ', 'UCv-Pg0rAiiEhKYFVKg8E7JQ',
               'UCwcH94eG9YYAn1G2nfw3j-w', 'UCx8Z14PpntdaxCt2hakbQLQ', 'UCAUNFgpgVisKPL3yq_-Nj-Q',
               'UCSKgOW8Pg_eZymYJyJc432g', 'UC7T5n_yWsI3T8EF7Auf6hfQ', 'UC6ZMkiLxLEtz2P3x57nLfhg',
               'UCAF1vryFYnaYhbpvL351ntA', 'UCIsotCCdTjStPCLyD9iBCCQ', 'UCuN4A3GCUq5-0wJDSiJoxRQ',
               'UC7IMq6lLHbptAnSucW1pClA', 'UCXyq6UjvT4dWjMOOiKuBncA', 'UCgKscN1f8lljHfGpoDQiaWQ',
               'UCCsDzd6T8aj4U7h1qEa3dgw', 'UCaAzhBJP6wnWL-s0SwD0vOA', 'UCurIxpcdDLU8sgHKN2TXJHw',
               'UCx2JNgWJPPTBODN-U56O-tg', 'UCa25AzX_P5GE4QW8NYkjzHA', 'UC2QBCIyo_FbTlm0Rfh_6wkw',
               'UCF8EOH5Ty0kk76qpK4-JNxQ', 'UC8xgalHz7T48nKqblwMTpuQ', 'UCj_eYKbmUXvsbDssD3pQBXA',
               'UCdxbhKxr8pyWTx1ExCSmJRw', 'UCDTM4H2YzteOzGMzoGSJwAQ', 'UCoaH2UtB1PsV7av17woV1BA',
               'UCI_6AcJI1sKexCLb2NAYTOQ', 'UCam4pU1NFraGs5Ng_SS9h9Q', 'UCTVO-4rIuzg_rCQP37jjSTA',
               'UCGNXTzswk02WtPilaK9AH1A', 'UCZqOZwn8I3luPzLO8lMkOBg', 'UCMbwFLJf2buLbHz-R2H_-aw',
               'UCT2TvAK2rP1yA4FGZMh_HWQ', 'UCMDPvpD-R1_isb1o0ua15ig', 'UCuHH7aX-2Djh9JnEUNwWOtw',
               'UCtLpPnn0bg52IPFhzMwkvBw', 'UC0tLWfVWGXrz3_LPTuU88VA', 'UCwPac8AUZHTcUXfnPQAtvFA',
               'UCWjwhAbIFpF9v0B1pCMyjrw', 'UCfqg37Xt-ezCSArBy0aTGMQ', 'UCdILIiuXr2UQUhc5SFVbKtw',
               'UCZT_OWY9abyO3V9Rmdvh02w', 'UCU5pQWmxvxqRP5u-SdoV8rQ', 'UCotI-SqRXnkAZX4bMqlRNjw',
               'UCkJZQddO3XTCdcLjN019FjA', 'UC7eZNThpjgXe5H0QwBN-XsQ', 'UCmDgY7fZKvgatlO0Tib2ggQ',
               'UCUld3m_rLi37zHm5XhbELPQ', 'UCdkyHgvtTdwEiGwsFTcdHVQ', 'UC3K_FQUZp_8heTLcQuQk-KA',
               'UCMsEZ4DiNPPHx-Fkc_ECuBg', 'UC3cOaBryqYJesL57W5fTE1A', 'UC9DDmuRH5lOngGe8po3Hk6Q',
               'UC7T-KouKpG3xwSH0WF6QJeg', 'UCNoV4HixKd8AoVt5paswr-w', 'UCx2HcmpB-UZGkMXOCJ4QIVA',
               'UCh29s-RNWI6z63670hUv71Q', 'UCcjwEAaBGu1yc89QE3DC1Tw', 'UCH_RhBu7gVUDictmhfL6DTg',
               'UCA1EgVvGBUzGviT6IhVCTVg', 'UCFeIPY_Xa9hddLEIQHCJVpg', 'UC0KNDUO6xcWY1vwNPFSsx4Q',
               'UC9YX5x_VU8gfe0Oui0TaLJg', 'UCCgbPq9LSkWX6XcbCpZfKfQ', 'UCs_lbrYHahO5gv7B27kkqBA',
               'UCZuM_9OulAMXu_0o0GMsjig', 'UCG9nVSLp4nQlW79sY5ihRrg', 'UCKrrATalRpJ-H5ltP788DkA',
               'UCs7hbbgJFer2AcKtwAV5Bwg', 'UCRXiA3h1no_PFkb1JCP0yMA', 'UCRSvEADlY-caz3sfDNwvR1A',
               'UC0HVR9T6oFS3veefhGCGEsw', 'UC659HgEmH8-mBpnVwZazuNQ', 'UCk1C8yXuiOvgXDgbz7hdf2w',
               'UCmuCmtaKmn5ouwaOYl3phlw', 'UC-DuRqsBQOEk_5o1q4Ze-Fg', 'UCqCDO1Ug87LHHG6PHgfPgGw',
               'UCfo0AlxPIc6LLIF8PWhVbVg', 'UCyzF_Sq8Sgc3iBl3poNi3xQ', 'UCV-oK1xe9nSAUUrNI-TrS0g',
               'UCx922skTDrcpo96rXUhb-AQ', 'UCgHSZo9QJoJvnL-0CZZDNyg', 'UCbOe5GzWdxrbZYiHa-wF98w',
               'UCMk9Tdc-d1BIcAFaSppiVkw', 'UCmphdqZNmqL72WJ2uyiNw5w', 'UCvxi7_X1VSaMr_osIamYNaw',
               'UCWJedqvdP2gB2f3qaJCO05g', 'UCu_q3VCbMiq5bYaBql_NmSQ', 'UCjWs7BxyjO5SLqevxSmp4vQ',
               'UCHcRjWV5XOGdYkM7ruvpgxw', 'UCN7B-QD0Qgn2boVH5Q0pOWg', 'UC7pluR6rB5KZIbN2IxamzxQ',
               'UC16niRr50-MSBwiO3YDb3RA', 'UCkul0EjxFOR2hS_0tOokGag', 'UCzR7770PbrKcG9OYGzBep9w',
               'UCQsob4fGjHWhYHW0OLb6rew', 'UCupvZG-5ko_eiXAupbDfxWw', 'UCVZ57OkKPAuRJ_wA_Rt4XFg',
               'UCXpIVRiyfsMB2aXXn2yq6LA', 'UCA7mhLMqut0YGwRcEP-qvEg', 'UCIRAYFbJmrP--jyrC9MAIWQ',
               'UCD3CdwT8lTCe5ZGHbUBxmWA', 'UC500dYMU9OMJdJKWRqGlhog', 'UCoUxsWakJucWg46KW5RsvPw',
               'UCAnVEUKmO_0BFbaneK4D5CA', 'UCx5e1u7BX0aKwEj3sdYXdXg', 'UC4Gez3KnLoxurTz5W3U6SuA',
               'UChYpOn9tfcJkfjq7e0KaAhw', 'UCNY_BMBmlEzCZHIeaTUqZnw', 'UCPPIsrNlEkaFQBk-4uNkOaw',
               'UCUmhg3QldZiYoO9mKGpumBA', 'UCqJkAAmi4QKCPCF62r_-BhQ', 'UCW-7nV1aZu0PKqj4c8SGP4A',
               'UCmBLIcqOLWwrrAJ-jaCttqw', 'UCICJyUwBYgUOE35F8klEoYw', 'UCChqsCRFePrP2X897iQkyAA',
               'UC5MHSwQ2menwYF2DqKoJsZg', 'UCwh_2Tw0JzkOEv74yvc20IQ', 'UC2ut_DrUZvO0BzUxZ_g9fXQ',
               'UCjD_oFXbVfCDhGEsHunI-hg', 'UCLIMdoS-InL1-GIBUKSvydg', 'UCWCEYVwSqr7Epo6sSCfUgiw',
               'UCrdPiSPVW0rtRsI002BX8iw', 'UC-sc9TAs4scvSUMpixFhcLA', 'UCQJTzCJD3HaB7NOc9gnd3bg',
               'UCytSP0M0Jdnw6qIy3Y-nTig', 'UC_5Xhfmz-KE3VyVuhwJ57hw', 'UCdegm7Y2AePJhkkmWCyYEwg',
               'UC-AY02VQo-F_trSeTKUu_QQ', 'UCl8J8ntI5SJkwf4V32ap-4Q', 'UCycyxZMoPwg9cuRDMyQE7PQ',
               'UCVGiG3GcphUyMzNXg2iMN3A', 'UCvwDYLQmuPo8tYk7GZ6c5pw', 'UC94m1fN8ypwZpZWs4E5mEcg',
               'UCJFOER35ggIWwsXh2ZDnqyg', 'UCmk6ZFMy1CT80orXca4tKew', 'UC6y7RE7foidnDJk05ERl7eA',
               'UCPgLNge0xqQHWM5B5EFH9Cg', 'UCJ04m06GHw3DUDQG6chCAUA', 'UCjkLYVF8Up8zt9ZQNLpR_TQ',
               'UC_gUM8rL-Lrg6O3adPW9K1g', 'UCe2JAC5FUfbxLCfAvBWmNJA', 'UCjmjWp38PCg15Z5ZS-tmpfw',
               'UCzeTWUnXO5jcxd6snccRJ6g', 'UCBNkm8o5LiEVLxO8w0p2sfQ', 'UC_DmCvOP5Q_eBMRDvqqRXjg',
               'UCc4sOlYxxoXNba2CYCRQJGw', 'UC0Hh-U1C8MjAIfvzzBUbr4A', 'UCQYZICqVoX2KKGxVOHZnfbg',
               'UCmc8RQmCy5F01IzbHaTnBSQ', 'UCPbUTEjdU04ChIZfi53jczw', 'UCwd-Oa3SpLAwZWqZJBQdEpw',
               'UCuhKTYX1J1M5XUX0cMbG3vQ', 'UCT0EeR8SMCItnnOgmLfEjoQ', 'UCXcygMj8Ezhq2EJ19QTzdww',
               'UCSbvkOATQLgqC_BCxcz7ZZA', 'UCCKjHsAIxvjtWG8KOcLuG8Q', 'UC2ZMJOSJNedeE648KDvIAbA',
               'UCwGHTKlxg__AW5Jae6sf1rw', 'UCNJcSUSzUeFm8W9P7UUlSeQ', 'UC8kB1Yfsq0WxJ9AFSkj7VHg',
               'UCAdq6Ae1SlV23DXqhmAw6UA']
    batches = batchify(handles, max_len=40)

    refresh_profiles_flow = 'refresh_yt_profiles'

    beat_connection = BaseHook.get_connection("beat")
    base_url = beat_connection.host

    refresh_profiles_url = base_url + "/scrape_request_log/flow/%s" % refresh_profiles_flow

    for batch in batches:
        refresh_profile_scl = {"flow": refresh_profiles_flow, "platform": "YOUTUBE", "params": {"channel_ids": batch}}
        requests.post(url=refresh_profiles_url, data=json.dumps(refresh_profile_scl))


with DAG(
        dag_id='sync_yt_critical_daily',
        start_date=days_ago(1),
        catchup=False,
        max_active_runs=1,
        concurrency=1,
        on_failure_callback = SlackNotifier.slack_fail_alert,
        on_success_callback = SlackNotifier.slack_fail_alert,
        schedule_interval="0 0 * * *",
) as dag:
    task = PythonOperator(
        task_id='sync_yt_critical_daily_scl',
        python_callable=create_scrape_request_log
    )
    task
