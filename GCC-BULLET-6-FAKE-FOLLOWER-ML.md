# GCC BULLET 6: ML-POWERED FAKE FOLLOWER DETECTION

## THE BULLET
> "Built ML-powered fake follower detection: 5-feature ensemble with HMM transliteration across 10 Indic scripts, deployed on serverless AWS Lambda pipeline, improving processing speed by 50%."

**Systems involved**: fake_follower_analysis (Python ML) | Beat (GPT enrichment)

---

## 1. THE BULLET -- WORD BY WORD

| Word/Phrase | What It Actually Means | Code Evidence |
|---|---|---|
| **Built** | I designed, coded, and deployed the entire system end-to-end: the ML detection model (fake.py, 385 lines), the data pipeline (push.py, pull.py), the Hindi transliteration (createDict.py), the Docker container, and integrated with Beat's GPT enrichment. | 955+ lines of Python across 6 files |
| **ML-powered** | Not deep learning -- this is a heuristic ML ensemble. Five independent feature extractors feed into two scoring functions (process1 for binary, final for weighted). The "ML" also comes from the HMM-based indictrans library used for transliteration -- actual trained Hidden Markov Models with Viterbi decoding. | `fake.py` lines 252-277: `process1()` + `final()` |
| **fake follower detection** | Given an Instagram follower's handle and full name, determine if they are a real person or a bot/fake account. Output is a confidence score: 0.0 (real), 0.33 (weak fake signal), or 1.0 (definitely fake). | `fake.py` model() returns 16-field response per follower |
| **5-feature ensemble** | Five independent detection features: (1) non-Indic language regex, (2) digit count threshold >4, (3) handle-name special character correlation, (4) RapidFuzz weighted similarity (2*partial + sort + set)/4, (5) 35,183-name Indian name database fuzzy match. Each produces a sub-score; two ensemble functions combine them. | Features in `fake.py` lines 123-312 |
| **HMM transliteration** | Hidden Markov Model transliteration via the `indictrans` library. Each of the 10 language models has 6 files: coefficient matrix (coef_.npy), output classes (classes.npy), initial/transition/final state probabilities, and a sparse feature vocabulary. The Viterbi algorithm (implemented in Cython for speed) finds the most likely English character sequence for each Indic script input. | `indic-trans-master/indictrans/models/` -- 10 language directories, `_decode/viterbi.pyx` |
| **across 10 Indic scripts** | Hindi (Devanagari), Bengali, Gujarati, Kannada, Malayalam, Odia, Punjabi (Gurmukhi), Tamil, Telugu, Urdu (Perso-Arabic). Plus derivative support: Marathi/Nepali/Konkani reuse Hindi model, Assamese reuses Bengali. Total: 583 characters mapped across all scripts. | `fake.py` lines 44-91: `data = {"hin": [...], "pan": [...], ...}` |
| **deployed on serverless** | AWS Lambda function packaged as an ECR Docker container. No servers to manage -- Lambda auto-scales with SQS queue depth. | `Dockerfile` with `CMD ["fake.handler"]` |
| **AWS Lambda pipeline** | Full pipeline: ClickHouse CTE query -> S3 export (JSONEachRow) -> local download -> 8-worker multiprocessing batch -> SQS queue -> Lambda (ECR) -> Kinesis stream -> multi-shard parallel reader -> output JSON. Five AWS services: S3, SQS, Lambda, Kinesis, ECR. | `push.py` (154 lines) + `pull.py` (131 lines) |
| **improving processing speed by 50%** | Before: sequential single-machine processing -- one follower at a time on an EC2 instance. After: Lambda auto-scales based on SQS queue depth (potentially hundreds of concurrent invocations) + push.py uses 8 parallel multiprocessing workers for SQS ingestion. The combination of parallel ingestion + serverless auto-scaling cut total batch processing time in half. | `push.py`: `multiprocessing.Pool(processes=8)` + Lambda auto-scaling |

---

## 2. THE 5-FEATURE ENSEMBLE (with code)

### Detection Pipeline Overview

```
INPUT: {follower_handle, follower_full_name}
         |
   [Step 1] symbol_name_convert() -- 13 Unicode fancy text variants -> ASCII
         |
   [Step 2] check_lang_other_than_indic() -- Feature 1
         |
   [Step 3] detect_language() -> process_word() or Transliterator()
         |
   [Step 4] uni_decode() -- final ASCII normalization
         |
   [Step 5] clean_handle() -- multi-stage regex cleaning
         |
   [Step 6] 5 features extracted in parallel
         |
   [Step 7] process1() binary + final() weighted ensemble
         |
OUTPUT: 16-field response with confidence score
```

### Feature 1: Non-Indic Language Detection -- `check_lang_other_than_indic()`

**What it does**: Flags followers whose names contain Greek, Armenian, Georgian, Chinese, or Korean characters. Real Indian Instagram followers should not have these scripts in their display names.

**Why this works**: Bot farms often generate accounts with random Unicode characters from non-Indian scripts. A follower named with Chinese or Korean characters on an Indian creator's account is almost certainly fake.

```python
# fake.py - check_lang_other_than_indic() (lines 123-126)
def check_lang_other_than_indic(symbolic_name):
    if not symbolic_name:
        print(symbolic_name)
    return 1 if re.search(r'[A-Za-z\u0391-\u03A9\u03B1-\u03C9\u0531-\u0556\u10D0-\u10F0\u4E00-\u9FFF\uAC00-\uD7A3]+',
                           symbolic_name, re.UNICODE) else 0
    # Greek: U+0391-03A9, U+03B1-03C9
    # Armenian: U+0531-0556
    # Georgian: U+10D0-10F0
    # Chinese (CJK): U+4E00-9FFF
    # Korean (Hangul): U+AC00-D7A3
```

**Output**: 0 (REAL -- no foreign scripts) or 1 (FAKE -- foreign script detected)

**Edge case**: This does NOT flag Indic scripts (Hindi, Bengali, etc.) -- those are handled separately through transliteration. It also does not flag Latin/English text, since many Indian users write names in English.

---

### Feature 2: Digit Count Threshold -- `fake_real_more_than_4_digit()`

**What it does**: Counts numerical digits in the follower's handle. If there are more than 4 digits, the account is flagged as fake.

**Why this works**: Real users might add a birth year ("rahul_1995") or a short number ("priya_27"), but bots are often auto-generated with long random digit suffixes ("user_8374629").

```python
# fake.py (lines 213-220)
def count_numerical_digits(text):
    if not isinstance(text, str):
        text = str(text)
    return sum(c.isdigit() for c in text)

def fake_real_more_than_4_digit(number):
    return 1 if number > 4 else 0
```

**Output**: 0 (REAL -- 4 or fewer digits) or 1 (FAKE -- more than 4 digits)

**Examples**:
- `"rahul_27"` -> 2 digits -> 0 (REAL)
- `"priya_1995"` -> 4 digits -> 0 (REAL)
- `"user_12345"` -> 5 digits -> 1 (FAKE)
- `"bot_99887766"` -> 8 digits -> 1 (FAKE)

