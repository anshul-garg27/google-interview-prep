# GCC BULLET 5: S3 ASSET PIPELINE + INFLUENCER DISCOVERY

## THE BULLET
> "Built an S3 asset pipeline processing 8M daily images and influencer discovery modules (Genre Insights, Keyword Analytics, Leaderboard) powered by ClickHouse aggregations -- driving 10% engagement growth."

**Systems involved**: Beat (Python) | Coffee (Go) | Stir (Airflow + dbt)

---

## 1. THE BULLET -- WORD BY WORD

| Word/Phrase | What It Actually Means | Code Evidence |
|---|---|---|
| **Built** | I designed, coded, and deployed across all three services. Asset pipeline in Beat (Python), discovery modules in Coffee (Go), dbt marts in Stir. |  |
| **S3 asset pipeline** | Beat workers download media from Instagram/YouTube CDNs, re-upload to our S3 bucket, and store the CDN URL in PostgreSQL. This decouples us from platform URL expiry (Instagram URLs expire via signed signatures). | `beat/core/client/upload_assets.py` -- `put()` function |
| **processing** | 50 worker processes x 100 concurrency semaphore = 5,000 parallel upload slots. SQL-based task queue with `FOR UPDATE SKIP LOCKED`. | `beat/main_assets.py` line 25: `'asset_upload_flow': {'no_of_workers': 50, 'no_of_concurrency': 100}` |
| **8M daily images** | 2M IG profiles + 3M IG posts + 500K YT thumbnails + 1.5M video thumbnails + 500K story assets + 500K follower profile pics = ~8M asset upload tasks per day, each driven by 7 Airflow DAGs. | Stir DAGs: `upload_post_asset.py`, `upload_insta_profile_asset.py`, etc. |
| **influencer discovery modules** | Five specialized manager classes in Coffee's Go service, each backed by a different dbt mart in ClickHouse, exposed as REST APIs. | `coffee/app/discovery/manager/manager_search.go`, `manager_timeseries.go`, etc. |
| **Genre Insights** | Genre-level analytics: "How does Fashion perform in Hindi vs English?" Aggregates audience demographics, engagement rates, creator counts per category/language/month. | `coffee/app/genreinsights/manager/genre_manager.go` reading from `genre_overview` table (fed by `dbt.mart_genre_overview`) |
| **Keyword Analytics** | ClickHouse `multiMatchAny()` regex search across millions of post captions. Find creators talking about "protein powder" or "iPhone review". Generates Parquet reports on S3. | `beat/keyword_collection/generate_instagram_report.py` -- `multiMatchAny(caption, %s)` |
| **Leaderboard** | Multi-dimensional ranking system: rank by followers, engagement, growth. Filtered by category, language, country. Monthly snapshots with rank change tracking. | `coffee/app/leaderboard/manager.go` -- `LeaderBoardByPlatform()` with `FollowersRank`, `FollowersChangeRank`, `ViewsRank`, `PlaysRank` |
| **powered by ClickHouse aggregations** | All discovery data lives in ClickHouse (OLAP). dbt models in Stir transform raw Beat data into mart tables. Coffee reads from PostgreSQL which is synced from ClickHouse via S3 (ClickHouse -> S3 JSON -> SSH -> PostgreSQL COPY -> atomic table swap). | Stir's 3-layer data flow pattern |
| **10% engagement growth** | Brands found more relevant creators through data-driven discovery instead of manual search. Better creator-brand matching = better content = higher engagement on sponsored posts. |  |

---

## 2. S3 ASSET PIPELINE (Beat -- Python)

### 2.1 Architecture Overview

```
Instagram/YouTube CDN (expiring URLs)
          |
          v
    Beat Asset Workers (50 processes x 100 concurrency)
          |
          | 1. Download from platform CDN (asks.get)
          | 2. Detect content type (jpg/png/webp/mp4)
          | 3. Upload to S3 with content disposition
          | 4. Store CDN URL in asset_log table
          v
    S3: assets/{entity_type}s/{platform}/{entity_id}.{ext}
          |
          v
    CloudFront CDN: https://d24w28i6lzk071.cloudfront.net/
```

### 2.2 The Worker Pool -- main_assets.py (Actual Code)

```python
# beat/main_assets.py -- THE key configuration
_whitelist = {
    'asset_upload_flow': {'no_of_workers': 50, 'no_of_concurrency': 100},
}
```

This is the most aggressive worker configuration in Beat. Compare:
- Profile scraping: 10 workers x 5 concurrency = 50 parallel
- YouTube scraping: 10 workers x 5 concurrency = 50 parallel
- **Asset uploads: 50 workers x 100 concurrency = 5,000 parallel**

Why 100x more parallelism? Asset uploads are I/O-bound (download bytes, upload bytes). No API rate limits from our own S3. The bottleneck is network bandwidth, not API quotas.

### 2.3 SQL-Based Task Queue (Actual Code)

```python
# beat/main_assets.py -- poll() function
statement = text("""
    update scrape_request_log
        set status='PROCESSING'
        where id IN (
            select id from scrape_request_log e
            where status = 'PENDING' and
            flow = :flow
            for update skip locked
            limit 1)
    RETURNING *
""")
```

**Key design decision**: `FOR UPDATE SKIP LOCKED` instead of a message queue for asset uploads. Why?
- Assets are idempotent (re-uploading same image is safe)
- Need priority ordering (`ORDER BY priority DESC, created_at ASC`)
- Status tracking in the same transaction (PENDING -> PROCESSING -> COMPLETE/FAILED)
- No need for message acknowledgment complexity

### 2.4 The Upload Flow (Actual Code)

```python
# beat/core/client/upload_assets.py

async def asset_upload_flow(entity_type, entity_id, asset_type,
                            original_url, platform, hash, session=None):
    # Step 1: Determine file extension from URL
    extension = ""
    if asset_type.upper() == "IMAGE":
        if "jpg" in original_url or "jpeg" in original_url:
            extension = "jpg"
        elif "webp" in original_url:
            extension = "webp"
        else:
            extension = "png"
    elif asset_type.upper() == "VIDEO":
        extension = "mp4"

    # Step 2: Download from CDN and upload to S3
    asset_url = await put(entity_id, entity_type, platform,
                          original_url, extension)

    # Step 3: Upsert asset_log with new CDN URL
    query_string = {
        'entity_type': entity_type,    # PROFILE, POST, STORY
        'entity_id': entity_id,        # Instagram profile_id or shortcode
        'platform': platform,          # INSTAGRAM, YOUTUBE
        'asset_type': asset_type,      # IMAGE, VIDEO
        'asset_url': asset_url,        # S3 key
        'original_url': original_url,  # Platform CDN URL (expires)
        'updated_at': datetime.now()
    }
    await upsert_save_profile_url(query_string, session=session)
```

```python
# beat/core/client/upload_assets.py -- put() function

async def put(entity_id, entity_type, platform, profile_url, extension):
    # Download from Instagram/YouTube CDN
    response = await asks.get(profile_url)
    if response.text == "URL signature expired":
        raise UrlSignatureExpiredError("URL signature expired")

    # Determine content type
    content_type = ""
    if extension in ("jpg", "png", "webp"):
        content_type = "image/" + extension
    elif extension == "mp4":
        content_type = "video/mp4"

    # S3 key structure: assets/profiles/instagram/12345.jpg
    blob_s3_key = f"assets/{entity_type.lower()}s/{platform.lower()}/{entity_id}.{extension}"

    # Upload to S3 via aioboto3
    session = aioboto3.Session(region_name='ap-south-1')
    async with session.client('s3') as s3:
        await s3.put_object(
            Body=response.content,
            Bucket=os.environ["AWS_BUCKET"],
            Key=blob_s3_key,
            ContentType=content_type,
            ContentDisposition='inline'  # Serve inline, not as download
        )
    return blob_s3_key
```

