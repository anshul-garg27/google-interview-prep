# Complete PR List - All Repositories

> **Purpose:** Complete list of all PRs (Open + Merged) for deep analysis
> **Last Updated:** February 2026
> **Total PRs:** 1000+ across 5 repos

---

## Table of Contents
1. [Summary Statistics](#summary-statistics)
2. [audit-api-logs-gcs-sink (Bullet 1, 10: Kafka Audit & Multi-Region)](#audit-api-logs-gcs-sink)
3. [dv-api-common-libraries (Bullet 2: Common Library JAR)](#dv-api-common-libraries)
4. [inventory-status-srv (Bullet 5, 6, 8: OpenAPI, DC Inventory, Auth)](#inventory-status-srv)
5. [cp-nrti-apis (Bullet 3, 4, 9: DSD, Spring Boot 3, Observability)](#cp-nrti-apis)
6. [inventory-events-srv (Bullet 7: Transaction Event History)](#inventory-events-srv)
7. [Cross-Repo Audit Logging Integration](#cross-repo-audit-logging-integration)

---

## Summary Statistics

| Repository | Total PRs | Merged | Open | Closed |
|------------|-----------|--------|------|--------|
| `audit-api-logs-gcs-sink` | 87 | 72 | 0 | 15 |
| `dv-api-common-libraries` | 16 | 8 | 3 | 5 |
| `inventory-status-srv` | 359+ | 300+ | 25+ | 30+ |
| `cp-nrti-apis` | 1600+ | 1400+ | 50+ | 150+ |
| `inventory-events-srv` | 143 | 100+ | 15+ | 25+ |

---

## audit-api-logs-gcs-sink

**Related Bullets:** 1 (Kafka Audit Logging), 10 (Multi-Region Architecture)

### MERGED PRs (72 total)

| PR# | Title | Branch | Created | Merged |
|-----|-------|--------|---------|--------|
| #87 | Vulnerability update and Start up probe changes | feature/start-up-probe-main | 2026-01-27 | 2026-01-29 |
| #86 | Adding dummy commit | feature/jan-dummy-commit-main | 2026-01-14 | 2026-01-14 |
| #85 | **GCS Sink Changes to support different site ids** | feature/single-deployment-main | 2025-11-19 | 2025-12-02 |
| #83 | CCM Changes for country code to site id in non-prod | feature/add-wm-site-id-config-main | 2025-11-03 | 2025-11-03 |
| #81 | Updating the connector version | feature/upgrade-connector-version | 2025-06-18 | 2025-06-18 |
| #80 | Adding dummy commit | feature/dummy-commit-main | 2025-06-17 | 2025-06-17 |
| #77 | Updating resolution path to have /* as override path | feature/update-resolution-path-ccm | 2025-06-13 | 2025-06-13 |
| #75 | Updating the data type as DECIMAL for timeout multiplier | feature/update-data-type-decimal-ccm | 2025-06-11 | 2025-06-11 |
| #73 | Adding all the worker and connector property in NON-PROD CCM | feature/add-ccm-properties | 2025-06-11 | 2025-06-11 |
| #72 | Updating README.md with local set up | feature/local-set-up-stage | 2025-06-09 | 2025-06-18 |
| #70 | **[CHG3235845] : Upgrade connector version. Remove custom header. Include retry policy** | release/connector-upgrade-retry-prod | 2025-06-02 | 2025-06-09 |
| #69 | Upgrade connector version. Remove custom header. Include retry policy | feature/connector-upgrade-retry-stage | 2025-05-30 | 2025-05-30 |
| #68 | Adding tunr names in CCM | feature/tunr-name-ccm | 2025-05-27 | 2025-05-27 |
| #66 | Adding names for config overrides | feature/include-stage-names-ccm | 2025-05-27 | 2025-05-27 |
| #65 | Adding gcs sink intl app key in sr.yaml | feature/sr-changes-intl | 2025-05-26 | 2025-05-26 |
| #64 | intl env changes in non prod ccm | feature/ccm-changes-intl-non-prod | 2025-05-22 | 2025-05-22 |
| #63 | Renaming stage as release/stage and updating connector configs | feature/stage-rename | 2025-05-19 | 2025-05-19 |
| #61 | **Increasing heap max** | hotfix/config-update-10 | 2025-05-15 | 2025-05-16 |
| #60 | Adding header converter at worker | hotfix/config-update-9 | 2025-05-15 | 2025-05-15 |
| #59 | **Adding the filter and custom header converter** | hotfix/config-update-8 | 2025-05-15 | 2025-05-15 |
| #58 | Removing unit test changes and filter class | hotfix/update-config-7 | 2025-05-15 | 2025-05-15 |
| #57 | **Overriding KAFKA_HEAP_OPTS** | feature/kafka-heap-opts | 2025-05-15 | 2025-05-15 |
| #55 | update latest version | hotfix/sinkissue | 2025-05-15 | 2025-05-15 |
| #53 | update kc config | hotfix/config-6 | 2025-05-14 | 2025-05-14 |
| #52 | **Setting error tolerance to none** | hotfix/update-config-4 | 2025-05-14 | 2025-05-14 |
| #51 | Setting error tolerance to none | hotfix/update-config-3 | 2025-05-14 | 2025-05-14 |
| #50 | Reducing the tasks to 1 | hotfix/update-configs-2 | 2025-05-14 | 2025-05-14 |
| #49 | Hotfix/update configs 1 | hotfix/update-configs-1 | 2025-05-14 | 2025-05-14 |
| #47 | **Increasing the poll timeout** | hotfix/increase-poll-interval | 2025-05-14 | 2025-05-14 |
| #46 | Decreasing the session time out | hotfix/decrease-session-time-out | 2025-05-14 | 2025-05-14 |
| #43 | Increasing the session time out | hotfix/increse-session-time-out | 2025-05-14 | 2025-05-14 |
| #42 | **Removing scaling on lag. Updating the heartbeat** | hotfix/remove-autoscaling-update-heartbeat | 2025-05-14 | 2025-05-14 |
| #41 | Updating the heartbeat | hotfix/update-heartbeat-prod | 2025-05-14 | 2025-05-14 |
| #40 | Updating CPU and ephemeral storage | hotfix/update-cpu-prod | 2025-05-13 | 2025-05-13 |
| #39 | Changing to standard consumer configs | hotfix/update-consumer-configs | 2025-05-13 | 2025-05-13 |
| #38 | **Adding consumer poll config** | feature/consumer-poll-config | 2025-05-13 | 2025-05-13 |
| #37 | Adding loggers | hotfix/enable-log-update-java | 2025-05-13 | 2025-05-13 |
| #35 | **Adding try catch block in filter** | hotfix/fix-2025-05-12-05-47-32 | 2025-05-12 | 2025-05-13 |
| #34 | Increasing the min and max replicas | hotfix/fix-2025-05-12-05-47-32 | 2025-05-12 | 2025-05-12 |
| #32 | Adding metadata as parent of labels in helm values | fix-2025-05-12-05-47-32 | 2025-05-12 | 2025-05-12 |
| #30 | **adding prod-eus2 env** | fix-2025-05-12-05-47-32 | 2025-05-12 | 2025-05-12 |
| #29 | **Adding autoscaling profile and related changes** | feature/autoscaling-flush-changes-stage | 2025-05-05 | 2025-05-07 |
| #28 | **CHG3180819 - Updating flush count, max.poll.interval.ms and offset assignor** | release/consumer-poll-config-changes-prod | 2025-05-02 | 2025-05-12 |
| #27 | Updating flush count, max.poll.interval.ms and offset assignor | feature/consumer-poll-config-changes-stage | 2025-04-28 | 2025-05-02 |
| #26 | Adding insights.yml | feature/update-bucket-name-stage | 2025-04-03 | 2025-04-03 |
| #25 | CRQ : CHG3041101 Dummy commit | release/gcs-sink-us-prod | 2025-04-03 | 2025-04-03 |
| #24 | Dummy commit | feature/update-bucket-name-stage | 2025-04-03 | 2025-04-03 |
| #23 | Updating KCQL of prod SCUS | release/gcs-sink-us-prod | 2025-04-03 | 2025-04-03 |
| #22 | Updating KCQL of prod SCUS | feature/update-bucket-name-stage | 2025-04-03 | 2025-04-03 |
| #21 | Updating the topic names | release/gcs-sink-us-prod | 2025-04-03 | 2025-04-03 |
| #20 | Updating the topic names | feature/update-bucket-name-stage | 2025-04-03 | 2025-04-03 |
| #19 | [CRQ:CHG3041101] Dummy commit to reflect non prod ccm and sr changes | release/gcs-sink-us-prod | 2025-04-02 | 2025-04-02 |
| #18 | Dummy commit to reflect non prod ccm and sr changes | feature/update-bucket-name-stage | 2025-04-02 | 2025-04-02 |
| #15 | Updating the notification email, including main branch in stage deployment | feature/update-bucket-name-stage | 2025-04-02 | 2025-04-02 |
| #14 | Updating the schema name folder | feature/update-bucket-name-stage | 2025-04-02 | 2025-04-02 |
| #13 | Removing ephemeral storage | feature/update-bucket-name-stage | 2025-04-01 | 2025-04-01 |
| #12 | **[CRQ:CHG3041101] GCS Sink changes for US records** | release/gcs-sink-us-prod | 2025-03-28 | 2025-04-02 |
| #11 | Updating the project id of new buckets | feature/update-bucket-name-stage | 2025-03-28 | 2025-03-28 |
| #10 | Storing in new test bucket | feature/update-bucket-name-stage | 2025-03-27 | 2025-03-27 |
| #9 | **Prod config, renaming class, updating field** | prod-config-only-stage | 2025-02-26 | 2025-03-25 |
| #5 | Updating the repo name - audit-api-logs-gcs-sink | repo-name-change-main | 2025-02-24 | 2025-02-26 |
| #4 | Updating the repo name as audit-api-logs-gcs-sink | repo-name-change-stage | 2025-02-24 | 2025-02-25 |
| #3 | **Deployment changes main** | deployment-changes-main | 2025-01-22 | 2025-02-20 |
| #1 | **Updating readme file, kafka topic and gcs bucket name** | deployment-changes-main | 2025-01-07 | 2025-01-07 |

### CLOSED PRs (Not Merged)

| PR# | Title | Branch | Created | Status |
|-----|-------|--------|---------|--------|
| #82 | Feature/single deployment stage | feature/single-deployment-stage | 2025-10-29 | CLOSED |
| #79 | Removing the envProfile from pathElements | feature/remove-envProfile-ccm | 2025-06-13 | CLOSED |
| #71 | Increasing the flush record count | feature/increase-record-count-stage | 2025-06-04 | CLOSED |
| #62 | GCP connection retry | hotfix/config-update-11 | 2025-05-16 | CLOSED |
| #56 | Adding connector properties | hotfix/update-config-5 | 2025-05-15 | CLOSED |
| #54 | Hotfix/config 6 | hotfix/config-6 | 2025-05-15 | CLOSED |
| #48 | Updating consumer configs | hotfix/update-configs-1 | 2025-05-14 | CLOSED |
| #45 | Decrease session time out | hotfix/increse-session-time-out | 2025-05-14 | CLOSED |
| #44 | HeartBeat and session Timeout | feature/update-heartbeat-and-session-timeout | 2025-05-14 | CLOSED |
| #36 | adding test | feature/adding-logger-stage | 2025-05-13 | CLOSED |
| #33 | Increasing min and max replicas | fix-2025-05-12-05-47-32 | 2025-05-12 | CLOSED |
| #6 | Prod configs changes | prod-configs-stage | 2025-02-25 | CLOSED |

### Key Interview PRs (Highlighted)
- **#1** - Initial setup (architecture decisions)
- **#9** - Prod config, class design
- **#12** - Multi-region US deployment [CRQ]
- **#28** - Performance tuning (flush count, poll interval)
- **#29** - Autoscaling profile
- **#35-61** - Hotfix series (great debugging story)
- **#70** - Retry policy, connector upgrade [CRQ]
- **#85** - Site-based routing

---

## dv-api-common-libraries

**Related Bullets:** 2 (Common Library JAR)

### MERGED PRs (8 total)

| PR# | Title | Branch | Created | Merged |
|-----|-------|--------|---------|--------|
| #15 | Version update for gtp bom | release/gtp-bom-changes-java17 | 2025-09-19 | 2025-09-19 |
| #12 | Create CODEOWNERS | feature/adding-codeowners | 2025-04-11 | 2025-04-11 |
| #11 | **enable endpoints with regex** | release/DSIDVL-42063 | 2025-04-08 | 2025-04-11 |
| #6 | **Develop** | develop | 2025-03-20 | 2025-03-20 |
| #5 | Update pom.xml | develop | 2025-03-20 | 2025-03-20 |
| #4 | snkyfix | release/snky_fix | 2025-03-20 | 2025-03-20 |
| #3 | **adding webclient in jar** | release/develop-copy | 2025-03-20 | 2025-03-20 |
| #1 | **Audit log service** | develop | 2025-03-07 | 2025-03-20 |

### OPEN PRs

| PR# | Title | Branch | Created | Status |
|-----|-------|--------|---------|--------|
| #16 | Release/gtp bom changes java17 | release/gtp-bom-changes-java17 | 2025-09-22 | OPEN |
| #7 | FIx review bugs | develop | 2025-03-21 | OPEN |
| #2 | Release/develop sonar fix | release/develop-sonar-fix | 2025-03-19 | OPEN |

### CLOSED PRs

| PR# | Title | Branch | Created | Status |
|-----|-------|--------|---------|--------|
| #14 | Add.insights.yml | release/add-insights.yml | 2025-04-11 | CLOSED |
| #13 | Update README.md | feature/testing | 2025-04-11 | CLOSED |
| #10 | Upgrade to jdk17 | release/jdk17_version | 2025-04-04 | CLOSED |
| #9 | upgrade to jdk17 | release/DSIDVL-41613-1 | 2025-04-03 | CLOSED |
| #8 | Release/dsidvl 41613 : upgrade to jdk 17 | release/DSIDVL-41613 | 2025-04-02 | CLOSED |

### Key Interview PRs
- **#1** - THE MAIN PR - Initial audit log library
- **#3** - WebClient integration
- **#6** - Core functionality
- **#11** - Regex endpoint filtering (flexibility feature)

---

## inventory-status-srv

**Related Bullets:** 5 (OpenAPI), 6 (DC Inventory), 8 (Auth Framework)

### MERGED PRs (Recent 100)

| PR# | Title | Branch | Created | Merged |
|-----|-------|--------|---------|--------|
| #357 | DSIDVL-51809 adding salchichas mx supplier | feature/salchichas-mx-supplier | 2026-01-28 | 2026-01-29 |
| #354 | Modified consumerId type to UUID | feature/search-inbound-dummy-commit | 2026-01-20 | 2026-01-20 |
| #352 | Enabled Search Inbound Endpoint [CRQ: CHG3699850] | feature/store-inbound-ccm-change | 2026-01-14 | 2026-01-14 |
| #350 | Updated ccm for featureFlag [CRQ: CHG3699539] | feature/fix-featureflag-issue | 2026-01-14 | 2026-01-14 |
| #348 | Dummy commit to Publish ccm [CRQ: CHG3699327] | feature/dummy-pr | 2026-01-14 | 2026-01-14 |
| #347 | Dummy Commit [CRQ: CHG3699327] | feature/dummy-pr | 2026-01-14 | 2026-01-14 |
| #345 | Removed fallback resolution path [CRQ: CHG3699327] | feature/ccm-fix-prod | 2026-01-14 | 2026-01-14 |
| #343 | Updated Prod CCM with default values [CRQ: CHG3699327] | feature/prod-ccm-fix | 2026-01-14 | 2026-01-14 |
| #341 | Prod CCM for Feature Flag config [CRQ: CHG3699327] | feature/prod-ccm-feature-flag-change | 2026-01-13 | 2026-01-14 |
| #339 | **Added feature flag to disable/enable endpoints** | feature/feature-flag-config | 2026-01-13 | 2026-01-13 |
| #338 | **Container test for dc inventory from main** | feature/container_test_for_dc_inventory_from_main | 2026-01-13 | 2026-01-13 |
| #337 | Added resiliency test | feature/resiliency-test | 2026-01-09 | 2026-01-12 |
| #335 | Added probe config | feature/add-probes | 2026-01-07 | 2026-01-08 |
| #334 | Fixed Linter Issues (modified type number to string) | feature/api-linter-fix | 2026-01-06 | 2026-01-06 |
| #333 | Using Main branch kitt-common in profiles | feature/kitt-common-profile-change | 2026-01-06 | 2026-01-06 |
| #332 | Added AutoCRQ config in kitt | feature/autocrq-kitt-changes | 2026-01-05 | 2026-01-06 |
| #331 | **Update Contract test gate for us** | feature/DSIDVL-50579 | 2025-12-15 | 2026-01-29 |
| #330 | **Feature/dc inv error fix** | feature/dc_inv_error_fix | 2025-12-03 | 2026-01-12 |
| #329 | adding hook | feature/add_intl_regression_hook | 2025-12-02 | 2025-12-03 |
| #328 | Modified Contract tests to Passive mode | feature/change-r2c-passive | 2025-12-02 | 2025-12-02 |
| #327 | **Added metrics in Kitt** | feature/add-metrics-kitt | 2025-12-02 | 2025-12-03 |
| #325 | Automaton hook for INTL and US | feature/automaton-hook-main | 2025-11-26 | 2025-12-02 |
| #323 | Update error messages for store inbound | feature/store-inbound-error-msg-fix | 2025-11-21 | 2025-11-24 |
| #322 | **V2 DSIDVL-48592 error handling dc inv** | feature/V2_DSIDVL-48592-error-handling_dc_inv | 2025-11-21 | 2025-11-25 |
| #320 | **Adding CCM-ANY value and siteId as tenantId** | feature/tenant-id-ccm-main | 2025-11-20 | 2025-11-26 |
| #318 | Update NON-PROD-1.0-ccm.yml | feature/ccm-fix- | 2025-11-18 | 2025-11-18 |
| #316 | Update kitt for trunk based deployment | feature/trunk-deployment-fix | 2025-11-18 | 2025-11-18 |
| #313 | Feature/fix deployment | feature/fix-deployment | 2025-11-18 | 2025-11-18 |
| #312 | **fix dc inventory api spec and remove gtin in response** | feature/dc-api-spec-fix | 2025-11-17 | 2025-11-20 |
| #304 | fix ccm | feature/ccm-updates | 2025-11-11 | 2025-11-11 |
| #302 | Fixed Container Tests | release/fix-item-seach-container-test | 2025-11-10 | 2025-11-11 |
| #301 | update ccm | feature/ccm-updates | 2025-11-07 | 2025-11-11 |
| #299 | Sandbox code for StoreInbound API (#276) | release/storeInbound-sandbox-deployment | 2025-11-07 | 2025-12-02 |
| #298 | Update refs for trunk based deployment | feature/fix-ccm | 2025-11-06 | 2025-11-07 |
| #296 | Update default Values in ccm | feature/fix-ccm | 2025-11-06 | 2025-11-06 |
| #295 | Container tests Stage to Main | release/container-tests | 2025-11-06 | 2025-11-10 |
| #292 | fix ccm | feature/fix-ccm | 2025-11-05 | 2025-11-06 |
| #287 | Update gtp bom | feature/gtp-bom-update | 2025-11-04 | 2025-11-05 |
| #285 | Adding mondelezinternationalinc to sr.yaml [skip ci] | feature/mondelezinternationalinc-onboarding | 2025-11-03 | 2025-11-03 |
| #283 | testing regression hook on inventory status main | feature/add_regression_main | 2025-10-07 | 2025-11-04 |
| #281 | **Update error handling for store inbound** | feature/DSIDVL-48592-error-handling | 2025-09-22 | 2025-11-05 |
| #277 | minor fix in wm_item_nbr variable name | feature/dc-inventory-api-spec-fix | 2025-08-21 | 2025-11-11 |
| #276 | **Sandbox code for StoreInbound API** | feature/sandbox-api-store-inbound | 2025-08-21 | 2025-11-06 |
| #271 | **Feature/dc inventory development** | feature/dc-inventory-development | 2025-07-24 | 2025-11-21 |
| #268 | Feature/update secret name | feature/update-secret-name | 2025-07-15 | 2025-07-15 |
| #267 | updating secret name | feature/DSIDVL-46059 | 2025-07-14 | 2025-07-15 |
| #266 | Container Tests (#259) | feature/containerTests | 2025-07-14 | 2025-07-16 |
| #265 | removing post deploy hooks for readme | feature/DSIDVL-46059 | 2025-07-14 | 2025-07-14 |
| #262 | **Feature/dsidvl 46057 api publishing automation** | feature/DSIDVL-46057-API-publishing-automation | 2025-07-09 | 2025-07-10 |
| #261 | updating akeyless path | feature/DSIDVL-46059-Updating-akeyless-path | 2025-07-08 | 2025-07-08 |
| #260 | **Feature/dc inventory api spec** | feature/dc-inventory-api-spec | 2025-07-08 | 2025-08-19 |
| #259 | Container Tests | feature/container-test | 2025-07-07 | 2025-07-14 |
| #257 | deleting files which are not needed | feature/DSIDVL-46059-api-automation-status-changes | 2025-07-03 | 2025-07-04 |
| #255 | onboarding ganaderos mx to inventory status | feature/onboarding-ganaderos-mx | 2025-06-23 | 2025-06-23 |
| #254 | adding snyk pr checker [skip-ci] | feature/add_snyk_integration_main | 2025-06-17 | 2025-06-30 |
| #252 | Cocacola CA consumer subscription | feature/cocacola_ca_consumer_subscription_intl | 2025-06-11 | 2025-06-11 |
| #250 | Cocacola CA consumer subscription | feature/cocacola_ca_consumer_subscription | 2025-06-11 | 2025-06-11 |
| #249 | **DSIDVL-44382-api spec development (Store Inbound)** | feature/DSIDVL-44382-Store-Inbound-Api-contract | 2025-06-06 | 2025-08-14 |
| #248 | snyk fixes and plugin | feature/stage_add_snyk_integration | 2025-06-03 | 2025-06-06 |
| #247 | Snyk Fixes | feature/snyk_fix | 2025-05-29 | 2025-05-29 |
| #246 | adding snyk tool--testing | feature/develop_add_snyk_integration | 2025-05-15 | 2025-06-03 |
| #244 | Updated ccm with latest domains | feature/latestDomain_ccm | 2025-05-09 | 2025-05-18 |
| #243 | integrate api automation process | feature/DSIDVL-43853 | 2025-05-09 | 2025-06-19 |
| #241 | Load Testing | feature/performance_test | 2025-05-07 | 2025-05-12 |
| #239 | Added sandbox unilever, groupbimbo consumerIds | feature/srChanges | 2025-05-06 | 2025-05-07 |
| #235 | Stage Deployment for Contract Test | feature/contract-test | 2025-04-29 | 2025-04-30 |
| #234 | Dev Deployment for Contract Testing | feature/contractTest | 2025-04-29 | 2025-04-29 |
| #233 | **Contract Testing Integration (#219)** | feature/contract-test | 2025-04-28 | 2025-04-28 |
| #231 | Fixed CCM Synk Sandbox | feature/ccmChange | 2025-04-25 | 2025-04-25 |
| #229 | Fixed EI url | feature/ccm-change | 2025-04-25 | 2025-04-25 |
| #228 | Snyk and Sandbox Fix | feature/snykandsandboxFix | 2025-04-25 | 2025-04-25 |
| #227 | Fixed Sandbox Intl Issue | feature/sandbox-fix | 2025-04-24 | 2025-04-24 |
| #226 | Fixed Synk Issues | feature/snyk-issues | 2025-04-24 | 2025-04-24 |
| #223 | PROD CCM Changes [CRQ: CHG3162033] | feature/prodChanges | 2025-04-24 | 2025-05-07 |
| #221 | Sonar Coverage Improvement (#217) | feature/sonar-coverage | 2025-04-24 | 2025-04-24 |
| #219 | **Contract Testing Integration** | feature/contractTesting | 2025-04-23 | 2025-04-28 |
| #218 | Corrected DB configuration | feature/postgreChanges | 2025-04-23 | 2025-04-23 |
| #217 | Sonar Coverage Improvement | feature/sonar_coverage_improvement | 2025-04-23 | 2025-04-24 |
| #216 | **Feature/cosmos to postgres migration (#197)** | feature/cosmos-to-postgre-migration | 2025-04-22 | 2025-04-23 |
| #214 | Added finalised ddl tables | feature/ddl-tables | 2025-04-21 | 2025-04-22 |
| #213 | sr.yaml change for PROD | feature/ca_consumerid_subscription | 2025-04-21 | 2025-04-24 |
| #212 | **Stage to Main Intl PostgresChanges CCM CA Sonar Regression** | feature/DSIDVL-42039-stg-to-main | 2025-04-17 | 2025-04-24 |
| #210 | Intl Changes from develop to stage | feature/intlChange | 2025-04-16 | 2025-04-16 |
| #208 | Feature/dsidvl 42039 dev to stage | feature/DSIDVL-42039-dev-to-stage | 2025-04-16 | 2025-04-16 |
| #205 | Build fix | feature/DSIDVL-42039-Canada-Launch-API-config-changes | 2025-04-14 | 2025-04-15 |
| #203 | Added INTL namespace in SR | feature/sr-changes | 2025-04-14 | 2025-04-15 |
| #202 | Update .insights.yml | feature/updateing-team-id | 2025-04-11 | 2025-04-11 |
| #198 | adding grupo bimbo sandbox subscription [skip-ci] | feature/grupo-mx-onboarding-sandbox | 2025-04-09 | 2025-04-10 |
| #197 | **Feature/cosmos to postgres migration** | feature/Cosmos_to_Postgres_migration | 2025-04-09 | 2025-04-22 |
| #196 | Canada API launch config changes | feature/DSIDVL-42039-Canada-Launch-API-config-change | 2025-04-09 | 2025-04-14 |
| #193 | adding UNILEVER PLC to NRTI prod and sandbox | feature/DSIDVL-42018 | 2025-04-08 | 2025-04-09 |
| #182 | DSIDVL-42506 adding grupo bimbo prod subscription | feature/grupo-consumer-onboarding | 2025-04-07 | 2025-04-07 |
| #180 | adding regression suite hook and snyk fix | feature/add_regression_test_hook_to_stage | 2025-04-03 | 2025-04-03 |
| #178 | Update Public Key | feature/update-public-key | 2025-04-02 | 2025-04-02 |
| #177 | Added INTL namespace in kitt | feature/intl-changes | 2025-04-01 | 2025-04-16 |
| #167 | route policy added | feature/routingPolicy | 2025-03-18 | 2025-03-18 |
| #166 | update stage to stg in CCM for EI US config | ei-stage-to-stg-change | 2025-03-12 | 2025-03-12 |
| #165 | ccm update Ei stage to stg | ei-stage-to-stg | 2025-03-12 | 2025-03-12 |
| #163 | Pod profiler | pod_profiler | 2025-02-24 | 2025-04-23 |
| #161 | **inventory-status store-GTIN enhancements [CRQ:CHG3014193]** | inventory-status-store-gtin-enhancement | 2025-02-19 | 2025-02-19 |

### OPEN PRs

| PR# | Title | Branch | Created | Status |
|-----|-------|--------|---------|--------|
| #359 | External Postgres CCM Change | feature/postgres-ccm-changes | 2026-01-30 | OPEN |
| #356 | DSIDVL-51368 Enhanced Logging | feature/enhanced-logging | 2026-01-27 | OPEN |
| #355 | Postgres Schema Changes | feature/postgres-schema-changes | 2026-01-21 | OPEN |
| #311 | Feature/test deployment | feature/test-deployment | 2025-11-17 | OPEN |
| #310 | container test added | feature/container_test_for_dc_inventory | 2025-11-17 | OPEN |
| #309 | Update sr.yaml | feature/DSIDVL-50331 | 2025-11-14 | OPEN |
| #308 | containarised test for store inbound | feature/DSIDVL-50329 | 2025-11-13 | OPEN |
| #307 | Load test config for storeInbound | feature/load-test-config | 2025-11-11 | OPEN |
| #289 | API Spec for items assortment | feature/mfc-api-spec-main | 2025-11-05 | OPEN |
| #288 | Error Handling added for DC Inventory | feature/DSIDVL-48592-error-handling_dc_inv | 2025-11-05 | OPEN |
| #284 | [GBI] Golden base image zulu:17-jdk-alpine-main | zulu-e125-6ezrl | 2025-10-16 | OPEN |
| #272 | [GBI] Golden base image zulu:17-jdk-alpine-main | zulu-af6f-om2wt | 2025-07-26 | OPEN |
| #270 | [GBI] Golden base image zulu:17-jdk-alpine-main | zulu-cb4d-tavbp | 2025-07-24 | OPEN |
| #269 | Feature/fixing apitoken variable issue | feature/fixing-apitoken-variable-issue | 2025-07-16 | OPEN |
| #264 | Feature/dsidvl 46057 api publishing automation | release/DSIDVL-46057-Api-publishing-automation-to-main | 2025-07-11 | OPEN |
| #263 | Config for Load Test | feature/load-test | 2025-07-11 | OPEN |
| #258 | api-spec created for DC inventory | feature/api-spec-for-dc-inventory | 2025-07-04 | OPEN |
| #242 | [GBI] Golden base image | zulu-2262-mn9iq | 2025-05-09 | OPEN |
| #238 | Feature/dsidvl 43176 canada spec file | feature/DSIDVL-43176-CanadaSpecFile | 2025-05-01 | OPEN |
| #237 | Performance Testing | feature/perf_test | 2025-04-30 | OPEN |
| #164 | Dsidvl 38603 mexico spec file | DSIDVL-38603-MexicoSpecFile | 2025-02-25 | OPEN |

### Key Interview PRs
- **#161** - GTIN authorization enhancements [CRQ]
- **#197, #216** - Cosmos to Postgres migration
- **#219, #233** - Contract Testing (R2C)
- **#249** - Store Inbound API spec
- **#260** - DC Inventory API spec
- **#271** - DC Inventory main implementation
- **#322** - DC Inventory error handling
- **#320** - Multi-tenant (siteId as tenantId)
- **#327** - Metrics in Kitt
- **#338** - Container tests

---

## cp-nrti-apis

**Related Bullets:** 3 (DSD Notifications), 4 (Spring Boot 3 Migration), 9 (Observability)

### MERGED PRs (Recent 100)

| PR# | Title | Branch | Created | Merged |
|-----|-------|--------|---------|--------|
| #1608 | Dummy commit to publish CCM | feature/dummy-ccm-change | 2026-02-04 | 2026-02-04 |
| #1605 | Added eso flag to false | feature/dummy-commit-eso | 2026-02-04 | 2026-02-04 |
| #1604 | Dummy Commit to publish core services v2 changes | feature/dummy-commit | 2026-02-04 | 2026-02-04 |
| #1602 | Core Services v2 Prod CCM Change [CRQ: CHG3745080] | feature/cs-v2-prod-ccm | 2026-01-19 | 2026-02-04 |
| #1601 | DSIDVL-50594 Resilience Test Gate to stage | feature/resilience-test-gate | 2026-01-12 | 2026-01-15 |
| #1599 | Dummy PR | feature/dummy-pr | 2026-01-07 | 2026-01-07 |
| #1595 | Modified parameters for probes | feature/verify-probe-values | 2026-01-07 | 2026-01-07 |
| #1594 | Fixed startUp probe | feature/fix-start-probe | 2025-12-16 | 2026-01-06 |
| #1579 | Core Services v2 ccm change | feature/cs-v2-ccm-change | 2025-12-09 | 2025-12-16 |
| #1576 | Core Services API V2 Adoption | feature/core-services-v2-adoption | 2025-12-09 | 2025-12-16 |
| #1572 | test updated perf test config | feature/main_test_404-iac-dsc-loadtest_config | 2025-12-02 | 2025-12-03 |
| #1571 | Release/api spec changes | release/api-spec-changes | 2025-11-20 | 2025-11-21 |
| #1568 | Added Test Hooks in Kitt | release/postdeploy-hooks | 2025-11-11 | 2025-11-28 |
| #1564 | Added Timeout for downstream APIs (#1562) | release/downstream-timeout-change | 2025-11-03 | 2025-11-04 |
| #1562 | Added Timeout for downstream APIs | feature/add-downstream-timeout | 2025-10-28 | 2025-11-03 |
| #1557 | Modified order for NumberFormatException (#1543) | release/5xx-vendorId-fix | 2025-10-23 | 2025-10-24 |
| #1556 | Added environment specific Postgres overrides [CRQ:CHG3585360] | release/postgres-ccm-change | 2025-10-23 | 2025-10-24 |
| #1550 | CCM Change for Environment specific Postgres | feature/ccm-postgres-change | 2025-10-23 | 2025-10-23 |
| #1548 | completed event type changes | feature/completed-event-type-main | 2025-10-21 | 2025-10-23 |
| #1546 | update automaton files to main branch | feature/automaton | 2025-10-17 | 2025-11-11 |
| #1545 | **Added Dynatrace Saas** | release/add-dynatrace-saas | 2025-10-17 | 2025-10-24 |
| #1543 | **Modified order for NumberFormatException** | feature/fix-5xx-issue | 2025-10-16 | 2025-10-17 |
| #1541 | Modified vendorNbr to Integer [CRQ:CHG3544061] (#1540) | release/vendorId-issue-fix | 2025-10-08 | 2025-10-09 |
| #1540 | Modified vendorNbr to Integer | feature/vendorId-fix | 2025-10-08 | 2025-10-08 |
| #1537 | Trigger Deployment for Prod | release/trigger-deploymnet | 2025-10-08 | 2025-10-08 |
| #1533 | Added SiteIdCCM Config [CRQ: CHG3536699] | release/kitt-ccm-siteId | 2025-10-07 | 2025-10-07 |
| #1531 | Fixed CCM Version for Sandbox | release/kitt-ccm-sandbox | 2025-10-07 | 2025-10-07 |
| #1529 | Feature/heap out of memory fix and postgres changes stage | feature/heap_out_of_memory_fix_and_postgres_changes_stage | 2025-10-06 | 2025-10-06 |
| #1528 | Feature/heap out of memory fix and postgres changes | feature/heap_out_of_memory_fix_and_postgres_changes | 2025-10-06 | 2025-10-06 |
| #1527 | **Correlation Id Issue Fix** | feature/correlation-id-fix | 2025-10-06 | 2025-10-06 |
| #1524 | Feature/heap out of memory fix and postgres changes | feature/heap_out_of_memory_fix_and_postgres_changes | 2025-10-02 | 2025-10-06 |
| #1522 | Postgres CCM Configuration [CRQ: CHG3536647] | release/ccm-change-postgres | 2025-09-29 | 2025-10-07 |
| #1520 | Cosmos to Postgres Changes | release/cosmos-postgres-migration | 2025-09-29 | 2025-10-07 |
| #1518 | **trip id fix and dock required bug fix (#1516) (#1517)** | feature/trip-id-bug-fix-main | 2025-09-25 | 2025-09-29 |
| #1517 | trip id fix and dock required bug fix (#1516) | feature/trip-id-bug-fix-stage | 2025-09-25 | 2025-09-25 |
| #1516 | trip id fix and dock required bug fix for nrt | feature/trip-id-bug-fix-dev-2 | 2025-09-25 | 2025-09-25 |
| #1514 | dock required changes to main | feature/dock-required-to-main | 2025-09-24 | 2025-09-25 |
| #1513 | Added ESO Flag | feature/ccm-change-stage | 2025-09-24 | 2025-09-24 |
| #1510 | CCM Overrides for Postgres | feature/ccm-db-change | 2025-09-24 | 2025-09-24 |
| #1506 | ccm change [CRQ : CHG3508324] | feature/store_gtin_ccm_change | 2025-09-24 | 2025-09-24 |
| #1503 | Release/store gtin removal | release/store_gtin_removal | 2025-09-23 | 2025-09-24 |
| #1501 | Release/gtp bom prod changes | release/gtp-bom-prod-changes | 2025-09-23 | 2025-09-23 |
| #1498 | GTP-BOM Issue Fix | release/gtp-bom-changes-stage | 2025-09-22 | 2025-09-22 |
| #1491 | Removed Supplier Persona and Modified access type for volt (#1487) | feature/volt-psp-fix | 2025-09-17 | 2025-09-17 |
| #1490 | audit log config changes fix (#1489) | feature/kitt_fix_stage | 2025-09-16 | 2025-09-16 |
| #1489 | **audit log config changes fix** | feature/kitt_fix | 2025-09-16 | 2025-09-16 |
| #1487 | **Removed Supplier Persona and Modified access type for volt** | feature/volt-psp-issue-fix | 2025-09-16 | 2025-09-17 |
| #1486 | SCUS to USWEST Cluster update for Stage deployment | feature/kitt-scus-uswest-change | 2025-09-15 | 2025-09-15 |
| #1483 | Updated ccm version and kitt change | feature/ccm-update-and-kitt-change | 2025-09-15 | 2025-09-15 |
| #1480 | Updated ccm configuration | feature/ccm-update | 2025-09-15 | 2025-09-15 |
| #1478 | **Cosmos to Postgres Migration (#1427)** | feature/postgres-code-changes | 2025-09-15 | 2025-09-15 |
| #1474 | logs and ccm updated (#1469) | feature/logs_updated_stage | 2025-09-12 | 2025-09-12 |
| #1469 | logs and ccm updated | feature/logs_updated | 2025-09-12 | 2025-09-12 |
| #1468 | kitt changes related to scus cluster | feature/change_scus-stage-a3_to_uswest | 2025-09-12 | 2025-09-12 |
| #1465 | update store gitn validation removal for some suppliers (#1443) | feature/storegtin_validation_removal_for_some_supplier_stage | 2025-09-11 | 2025-09-11 |
| #1446 | **DSIDVL-47611 DockRequired request field changes (#1433)** | feature/DSIDVL-47611-stage | 2025-09-05 | 2025-09-05 |
| #1443 | update store gitn validation removal for some suppliers | feature/storegtin_validation_removal_for_some_supplier | 2025-09-04 | 2025-09-10 |
| #1433 | **DSIDVL-47611 isDockRequired request field changes** | feature/DSIDVL-47611 | 2025-08-22 | 2025-09-05 |
| #1428 | Dummy commit to publish ccm changes | feature/ccm-change | 2025-08-13 | 2025-08-13 |
| #1427 | **Cosmos to Postgres Migration** | feature/postgres-changes | 2025-08-12 | 2025-09-15 |
| #1426 | Dummy Commit to publish CCM Changes | feature/ccm-dummy-commit | 2025-08-12 | 2025-08-13 |
| #1425 | Updated gtp bom version (#1424) | feature/gtp-bom-update | 2025-08-07 | 2025-08-11 |
| #1424 | Updated gtp bom version | feature/gtp-bom-version-update | 2025-08-01 | 2025-08-07 |
| #1421 | DSIDVL-45625 upgrading walmart-cosmos version (#1419) (#1420) | feature/DSIDVL-45625-walmart-cosmos-upgrade-main | 2025-07-18 | 2025-07-22 |
| #1420 | DSIDVL-45625 upgrading walmart-cosmos version (#1419) | feature/DSIDVL-45625-walmart-cosmos-upgrade-stg | 2025-07-17 | 2025-07-17 |
| #1419 | upgrading walmart-cosmos version to the supported version | feature/DSIDVL-45625-walmart-cosmos-upgrade | 2025-07-16 | 2025-07-17 |
| #1417 | **adding packNbrQty to kafka payload (#1405) (#1406)** | feature/DSIDVL-45265-dsc-packqty-kafka-main | 2025-07-09 | 2025-07-15 |
| #1415 | Feature/sandbox product metrics fix sandbox | feature/sandbox_product_metrics_fix-sandbox | 2025-07-04 | 2025-07-07 |
| #1414 | fix_kitt (#1413) | feature/sandbox_product_metrics_fix-stage | 2025-07-04 | 2025-07-04 |
| #1413 | fix_kitt | feature/sandbox_product_metrics_fix | 2025-07-04 | 2025-07-04 |
| #1410 | change stage to stg per ei update | hotfix/develop_update_ei_ccm_env_property | 2025-07-03 | 2025-07-03 |
| #1408 | kitt changes for sandbox product metrics (#1407) | feature/sandbox_product_metrics_kitt_change-stage | 2025-07-02 | 2025-07-03 |
| #1407 | kitt changes for sandbox product metrics | feature/sandbox_product_metrics_kitt_change | 2025-07-02 | 2025-07-02 |
| #1406 | **adding packNbrQty to kafka payload (#1405)** | feature/DSIDVL-45265-dsc-packqty-kafka-stg | 2025-06-26 | 2025-06-26 |
| #1405 | **adding packNbrQty to kafka payload** | feature/DSIDVL-45265-dsc-packqty-kafka | 2025-06-25 | 2025-06-26 |
| #1404 | testing snyk fix | feature/snyk-test | 2025-06-20 | 2025-06-20 |
| #1401 | changes for snyk fixes and plugin | feature/stage_add_snyk_integration | 2025-06-03 | 2025-06-06 |
| #1400 | **DSIDVL-44428 adding driver id to DSC request and kafka payload (#1397)** | feature/DSIDVL-44428-dsc-driver-id-main | 2025-06-03 | 2025-06-20 |
| #1399 | DSIDVL-44428 DSC Driver ID Stage | feature/DSIDVL-44428-dsc-driver-id-stage | 2025-06-02 | 2025-06-02 |
| #1397 | **DSIDVL-44428 adding driver id to DSC request and kafka payload** | feature/DSIDVL-44428-dsc-driver-id | 2025-05-27 | 2025-06-02 |
| #1394 | snyk fixes and testing new plugin | feature/develop_add_snyk_integration | 2025-05-20 | 2025-06-03 |
| #1391 | **DSIDVL-43873 removing inventory snapshot code (#1389) (#1390)** | feature/DSIDVL-43873-deprecate-inv-snap-main | 2025-05-15 | 2025-05-16 |
| #1390 | DSIDVL-43873 removing inventory snapshot code (#1389) | feature/DSIDVL-43873-deprecate-inv-snap-stage | 2025-05-15 | 2025-05-15 |
| #1389 | DSIDVL-43873 removing inventory snapshot code | feature/DSIDVL-43873-deprecate-inv-snap | 2025-05-14 | 2025-05-15 |
| #1381 | Updating latest domain | feature/latest_origin | 2025-05-09 | 2025-05-19 |
| #1375 | Update kitt.yml | feature/updating_ad_group | 2025-04-29 | 2025-04-29 |
| #1373 | **Update audit log jar version** | feature/audit-log-jdk17 | 2025-04-29 | 2025-04-29 |
| #1371 | fix stage | feature/fix-stage | 2025-04-29 | 2025-04-29 |
| #1368 | update ccm for audit log endpoints | feature/update-ccm-for-endpoints | 2025-04-29 | 2025-04-29 |
| #1363 | Update ccm.yml | feature/ccm_audit_logging_changes | 2025-04-29 | 2025-04-29 |
| #1356 | adding feeding america and core mark to commodityTypeMappings [CRQ:CHG3172794] | feature/DSIDVL-43351-feeding-america-update-main | 2025-04-28 | 2025-05-13 |
| #1353 | Feeding America Commodity Type Mappings Update [skip-ci] | feature/DSIDVL-43351-feeding-america-update | 2025-04-28 | 2025-04-28 |
| #1351 | update audit log jar | feature/audit-log-update-jdk17 | 2025-04-25 | 2025-04-29 |
| #1344 | Feature/ccm update audit log | feature/ccm-update-audit-log | 2025-04-24 | 2025-04-24 |
| #1340 | added tag for ignore empty config definition values [CRQ:CHG3172179] | release/ccm_update_lm_clients | 2025-04-24 | 2025-04-24 |
| #1339 | reduced scaling for prod and IAC prod | release/reduce_scaling_for_prod_IAC | 2025-04-24 | 2025-04-24 |
| #1337 | **Feature/springboot 3 migration (#1312)** | release/spring-boot-3-migration-main | 2025-04-24 | 2025-04-24 |
| #1335 | one major sonar issue fixed | feature/mutable_issue__fixed | 2025-04-23 | 2025-04-23 |
| #1334 | snyk issues fixed and mutable list issue fixed (#1332) | feature/mutable_issue_fixed | 2025-04-23 | 2025-04-23 |
| #1332 | snyk issues fixed and mutable list issue fixed | release/sonar_fixes_hotfix | 2025-04-23 | 2025-04-23 |
| #1324 | Feature/dsidvl update ccm.yml | feature/dsidvl-update-ccm.yml | 2025-04-22 | 2025-04-22 |
| #1313 | Release/dsidvl ccm.yml | release/DSIDVL-ccm.yml | 2025-04-22 | 2025-04-22 |
| #1312 | **Feature/springboot 3 migration** | feature/springboot_3_migration | 2025-04-17 | 2025-04-17 |

### Key Interview PRs (DSD, Spring Boot 3, Observability)

**DSC (Direct Shipment Capture) - Bullet 3:**
- #1397, #1400 - Driver ID to DSC
- #1405, #1406, #1417 - packNbrQty to Kafka
- #1433, #1446 - isDockRequired field
- #1516, #1517, #1518 - Trip ID bug fix

**Spring Boot 3 Migration - Bullet 4:**
- #1312 - Main Spring Boot 3 migration PR
- #1337 - Production deployment
- #1332-1335 - Post-migration fixes

**Observability - Bullet 9:**
- #1545 - Dynatrace Saas
- #1527 - Correlation ID fix
- #1489 - Audit log config
- #1373 - Audit log jar update

**Audit Logging Integration:**
- #1363, #1368 - CCM for audit log endpoints
- #1373, #1351 - Audit log JAR updates
- #1489, #1490 - Audit log config fixes

---

## inventory-events-srv

**Related Bullets:** 7 (Transaction Event History)

### MERGED PRs (All)

| PR# | Title | Branch | Created | Merged |
|-----|-------|--------|---------|--------|
| #143 | Resiliency Test Gate | feature/resiliency-test | 2026-02-02 | 2026-02-04 |
| #141 | DSIDVL-51809 adding salchichas mx supplier | feature/salchichas-mx-supplier | 2026-01-28 | 2026-01-29 |
| #140 | Added autoCRQ config and probe config | feature/autocrq-changes | 2026-01-08 | 2026-01-15 |
| #139 | Added r2c hook in Kitt | feature/contract-test-intl | 2025-12-02 | 2025-12-18 |
| #137 | Added perf test gate and load test files | feature/perf-test | 2025-11-21 | 2025-12-02 |
| #136 | Added actuator dependency in pom | feature/actuator-fix | 2025-11-17 | 2025-11-17 |
| #134 | Removed develop, stage branch in refs from kitt | release/fix-kitt | 2025-11-17 | 2025-11-17 |
| #133 | Container Tests config for Inventory Events | release/container-tests-main | 2025-11-14 | 2025-11-17 |
| #131 | Adding mondelezinternationalinc to sr.yaml [skip ci] | feature/mondelezinternationalinc | 2025-11-03 | 2025-11-03 |
| #125 | **Feature/iac v2 to stage** | feature/iac_v2_to_stage | 2025-10-03 | 2025-10-03 |
| #123 | **Feature/develop add iac endpoint contract ext** | feature/develop_add_iac_endpoint_contract_ext | 2025-09-05 | 2025-09-30 |
| #122 | Removing post deploy hook | feature/DSIDVL-47299-Remove-reamde-post-deploy-hook | 2025-07-31 | 2025-08-11 |
| #118 | updating-devx-config-files | feature/updating-devx-config-files | 2025-07-23 | 2025-07-23 |
| #117 | secret update | feature/updating-secret-name | 2025-07-15 | 2025-07-16 |
| #116 | remove post deploy hooks of readme and secret update | feature/DSIDVL-46059 | 2025-07-15 | 2025-07-15 |
| #115 | remove post deploy hooks of readme and secret update | feature/DSIDVL-46059-Remove-hooks-and-update-secret | 2025-07-14 | 2025-07-15 |
| #114 | Config for Container Tests (#82) | feature/container-test | 2025-07-14 | 2025-07-15 |
| #112 | **Updating kitt.yml changes and market Final changes to main** | release/DSIDVL-46057-API-publishing-automation-to-main-for-events | 2025-07-11 | 2025-07-14 |
| #111 | Akeyless-path-update-develop | feature/DSIDVL-46059-Akeyless-path-update-develop | 2025-07-09 | 2025-07-10 |
| #110 | Feature/dsidvl 46059 api publishing automation to stage | feature/DSIDVL-46059-API-publishing-automation-to-stage | 2025-07-09 | 2025-07-10 |
| #109 | dev-config changes for ca, and repo name update-fix | feature/DSIDVL-46059-API-publishing-automation-events-srv | 2025-07-03 | 2025-07-04 |
| #106 | onboarding ganaderos mx to inventory activity | feature/onboarding-ganaderos-mx | 2025-06-23 | 2025-06-23 |
| #104 | Updating request Timeout to 2 seconds | feature/request-timeout-update | 2025-06-20 | 2025-06-20 |
| #103 | Feature/dsidvl 43853 integrating events repo to automation process (#81) | feature/api-automation-process-changes-to-develop | 2025-06-18 | 2025-06-19 |
| #102 | Akeyless Path Change dev to stage | feature/akeylessChange | 2025-06-16 | 2025-06-17 |
| #101 | Akeyless Path Change | feature/akeyless-change | 2025-06-16 | 2025-06-16 |
| #99 | Cocacola CA consumer subscription | feature/cocacola_ca_consumer_subscription | 2025-06-11 | 2025-06-11 |
| #97 | Added CA test consumer Id | feature/sr-prod-changes | 2025-06-11 | 2025-06-11 |
| #96 | **Inventory Events CA Launch Cosmos-Postgres CCM Snyk [CRQ:CHG3238236]** | release/ca-launch | 2025-06-10 | 2025-06-11 |
| #95 | CA Config Changes Dev to Stage | feature/ca-ccm-change | 2025-06-04 | 2025-06-04 |
| #94 | snyk PR checker and fixes (#88) | feature/stage_add_snyk_integration | 2025-06-03 | 2025-06-06 |
| #93 | CCM Issue Fix | feature/ccm-fix | 2025-06-03 | 2025-06-03 |
| #92 | non prod ccm changes | feature/non_prod_ccm_change | 2025-06-02 | 2025-06-02 |
| #91 | non prod ccm changes | feature/non_prod_ccm_changes--develop | 2025-06-02 | 2025-06-02 |
| #89 | non prod ccm changes | feature/non_prod_ccm_changes | 2025-06-02 | 2025-06-02 |
| #88 | snyk PR checker and fixes | feature/develop_add_snyk_integration | 2025-05-30 | 2025-06-03 |
| #87 | Feature/cosmos to postgres migration codegate issue fixed stage | feature/cosmos_to_postgres_migration_codegate_issue_fixed--stage | 2025-05-30 | 2025-05-30 |
| #86 | Feature/cosmos to postgres migration codegate issue fixed | feature/cosmos_to_postgres_migration_codegate_issue_fixed | 2025-05-29 | 2025-05-29 |
| #85 | Feature/cosmos to postgres migration codegate issue fixed | feature/cosmos_to_postgres_migration_codegate_issue_fixed | 2025-05-29 | 2025-05-29 |
| #83 | code gate issue resolved | feature/cosmos_to_postgres_migration_codegate_issue_fixed | 2025-05-28 | 2025-05-28 |
| #82 | Config for Container Tests | feature/integration-test | 2025-05-23 | 2025-07-10 |
| #81 | Feature/dsidvl 43853 integrating events repo to automation process | feature/DSIDVL-43853-INTEGRATING-EVENTS-REPO-TO-AUTOMATION-PROCESS | 2025-05-19 | 2025-06-18 |
| #80 | **cosmos to postgres migration** | feature/Cosmos_to_Postgres_migration | 2025-05-14 | 2025-05-28 |
| #79 | Updating latest domain | feature/latestDomain_ccm | 2025-05-09 | 2025-05-18 |
| #77 | Added CA consumers for stg and sandbox | feature/ca-consumers | 2025-05-09 | 2025-05-13 |
| #67 | grupo bimbo sandbox onboarding [skip-ci] | feature/grupo-mx-sandbox-onboarding | 2025-04-09 | 2025-04-10 |
| #66 | Canada launch api config changes | feature/DSIDVL-42039-Canada-Launch-API-config-changes | 2025-04-09 | 2025-06-02 |
| #65 | adding UNILEVER PLC to NRTI prod and sandbox | feature/DSIDVL-42018_add_UNILEVER_PLC_to_NRTI_mx | 2025-04-08 | 2025-04-09 |
| #64 | grupo bimbo onboarding [skip-ci] | feature/mx-grupo-onboarding | 2025-04-08 | 2025-04-10 |
| #63 | adding regression suite hook to stage | feature/add_regression_suite_concord_to_stage | 2025-04-03 | 2025-04-03 |
| #61 | Update openapi.json | feature/k0g0grg-patch-1-1 | 2025-04-02 | 2025-04-02 |
| #60 | CRQ : CHG3095678 dummy commit [CRQ : CHG3095678] | feature/nameUpdate2 | 2025-04-02 | 2025-04-02 |
| #59 | CRQ : CHG3095678 name update [CRQ : CHG3095678] | feature/nameUpdate1 | 2025-04-02 | 2025-04-02 |
| #58 | CRQ : CHG3095678 name update [CRQ : CHG3095678] | feature/nameUpdate | 2025-04-02 | 2025-04-02 |
| #57 | CRQ : CHG3095678 Feature/team roster id and apm id update [CRQ:CHG3095678] | feature/teamRosterIdAndApmIdUpdate | 2025-04-02 | 2025-04-02 |
| #56 | CRQ : CHG3095678 CCM Update [CRQ : CHG3095678] | feature/k0g0grg-patch-1 | 2025-04-02 | 2025-04-02 |
| #55 | CRQ : CHG3095678 Dummy Commit [CRQ : CHG3095678] | feature/dummyCommit | 2025-04-02 | 2025-04-02 |
| #54 | CRQ : CHG3095678 Update ccm dummy changes [CRQ:CHG3095678] | feature/ccmDummyUpdate | 2025-04-01 | 2025-04-02 |
| #51 | Added multiple headers for Sandbox(#50) | feature/sandboxChanges | 2025-04-01 | 2025-04-01 |
| #50 | Added multiple headers | feature/sandbox-changes | 2025-04-01 | 2025-04-01 |
| #49 | SR Intl App changes (#47) | feature/DSIDVL-42274-sr-intl-changes-stg | 2025-04-01 | 2025-04-01 |
| #47 | SR Intl App changes | feature/DSIDVL-42274-sr-intl-changes | 2025-03-28 | 2025-04-01 |
| #45 | Added sr intl app | feature/add-sr-intl-app | 2025-03-28 | 2025-03-28 |
| #42 | Feature/kitt change linter fix load test | feature/kittChange-LinterFix-LoadTest | 2025-03-26 | 2025-03-26 |
| #39 | Update kitt.intl.yml | feature/kittIntlChange | 2025-03-26 | 2025-03-26 |
| #38 | **CRQ CHG3095678 TransactionHistoryApi Changes Stage to Main** | feature/transactionhistoryApiChanges | 2025-03-26 | 2025-04-01 |
| #37 | 405 error code fix (#36) | feature/bugFix-develop-to-stage | 2025-03-25 | 2025-03-25 |
| #36 | 405 error code fix | feature/bugFix | 2025-03-25 | 2025-03-25 |
| #35 | fixed stage contract testing. (#34) | feature/contractTestFix-dev-to-stage | 2025-03-25 | 2025-03-25 |
| #34 | fixed stage contract testing. | feature/contractTestingFix | 2025-03-25 | 2025-03-25 |
| #32 | Fixed sonar issues (#24) | feature/sonar-fixes | 2025-03-24 | 2025-03-24 |
| #31 | removed pattern | feature/patternIssue | 2025-03-24 | 2025-03-24 |
| #30 | Feature/develop to stage | feature/develop-to-stage | 2025-03-24 | 2025-03-24 |
| #29 | Feature/test | feature/test | 2025-03-24 | 2025-03-24 |
| #28 | Feature/synk fix bug fix | feature/synk-fix-bug-fix | 2025-03-24 | 2025-03-24 |
| #27 | sr.yaml update | feature/srYamlUpdate | 2025-03-24 | 2025-03-24 |
| #26 | **Config for Load Test** | feature/load_test | 2025-03-24 | 2025-03-26 |
| #25 | **Feature/contract testing** | feature/contractTesting | 2025-03-21 | 2025-03-24 |
| #24 | Fixed sonar issues | feature/sonar_issues | 2025-03-21 | 2025-03-24 |
| #22 | **Dsidvl 40282 transaction history (#8)** | feature/develop_to_stage | 2025-03-21 | 2025-03-21 |
| #17 | Reverted transactions in uri-mapping-policy | uri-tranformation | 2025-03-13 | 2025-03-13 |
| #15 | Add transaction in uri transform policy | DSIDVL-41572-URI-TRANSFORM-POLICY | 2025-03-13 | 2025-03-13 |
| #14 | sr.yaml configuration | DSIDVL-40834 | 2025-03-10 | 2025-03-26 |
| #13 | updated kitt config | DSIDVL-40282_transaction_history_multi_kitt_config | 2025-03-07 | 2025-03-07 |
| #10 | **DSIDVL-40848 Junit Sonar Snyk Codegate** | DSIDVL-40848-Junits | 2025-02-28 | 2025-03-20 |
| #9 | Update sr.yaml | k0g0grg-patch-1 | 2025-02-28 | 2025-02-28 |
| #8 | **Dsidvl 40282 transaction history** | DSIDVL-40282_transaction_history | 2025-02-25 | 2025-03-21 |
| #7 | Dsidvl 39961 sandbox development | DSIDVL-39961-sandbox-development | 2025-02-24 | 2025-03-20 |
| #5 | **API spec development** | DSIDVL-38198-API-spec-development | 2025-02-12 | 2025-02-17 |

### OPEN PRs

| PR# | Title | Branch | Created | Status |
|-----|-------|--------|---------|--------|
| #130 | API spec | feature/IACV2_api_spec | 2025-10-27 | OPEN |
| #129 | IACv2 to main | feature/IACv2_to_main | 2025-10-22 | OPEN |
| #128 | [GBI] Golden base image | zulu-e125-9f5rw | 2025-10-16 | OPEN |
| #121 | updating token variable | feature/updating-api-token-variable | 2025-07-29 | OPEN |
| #120 | [GBI] Golden base image | zulu-af6f-yj7ai | 2025-07-26 | OPEN |
| #119 | [GBI] Golden base image | zulu-cb4d-buf5f | 2025-07-23 | OPEN |
| #108 | initial changes for IAC API Spec Development | feature/develop_add_iac_endpoint | 2025-07-01 | OPEN |
| #84 | Performance Testing | feature/perf_test | 2025-05-28 | OPEN |
| #78 | [GBI] Golden base image | zulu-2262-ea9hq | 2025-05-09 | OPEN |
| #75 | Sample Load test for 4xx errors | feature/loadtest | 2025-04-15 | OPEN |
| #48 | Feature/DSIDVL 41618 transaction api spec file | feature/DSIDVL-41618-TransactionApiSpecFile | 2025-03-31 | OPEN |

### Key Interview PRs (Transaction History)
- **#5** - API spec development (initial design)
- **#8** - Main transaction history implementation
- **#10** - Testing, Sonar, Snyk, Codegate
- **#22** - Dev to stage deployment
- **#25** - Contract testing
- **#38** - Production deployment [CRQ]
- **#80** - Cosmos to Postgres migration
- **#96** - CA Launch with all changes [CRQ]
- **#123** - IAC endpoint contract

---

## Cross-Repo Audit Logging Integration

This section maps how audit logging (Bullet 1, 2) integrates across repos.

### Flow: API Request → Audit Log → GCS

```
┌─────────────────────────────────────────────────────────────────┐
│                     APPLICATION LAYER                            │
│  ┌─────────────────┐   ┌─────────────────┐   ┌────────────────┐│
│  │ cp-nrti-apis    │   │inventory-status │   │inventory-events││
│  │                 │   │     -srv        │   │     -srv       ││
│  │ PRs: #1373,     │   │ PRs: #327       │   │ PRs: Related   ││
│  │ #1368, #1351,   │   │ (metrics)       │   │ to events      ││
│  │ #1363 (CCM)     │   │                 │   │                ││
│  └────────┬────────┘   └────────┬────────┘   └────────┬───────┘│
│           │                     │                     │         │
│           └──────────┬──────────┴──────────┬──────────┘         │
│                      ▼                     ▼                    │
│           ┌─────────────────────────────────────────┐           │
│           │     dv-api-common-libraries             │           │
│           │     PRs: #1 (main), #3, #6, #11         │           │
│           │     - AuditFilter (HTTP interceptor)    │           │
│           │     - KafkaAuditProducer               │           │
│           │     - AsyncConfig (thread pool)         │           │
│           └─────────────────┬───────────────────────┘           │
└─────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     KAFKA LAYER                                  │
│           ┌─────────────────────────────────────────┐           │
│           │     Kafka Topics (EUS2 / SCUS)          │           │
│           │     audit-api-logs topic                │           │
│           └─────────────────┬───────────────────────┘           │
└─────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     SINK LAYER                                   │
│           ┌─────────────────────────────────────────┐           │
│           │     audit-api-logs-gcs-sink             │           │
│           │     PRs: #1 (initial), #9, #12, #28,    │           │
│           │     #29, #70, #85                       │           │
│           │     - Kafka Connect sink connector      │           │
│           │     - SMT filters                       │           │
│           │     - Site-based routing                │           │
│           └─────────────────┬───────────────────────┘           │
└─────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     STORAGE LAYER                                │
│           ┌─────────────────────────────────────────┐           │
│           │     GCS Bucket (Parquet files)          │           │
│           │     - Partitioned by date/site          │           │
│           │     - 2M+ events/day                    │           │
│           └─────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────────────┘
```

### Related PRs by Component

**Library Creation (dv-api-common-libraries):**
- #1 - Audit log service (initial)
- #3 - WebClient integration
- #11 - Regex endpoint filtering

**NRT Integration (cp-nrti-apis):**
- #1363 - CCM audit logging changes
- #1368 - CCM for audit log endpoints
- #1373 - Audit log JAR version update
- #1351 - Audit log JAR update
- #1489, #1490 - Audit log config fixes

**GCS Sink (audit-api-logs-gcs-sink):**
- #1 - Initial setup
- #9 - Prod config
- #12 - US records [CRQ]
- #28 - Performance tuning
- #29 - Autoscaling
- #70 - Connector upgrade, retry [CRQ]
- #85 - Site-based routing

---

## Quick Commands to View PRs

```bash
# View any PR details
GH_HOST=gecgithub01.walmart.com gh pr view <PR#> --repo dsi-dataventures-luminate/<repo-name>

# View PR diff
GH_HOST=gecgithub01.walmart.com gh pr diff <PR#> --repo dsi-dataventures-luminate/<repo-name>

# Examples:
# Audit logging initial PR
GH_HOST=gecgithub01.walmart.com gh pr view 1 --repo dsi-dataventures-luminate/audit-api-logs-gcs-sink

# Common library main PR
GH_HOST=gecgithub01.walmart.com gh pr view 1 --repo dsi-dataventures-luminate/dv-api-common-libraries

# Spring Boot 3 migration
GH_HOST=gecgithub01.walmart.com gh pr view 1312 --repo dsi-dataventures-luminate/cp-nrti-apis

# Transaction history main PR
GH_HOST=gecgithub01.walmart.com gh pr view 8 --repo dsi-dataventures-luminate/inventory-events-srv
```

---

*Generated for deep dive interview preparation*