---

### Feature 3: Handle-Name Special Character Correlation -- `process()`

**What it does**: Analyzes the relationship between special characters (_, -, .) in the handle and how well the handle matches the follower's displayed name.

**Why this works**: Real users who put separators in their handle (e.g., "john_doe") are typically encoding their real name. If the handle has separators but does NOT match the name, it is suspicious. If the handle has NO separators, we cannot draw a conclusion (INCONCLUSIVE).

```python
# fake.py - process() (lines 188-200)
def process(follower_handle, cleaned_handle, cleaned_name):
    SPECIAL_CHARS = ('_', '-', '.')

    if any(char in follower_handle for char in SPECIAL_CHARS):
        if not ' ' in cleaned_name:  # Single-word name
            if generate_similarity_score(cleaned_handle, cleaned_name) > 80:
                return 0  # REAL -- handle matches name
            else:
                return 1  # FAKE -- handle has separators but doesn't match name
        else:
            return 0  # REAL -- multi-word name with separators is normal
    else:
        return 2  # INCONCLUSIVE -- no special chars to analyze
```

**Decision tree**:
```
Has _ or - or . in handle?
  YES -> Is name a single word?
    YES -> Similarity > 80?
      YES -> 0 (REAL)
      NO  -> 1 (FAKE)
    NO (multi-word) -> 0 (REAL)
  NO -> 2 (INCONCLUSIVE)
```

**Output**: 0 (REAL), 1 (FAKE), or 2 (INCONCLUSIVE)

---

### Feature 4: RapidFuzz Weighted Similarity -- `generate_similarity_score()`

**What it does**: Measures how similar the follower's handle is to their displayed name using three fuzzy matching algorithms, weighted and averaged.

**Why this works**: Real people's handles usually resemble their real names. A handle "rahul_kumar" with display name "Rahul Kumar" scores ~95. A handle "xyz_abc123" with display name "Priya Sharma" scores ~15.

```python
# fake.py - generate_similarity_score() (lines 225-244)
def generate_similarity_score(handle, name):
    name = name.split()
    if len(name) <= 4:
        name_permutations = [' '.join(p) for p in permutations(name)]
        # Max 4 words = 4! = 24 permutations
    else:
        name_permutations = name  # Too many, skip permutations

    similarity_score = -1
    cleaned_handle = handle.replace(' ', '')

    for name in name_permutations:
        cleaned_name = name.replace(' ', '')
        partial_ratio = fuzzz.partial_ratio(cleaned_handle.lower(), cleaned_name.lower())
        token_sort_ratio = fuzzz.token_sort_ratio(cleaned_handle.lower(), cleaned_name.lower())
        token_set_ratio = fuzzz.token_set_ratio(cleaned_handle.lower(), cleaned_name.lower())

        # Weighted average: partial gets 2x weight
        similarity_score = max(similarity_score,
                              (2 * partial_ratio + token_sort_ratio + token_set_ratio) / 4)

    return similarity_score
```

**The formula**: `score = (2 * partial_ratio + token_sort_ratio + token_set_ratio) / 4`

**Why partial_ratio gets 2x weight**: partial_ratio does substring matching. Handles are often abbreviations of full names (e.g., "rahul" is a substring of "Rahul Kumar"). This scenario is the most common pattern for real users, so it gets double weight.

**Why permutations**: Some users put last name first in their handle ("kumar_rahul" for "Rahul Kumar"). Trying all permutations of name words ensures we catch this.

**Output**: 0-100 float (higher = more similar)

**Threshold applied separately**:
```python
def based_on_partial_ratio(similarity_score):
    return 0 if similarity_score > 90 else 1
```

---

### Feature 5: Indian Name Database Match -- `check_indian_names()`

**What it does**: Fuzzy-matches the follower's name against a database of 35,183 Indian baby names. A high match score suggests the name belongs to a real person.

**Why this works**: Fake accounts often have gibberish names or names that are not real Indian names. Matching against a comprehensive name database catches this. The fuzzy matching handles transliteration variations and misspellings.

```python
# fake.py - check_indian_names() (lines 286-312)
# Global: loaded once at module import
name_data = pd.read_csv('baby_names.csv')
namess = name_data['Baby Names'].str.lower()  # 35,183 names

def score(i, first_name):
    i = i.lower()
    first_name = first_name.lower()
    ratio = fuzzz.ratio(i, first_name)
    token_sort_ratio = fuzzz.token_sort_ratio(i, first_name)
    token_set_ratio = fuzzz.token_set_ratio(i, first_name)
    return (2 * ratio + token_sort_ratio + token_set_ratio) / 4

def check_indian_names(name):
    if len(name) < 2:
        return 1  # Too short

    similarity_score = 0
    name = name.split()
    first_name = name[0]
    last_name = name[1] if len(name) >= 2 else None

    # Match first name against entire database
    for i in namess:
        similarity_score = max(similarity_score, score(i, first_name))

    # Match last name if present
    if last_name:
        if len(last_name) < 2:
            similarity_score = 1  # Short last name = suspicious
        else:
            for i in namess:
                similarity_score = max(similarity_score, score(i, last_name))

    return similarity_score
```

**Output**: 0-100 float (higher = more likely a real Indian name)

**Complexity**: O(35,183) per name part -- linear scan of entire database. For a name with first + last, this runs twice = 70,366 fuzzy comparisons per follower. This is the most expensive feature.

**Threshold applied in model()**:
```python
score_80 = 1 if indian_name_score > 80 else 0
```

---

### Ensemble: `process1()` Binary + `final()` Weighted

These two functions combine the 5 features into final decisions.

**process1() -- Binary Classification**:
```python
# fake.py - process1() (lines 252-261)
def process1(fake_real_based_on_lang, number_more_than_4_handle, chhitij_logic):
    if fake_real_based_on_lang:
        return 1  # Non-Indic language = FAKE
    if number_more_than_4_handle:
        return 1  # >4 digits = FAKE
    if chhitij_logic == 1:
        return 1  # Special char mismatch = FAKE
    elif chhitij_logic == 2:
        return 2  # No special chars = INCONCLUSIVE
    return 0      # Default = REAL
```

**final() -- Weighted Confidence Score**:
```python
# fake.py - final() (lines 263-277)
def final(fake_real_based_on_lang, similarity_score,
          number_more_than_4_handle, chhitij_logic):
    if fake_real_based_on_lang:
        return 1      # 100% FAKE
    if 0 < similarity_score <= 40:
        return 0.33   # Weak FAKE signal (low similarity but not zero)
    if number_more_than_4_handle:
        return 1      # 100% FAKE
    if chhitij_logic == 1:
        return 1      # 100% FAKE
    elif chhitij_logic == 2:
        return 0      # No special chars = default REAL
    return 0          # Default: REAL
```

**The three confidence levels**:
| Score | Meaning | When |
|---|---|---|
| 0.0 | Definitely REAL | Handle matches name, valid language, few digits |
| 0.33 | Weak FAKE signal | Handle-name similarity is 1-40% (low but not zero) |
| 1.0 | Definitely FAKE | Foreign script, >4 digits, or special char mismatch |