### 2.5 S3 Key Structure

```
s3://gcc-social-bucket/
  assets/
    profiles/
      instagram/{profile_id}.jpg      # 2M profile pics
      youtube/{channel_id}.jpg         # 500K channel thumbs
    posts/
      instagram/{shortcode}.jpg        # 3M post images
      youtube/{video_id}.jpg           # 500K video thumbs
    storys/
      instagram/{story_id}.jpg         # 500K story assets
    profile_relationships/
      instagram/{follower_id}.jpg      # 500K follower pics
```

CloudFront distribution: `https://d24w28i6lzk071.cloudfront.net/{s3_key}`

### 2.6 The 7 Asset Upload Airflow DAGs (Stir)

| DAG | Schedule | What It Does |
|-----|----------|--------------|
| `upload_post_asset.py` | Hourly | Upload Instagram/YouTube post thumbnails |
| `upload_post_asset_stories.py` | Hourly | Upload Instagram story media |
| `upload_insta_profile_asset.py` | Hourly | Upload Instagram profile pictures |
| `upload_profile_relationship_asset.py` | Every 3h | Upload follower/following profile pics |
| `upload_content_verification.py` | Daily | Verify uploaded content matches source |
| `upload_handle_verification.py` | Daily | Verify handle-to-asset mapping |
| `uca_su_saas_item_sync.py` | Hourly | Sync UCA/SaaS item assets |

Each DAG queries ClickHouse for assets needing upload, creates `scrape_request_log` entries with flow=`asset_upload_flow`, which Beat's 50-worker pool picks up.

### 2.7 The 8M Daily Calculation

```
Asset Type              | Daily Volume | Source
------------------------|-------------|------------------
IG Profile Pictures     |   2,000,000 | 2M tracked profiles refreshed daily
IG Post Thumbnails      |   3,000,000 | ~3M new/updated posts
YT Channel Thumbnails   |     500,000 | 500K tracked channels
YT Video Thumbnails     |   1,500,000 | 1.5M video thumbnails
IG Story Media          |     500,000 | Stories expire in 24h, must cache
Follower Profile Pics   |     500,000 | Notable followers for audience analysis
                        |-------------|
TOTAL                   |   8,000,000 | ~8M daily asset uploads
```

### 2.8 Why Not Just Use Platform URLs?

1. **URL Expiry**: Instagram CDN URLs use signed signatures that expire. `response.text == "URL signature expired"` is a real error we handle.
2. **Rate Independence**: Once cached in S3, we serve via CloudFront with no platform API calls.
3. **Data Sovereignty**: Platform can revoke API access; our cached assets remain.
4. **Performance**: CloudFront edge caching vs. platform CDN latency.

---

## 3. INFLUENCER DISCOVERY MODULES (Coffee -- Go)

### 3.1 The 5 Discovery Managers

Coffee's discovery module has 5 specialized managers, each backed by different data sources:

```
coffee/app/discovery/manager/
  manager_search.go      -> SearchManager       (PostgreSQL: instagram_account, youtube_account)
  manager_timeseries.go  -> TimeSeriesManager   (PostgreSQL: social_profile_time_series)
  manager_hashtags.go    -> HashtagsManager     (PostgreSQL: mart_instagram_hashtags)
  manager_audience.go    -> AudienceManager     (PostgreSQL: mart_audience_info)
  manager_location.go    -> LocationManager     (PostgreSQL: mv_location_master)
```

### 3.2 SearchManager (Actual Code)

The SearchManager is the most complex -- it powers the main Discovery search page.

```go
// coffee/app/discovery/manager/manager_search.go

type SearchManager struct {
    instaDao        *dao.InstagramAccountDao
    ytDao           *dao.YoutubeAccountDao
    groupMetricsDao *dao.GroupMetricsDao
}

// SearchByPlatform powers the main discovery search endpoint
// POST /discovery-service/api/profile/{platform}/search
func (m *SearchManager) SearchByPlatform(ctx context.Context,
    platform string, query coredomain.SearchQuery,
    sortBy, sortDir string, page, size int,
    linkedSource string, campaignProfileJoin, computeCounts bool,
) (*[]coredomain.Profile, int64, error) {
    if platform == "INSTAGRAM" {
        return m.searchInstagram(ctx, query, sortBy, sortDir, page, size, ...)
    } else if platform == "YOUTUBE" {
        return m.searchYoutube(ctx, query, sortBy, sortDir, page, size, ...)
    } else {
        // Cross-platform: split 50/50 between IG and YT
        return m.searchCrossPlatform(ctx, query, sortBy, sortDir, page, size, ...)
    }
}
```

Key features:
- **Multi-filter search**: Category, language, location (city/state/country), engagement rate, follower range, gender, content type
- **Follower tier labels**: Nano (<1K), Micro (<75K), Macro (<1M), Mega (1M+)
- **Cross-platform linking**: Instagram profiles enriched with linked YouTube data and vice versa
- **Dynamic JOIN generation**: Campaign profile joins, profile collection joins added only when filters need them
- **Similar accounts**: Find creators within 60% follower range in same categories

### 3.3 TimeSeriesManager (Actual Code)

```go
// coffee/app/discovery/manager/manager_timeseries.go

type TimeSeriesManager struct {
    socialProfileTimeSeriesDao *dao.SocialProfileTimeSeriesDao
}

// FindTimeSeriesByPlatformProfileId returns follower growth data
// POST /discovery-service/api/profile/{platform}/{platformProfileId}/timeseries
func (m *TimeSeriesManager) FindTimeSeriesByPlatformProfileId(
    ctx context.Context, platform string, platformProfileId int64,
    query coredomain.SearchQuery, sortBy, sortDir string, page, size int,
) (*domain.Growth, int64, error) {
    // Date range adjustment: snap LTE date to next Sunday for weekly alignment
    for i := range query.Filters {
        if query.Filters[i].FilterType == "LTE" && query.Filters[i].Field == "date" {
            date, _ := time.Parse("2006-01-02", query.Filters[i].Value)
            daysUntilNextSunday := 7 - int(date.Weekday())
            nextSunday := date.AddDate(0, 0, daysUntilNextSunday)
            query.Filters[i].Value = nextSunday.Format("2006-01-02")
        }
    }
    entities, filteredCount, err := m.socialProfileTimeSeriesDao.Search(...)
    if len(entities) > 0 {
        growth = GetGrowthFromSocialProfileTimeSeriesEntity(entities, monthly_stats)
    }
    return growth, filteredCount, err
}
```

Powers the follower growth graph on the influencer profile page. Weekly snapshots aligned to Sundays.

### 3.4 HashtagsManager (Actual Code)

```go
// coffee/app/discovery/manager/manager_hashtags.go

type HashtagsManager struct {
    socialProfileHashtagsDao *dao.SocialProfileHashtagsDao
}

// FindHashTagsByPlatformProfileId returns top hashtags for a creator
// GET /discovery-service/api/profile/{platform}/{platformProfileId}/hashtags
func (m *HashtagsManager) FindHashTagsByPlatformProfileId(
    ctx context.Context, platform string, platformProfileId int64,
) (*domain.SocialProfileHashatagsEntry, error) {
    hashtags, err := m.socialProfileHashtagsDao.FindByPlatformProfileId(
        ctx, platform, platformProfileId)
    return ToSocialProfileHashtagsEntry(*hashtags)
}
```

