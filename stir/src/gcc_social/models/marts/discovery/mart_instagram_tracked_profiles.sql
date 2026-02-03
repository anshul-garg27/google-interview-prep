{{ config(materialized = 'table', tags=["account_tracker_stats"], order_by='handle') }}
with
     handles_attempted as (select JSONExtractString(params, 'handle') handle,
                                  max(scraped_at)                     last_attempt
                           from dbt.stg_beat_scrape_request_log
                           where status IN ('FAILED', 'COMPLETE')
                             and ((data NOT LIKE '(204%'
                               and data NOT LIKE '(409%'
                               and data NOT LIKE '%NoneType%'
                               and data NOT LIKE '(429%') or data is null)
                             and flow = 'refresh_profile_by_handle'
                           group by handle
                           order by last_attempt desc),
     app_handles as (select sa.handle handle,
                            'app'     source,
                            1         priority
                     from dbt.stg_identity_social_account sa
                              left join handles_attempted ia on ia.handle = sa.handle
                     where sa.created_on >= now() - INTERVAL 7 DAY
                       and sa.platform = 'INSTAGRAM'
                       and (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 MONTH)
                     group by sa.handle),
     gcc_collection_handles as (select sia.handle                handle,
                                       'saas-profile-collection' source,
                                       5                         priority
                                from dbt.stg_coffee_profile_collection_item pci
                                         left join dbt.stg_coffee_profile_collection pc on pc.id = pci.profile_collection_id
                                         left join dbt.stg_coffee_view_instagram_account_lite sia
                                                   on sia.id = pci.platform_account_code
                                         left join handles_attempted ia on ia.handle = sia.handle
                                where profile_social_id != ''
                                  and pc.source = 'SAAS'
                                  and pci.platform = 'INSTAGRAM'
                                  and pci.updated_at >= now() - INTERVAL 30 DAY
                                  and (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 DAY)
                                group by handle),
     saas_collection_handles as (select sia.handle        handle,
                                        'saas-collection' source,
                                        4                 priority
                                 from dbt.stg_coffee_profile_collection_item pci
                                          left join dbt.stg_coffee_profile_collection pc on pc.id = pci.profile_collection_id
                                          left join dbt.stg_coffee_view_instagram_account_lite sia
                                                    on sia.id = pci.platform_account_code
                                          left join handles_attempted ia on ia.handle = sia.handle
                                 where profile_social_id != ''
                                   and pc.source = 'SAAS-AT'
                                   and pci.platform = 'INSTAGRAM'
                                   and (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 DAY)
                                 group by handle),
     leaderboard_handles as (select l.handle      handle,
                                    'leaderboard' source,
                                    4             priority
                             from dbt.mart_instagram_leaderboard l
                                      left join handles_attempted ia on ia.handle = l.handle
                             where month >= date('2023-02-01')
                               and (
                                 followers_rank < 1000 or
                                 followers_change_rank < 1000 or
                                 views_rank < 1000 or
                                 plays_rank < 1000 or
                                 followers_rank_by_cat < 1000 or
                                 followers_change_rank_by_cat < 1000 or
                                 views_rank_by_cat < 1000 or
                                 plays_rank_by_cat < 1000 or
                                 followers_rank_by_lang < 1000 or
                                 views_rank_by_lang < 1000 or
                                 followers_change_rank_by_lang < 1000 or
                                 plays_rank_by_lang < 1000 or
                                 views_rank_by_cat_lang < 1000 or
                                 followers_rank_by_cat_lang < 1000 or
                                 followers_change_rank_by_cat_lang < 1000 or
                                 plays_rank_by_cat_lang < 1000
                                 )
                               and (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 WEEK)),
     saas_handles as (select bia.handle handle,
                             'saas'     source,
                             5          priority
                      from dbt.mart_instagram_account bia
                               left join handles_attempted ia on ia.handle = bia.handle
                      where country = 'IN'
                        and (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 2 WEEK)
                      group by handle),
     cm_handles as (select wia.handle handle,
                           'cm'       source,
                           6          priority
                    from dbt.stg_coffee_campaign_profiles cp
                             left join dbt.stg_coffee_view_instagram_account_lite wia
                                       on wia.id = cp.platform_account_id and cp.platform = 'INSTAGRAM'
                             left join handles_attempted ia on ia.handle = wia.handle
                    where cp.platform_account_id > 0 and cp.gcc_user_account_id > 0
                      and (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 MONTH)
                    group by handle),
     rest_2k as (select bia.handle handle,
                        'rest_2k'  source,
                        7          priority
                 from dbt.mart_instagram_account bia
                          left join handles_attempted ia on ia.handle = bia.handle
                 where (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 MONTH)
                 group by handle
                 having max(bia.followers) > 2000),
     rest as (select bia.handle handle,
                     'rest'     source,
                     8          priority
              from dbt.mart_instagram_account bia
                       left join handles_attempted ia on ia.handle = bia.handle
              where (ia.last_attempt is null or ia.last_attempt <= now() - INTERVAL 1 MONTH)
              group by handle),
     all_handles as (select handle, source, priority
                     from app_handles
                     union all
                     select handle, source, priority
                     from gcc_collection_handles
                     union all
                     select handle, source, priority
                     from saas_collection_handles
                     union all
                     select handle, source, priority
                     from leaderboard_handles
                     union all
                     select handle, source, priority
                     from cm_handles
                     union all
                     select handle, source, priority
                     from saas_handles
                     union all
                     select handle, source, priority
                     from rest_2k
                     union all
                     select handle, source, priority
                     from rest),
    combined as (
                select
                    handle,
                    argMin(source, priority) source,
                    min(priority) p
                from all_handles
                group by handle
            )
select ifNull(handle, '') handle, source from combined order by p asc