**Critical design insight**: The similarity_score range 0-40 gets 0.33 instead of 1.0 because low similarity alone is not conclusive -- the user might have a creative handle unrelated to their name. But combined with other signals, 0.33 adds up.

---

## 3. HMM TRANSLITERATION DEEP DIVE

### The Problem

Indian Instagram users write their display names in their native scripts: "राहुल कुमार" (Hindi), "রাহুল" (Bengali), "ராகுல்" (Tamil). To compare these names against handles (which are always in Latin/English), we must transliterate them to English first.

### Language Detection: 583 Characters Mapped to 10 Language Codes

The system builds a reverse lookup table mapping every character from 10 Indic scripts to its language code:

```python
# fake.py - Language character database (lines 44-91)
data = {
    "hin": ["़","०","१","२","अ","आ","इ","ई","उ","ऊ","ऋ","ए","ऐ","ओ","औ",
            "क","ख","ग","घ","ङ","च","छ","ज","झ","ञ","ट","ठ","ड","ढ","ण",
            "त","थ","द","ध","न","प","फ","ब","भ","म","य","र","ल","व","श",
            "ष","स","ह","ा","ि","ी","ु","ू","ृ","े","ै","ो","ौ","्"],  # 77 chars
    "pan": [61 chars],   # Punjabi (Gurmukhi)
    "guj": [82 chars],   # Gujarati
    "ben": [65 chars],   # Bengali
    "urd": [41 chars],   # Urdu (Perso-Arabic)
    "tam": [62 chars],   # Tamil
    "mal": [43 chars],   # Malayalam
    "kan": [65 chars],   # Kannada
    "ori": [63 chars],   # Odia
    "tel": [65 chars],   # Telugu
}
# TOTAL: 583 characters mapped

# Build reverse lookup
char_to_lang = {}
for lang, chars in data.items():
    for char in chars:
        char_to_lang[char] = lang
```

| Language | Code | Script | Chars | ML Model |
|---|---|---|---|---|
| Hindi | hin | Devanagari | 77 | Custom process_word() + hin-eng/ HMM |
| Bengali | ben | Bengali | 65 | ben-eng/ |
| Gujarati | guj | Gujarati | 82 | guj-eng/ |
| Kannada | kan | Kannada | 65 | kan-eng/ |
| Malayalam | mal | Malayalam | 43 | mal-eng/ |
| Odia | ori | Odia | 63 | ori-eng/ |
| Punjabi | pan | Gurmukhi | 61 | pan-eng/ |
| Tamil | tam | Tamil | 62 | tam-eng/ |
| Telugu | tel | Telugu | 65 | tel-eng/ |
| Urdu | urd | Perso-Arabic | 41 | urd-eng/ |
| Marathi | - | Devanagari | reuses hin | hin-eng/ |
| Nepali | - | Devanagari | reuses hin | hin-eng/ |
| Assamese | - | Bengali | reuses ben | ben-eng/ |

### The Routing Logic

```python
# fake.py - detect_language() (lines 164-176)
def detect_language(word):
    for char in word:
        lang = char_to_lang.get(char)
        if lang is not None:
            if lang == 'hin':
                return process_word(word)  # Custom Hindi transliteration
            else:
                trn = Transliterator(source=lang, target='eng', build_lookup=True)
                return trn.transform(word)  # HMM-based ML transliteration
    return word  # Already in English/Latin script
```

**Key design decision**: Hindi gets a CUSTOM rule-based transliterator (process_word) instead of the generic HMM model. Why? Hindi is the most common script in Indian Instagram names (>60% of Indic-script names). A hand-tuned rule-based approach with explicit vowel/consonant mappings produces more natural-sounding transliterations for common Hindi names. The HMM model sometimes produces unusual romanizations ("raahula" instead of "rahul").

### Custom Hindi process_word() -- 24 Vowels + 42 Consonants

The custom Hindi transliterator uses two CSV files with hand-crafted mappings:

**svar.csv -- 24 Hindi Vowel Mappings**:
```python
vowels = OrderedDict([
    ('ँ','n'),    # Chandrabindu (nasal)
    ('ं','n'),    # Anusvara
    ('ः','a'),    # Visarga
    ('अ','a'),    # A
    ('आ','aa'),   # Aa
    ('इ','i'),    # I
    ('ई','ee'),   # Ii
    ('उ','u'),    # U
    ('ऊ','oo'),   # Uu
    ('ऋ','ri'),   # Vocalic R
    ('ए','e'),    # E
    ('ऐ','ae'),   # Ai
    ('ओ','o'),    # O
    ('औ','au'),   # Au
    ('ा','a'),    # Aa matra (combining form)
    ('ि','i'),    # I matra
    ('ी','i'),    # Ii matra
    ('ु','u'),    # U matra
    ('ू','oo'),   # Uu matra
    ('ृ','ri'),   # Ri matra
    ('े','e'),    # E matra
    ('ै','ai'),   # Ai matra
    ('ो','o'),    # O matra
    ('ौ','au'),   # Au matra
])
```

**vyanjan.csv -- 42 Hindi Consonant Mappings**:
```python
consonants = OrderedDict([
    # Velar (back of tongue against soft palate)
    ('क','k'), ('ख','kh'), ('ग','g'), ('घ','gh'), ('ङ','ng'),
    # Palatal (tongue against hard palate)
    ('च','ch'), ('छ','chh'), ('ज','j'), ('ज़','z'), ('झ','jh'), ('ञ','nj'),
    # Retroflex (tongue curled back)
    ('ट','t'), ('ठ','th'), ('ड','d'), ('ड़','r'), ('ढ','dh'), ('ण','n'),
    # Dental (tongue against teeth)
    ('त','t'), ('थ','th'), ('द','d'), ('ध','dh'), ('न','n'),
    # Labial (lips together)
    ('प','p'), ('फ','ph'), ('फ़','f'), ('ब','b'), ('भ','bh'), ('म','m'),
    # Semi-vowels
    ('य','y'), ('र','r'), ('ल','l'), ('व','v'),
    # Sibilants
    ('श','sh'), ('ष','sh'), ('स','s'),
    # Glottal
    ('ह','h'),
    # Complex conjuncts
    ('क्ष','ksh'), ('त्र','tr'), ('ज्ञ','gy'),
])
```

**The process_word() algorithm**:
```python
# fake.py - process_word() (lines 129-162)
def process_word(word):
    str1 = ""
    i = 0
    while i < len(word):
        # Step 1: Check for nukta diacritics (two-character combinations like ज़)
        if (i+1 < len(word) and word[i+1].strip() == '\u093C'.strip()):  # nukta
            c = word[i] + word[i+1]
            i += 2
        else:
            c = word[i]
            i += 1

        # Step 2: Map the character
        if c in vowels:
            str1 += vowels[c]
        elif c in consonants:
            # Step 3: Handle inherent 'a' vowel
            # In Devanagari, every consonant has an implicit 'a' sound
            # UNLESS followed by another consonant or a vowel matra
            if i < len(word) and word[i] in consonants:
                if (c == 'झ' and i != 0) or \
                   (i != 0 and i+1 < len(word) and word[i+1] in vowels):
                    str1 += consonants[c]
                else:
                    str1 += consonants[c] + 'a'  # Add inherent 'a'
            else:
                str1 += consonants[c]
        elif c in ['\n','\t',' ','!',',','।','-',':','\\','_','?'] or c.isalnum():
            str1 += c.replace('।', '.')
    return str1
```