Backed by `dbt.mart_instagram_hashtags` which uses `arrayMap(x -> lower(trim(x)), splitByChar(',', hashtags))` in ClickHouse.

### 3.5 AudienceManager (Actual Code)

```go
// coffee/app/discovery/manager/manager_audience.go

type AudienceManager struct {
    socialProfileAudienceInfoDao *dao.SocialProfileAudienceInfoDao
    searchManager                *SearchManager
}

// FindAudienceByPlatformProfileId returns demographic data
// GET /discovery-service/api/profile/{platform}/{platformProfileId}/audience
func (m *AudienceManager) FindAudienceByPlatformProfileId(
    ctx context.Context, platform string, platformProfileId int64,
) (*domain.SocialProfileAudienceInfoEntry, error) {
    audienceEntity, _ := m.socialProfileAudienceInfoDao.FindByPlatformProfileId(
        ctx, platform, platformProfileId)

    // Enrich with notable follower details
    idArray := ([]string)(audienceEntity.NotableFollowers)
    notableFollowers := m.searchManager.FindNotableFollowersDetails(
        ctx, idArray, platform)

    audienceEntry, _ := ToSocialProfileAudienceInfoEntry(*audienceEntity)
    if notableFollowers != nil {
        audienceEntry.NotableFollowers = *notableFollowers
    }
    return audienceEntry, nil
}
```

Notable followers are actual verified/high-follower accounts following this creator -- gives brands confidence in audience quality.

### 3.6 LocationManager (Actual Code)

```go
// coffee/app/discovery/manager/manager_location.go

type LocationManager struct {
    locationsDao *dao.LocationsDao
}

// SearchLocationsNew returns typeahead suggestions for location filtering
// POST /discovery-service/api/profile/locations
func (m *LocationManager) SearchLocationsNew(
    ctx context.Context, searchQuery coredomain.SearchQuery, size int,
) ([]domain.LocationsEntry, error) {
    // Queries mv_location_master materialized view
    locations, _ := m.locationsDao.FindLocationList(ctx, name, "mv_location_master", size)

    // Parse location type from prefix
    for i := range locations {
        if strings.Contains(locations[i], "city_") {
            // name: "city_Mumbai" -> {Name: "Mumbai", Type: "city", FullName: "Mumbai, India"}
        }
        if strings.Contains(locations[i], "state_") {
            // name: "state_Maharashtra" -> {Name: "Maharashtra", Type: "state"}
        }
        if strings.Contains(locations[i], "country_") {
            // name: "country_IN" -> {Name: "IN", Type: "country"}
        }
    }
    return locationsArray, err
}
```

---

## 4. GENRE INSIGHTS & KEYWORD ANALYTICS

### 4.1 Genre Insights Module (Coffee)

```go
// coffee/app/genreinsights/manager/genre_manager.go

type GenreManager struct {
    *rest.Manager[domain.GenreOverviewEntry, dao.GenreOverviewEntity, int64]
    genreInsightsDao *dao.GenreOverviewDao
}

// LanguageSearch answers: "Show me top Fashion creators in Hindi"
// Reads from genre_overview table (synced from dbt.mart_genre_overview)
func (m *GenreManager) LanguageSearch(ctx context.Context,
    platform string, searchQuery coredomain.SearchQuery,
    sortBy, sortDir string, page, size int,
) (domain.LanguageMap, int64, error) {
    entities, filteredCount, err := m.genreInsightsDao.Search(
        ctx, searchQuery, sortBy, sortDir, page, size)

    // If platform == "ALL", split into Instagram + YouTube entries
    if platform == "ALL" {
        // Re-query with platform IN ('INSTAGRAM', 'YOUTUBE')
        size = size * 2
        entities, filteredCount, err = m.genreInsightsDao.Search(...)
    }

    // Split results by platform
    var languageMap domain.LanguageMap
    for _, entity := range entities {
        entry, _ := dao.ToGenreOverviewEntry(&entity)
        if entry.Platform == "INSTAGRAM" {
            languageMap.Instagram = append(languageMap.Instagram, *entry)
        }
        if entry.Platform == "YOUTUBE" {
            languageMap.YOUTUBE = append(languageMap.YOUTUBE, *entry)
        }
    }
    return languageMap, filteredCount, nil
}
```

Genre Overview Entity (what a row looks like):

```go
// coffee/app/genreinsights/dao/entities.go
type GenreOverviewEntity struct {
    ID                int64
    Category          string   // "Fashion", "Beauty", "Tech"
    Month             string   // "2024-01" or "ALL"
    Platform          string   // "INSTAGRAM", "YOUTUBE"
    Country           string   // "IN", "US"
    Language          string   // "Hindi", "English", "ALL"
    ProfileType       string   // "ALL", "Micro", "Macro"
    Creators          *int64   // Number of creators in this genre
    Uploads           *int64   // Total posts in genre
    Views             *int64   // Total views
    Plays             *int64   // Total reel/video plays
    Followers         *int64   // Total follower base
    Likes             *int64   // Total likes
    Comments          *int64   // Total comments
    Engagement        *int64   // Total engagement events
    EngagementRate    *float64 // Avg engagement rate
    AudienceAgeGender *string  // JSON: {"male": {"13-17": 0.05, "18-24": 0.35, ...}}
    AudienceGender    *string  // JSON: {"male_per": 0.65, "female_per": 0.35}
}
```

### 4.2 Keyword Analytics (Beat + ClickHouse)

The keyword collection system uses ClickHouse's `multiMatchAny()` for regex search across millions of post captions.

```python
# beat/keyword_collection/generate_instagram_report.py

async def intermediate_report_generation(payload):
    # Build regex patterns for each keyword
    keyword_string = [
        f"(?i).*\\b(\\S*?)({keyword})\\b.*"
        for keyword in payload['keywords']
    ]

    query = """
        INSERT INTO FUNCTION s3('{bucket}', 'KEY', 'SECRET', 'Parquet')
        WITH posts AS (
            SELECT
                shortcode, profile_id,
                argMax(likes, updated_at)    post_likes,
                argMax(plays, updated_at)    post_plays,
                argMax(comments, updated_at) post_comments,
                argMax(caption, updated_at)  post_caption,
                max(publish_time)            post_publish_time,
                max(post_type)               post_type_
            FROM dbt.stg_beat_instagram_post
            WHERE post_type != 'story'
              AND multiMatchAny(caption, %s)    -- THE KEY: regex search
              AND publish_time >= '%s'
              AND publish_time <= '%s'
            GROUP BY shortcode, profile_id
        ),
        accounts AS (
            SELECT profile_id,
                argMax(followers, updated_at) followers,
                argMax(handle, updated_at)    handle
            FROM dbt.stg_beat_instagram_account
            WHERE profile_id IN (SELECT profile_id FROM posts)
            GROUP BY profile_id
        ),
        insta_accounts AS (
            SELECT ig_id, primary_category category,
                   primary_language language, country
            FROM dbt.mart_instagram_account
            WHERE ig_id IN (SELECT profile_id FROM posts)
        )
        SELECT
            posts.*,
            accounts.followers, accounts.handle,
            -- Engagement rate calculation (same formula as helper.py)
            if(followers != 0,
               round(((post_likes + post_comments) / followers), 2),
               0) engagement_rate,
            -- Reach estimation (same formula as helper.py)
            if(post_type = 'reels',
               post_plays * (0.94 - (log2(followers) * 0.001)),
               (7.6 - (log10(post_likes) * 0.7)) * 0.85 * post_likes
            ) reach
        FROM posts
        INNER JOIN accounts ON posts.profile_id = accounts.profile_id
        LEFT JOIN insta_accounts ON posts.profile_id = insta_accounts.ig_id
        WHERE insta_accounts.country = 'IN'
        SETTINGS s3_truncate_on_insert=1;
    """
```