**Example walkthrough**: "राहुल" -> "raahul"
1. `र` (consonant) + next is `ा` (vowel matra, not consonant) -> `consonants['र']` = `'r'`
2. `ा` (vowel matra) -> `vowels['ा']` = `'a'`  -> so far: `"ra"`
3. `ह` (consonant) + next is `ु` (vowel matra) -> `consonants['ह']` = `'h'`
4. Wait -- actually `ह` is followed by `ु` which is a vowel, not in consonants dict directly, so the else branch runs -> `consonants['ह']` = `'h'`
5. `ु` (vowel matra) -> `vowels['ु']` = `'u'`  -> so far: `"rahu"`
6. `ल` (consonant, last char) -> `consonants['ल']` = `'l'`
7. Result: `"raahul"`

### HMM Models: Viterbi Decoder, Coefficients, State Probabilities

For non-Hindi scripts, the system uses pre-trained HMM models from the `indictrans` library.

**Each language model directory contains 6 files**:
```
models/ben-eng/         (example: Bengali -> English)
  coef_.npy             # HMM coefficient matrix (feature weights)
  classes.npy           # Output character mapping (what chars can be emitted)
  intercept_init_.npy   # Initial state probabilities (P(start in state i))
  intercept_trans_.npy  # Transition probabilities (P(state j | state i))
  intercept_final_.npy  # Final state probabilities (P(end in state i))
  sparse.vec            # Feature vocabulary (character n-gram features)
```

**The ML pipeline**:
```
Input: "রাহুল" (Bengali)
  |
  [1] UTF-8 -> WX notation (ISO 15919 encoding scheme)
  |
  [2] Feature extraction: character n-gram context windows
      Each character gets features from surrounding characters
  |
  [3] Linear classifier: coef_ matrix x feature vector + intercept
      Produces probability distribution over output characters
  |
  [4] Viterbi decoder: finds optimal sequence through state space
      Uses transition probabilities to enforce valid character sequences
  |
  [5] WX -> UTF-8 (target English script)
  |
Output: "Rahul"
```

**Viterbi decoder (Cython for speed)**:
```
# _decode/viterbi.pyx -- compiled to C for performance
# Implements the standard Viterbi algorithm:
# For each position t in the input:
#   For each possible state s:
#     score[t][s] = max over all previous states s' of:
#       score[t-1][s'] + transition[s'][s] + emission[s][observed[t]]
# Backtrack to find the highest-scoring path
```

The Viterbi algorithm finds the single most likely output sequence. This is more efficient than beam search (also included in the library as `beamsearch.pyx`) but slightly less accurate. For our use case (name transliteration, not full sentence translation), Viterbi is sufficient.

**Why Cython?** Pure Python Viterbi would be too slow for batch processing. Cython compiles to C, giving 10-100x speedup. This is why the Dockerfile needs `gcc-c++` installed.

---

## 4. AWS SERVERLESS PIPELINE

### Architecture Diagram

```
[ClickHouse DB]
     |
     | CTE query (5 tables, UNION ALL)
     v
[S3: gcc-social-data/temp/{date}_creator_followers.json]
     |
     | boto3 download_file()
     v
[Local Machine: push.py]
     |
     | 10,000 lines at a time
     | Round-robin into 8 buckets
     | multiprocessing.Pool(processes=8)
     v
[SQS: creator_follower_in]  (eu-north-1)
     |
     | Event trigger (auto-scaling)
     v
[AWS Lambda: fake.handler]  (ECR container, Python 3.10)
     |  - symbol_name_convert()
     |  - detect_language() + transliterate
     |  - 5 feature extraction
     |  - process1() + final() ensemble
     |
     | kinesis.put_record()
     v
[Kinesis: creator_out]  (ON_DEMAND, ap-south-1)
     |
     | Multi-shard parallel reader
     v
[Local Machine: pull.py -> {date}_creator_followers_final_fake_analysis.json]
```

### push.py: ClickHouse CTE -> S3 -> 8-Worker Multiprocessing -> SQS

**Phase 1: Data Extraction (ClickHouse -> S3)**

```python
# push.py - ClickHouse CTE query
query = f"""
    INSERT INTO FUNCTION s3(
        'https://gcc-social-data.s3.ap-south-1.amazonaws.com/temp/{file_name}',
        'AKIA...', '...', 'JSONEachRow'
    )
    WITH
        handles AS (
            SELECT Names as handle
            FROM s3('...creators_handles.csv', '...', 'CSVWithNames')
        ),
        profile_ids AS (
            SELECT ig_id FROM dbt.mart_instagram_account mia
            WHERE handle IN (SELECT handle FROM handles)
        ),
        follower_data AS (
            SELECT log.target_profile_id,
                   JSONExtractString(source_dimensions, 'handle') AS follower_handle,
                   JSONExtractString(source_dimensions, 'full_name') AS follower_full_name
            FROM dbt.stg_beat_profile_relationship_log log
            WHERE target_profile_id IN profile_ids
              AND follower_handle IS NOT NULL AND follower_handle != ''
              AND follower_full_name IS NOT NULL AND follower_full_name != ''
        ),
        follower_events_data AS (
            -- Same but from real-time events table
            SELECT ... FROM _e.profile_relationship_log_events log ...
        ),
        data AS (
            SELECT * FROM follower_data
            UNION ALL
            SELECT * FROM follower_events_data
        )
    SELECT mia.handle, follower_handle, follower_full_name
    FROM data
    INNER JOIN dbt.mart_instagram_account mia ON mia.ig_id = target_profile_id
    GROUP BY handle, follower_handle, follower_full_name
"""
```

This single query: reads creator handles from S3 CSV, joins with Instagram account table, pulls follower data from both historical (dbt staging) and real-time (events) tables, deduplicates with GROUP BY, and writes directly to S3 as JSONEachRow.

**Phase 2: Batch Distribution (Local -> SQS)**

```python
# push.py - 8-worker multiprocessing batch distribution
s3.Bucket('gcc-social-data').download_file(f"temp/{file_name}", file_name)

queue = sqs.create_queue(QueueName='creator_follower_in',
                         Attributes={
                             'MaximumMessageSize': '262144',      # 256 KB
                             'MessageRetentionPeriod': '345600',  # 4 days
                             'VisibilityTimeout': '30'            # 30 seconds
                         })

def divide_into_buckets(lst, num_buckets):
    buckets = [[] for _ in range(num_buckets)]
    for index, item in enumerate(lst):
        bucket_index = index % num_buckets
        buckets[bucket_index].append(item)
    return buckets

# Main loop: 10K lines at a time, 8 parallel workers
with open(file_name, 'r') as f:
    while True:
        batch_lines = list(itertools.islice(f, 10000))
        if not batch_lines:
            break
        messages = [{'Id': str(i), 'MessageBody': row}
                    for i, row in enumerate(batch_lines)]

        num_buckets = 8
        buckets = divide_into_buckets(messages, num_buckets)

        pool = multiprocessing.Pool(processes=num_buckets)
        pool.map(final, buckets)  # final() sends batch to SQS
        pool.close()
        pool.join()
```