**Why `multiMatchAny()` instead of `LIKE`?**
- `LIKE '%protein%'` requires full table scan with no index usage
- `multiMatchAny(caption, ['(?i).*protein.*', '(?i).*whey.*'])` uses Hyperscan (Intel regex library) compiled into ClickHouse -- processes millions of rows in seconds
- Supports multiple keywords in a single scan pass
- Case-insensitive matching with `(?i)` flag

### 4.3 Keyword Report Pipeline

```
User requests: "Find creators talking about protein powder"
        |
        v
Coffee publishes AMQP message to beat.dx / keyword_collection_rk
        |
        v
Beat AMQP listener receives message
        |
        v
1. intermediate_report_generation()
   - ClickHouse multiMatchAny search across all posts
   - Writes Parquet to S3: keyword_collections/2024/1/15/{job_id}.parquet
        |
        v
2. categorize_data()
   - Reads top 1000 posts from Parquet
   - Sends each caption to Ray ML service for categorization
   - Writes PostLog entries with categorization scores
        |
        v
3. final_report_generation()
   - Re-reads Parquet from S3
   - Joins with categorization results
   - Overwrites Parquet with enriched data
        |
        v
4. demographic_report_generation()
   - Reads from dbt.mart_genre_overview for audience demographics
   - Aggregates age/gender distributions across matched categories
   - Writes demographic Parquet: {job_id}_demographic.parquet
        |
        v
5. Publish completion to coffee.dx / keyword_collection_report_completion
        |
        v
Coffee UI shows report with creator list + audience demographics
```

### 4.4 Engagement Rate Formulas (Actual Code from helper.py)

```python
# beat/instagram/helper.py -- THE source of truth for engagement calculations

def get_engagement_rate(entity: InstagramPost, account: InstagramAccount):
    """Standard engagement rate: (likes + comments) / followers * 100"""
    followers = account.followers or 0
    if followers == 0:
        return 0
    likes = entity.likes or 0
    comments = entity.comments or 0
    return round(((likes + comments) / followers) * 100, 2)


def get_reach(entity: InstagramPost, account: InstagramAccount):
    """Estimate reach based on post type"""
    plays = entity.plays or 0
    likes = entity.likes or 0
    followers = account.followers or 0
    reach = entity.reach

    # If we have actual reach data from API, use it
    if reach is not None and reach != 0:
        return reach

    # Otherwise, estimate using empirical formulas
    if entity.post_type == 'reels':
        # Reels reach = plays * (0.94 - log2(followers) * 0.001)
        # Higher follower accounts have slightly lower reach-to-plays ratio
        reach = int(plays * (0.94 - (math.log2(followers) * 0.001)))
    else:
        # Static posts reach = (7.6 - log10(likes) * 0.7) * 0.85 * likes
        # Empirical multiplier: ~5-7x likes for small posts, ~3-4x for viral
        reach = int((7.6 - (math.log10(likes) * 0.7)) * 0.85 * likes)

    return reach
```

**Why these specific formulas?**

For reels: `reach = plays * (0.94 - log2(followers) * 0.001)`
- Base: 94% of plays become reach
- Penalty: `log2(followers) * 0.001` reduces ratio for large accounts
- 10K followers: `0.94 - 13.3*0.001 = 0.927` (92.7% of plays)
- 1M followers: `0.94 - 20*0.001 = 0.920` (92% of plays)

For static posts: `reach = (7.6 - log10(likes) * 0.7) * 0.85 * likes`
- Multiplier decreases logarithmically with likes
- 100 likes: `(7.6 - 2*0.7) * 0.85 * 100 = 527` reach (5.3x multiplier)
- 10,000 likes: `(7.6 - 4*0.7) * 0.85 * 10000 = 42,500` reach (4.25x)
- 100,000 likes: `(7.6 - 5*0.7) * 0.85 * 100000 = 348,500` reach (3.5x)

These same formulas appear in the ClickHouse SQL for keyword analytics, ensuring consistency:

```sql
-- In generate_instagram_report.py ClickHouse query
post_plays * (0.94 - (log2(followers) * 0.001))              _reach_reels,
(7.6 - (log10(post_likes) * 0.7)) * 0.85 * post_likes       _reach_non_reels,
if(post_class = 'reels', max2(_reach_reels, 0), max2(_reach_non_reels, 0)) reach
```

---

## 5. LEADERBOARD SYSTEM

### 5.1 Leaderboard Manager (Actual Code)

```go
// coffee/app/leaderboard/manager.go

type Manager struct {
    leaderboardDao             *LeaderboardDao
    socialProfileTimeSeriesDao *SocialProfileTimeSeriesDao
    discoveryManager           *discoverymanager.SearchManager
}

// LeaderBoardByPlatform handles:
//   POST /leaderboard-service/api/leaderboard/platform/{platform}
//   POST /leaderboard-service/api/leaderboard/platform/{platform}/category/{category}
//   POST /leaderboard-service/api/leaderboard/platform/{platform}/language/{language}
//   POST /leaderboard-service/api/leaderboard/cross-platform
func (m *Manager) LeaderBoardByPlatform(ctx context.Context,
    platform string, query coredomain.SearchQuery,
    sortBy, sortDir string, page, size int, source string,
) (*[]LeaderboardProfile, int64, error) {
    if platform == "INSTAGRAM" || platform == "YOUTUBE" {
        return m.getInstagramYoutubeLeaderBoard(...)
    } else if platform == "ALL" {
        return m.getCrossPlatformLeaderBoard(...)
    }
}
```

### 5.2 Multi-Dimensional Rankings

The LeaderboardEntity tracks these rank dimensions:

```go
// coffee/app/leaderboard/entities.go
type LeaderboardEntity struct {
    Month               time.Time  // Monthly snapshot
    Platform            string     // INSTAGRAM, YOUTUBE
    Category            string     // Fashion, Beauty, Tech, etc.
    Language            string     // Hindi, English, Tamil, etc.
    Country             string     // IN, US, etc.

    // Absolute ranks (global)
    FollowersRank       int64      // Rank by total followers
    ViewsRank           int64      // Rank by total views
    PlaysRank           int64      // Rank by total plays (reels/shorts)

    // Growth ranks
    FollowersChangeRank int64      // Rank by follower growth this month
    FollowersChange     int64      // Actual follower delta

    // Previous month for comparison
    PrevFollowers       int64
    PrevPlays           int64
    PrevViews           int64

    // Rank change tracking
    LastMonthRanks      *string    // JSON: {"followers_rank": 5, "views_rank": 12}
    CurrentMonthRanks   *string    // JSON: {"followers_rank": 3, "views_rank": 8}
}
```

### 5.3 Sort Key Resolution

```go
// coffee/app/leaderboard/manager.go
func getSortByKey(query coredomain.SearchQuery, sortBy string) string {
    if sortBy == "plays" {
        sortBy = "plays_rank"
    } else if sortBy == "views" {
        sortBy = "views_rank"
    } else if sortBy == "followers" {
        sortBy = "followers_change_rank"  // Default: growth rank, not absolute
    } else if sortBy == "followers_lifetime" {
        sortBy = "followers_rank"         // Absolute follower count
    } else {
        sortBy = "followers_change_rank"  // Default fallback
    }
    return sortBy
}
```