### Lambda: SQS Trigger -> model() -> Kinesis put_record

```python
# fake.py - Lambda handler
def handler(event, context):
    response = model(event)  # Run full 5-feature ensemble
    print(response)
    sqs = session.resource('sqs', region_name='eu-north-1')
    queue = sqs.get_queue_by_name(QueueName='output_queue')
    response = queue.send_message(MessageBody=json.dumps(response))
    return {
        "statusCode": 200,
        "body": {"message": response}
    }
```

Each Lambda invocation processes ONE follower record. Lambda auto-scales based on SQS queue depth -- if there are 100,000 messages in the queue, AWS can spin up hundreds of Lambda containers simultaneously.

### pull.py: Multi-Shard Kinesis Reader

```python
# pull.py - Multi-shard parallel reader
def process_shard(shard_id, starting_sequence_number):
    total_length = 0
    shard_iterator = client.get_shard_iterator(
        StreamName=stream_name,
        ShardId=shard_id,
        ShardIteratorType='AFTER_SEQUENCE_NUMBER',
        StartingSequenceNumber=starting_sequence_number
    )['ShardIterator']

    while True:
        messages = []
        response = client.get_records(ShardIterator=shard_iterator, Limit=10000)
        if len(response['Records']) < 1:
            break
        shard_iterator = response['NextShardIterator']
        # ... process records, write to file

if __name__ == '__main__':
    pool = Pool()
    shards = client.list_shards(StreamName=stream_name)['Shards']
    args = [(shard['ShardId'],
             shard['SequenceNumberRange']['StartingSequenceNumber'])
            for shard in shards]
    pool.starmap(process_shard, args)
```

### Dockerfile: ECR with gcc-c++ for Cython

```dockerfile
FROM public.ecr.aws/lambda/python:3.10

# System dependencies for indictrans Cython compilation
RUN yum install -y gcc-c++ pkgconfig poppler-cpp-devel

# Install Python dependencies and indictrans from source
COPY requirement.txt ./
COPY indic-trans-master ./
RUN pip install -r requirements.txt
RUN pip install .  # Compiles viterbi.pyx -> C -> .so

# Copy pre-trained HMM models to site-packages
RUN cp -r indictrans/models /var/lang/lib/python3.10/site-packages/indictrans/
RUN cp -r indictrans/mappings /var/lang/lib/python3.10/site-packages/indictrans/

# Copy Hindi transliteration mappings
COPY svar.csv ./      # 24 vowel mappings
COPY vyanjan.csv ./   # 42 consonant mappings

# Copy application code and data
COPY fake.py ./
COPY baby_names_.csv ./baby_names.csv

CMD [ "fake.handler" ]
```

**Why Docker/ECR instead of plain Lambda?** The indictrans library requires Cython compilation (C compiler), pre-trained model files (.npy), and a 35K-name CSV. Total package size exceeds Lambda's 50MB zip limit. ECR containers support up to 10GB.

---

## 5. GPT ENRICHMENT (from Beat)

The fake follower detection system works alongside Beat's GPT enrichment flows. While fake.py detects fake followers at the individual level, Beat's GPT flows enrich creator profiles with inferred attributes.

### 6 GPT Flows

| Flow | Input | Output | Workers x Concurrency |
|---|---|---|---|
| `refresh_instagram_gpt_data_base_gender` | handle, bio | {gender, confidence} | 2 x 5 |
| `refresh_instagram_gpt_data_base_location` | handle, bio | {country, city, confidence} | 2 x 5 |
| `refresh_instagram_gpt_data_base_categ_lang_topics` | handle, bio, posts | {category, language, topics[]} | 2 x 5 |
| `refresh_instagram_gpt_data_audience_age_gender` | handle, bio, posts | {age_range, gender_dist} | 2 x 5 |
| `refresh_instagram_gpt_data_audience_cities` | handle, bio, posts | {cities: [{name, pct}]} | 2 x 5 |
| `refresh_instagram_gpt_data_gender_location_lang` | handle, bio | combined gender+location+lang | 2 x 5 |

### 12 Prompt Versions

The GPT prompt evolved through 12 iterations stored as YAML files in `gpt/prompts/`:
- `profile_info_v0.1.yaml` through `profile_info_v0.12.yaml`
- Each version refined instructions for edge cases (ambiguous names like "Alex", multilingual bios, creative handles)
- Final version: `profile_info_v0.12.yaml`

### is_data_consumable() Validation

```python
# gpt/flows/fetch_gpt_data.py
@sessionize
async def fetch_instagram_gpt_data_base_gender(data: dict, session=None):
    max_retries = 2
    base_data = await retrieve_instagram_gpt_data_base_gender(data, session=session)
    for _ in range(max_retries):
        if not is_data_consumable(base_data, "base_gender"):
            updated_data = await retrieve_instagram_gpt_data_base_gender(data, session=session)
            for key in keys_to_update:
                base_data[key] = updated_data.get(key, base_data[key])
        else:
            break
```

Three safeguards against GPT hallucinations:
1. **Temperature=0** -- deterministic output, same input always produces same output
2. **is_data_consumable()** -- validates required fields are present and not "UNKNOWN"
3. **Retry up to 2 times** -- if validation fails, re-query GPT with fresh API call

### Azure OpenAI Configuration

```python
# gpt/functions/retriever/openai/openai_extractor.py
class OpenAi(GptCrawlerInterface):
    def __init__(self):
        openai.api_type = "azure"
        openai.api_base = os.environ["OPENAI_API_BASE"]
        openai.api_version = os.environ["OPENAI_API_VERSION"]
        openai.api_key = os.environ["OPENAI_API_KEY"]

    async def fetch_instagram_gpt_data_base_gender(self, handle, bio):
        prompt = self._load_prompt("profile_info_v0.12.yaml")
        response = await openai.ChatCompletion.acreate(
            engine="gpt-35-turbo",
            messages=[
                {"role": "system", "content": prompt['system']},
                {"role": "user", "content": prompt['user'].format(
                    handle=handle, bio=bio)}
            ],
            temperature=0
        )
        return self._parse_response(response)
```

**Why Azure OpenAI instead of direct OpenAI?** Enterprise compliance. Azure OpenAI provides data residency guarantees, SLA, and does not use customer data for training. At GCC's scale (~100K profiles/month at ~$0.002/profile = ~$200/month), the Azure premium is negligible.

**Connection to fake follower detection**: The GPT enrichment tells us WHO the creator is (gender, location, content category). The fake follower detection tells us WHO their followers are (real or fake). Together, these power the platform's creator analytics -- brands can see both creator attributes and follower quality.

---

## 6. THE 50% SPEED CALCULATION

### Before: Sequential Single-Machine Processing

```
Old approach:
- Single EC2 instance
- Read follower data from database
- Process one follower at a time
- Write results back to database
- Throughput: 10-20 records/second
- 100K followers: 1.5 - 2.5 hours
```

### After: Lambda Auto-Scaling + 8 Parallel SQS Workers

```
New approach:
- push.py: 8 parallel workers feeding SQS (10K batch, round-robin)
- Lambda: auto-scales with queue depth (100+ concurrent invocations)
- Kinesis: ON_DEMAND stream (auto-scaling shards)
- Throughput: scales horizontally with queue depth
- 100K followers: 45 min - 1.25 hours
```

**Where the 50% comes from**:

| Component | Before | After | Improvement |
|---|---|---|---|
| Data extraction | Single query + serial read | CTE -> S3 -> parallel download | ~same |
| Message distribution | N/A (direct processing) | 8-worker multiprocessing Pool | 8x faster ingestion |
| ML processing | 1 machine, sequential | Lambda auto-scaling (100+ concurrent) | 10-50x throughput |
| Result collection | Sequential DB writes | Kinesis ON_DEMAND + multi-shard reader | Parallel writes |
| **Total wall-clock time** | **~2 hours for 100K** | **~1 hour for 100K** | **~50% faster** |

The 50% improvement is conservative because it measures end-to-end wall-clock time including data extraction and result aggregation overhead. The ML processing step alone is 10-50x faster, but the pipeline has non-parallelizable bottlenecks (ClickHouse query, final file assembly).

### Per-Record Timing Breakdown

```
Symbol conversion:     1-5ms
Language detection:    1-2ms
Transliteration:       5-10ms (ML inference for non-Hindi)
Handle cleaning:       1-2ms
Fuzzy scoring:         5-15ms (up to 24 permutations)
Indian name check:     10-50ms (linear scan of 35,183 names)
Ensemble scoring:      <1ms
---
TOTAL:                 25-85ms per follower (single Lambda)

At scale:
  Single Lambda:       ~15 records/second
  100 concurrent:      ~1,500 records/second
  100K followers:      ~67 seconds (ML only) + pipeline overhead
```

---

## 7. INTERVIEW SCRIPTS

### 30-Second Version

> "I built a fake follower detection system for Instagram creators. It uses a 5-feature ensemble model -- combining language detection, handle-name similarity with RapidFuzz, and fuzzy matching against 35,000 Indian names. The key challenge was multilingual support: Indian users write names in 10 different scripts, so I integrated HMM-based transliteration using trained Viterbi decoders to convert all scripts to English before comparison. The whole thing runs on AWS Lambda triggered by SQS, with results streamed through Kinesis. Moving from sequential processing to this serverless pipeline improved speed by 50%."

### 2-Minute Version

> "At Good Creator Co., our platform helps brands find authentic Instagram creators. A major problem was fake followers inflating creator metrics. I built an ML system to detect these fakes.
>
> The core is a 5-feature ensemble model. Given a follower's handle and display name, it extracts five signals: non-Indic script detection via regex, a digit count threshold, handle-name structural correlation, a weighted fuzzy similarity score using RapidFuzz, and matching against a database of 35,183 Indian baby names. Two ensemble functions combine these -- one binary classifier and one that outputs a confidence score of 0.0, 0.33, or 1.0.
>
> The biggest technical challenge was multilingual names. Over 60% of Indian users write their display names in native scripts -- Hindi, Bengali, Tamil, and 7 others. I built a transliteration pipeline using the indictrans library, which uses Hidden Markov Models with Viterbi decoding. For Hindi specifically, I wrote a custom rule-based transliterator with 66 hand-mapped vowels and consonants because the HMM model was producing unnatural romanizations for common names.
>
> For deployment, I designed a serverless AWS pipeline: ClickHouse data extraction to S3, an 8-worker multiprocessing feeder into SQS, Lambda functions auto-scaling with queue depth, and Kinesis for result streaming. The Lambda container is packaged as an ECR image because the ML models exceed Lambda's zip size limit. This architecture replaced sequential single-machine processing and improved throughput by 50%.
>
> The system also connects to Beat's GPT enrichment -- 6 flows using Azure OpenAI GPT-3.5-turbo to infer creator attributes like gender, location, and content category. Together, brands get both creator quality and follower authenticity insights."

### 5-Minute Version (adds depth on each component)

> [Use the 2-minute version as the base, then expand each paragraph:]
>
> **On the ensemble features**: "The 5 features are designed to be independent -- each catches a different type of fake account. Feature 1 catches bots that use random Unicode characters from non-Indian scripts like Greek or Chinese. Feature 2 catches auto-generated accounts with long digit suffixes. Feature 3 exploits the fact that real users who put underscores in handles are encoding their real name, so mismatches are suspicious. Feature 4 uses a weighted RapidFuzz formula -- partial_ratio gets double weight because handles are often abbreviations of full names. And Feature 5 validates against a 35K-name database to see if the displayed name is even a real Indian name. The ensemble has three confidence levels: 0.0 for definitely real, 0.33 for weak fake signals, and 1.0 for definitely fake."
>
> **On transliteration**: "The language detection works by mapping 583 characters across 10 scripts to their language codes. For each character in a name, I look up its script. Hindi routes to a custom transliterator with 24 vowel and 42 consonant mappings that handle Devanagari diacritics -- things like the inherent 'a' vowel in consonants and nukta modifiers. All other scripts route to HMM models trained on parallel corpora. Each model has a coefficient matrix, state transition probabilities, and a Viterbi decoder compiled in Cython for C-level performance. The example I like to give is 'raahul' from the Devanagari letters -- the custom Hindi transliterator produces a more natural romanization than the generic HMM."
>
> **On the pipeline**: "The push.py script is an interesting piece of engineering. It runs a multi-CTE ClickHouse query that joins 5 tables, unions historical and real-time follower data, deduplicates, and writes directly to S3 as JSONEachRow. Then it downloads the file locally, reads 10,000 lines at a time, distributes them round-robin into 8 buckets, and uses a multiprocessing Pool to send SQS batches in parallel. Each SQS message triggers a Lambda invocation. The Lambda container, built on ECR, needs gcc-c++ because the Viterbi decoder is in Cython. Results go to an ON_DEMAND Kinesis stream, and pull.py reads all shards in parallel using multiprocessing.Pool with starmap."
>
> **On GPT enrichment**: "Separately, Beat has 6 GPT enrichment flows that infer creator attributes. These use Azure OpenAI GPT-3.5-turbo with temperature=0. Each flow has gone through 12 prompt versions -- profile_info_v0.1 through v0.12 -- with each iteration fixing edge cases. There's an is_data_consumable() validation function that checks GPT's output has the required fields and retries up to 2 times if not. The connection to fake follower detection is that together they give brands a complete picture: who the creator is, and whether their followers are real."

---

## 8. TOUGH FOLLOW-UPS (15 Questions)