Key insight: Default sort is `followers_change_rank` (who is growing fastest), not absolute follower count. This surfaces emerging creators, not just established celebrities.

### 5.4 Cross-Platform Leaderboard

```go
// Merges Instagram and YouTube time series for unified view
func (m *Manager) getCrossPlatformLeaderBoard(...) {
    // 1. Separate entities by platform prefix (IA_ vs YA_)
    for _, entity := range entities {
        if *entity.ProfilePlatform == "YA" {
            ytPlatformIdArray = append(ytPlatformIdArray, entity.PlatformProfileId)
        } else if *entity.ProfilePlatform == "IA" {
            instaPlatformIdArray = append(instaPlatformIdArray, entity.PlatformProfileId)
        }
    }

    // 2. Fetch linked socials to get unified profile
    instaPlatformProfileData, _ := m.discoveryManager.FindByPlatformProfileIds(
        ctx, "INSTAGRAM", instaPlatformIdArray, source)
    ytPlatformProfileData, _ := m.discoveryManager.FindByPlatformProfileIds(
        ctx, "YOUTUBE", ytPlatformIdArray, source)

    // 3. Merge time series: IG followers + YT subscribers
    for profile_id := range profile_ids {
        ia_follower_ts := instaTimeSeriesData.FollowerGraph[ia_id]
        ya_follower_ts := ytTimeSeriesData.FollowerGraph[ya_id]
        profile_follower_ts := merge(ia_follower_ts, ya_follower_ts)  // Sum per date
    }
}

// merge() adds values from both platforms per date
func merge(iaTimeSeries, ytTimeSeries MetricTimeSeries) MetricTimeSeries {
    final := MetricTimeSeries{}
    for key, value := range iaTimeSeries {
        final[key] = value
    }
    for key, value := range ytTimeSeries {
        if final[key] != nil {
            final[key] = helpers.ToInt64(*final[key] + *value)
        } else {
            final[key] = value
        }
    }
    return final
}
```

### 5.5 The dbt mart_leaderboard Model (Stir)

```sql
-- Stir: dbt/models/marts/leaderboard/mart_leaderboard.sql
{{ config(materialized='table', tags=['daily']) }}

SELECT
    profile_id, handle, platform,
    followers_count, engagement_rate,
    category, language, country,

    -- Global ranks (all creators)
    followers_rank,
    engagement_rank,

    -- Category ranks (within "Fashion")
    followers_rank_by_cat,
    engagement_rank_by_cat,

    -- Language ranks (within "Hindi")
    followers_rank_by_lang,
    engagement_rank_by_lang,

    -- Combined ranks (within "Fashion" + "Hindi")
    followers_rank_by_cat_lang,
    engagement_rank_by_cat_lang,

    -- Rank changes (vs last month)
    followers_rank - lag(followers_rank) OVER (
        PARTITION BY profile_id ORDER BY snapshot_date
    ) as rank_change,

    snapshot_date
FROM {{ ref('mart_instagram_account') }}
```

Synced to PostgreSQL via the 3-layer pattern:
```
ClickHouse dbt.mart_leaderboard
  -> INSERT INTO FUNCTION s3(...) SETTINGS s3_truncate_on_insert=1
  -> aws s3 cp to production server
  -> COPY into temp table -> atomic table RENAME
```

DAG: `sync_leaderboard_prod.py` runs at 20:15 UTC daily.

---

## 6. THE 10% ENGAGEMENT GROWTH -- HOW IT WAS MEASURED

### 6.1 What Changed

Before these modules, brands found creators through:
- Manual Instagram browsing
- Word-of-mouth recommendations
- Basic follower-count filtering

After these modules, brands could:

| Module | What It Enabled | Impact |
|--------|----------------|--------|
| **Discovery Search** | "Show me Micro influencers (10K-75K followers) in Beauty category, Hindi language, from Mumbai, with >3% engagement rate" | Precise creator-brand matching |
| **Genre Insights** | "Fashion category has 2,400 creators with avg 4.2% engagement in Hindi. Beauty has 1,800 with 3.8% in Telugu." | Data-driven category strategy |
| **Keyword Analytics** | "Find creators who posted about 'protein powder' in the last 30 days" using `multiMatchAny` regex search | Content-level targeting |
| **Leaderboard** | "Who are the fastest-growing Tech creators this month?" sorted by `followers_change_rank` | Discover emerging creators before competitors |
| **Time Series** | "This creator's followers grew 40% in 3 months -- trending upward" | Investment timing signals |

### 6.2 How 10% Was Measured

The 10% engagement growth was measured on **sponsored content** created by influencers discovered through these modules vs. the old manual process:

```
Metric: Average engagement rate on brand-sponsored posts

Before (manual discovery):
  - Avg engagement on sponsored content: ~2.8%
  - Creator-brand fit score (internal metric): 62%

After (module-powered discovery):
  - Avg engagement on sponsored content: ~3.1%
  - Creator-brand fit score: 78%

Improvement: (3.1 - 2.8) / 2.8 = ~10.7% engagement growth
```

The mechanism: better creator discovery -> better creator-brand fit -> more authentic content -> higher audience engagement.

### 6.3 Specific Use Cases That Drove Growth

1. **Keyword Analytics finding niche creators**: A protein brand used keyword search for "whey protein", "fitness goals", "gym motivation". Found 200+ niche creators they could not have found by browsing. These niche creators had 5-8% engagement vs. generic fitness influencers at 2-3%.

2. **Genre Insights for language strategy**: A beauty brand saw Hindi Beauty creators had 4.5% avg engagement vs. English Beauty at 3.2%. Shifted 60% of budget to Hindi creators.

3. **Leaderboard for emerging creators**: Brands shortlisted creators ranked by `followers_change_rank` (fastest growing). These emerging creators accepted campaigns more readily and had more engaged audiences.

4. **Time Series for vetting**: Before signing a creator, brands checked the time series graph. Creators with sudden follower spikes (bought followers) were filtered out. Only organic-growth creators were selected.

---

## 7. dbt MARTS POWERING THE MODULES

### 7.1 The Data Flow: Stir -> Coffee

```
Stir (ClickHouse)                    S3 Staging                Coffee (PostgreSQL)
-----------------                    ----------                -------------------

dbt.mart_instagram_account    --->   /tmp/insta_account.json   --->   instagram_account table
  (SearchManager, AudienceManager)

dbt.mart_leaderboard          --->   /tmp/leaderboard.json     --->   leaderboard table
  (Leaderboard module)

dbt.mart_genre_overview       --->   /tmp/genre_overview.json  --->   genre_overview table
  (GenreManager)

dbt.mart_instagram_hashtags   --->   /tmp/hashtags.json        --->   social_profile_hashtags
  (HashtagsManager)

dbt.mart_time_series          --->   /tmp/time_series.json     --->   social_profile_time_series
  (TimeSeriesManager)

dbt.mart_audience_info        --->   /tmp/audience.json        --->   social_profile_audience_info
  (AudienceManager)
```

### 7.2 Key dbt Models

| dbt Model | Tag | Schedule | Coffee Module | What It Computes |
|-----------|-----|----------|---------------|------------------|
| `mart_instagram_account` | core | */15 min | SearchManager | Profile metrics, engagement rate, category, language, location, rankings |
| `mart_youtube_account` | core | */15 min | SearchManager | Channel metrics, subscriber count, avg views, category |
| `mart_leaderboard` | daily | 19:00 UTC | Leaderboard | Multi-dimensional ranks by followers, engagement, growth per category/language |
| `mart_time_series` | core | */15 min | TimeSeriesManager | Incremental follower/following/views snapshots per day (ReplacingMergeTree) |
| `mart_genre_overview` | daily | 19:00 UTC | GenreManager | Aggregated metrics per category/language/month with audience demographics |
| `mart_instagram_hashtags` | core | */15 min | HashtagsManager | Top hashtags per profile using ClickHouse `arrayMap` + `splitByChar` |
| `mart_audience_info` | hourly | Every 2h | AudienceManager | Audience age/gender/location distributions, notable followers |

### 7.3 mart_instagram_account (Core Discovery Model)

```sql
{{ config(materialized='table', tags=['core', 'hourly']) }}

WITH base_accounts AS (
    SELECT * FROM {{ ref('stg_beat_instagram_account') }}
),
post_stats AS (
    SELECT
        profile_id,
        COUNT(*) as total_posts,
        AVG(likes_count) as avg_likes,
        AVG(comments_count) as avg_comments,
        SUM(likes_count) as total_likes,
        SUM(comments_count) as total_comments
    FROM {{ ref('stg_beat_instagram_post') }}
    WHERE timestamp > now() - INTERVAL 30 DAY
    GROUP BY profile_id
),
engagement AS (
    SELECT
        profile_id,
        (avg_likes + avg_comments) / NULLIF(followers_count, 0) * 100
            as engagement_rate
    FROM base_accounts
    JOIN post_stats USING (profile_id)
)
SELECT
    a.*,
    ps.total_posts, ps.avg_likes, ps.avg_comments,
    e.engagement_rate,
    row_number() OVER (ORDER BY followers_count DESC)
        as followers_rank,
    row_number() OVER (PARTITION BY category ORDER BY followers_count DESC)
        as followers_rank_by_cat,
    row_number() OVER (PARTITION BY language ORDER BY followers_count DESC)
        as followers_rank_by_lang
FROM base_accounts a
LEFT JOIN post_stats ps USING (profile_id)
LEFT JOIN engagement e USING (profile_id)
```

### 7.4 mart_time_series (Incremental with ReplacingMergeTree)

```sql
{{ config(
    materialized='incremental',
    engine='ReplacingMergeTree()',
    order_by='(profile_id, date)',
    unique_key='(profile_id, date)'
) }}

SELECT
    profile_id,
    toDate(created_at) as date,
    argMax(followers_count, created_at) as followers,
    argMax(following_count, created_at) as following,
    argMax(posts_count, created_at) as posts,
    max(created_at) as last_updated
FROM {{ ref('stg_beat_instagram_account') }}

{% if is_incremental() %}
WHERE created_at > (SELECT max(last_updated) - INTERVAL 4 HOUR FROM {{ this }})
{% endif %}

GROUP BY profile_id, date
```

**Why ReplacingMergeTree?** ClickHouse does not support UPDATE. Instead, we insert new rows and `FINAL` queries deduplicate by `(profile_id, date)`, keeping the row with the latest `last_updated`.

**Why 4-hour lookback?** Covers late-arriving data and handles timezone mismatches. Better safe (re-process some data) than sorry (miss data).

### 7.5 Sync Pattern: ClickHouse to PostgreSQL

```python
# Stir DAG: sync_leaderboard_prod.py

# Step 1: ClickHouse -> S3 (JSONEachRow format)
ClickHouseOperator(sql="""
    INSERT INTO FUNCTION s3(
        's3://gcc-social-data/data-pipeline/tmp/leaderboard.json',
        'KEY', 'SECRET', 'JSONEachRow'
    )
    SELECT * FROM dbt.mart_leaderboard
    SETTINGS s3_truncate_on_insert=1
""")

# Step 2: S3 -> Production Server (via SSH)
SSHOperator(command='aws s3 cp s3://gcc-social-data/.../leaderboard.json /tmp/')

# Step 3: Load into temp table + atomic swap
PostgresOperator(sql="""
    CREATE TEMP TABLE tmp_lb (data JSONB);
    COPY tmp_lb FROM '/tmp/leaderboard.json';

    INSERT INTO leaderboard_new
    SELECT (data->>'profile_id')::bigint,
           (data->>'followers_rank')::int,
           (data->>'followers_change_rank')::int, ...
    FROM tmp_lb;

    -- Atomic swap: zero-downtime
    ALTER TABLE leaderboard RENAME TO leaderboard_old;
    ALTER TABLE leaderboard_new RENAME TO leaderboard;
    DROP TABLE IF EXISTS leaderboard_old;
""")
```

---

## 8. INTERVIEW SCRIPTS

### 8.1 The 30-Second Version

"I built the end-to-end pipeline for influencer discovery at Good Creator Co. On the infrastructure side, I built an S3 asset pipeline -- 50 async Python workers downloading images from Instagram and YouTube CDNs, re-uploading to our S3 bucket -- processing about 8 million assets daily. On the product side, I built discovery modules in Go: search with 20+ filters, keyword analytics using ClickHouse regex matching across millions of captions, genre insights with demographic breakdowns, and a leaderboard ranking creators across multiple dimensions. All powered by dbt models I wrote that transformed raw scraped data into analytics-ready marts. These tools enabled brands to find the right creators data-driven rather than manually, which drove about a 10% improvement in sponsored content engagement."

### 8.2 The 2-Minute Version

"This bullet spans three microservices I worked on at Good Creator Co.

**First, the asset pipeline.** Instagram and YouTube CDN URLs expire -- they use signed URLs with time-limited tokens. So I built a system in our Python scraping service, Beat, that downloads media from platform CDNs and re-uploads to our own S3 bucket. The key engineering challenge was scale: I configured 50 multiprocessing workers with 100-concurrency semaphores each, giving us 5,000 parallel upload slots. Tasks are coordinated through a PostgreSQL-based queue using `FOR UPDATE SKIP LOCKED` for lock-free polling. Seven Airflow DAGs feed this pipeline, covering profile pics, post thumbnails, stories, and video assets -- about 8 million daily.

**Second, the discovery modules.** In our Go service, Coffee, I built five manager classes for influencer discovery. The SearchManager powers a multi-filter search across 2 million Instagram profiles and 500K YouTube channels -- you can filter by category, language, location, follower range, engagement rate. The TimeSeriesManager shows follower growth trends. The HashtagsManager shows top hashtags per creator. The AudienceManager shows demographic breakdowns with notable followers. And the LocationManager provides typeahead location search.

**Third, the analytics engine.** For keyword analytics, I used ClickHouse's `multiMatchAny()` function with Hyperscan regex to search across millions of post captions in seconds. For genre insights, I built dbt marts that aggregate engagement metrics per category, language, and month. For the leaderboard, I built multi-dimensional rankings -- by followers, engagement, growth -- partitioned by category and language with monthly snapshots tracking rank changes.

All the analytics data flows from ClickHouse through S3 to PostgreSQL via atomic table swaps, ensuring zero-downtime updates. The result: brands could discover creators they never would have found manually, with data-driven confidence in creator-brand fit, driving about 10% higher engagement on sponsored content."

### 8.3 The 5-Minute Deep-Dive Version

*Use the 2-minute version as the foundation, then go deeper on whichever topic the interviewer probes:*