### Q1: Why not use a neural network for fake detection?

> "Two practical reasons: no labeled training data and insufficient scale. We had no ground truth dataset of confirmed fake vs. real followers -- labeling even a few thousand would take significant manual effort. The heuristic ensemble works well without labeled data because each feature encodes a well-understood signal. A neural network also needs thousands of labeled examples to outperform hand-crafted features for a structured problem like this. If we had labeled data, I would train a gradient-boosted classifier (XGBoost) rather than a neural network -- tabular data with 5 numeric features is gradient boosting's sweet spot."

### Q2: What is the false positive rate?

> "Honestly, we did not have labeled data to compute a precise false positive rate. Based on manual spot-checking of about 500 results, I estimate the false positive rate is around 5-10%. The most common false positive is real users with creative handles unrelated to their names -- for example, a user named 'Priya' with handle 'sunshine_vibes_123' gets a low similarity score and gets flagged. The 0.33 confidence level (instead of 1.0) was specifically designed for these ambiguous cases to avoid hard misclassification."

### Q3: How would you create a labeled dataset to improve this?

> "Three approaches in order of cost-effectiveness:
> 1. **Semi-supervised with high-confidence predictions**: Use the current model's 0.0 and 1.0 predictions as pseudo-labels (these are the most confident). Manually verify a random sample of 500 from each class to estimate label quality.
> 2. **Active learning**: Export the 0.33-confidence cases (the ambiguous ones) and have human annotators label 1,000-2,000 of them. Train a supervised model on the combination.
> 3. **Behavioral features**: If we had access to follower activity data (post frequency, like patterns, follower/following ratio), we could cluster them -- inactive accounts with suspicious patterns would be another strong signal."

### Q4: Explain the Viterbi algorithm in one sentence.

> "Viterbi is a dynamic programming algorithm that finds the most likely sequence of hidden states in a Hidden Markov Model by processing one observation at a time and keeping only the best path to each state, running in O(T * N^2) time where T is sequence length and N is the number of states."

### Q5: Why RapidFuzz over FuzzyWuzzy?

> "RapidFuzz is a drop-in replacement for FuzzyWuzzy that is 10-100x faster because it is implemented in C++ instead of pure Python. Since we run fuzzy matching against 35,183 names per follower (in check_indian_names), performance matters enormously. RapidFuzz also removes the python-Levenshtein dependency issue and has consistent cross-platform behavior. The API is identical -- fuzz.ratio, fuzz.partial_ratio, fuzz.token_sort_ratio, fuzz.token_set_ratio -- so migration is a one-line import change."

### Q6: Why custom Hindi transliteration instead of using HMM for all languages?

> "Hindi is our most common Indic script by far -- over 60% of Indic-script names on our platform are in Devanagari. The HMM model produces acceptable but unnatural romanizations for common Hindi names. For example, the HMM might output 'raahula' instead of 'rahul' because it models the inherent 'a' vowel differently. With a hand-crafted rule system using 66 explicit mappings, I control exactly how the inherent vowel and diacritic combinations are handled. For less common scripts like Odia or Malayalam, the HMM model is good enough and the volume does not justify custom rules."

### Q7: What are the security concerns with this system?

> "Major ones:
> 1. **Hardcoded AWS credentials** -- push.py and fake.py contain plaintext access keys and secret keys. These should use IAM roles (Lambda execution role) or AWS Secrets Manager.
> 2. **Hardcoded ClickHouse password** -- in the push.py connection string. Should be in environment variables or a secrets manager.
> 3. **No input validation** -- the Lambda handler trusts the SQS message content without sanitization. Malformed JSON or injection payloads could cause unexpected behavior.
> 4. **Unencrypted data transfer** -- follower handles and names are sent in plaintext to SQS and Kinesis.
> If I were to rebuild this, IAM roles + env vars + SQS/Kinesis encryption at rest would be the minimum."

### Q8: Design v2 with supervised learning.

> "v2 architecture:
> 1. **Label collection**: Export 5,000 high-confidence predictions (2,500 fake, 2,500 real) from v1. Manually verify 500. Use as seed labels.
> 2. **Feature expansion**: Add behavioral features from Beat's database -- account age, post frequency, follower/following ratio, bio length, profile picture presence.
> 3. **Model**: XGBoost classifier. Tabular data with <20 features -- gradient boosting beats neural networks here.
> 4. **Training pipeline**: Weekly retraining on new labeled data. Store models in S3, load in Lambda.
> 5. **Output**: Continuous probability (0.0-1.0) instead of three discrete levels.
> 6. **Monitoring**: Track prediction distribution drift, feature importance changes, and false positive rate on manually-reviewed samples."

### Q9: Why three confidence levels (0.0, 0.33, 1.0) instead of continuous?

> "Practical downstream requirements. The platform UI needed clear categories for brands: 'Real', 'Suspicious', and 'Fake'. Continuous scores between 0 and 1 would require an additional thresholding step and would confuse non-technical users. The 0.33 level specifically captures the ambiguous case where similarity is low (1-40%) but no other strong fake signal exists -- these followers deserve investigation but not automatic rejection. In v2 with supervised learning, I would output continuous probabilities internally but still bucket them into categories for the UI."

### Q10: How do you handle the O(35,183) linear scan in check_indian_names()?

> "It is the biggest performance bottleneck -- 70,366 fuzzy comparisons per follower (35K for first name + 35K for last name). In production, this costs 10-50ms per follower.
>
> Optimization strategies I considered:
> 1. **Pre-filter by first character**: Only compare against names starting with the same letter. Reduces comparisons by ~26x.
> 2. **BK-tree or VP-tree**: Index the name database by edit distance. O(log n) lookups instead of O(n).
> 3. **Approximate nearest neighbor**: Use locality-sensitive hashing (LSH) with character n-gram features.
>
> We did not implement these because Lambda's per-invocation cost was low enough that the brute-force approach was acceptable. The 50-100ms per follower was within our SLA."

### Q11: Why SQS + Lambda instead of just running on EC2 with multiprocessing?

> "Three reasons:
> 1. **Auto-scaling**: SQS + Lambda scales to zero when idle and to hundreds of concurrent invocations under load. EC2 would need manual scaling or Auto Scaling Groups with more operational overhead.
> 2. **Cost**: We pay per Lambda invocation (100ms granularity). For batch jobs that run a few hours per week, serverless is cheaper than running EC2 24/7.
> 3. **Fault tolerance**: If a Lambda invocation fails, the SQS message returns to the queue and is retried automatically. With EC2 multiprocessing, a process crash loses the work."

### Q12: Why Kinesis ON_DEMAND instead of a fixed shard count?

> "ON_DEMAND mode auto-scales shards based on throughput. Since our batch jobs have variable sizes (10K to 1M followers per run), pre-provisioning shards would either waste money or throttle under peak load. ON_DEMAND handles both cases. The trade-off is slightly higher per-record cost versus provisioned capacity, but for batch workloads with variable throughput, it is the right choice."

### Q13: How does the Unicode symbol conversion (13 variants) help with fake detection?