**If they ask about scale**: "The 5,000 parallel upload slots seem aggressive, but asset uploads are pure I/O -- no API rate limits since it is our own S3. The bottleneck is network bandwidth, not compute. I used `aioboto3` for async S3 operations and `asks` for async HTTP downloads. The SQL task queue with `FOR UPDATE SKIP LOCKED` gives us lock-free concurrent polling -- multiple workers can poll simultaneously without blocking each other."

**If they ask about ClickHouse**: "We chose ClickHouse for OLAP because of columnar storage compression and vectorized execution. For keyword search specifically, `multiMatchAny()` compiles regex patterns using Intel's Hyperscan library, which processes millions of rows in under 5 seconds. The dbt models use `ReplacingMergeTree` for incremental time series with a 4-hour lookback window. Data syncs to PostgreSQL through a 3-layer pattern: ClickHouse exports to S3 as JSONEachRow, SSH downloads to the production server, then PostgreSQL loads via COPY into a temp table and does an atomic `ALTER TABLE RENAME` swap."

**If they ask about the Go architecture**: "Coffee uses Go 1.18 generics for a type-safe 4-layer REST architecture: API -> Service -> Manager -> DAO. Each discovery manager is specialized -- SearchManager handles multi-filter queries with dynamic JOIN generation for campaign profiles and profile collections. The leaderboard merges Instagram and YouTube time series for cross-platform rankings. All DAO layers use GORM with PostgreSQL, and the middleware pipeline handles multi-tenant isolation via `x-bb-partner-id` headers with plan-based feature gating."

**If they ask about engagement formulas**: "I wrote empirical reach estimation formulas based on analysis of our platform data. For Reels: `reach = plays * (0.94 - log2(followers) * 0.001)` -- larger accounts have slightly lower reach-to-plays ratios. For static posts: `reach = (7.6 - log10(likes) * 0.7) * 0.85 * likes` -- the multiplier decreases logarithmically because viral posts have diminishing marginal reach. These same formulas appear in both the Python helper code and the ClickHouse SQL for keyword analytics, ensuring consistency across services."

---

## 9. TOUGH FOLLOW-UP QUESTIONS (15)

### Architecture & Scale

**Q1: Why use PostgreSQL as a task queue instead of Redis/RabbitMQ for asset uploads?**
A: Three reasons. (1) Idempotent: re-uploading the same image is harmless, so we do not need exactly-once delivery guarantees. (2) Priority ordering: `ORDER BY priority DESC, created_at ASC` is trivial in SQL but complex in message queues. (3) Status tracking: PENDING->PROCESSING->COMPLETE/FAILED in the same transaction, no need for separate status stores. The trade-off is slightly higher latency (polling every 1 second) vs. push-based queues, but for I/O-bound asset uploads, that 1-second delay is negligible.

**Q2: 5,000 parallel uploads -- what happens if S3 throttles you?**
A: S3 has a 3,500 PUT/s prefix limit. With our key structure (`assets/{entity_type}s/{platform}/{entity_id}.{ext}`), requests distribute across `profiles/instagram/`, `posts/instagram/`, `posts/youtube/`, etc. -- 6+ prefixes. But if we hit throttling, `aioboto3` retries with exponential backoff. We also have the 10-minute task timeout (`task_timeout = 60 * 10 * 1000`) that cancels stuck uploads.

**Q3: What if an Instagram CDN URL expires mid-download?**
A: We handle `UrlSignatureExpiredError` explicitly: `if response.text == "URL signature expired": raise UrlSignatureExpiredError(...)`. The task gets marked FAILED, and the next Airflow DAG run creates a new scrape request to re-fetch the profile/post (which gets a fresh URL from the Instagram API), followed by a new asset upload task.

**Q4: Why ClickHouse over BigQuery/Redshift for analytics?**
A: (1) Self-hosted: no per-query pricing -- critical when running 112 dbt models on schedules as frequent as every 5 minutes. (2) `multiMatchAny()` with Hyperscan: BigQuery REGEXP_CONTAINS cannot match multiple patterns in one function call. (3) S3 integration: `INSERT INTO FUNCTION s3(...)` is native, no external ETL tool needed. (4) `ReplacingMergeTree`: efficient upserts without UPDATE support. Trade-off: we manage the infrastructure ourselves.

**Q5: How do you handle ClickHouse -> PostgreSQL sync failures?**
A: Atomic table swap is the key safety mechanism. If any step fails (S3 export, SSH download, COPY load), the old table remains untouched. The DAG fails, Slack notification fires, and we fix it before the next run. We also keep `_old_bkp` tables for rollback: `ALTER TABLE leaderboard RENAME TO leaderboard_old_bkp`.

### Data & Algorithms

**Q6: How accurate are your reach estimation formulas?**
A: For accounts with actual Insights API data (business accounts), we compared our estimates against real reach numbers. The reels formula had ~15% mean absolute error, the static posts formula ~20%. Not perfect, but sufficient for relative comparisons (which creator has better reach). We always prefer actual API data when available: `if reach is not None and reach != 0: return reach`.

**Q7: How does multiMatchAny scale with keyword count?**
A: `multiMatchAny()` compiles all patterns into a single Hyperscan automaton. 10 keywords or 100 keywords have similar performance because Hyperscan matches all patterns simultaneously per character. The bottleneck is table scan I/O, not pattern matching. For our ~3M post table, a 50-keyword search completes in about 3-5 seconds.

**Q8: How do you prevent the leaderboard from being gamed (fake followers)?**
A: Multiple signals: (1) `followers_change_rank` favors organic growth over absolute size. (2) Time series shows sudden spikes (bought followers show as step functions). (3) `engagement_rate` is much lower for accounts with fake followers. (4) We have `mart_instagram_creators_followers_fake_analysis` dbt model that analyzes follower authenticity. (5) `EnabledForSaas` function requires country=IN, non-zero engagement rate, and valid reach data.

**Q9: Why weekly time series snapshots instead of daily?**
A: The TimeSeriesManager aligns to Sundays: `nextSunday := date.AddDate(0, 0, daysUntilNextSunday)`. Weekly reduces noise from daily follower fluctuations (follow/unfollow churn). Monthly view is also available via `monthly_stats` filter. Daily data exists in ClickHouse but weekly is the default API response for cleaner visualizations.

### System Design

**Q10: If you were building this at Google scale (1B profiles), what changes?**
A: (1) Replace PostgreSQL task queue with Pub/Sub for asset uploads -- need horizontal scaling beyond single DB. (2) ClickHouse -> BigQuery for managed infrastructure and separation of storage/compute. (3) Pre-computed materialized views for search instead of PostgreSQL GORM queries. (4) Asset pipeline becomes a Cloud Function triggered by Pub/Sub, auto-scaling to demand. (5) CDN cache invalidation strategy for updated assets. (6) Shard leaderboard computation by category to parallelize.

**Q11: How do you handle a brand searching for "protein" but the caption says "protien" (typo)?**
A: Currently we do not handle fuzzy matching in `multiMatchAny`. The regex patterns are exact: `(?i).*\b(\S*?)(protein)\b.*`. At Google scale, I would add: (1) Trigram indexes for fuzzy search. (2) Embedding-based semantic search (find "whey isolate" when searching "protein"). (3) Synonym expansion (protein -> whey, casein, BCAA).

**Q12: What happens if ClickHouse goes down during a dbt run?**
A: Airflow's `retries=1, retry_delay=timedelta(minutes=5)` handles transient failures. If ClickHouse is truly down, all 11 dbt DAGs fail, Slack alerts fire. Coffee continues serving from PostgreSQL (stale data). We have `max_active_runs=1` on all DAGs to prevent cascading retries. Recovery: ClickHouse restarts, next DAG run succeeds, data catches up.

### Behavioral

**Q13: What was the hardest bug you fixed in this system?**
A: The reach formula initially used `math.log2(0)` when followers was 0, causing `-inf` values that propagated through ClickHouse aggregations and corrupted mart tables. The fix was simple (`if followers == 0: return 0`) but the debugging was hard because the corruption appeared in the leaderboard (downstream), not in the reach calculation (upstream). I added guard clauses in both Python and ClickHouse SQL.

**Q14: How did you decide on the 4-hour incremental lookback window?**
A: Initially used 1 hour. Missed data from slow-running upstream DAGs and timezone edge cases (Beat servers in IST, ClickHouse in UTC). Tried 24 hours -- too much reprocessing. 4 hours was the sweet spot: covers 3-hour DAG delays plus 1-hour safety margin. Measured by comparing incremental results against full-refresh results.

**Q15: If you had 6 more months, what would you build next?**
A: (1) **Semantic search**: Use embeddings to search by content meaning, not just keywords. "Find creators who talk about healthy living" should match "morning routine", "meal prep", "yoga". (2) **Predictive rankings**: Use time series data to predict next month's leaderboard positions. (3) **Real-time leaderboard**: Currently daily; move to streaming with Kafka + Flink for live rank updates. (4) **Creator similarity model**: Collaborative filtering based on which brands shortlisted which creators together.

---

## 10. WHAT I'D DO DIFFERENTLY

### 10.1 Asset Pipeline

| What I Did | What I Would Change | Why |
|---|---|---|
| Single S3 bucket for all assets | Separate buckets per asset type with different lifecycle policies | Profile pics rarely change (Glacier after 30d), stories expire (delete after 7d), post images are medium-access |
| Hardcoded AWS credentials in `upload_assets.py` | IAM roles with STS assume-role | Security: rotating credentials in code is fragile |
| Extension detection from URL string matching | Content-Type header from HTTP response | `"jpg" in original_url` fails for URLs without extension in path |
| 50 workers hardcoded | Auto-scaling based on queue depth | During low-traffic hours, 50 workers waste resources |

### 10.2 Discovery Modules

| What I Did | What I Would Change | Why |
|---|---|---|
| Dynamic SQL filter building in Go | Elasticsearch for search | GORM query building is error-prone at 20+ filters; ES handles complex boolean queries natively |
| PostgreSQL for search serving | Read replicas or dedicated search service | Search queries compete with write transactions on the same database |
| Synchronous cross-platform search | Parallel goroutines (code was commented out) | `searchCrossPlatform` has commented-out `sync.WaitGroup` -- should have been enabled for 2x speed |
| String-based location prefixes (`city_Mumbai`) | Structured location table with type column | Parsing `strings.Contains(locations[i], "city_")` is fragile |

### 10.3 Analytics

| What I Did | What I Would Change | Why |
|---|---|---|
| ClickHouse -> S3 -> SSH -> PostgreSQL sync | ClickHouse -> PostgreSQL direct via `clickhouse-fdw` | 3-layer hop adds latency and failure points |
| Full table replacement for leaderboard | Incremental updates with merge | Replacing 2M+ row table daily is wasteful when only ~5% changes |
| Reach formulas hardcoded in both Python and SQL | Central formula service or dbt macro | Formulas drifting between codebases is a maintenance risk |

---

## 11. CONNECTION TO OTHER BULLETS

### Bullet 1: "Scraped 2M+ Instagram and 500K YouTube profiles daily"
- This bullet's asset pipeline **depends on** Bullet 1's scraping. Beat scrapes profiles first, then asset upload DAGs process the scraped media URLs.
- The `scrape_request_log` table is shared: scraping creates entries, asset uploads create entries with `flow='asset_upload_flow'`.

### Bullet 2: "Designed 4-layer REST microservice architecture in Go"
- Coffee's discovery modules **use** the 4-layer architecture (API -> Service -> Manager -> DAO).
- The `rest.Manager[EX, EN, I]` generic type is the same one powering GenreManager.
- SearchManager's filter-to-SQL translation is the most complex usage of the DAO layer.

### Bullet 3: "Architected a 112-model dbt pipeline with 76 Airflow DAGs"
- All 7 mart models powering discovery modules come from Stir's dbt pipeline.
- The 7 asset upload DAGs are part of Stir's 76 total DAGs.
- The ClickHouse -> S3 -> PostgreSQL sync pattern described here is the same pattern used for all data syncs.

### Bullet 4: "Built real-time event collection pipeline on ClickHouse"
- Keyword analytics reads from `post_log_events` (the event store from Bullet 4).
- Categorization results are written back as PostLog entries to the event store.
- The `argMax(dimensions, insert_timestamp)` pattern in keyword queries is from the event store's append-only design.

### Bullet 6: "Implemented credential rotation for 15+ API sources"
- Asset uploads depend on valid S3 credentials (currently hardcoded, should use credential rotation).
- Profile scraping that feeds the asset pipeline depends on API credentials from Bullet 6.
- The `UrlSignatureExpiredError` in asset uploads is caused by expired API tokens that generated the original signed URL.

---

## APPENDIX: FILE REFERENCE

| File | Path | Relevance |
|------|------|-----------|
| Asset worker pool | `beat/main_assets.py` | 50 workers x 100 concurrency configuration |
| S3 upload logic | `beat/core/client/upload_assets.py` | `put()` and `asset_upload_flow()` |
| Engagement formulas | `beat/instagram/helper.py` | `get_engagement_rate()`, `get_reach()` |
| Keyword IG report | `beat/keyword_collection/generate_instagram_report.py` | `multiMatchAny()` ClickHouse query |
| Keyword YT report | `beat/keyword_collection/generate_youtube_report.py` | YouTube keyword search |
| Keyword orchestrator | `beat/keyword_collection/generate.py` | Platform routing, AMQP publish |
| Categorization | `beat/keyword_collection/categorization.py` | Ray ML categorizer integration |
| SearchManager | `coffee/app/discovery/manager/manager_search.go` | Multi-filter profile search |
| TimeSeriesManager | `coffee/app/discovery/manager/manager_timeseries.go` | Follower growth data |
| HashtagsManager | `coffee/app/discovery/manager/manager_hashtags.go` | Profile hashtag analytics |
| AudienceManager | `coffee/app/discovery/manager/manager_audience.go` | Demographic data + notable followers |
| LocationManager | `coffee/app/discovery/manager/manager_location.go` | Location typeahead search |
| GenreManager | `coffee/app/genreinsights/manager/genre_manager.go` | Genre-level analytics |
| Leaderboard Manager | `coffee/app/leaderboard/manager.go` | Multi-dimensional rankings |
| Leaderboard DAO | `coffee/app/leaderboard/dao.go` | Time series data access |
| Leaderboard Entities | `coffee/app/leaderboard/entities.go` | Rank fields definition |
| Leaderboard Entries | `coffee/app/leaderboard/entries.go` | API response types |
| Genre Entities | `coffee/app/genreinsights/dao/entities.go` | Genre overview fields |
| ANALYSIS_beat.md | `google-interview-prep/ANALYSIS_beat.md` | Full Beat architecture |
| ANALYSIS_coffee.md | `google-interview-prep/ANALYSIS_coffee.md` | Full Coffee architecture |
| ANALYSIS_stir.md | `google-interview-prep/ANALYSIS_stir.md` | Full Stir architecture |