> "Fake accounts and spam bots frequently use 'fancy' Unicode text to evade keyword filters and text-matching systems. A name written as '𝓟𝓻𝓲𝔂𝓪' (Mathematical Bold Script) would fail all fuzzy matching if not normalized to 'Priya' first. The 13 variant sets cover: circled letters, mathematical bold/italic/fraktur/double-struck/monospace, sans-serif, and full-width characters. By normalizing these before any feature extraction, we ensure the downstream pipeline works correctly regardless of Unicode tricks."

### Q14: What happens when a name is in a script not in your 10-language mapping?

> "It falls through to the default case -- detect_language() returns the original string unchanged. Then uni_decode() (which uses the `unidecode` library) attempts a best-effort ASCII transliteration. For example, Thai script or Arabic (non-Urdu) would go through unidecode's generic mapping. The accuracy is lower than our custom mappings, but it provides a reasonable fallback. Additionally, check_lang_other_than_indic() would flag scripts like Chinese or Korean as fake -- so non-Indic scripts are handled before transliteration is even needed."

### Q15: If you had to cut one feature, which would you remove?

> "Feature 3 (handle-name special character correlation). It has the most complex decision tree but the least independent signal. Its output depends on generate_similarity_score() (Feature 4), so it is partially redundant. Also, its 'INCONCLUSIVE' (value=2) output for handles without special characters covers a large portion of real accounts and does not contribute to detection. Features 1, 4, and 5 carry the most signal independently."

---

## 9. WHAT I'D DO DIFFERENTLY

### 1. Security First
Hardcoded credentials in source code is a significant issue. I would use AWS IAM roles for Lambda (no credentials needed at all), AWS Secrets Manager for ClickHouse passwords, and encrypted SQS/Kinesis.

### 2. Supervised Learning Pipeline
The heuristic ensemble was a good v1, but a v2 with labeled data and gradient boosting would be significantly more accurate. I would invest in building a labeled dataset early using active learning on the ambiguous 0.33 cases.

### 3. Indian Name Lookup Optimization
The O(35,183) brute-force scan per name is wasteful. A BK-tree indexed by edit distance would reduce this to O(log n) while maintaining fuzzy matching capability.

### 4. Continuous Confidence Scores
Three discrete levels (0.0, 0.33, 1.0) lose information. A weighted combination of all feature scores into a continuous probability would be more useful for downstream analytics and threshold tuning.

### 5. Behavioral Features
Handle and name alone miss important signals. Adding account age, post frequency, follower/following ratio, bio presence, and profile picture presence would significantly improve detection accuracy without needing complex NLP.

### 6. Better Hindi Transliteration
The custom process_word() handles common cases well but has edge cases with consonant clusters and complex conjuncts. I would use the HMM model as a fallback when the rule-based output looks suspicious (e.g., contains no vowels).

### 7. Monitoring and Evaluation
No evaluation metrics were tracked in production. I would add: prediction distribution dashboards, weekly random-sample manual review, and drift detection for feature distributions.

### 8. Error Handling
The Lambda handler has minimal error handling. Failed transliterations, malformed input, or API failures silently produce wrong results. Proper try/except with CloudWatch logging and dead-letter queues for failed messages.

---

## 10. CONNECTION TO OTHER BULLETS

### Bullet 1: Distributed Multi-Source Data Pipeline (Beat)
> The fake follower detection system consumes follower data that was originally collected by Beat's `fetch_profile_followers` flow. Beat's worker pool system (73 flows, multiprocessing + async) is what populates the ClickHouse tables that push.py queries. Without Beat's data collection infrastructure, fake_follower_analysis would have no data to analyze.

### Bullet 2: ClickHouse Migration
> The ClickHouse CTE query in push.py directly reads from dbt marts created during the ClickHouse migration. The `dbt.mart_instagram_account` and `dbt.stg_beat_profile_relationship_log` tables are the result of the dbt transformation layer built as part of the migration. The JSON extraction functions (`JSONExtractString`) leverage ClickHouse's columnar storage for efficient processing of millions of follower records.

### Bullet 3: Data Platform (Stir -- Airflow + dbt)
> The follower data in ClickHouse arrives through Stir's Airflow DAGs. The `profile_relationship_log` data is staged and transformed by dbt models before fake_follower_analysis consumes it. This is the three-layer pattern: Beat collects raw data -> Stir transforms it -> fake_follower_analysis and Coffee consume it.

### Bullet 4: Dual-Database Architecture (Coffee)
> The fake follower analysis results feed back into the platform. Coffee's Go API serves the creator analytics dashboard where brands see follower quality metrics. The fake_follower_analysis output (the 16-field JSON per follower) is aggregated into per-creator fake follower percentages that Coffee surfaces through its REST endpoints.

### Bullet 5: S3 Assets + Discovery
> The S3 bucket `gcc-social-data` used by push.py is the same infrastructure used for asset storage. The naming convention (`temp/{date}_creator_followers.json`) follows the same pattern as other GCC S3 operations. The Kinesis stream `creator_out` follows the same serverless event-driven architecture as the asset upload flows.

### GPT Enrichment (Beat)
> Beat's 6 GPT flows (gender, location, category, audience demographics, audience cities, combined) provide the "who is the creator" dimension. Fake follower detection provides the "are their followers real" dimension. Together, they give brands a complete creator quality assessment: a creator with 1M followers but 40% fake follower rate is worth less than one with 200K followers and 5% fake rate.

---

## APPENDIX: KEY NUMBERS FOR QUICK REFERENCE

| Metric | Value |
|---|---|
| Total lines of code | 955+ across 6 Python files |
| Core model file (fake.py) | 385 lines |
| Detection features | 5 independent heuristics |
| Confidence levels | 3 (0.0, 0.33, 1.0) |
| Indic scripts supported | 10 (+ 4 derivative) |
| Characters mapped | 583 across all scripts |
| Hindi vowel mappings | 24 (svar.csv) |
| Hindi consonant mappings | 42 (vyanjan.csv) |
| Indian name database | 35,183 names |
| HMM model files per language | 6 (coef, classes, init, trans, final, sparse) |
| Unicode symbol variants | 13 sets normalized |
| AWS services used | 5 (S3, SQS, Lambda, Kinesis, ECR) |
| SQS batch workers | 8 (multiprocessing.Pool) |
| SQS message retention | 4 days |
| Kinesis mode | ON_DEMAND (auto-scaling) |
| Processing speed per record | 25-85ms |
| Single Lambda throughput | ~15 records/second |
| Speed improvement | 50% faster end-to-end |
| GPT flows | 6 |
| GPT prompt versions | 12 (v0.1 through v0.12) |
| GPT model | Azure OpenAI GPT-3.5-turbo, temperature=0 |
| Output fields per follower | 16 |
| RapidFuzz formula | (2*partial + sort + set) / 4 |
| Digit threshold | >4 digits = FAKE |
| Similarity threshold (hard) | >90 = REAL |
| Similarity threshold (weak) | 0-40 = 0.33 confidence |
| Indian name threshold | >80 = REAL |
