package constants

const ApplicationName = "coffee"
const AppContextKey = "appContext"
const LogLevelDebug = "debug"
const LogLevelInfo = "info"
const LogLevelWarn = "warn"
const LogLevelError = "error"
const LogLevelConfig = "LOG_LEVEL"
const ServerPortConfig = "SERVER_PORT"
const PgHostConfig = "PG_HOST"
const PgDbConfig = "PG_DB"
const PgPortConfig = "PG_PORT"
const PgMaxIdleConnConfig = "PG_MAX_IDLE_CONN"
const PgMaxOpenConnConfig = "PG_MAX_OPEN_CONN"
const ChMaxIdleConnConfig = "CH_MAX_IDLE_CONN"
const ChMaxOpenConnConfig = "CH_MAX_OPEN_CONN"
const ServerTimeoutConfig = "SERVER_TIMEOUT"
const SentryDsn = "https://d3ac0488c9fe429181a95959027acad3@sentry.bulbul.tv/38"

var LanguageCodes = map[string]string{
	"Assamese":  "as",
	"Bengali":   "bn",
	"Bodo":      "brx",
	"Dogri":     "doi",
	"English":   "en",
	"Gujarati":  "gu",
	"Hindi":     "hi",
	"Kannada":   "kn",
	"Kashmiri":  "ks",
	"Konkani":   "gom",
	"Maithili":  "mai",
	"Malayalam": "ml",
	"Manipuri":  "mni",
	"Marathi":   "mr",
	"Nepali":    "ne",
	"Odia":      "or",
	"Punjabi":   "pa",
	"Sanskrit":  "sa",
	"Santali":   "sat",
	"Sindhi":    "sd",
	"Tamil":     "ta",
	"Telugu":    "te",
	"Urdu":      "ur",
}

var LanguageNames = map[string]string{
	"as":  "Assamese",
	"bn":  "Bengali",
	"brx": "Bodo",
	"doi": "Dogri",
	"en":  "English",
	"gu":  "Gujarati",
	"hi":  "Hindi",
	"kn":  "Kannada",
	"ks":  "Kashmiri",
	"gom": "Konkani",
	"mai": "Maithili",
	"ml":  "Malayalam",
	"mni": "Manipuri",
	"mr":  "Marathi",
	"ne":  "Nepali",
	"or":  "Odia",
	"pa":  "Punjabi",
	"sa":  "Sanskrit",
	"sat": "Santali",
	"sd":  "Sindhi",
	"ta":  "Tamil",
	"te":  "Telugu",
	"ur":  "Urdu",
}

var LanguageNamesExhaustive = map[string]string{
	"zxx": "Unknown",
	"aa":  "Afar",
	"ab":  "Abkhazian",
	"af":  "Afrikaans",
	"ak":  "Akan",
	"sq":  "Albanian",
	"am":  "Amharic",
	"ar":  "Arabic",
	"an":  "Aragonese",
	"hy":  "Armenian",
	"as":  "Assamese",
	"av":  "Avaric",
	"ae":  "Avestan",
	"ay":  "Aymara",
	"az":  "Azerbaijani",
	"bm":  "Bambara",
	"ba":  "Bashkir",
	"eu":  "Basque",
	"be":  "Belarusian",
	"bn":  "Bengali",
	"bh":  "Bihari languages",
	"bi":  "Bislama",
	"bs":  "Bosnian",
	"br":  "Breton",
	"bg":  "Bulgarian",
	"my":  "Burmese",
	"ca":  "Catalan",
	"ch":  "Chamorro",
	"ce":  "Chechen",
	"ny":  "Chichewa",
	"zh":  "Chinese",
	"cu":  "Church Slavic",
	"cv":  "Chuvash",
	"kw":  "Cornish",
	"co":  "Corsican",
	"cr":  "Cree",
	"hr":  "Croatian",
	"cs":  "Czech",
	"da":  "Danish",
	"dv":  "Dhivehi",
	"nl":  "Dutch",
	"dz":  "Dzongkha",
	"en":  "English",
	"eo":  "Esperanto",
	"et":  "Estonian",
	"ee":  "Ewe",
	"fo":  "Faroese",
	"fj":  "Fijian",
	"fi":  "Finnish",
	"fr":  "French",
	"ff":  "Fulah",
	"gl":  "Galician",
	"ka":  "Georgian",
	"de":  "German",
	"el":  "Greek",
	"gn":  "Guarani",
	"gu":  "Gujarati",
	"ht":  "Haitian",
	"ha":  "Hausa",
	"he":  "Hebrew",
	"hz":  "Herero",
	"hi":  "Hindi",
	"ho":  "Hiri Motu",
	"hu":  "Hungarian",
	"ia":  "Interlingua",
	"id":  "Indonesian",
	"ie":  "Interlingue",
	"ga":  "Irish",
	"ig":  "Igbo",
	"ik":  "Inupiaq",
	"io":  "Ido",
	"is":  "Icelandic",
	"it":  "Italian",
	"iu":  "Inuktitut",
	"ja":  "Japanese",
	"jv":  "Javanese",
	"kl":  "Kalaallisut",
	"kn":  "Kannada",
	"kr":  "Kanuri",
	"ks":  "Kashmiri",
	"kk":  "Kazakh",
	"km":  "Khmer",
	"ki":  "Kikuyu",
	"rw":  "Kinyarwanda",
	"ky":  "Kirghiz",
	"kv":  "Komi",
	"kg":  "Kongo",
	"ko":  "Korean",
	"ku":  "Kurdish",
	"kj":  "Kwanyama",
	"la":  "Latin",
	"lb":  "Luxembourgish",
	"lg":  "Ganda",
	"li":  "Limburgish",
	"ln":  "Lingala",
	"lo":  "Lao",
	"lt":  "Lithuanian",
	"lu":  "Luba-Katanga",
	"lv":  "Latvian",
	"gv":  "Manx",
	"mk":  "Macedonian",
	"mg":  "Malagasy",
	"ms":  "Malay",
	"ml":  "Malayalam",
	"mt":  "Maltese",
	"mi":  "Maori",
	"mr":  "Marathi",
	"mh":  "Marshallese",
	"mn":  "Mongolian",
	"na":  "Nauru",
	"nv":  "Navajo",
	"nd":  "North Ndebele",
	"ne":  "Nepali",
	"ng":  "Ndonga",
	"nb":  "Norwegian Bokmål",
	"nn":  "Norwegian Nynorsk",
	"no":  "Norwegian",
	"ii":  "Sichuan Yi",
	"nr":  "South Ndebele",
	"oc":  "Occitan",
	"oj":  "Ojibwa",
	"om":  "Oromo",
	"or":  "Oriya",
	"os":  "Ossetian",
	"pa":  "Punjabi",
	"pi":  "Pali",
	"fa":  "Persian",
	"pl":  "Polish",
	"ps":  "Pashto",
	"pt":  "Portuguese",
	"qu":  "Quechua",
	"rm":  "Romansh",
	"rn":  "Rundi",
	"ro":  "Romanian",
	"ru":  "Russian",
	"sa":  "Sanskrit",
	"sc":  "Sardinian",
	"sd":  "Sindhi",
	"se":  "Northern Sami",
	"sm":  "Samoan",
	"sg":  "Sango",
	"sr":  "Serbian",
	"gd":  "Gaelic",
	"sn":  "Shona",
	"si":  "Sinhalese",
	"sk":  "Slovak",
	"sl":  "Slovenian",
	"so":  "Somali",
	"st":  "Southern Sotho",
	"es":  "Spanish",
	"su":  "Sundanese",
	"sw":  "Swahili",
	"ss":  "Swati",
	"sv":  "Swedish",
	"ta":  "Tamil",
	"te":  "Telugu",
	"tg":  "Tajik",
	"th":  "Thai",
	"ti":  "Tigrinya",
	"bo":  "Tibetan",
	"tk":  "Turkmen",
	"tl":  "Tagalog",
	"tn":  "Tswana",
	"to":  "Tonga",
	"tr":  "Turkish",
	"ts":  "Tsonga",
	"tt":  "Tatar",
	"tw":  "Twi",
	"ty":  "Tahitian",
	"ug":  "Uighur",
	"uk":  "Ukrainian",
	"ur":  "Urdu",
	"uz":  "Uzbek",
	"ve":  "Venda",
	"vi":  "Vietnamese",
	"vo":  "Volapük",
	"wa":  "Walloon",
	"cy":  "Welsh",
	"wo":  "Wolof",
	"fy":  "Western Frisian",
	"xh":  "Xhosa",
	"yi":  "Yiddish",
	"yo":  "Yoruba",
	"za":  "Zhuang",
	"zu":  "Zulu",
}

var CategoryNames = map[string]string{
	"1":  "Food",
	"2":  "Fashion",
	"3":  "Beauty",
	"4":  "Lifestyle",
	"5":  "Travel",
	"6":  "Photography",
	"7":  "Art",
	"8":  "Make-up",
	"11": "Gaming",
	"12": "Fitness",
	"13": "Health",
	"14": "Entertainment",
	"15": "Gadgets & Tech",
	"16": "Finance",
	"17": "Luxury",
	"18": "Animal/Pet",
	"19": "Self Improvement",
	"20": "Parenting",
	"21": "Books",
	"22": "Sports",
	"23": "Education",
	"24": "Automobile",
	"25": "Business",
	"26": "Wedding",
	"27": "Decor",
}

var CreatorProgramsNameToId = map[string]int{
	"BEAUTY":         1,
	"GOOD_PARENTING": 2,
	"GOOD_LIFE":      3,
}

var CreatorProgramsIdToName = map[int]string{
	1: "BEAUTY",
	2: "GOOD_PARENTING",
	3: "GOOD_LIFE",
}

var CreatorCohortsNameToId = map[string]int{
	"CONTENT": 1,
	"ORDERS":  2,
	"VIEWS":   3,
}

var CreatorCohortsIdToName = map[int]string{
	1: "CONTENT",
	2: "ORDERS",
	3: "VIEWS",
}

var CountryCodes = map[string]string{
	"EN":    "England",
	"RU":    "Russia",
	"KP":    "North Korea",
	"DK":    "Denmark",
	"SV":    "El Salvador",
	"UAE":   "United Arab Emirates",
	"CZ":    "Czech Republic",
	"KR":    "South Korea",
	"JP":    "Japan",
	"VE":    "Venezuela",
	"UZ":    "Uzbekistan",
	"AU":    "Australia",
	"CL":    "Chile",
	"QA":    "Qatar",
	"GE":    "Georgia",
	"EE":    "Estonia",
	"AT":    "Austria",
	"VN":    "Vietnam",
	"NF":    "Norfolk Island",
	"KW":    "Kuwait",
	"AR":    "Argentina",
	"BR":    "Brazil",
	"PL":    "Poland",
	"MC":    "Monaco",
	"MT":    "Malta",
	"BY":    "Belarus",
	"EG":    "Egypt",
	"IN":    "India",
	"SA":    "Saudi Arabia",
	"DE":    "Germany",
	"SK":    "Slovakia",
	"LT":    "Lithuania",
	"DE-CH": "Switzerland (German-speaking part)",
	"CY":    "Cyprus",
	"SZ":    "Eswatini",
	"HU":    "Hungary",
	"UK":    "United Kingdom",
	"FI":    "Finland",
	"RS":    "Serbia",
	"EC":    "Ecuador",
	"AO":    "Angola",
	"IE":    "Ireland",
	"RE":    "Réunion",
	"BD":    "Bangladesh",
	"UA":    "Ukraine",
	"DO":    "Dominican Republic",
	"ES":    "Spain",
	"EN-CA": "Canada (English-speaking part)",
	"NO":    "Norway",
	"BA":    "Bosnia and Herzegovina",
	"TH":    "Thailand",
	"TN":    "Tunisia",
	"TR":    "Turkey",
	"BH":    "Bahrain",
	"MA":    "Morocco",
	"AL":    "Albania",
	"BE":    "Belgium",
	"CO":    "Colombia",
	"DZ":    "Algeria",
	"AZ":    "Azerbaijan",
	"GF":    "French Guiana",
	"GT":    "Guatemala",
	"PE":    "Peru",
	"PY":    "Paraguay",
	"SG":    "Singapore",
	"US":    "United States",
	"KE":    "Kenya",
	"YE":    "Yemen",
	"NZ":    "New Zealand",
	"HN":    "Honduras",
	"AE":    "United Arab Emirates",
	"IQ":    "Iraq",
	"IR":    "Iran",
	"NL":    "Netherlands",
	"CN":    "China",
	"ME":    "Montenegro",
	"GM":    "Gambia",
	"MY":    "Malaysia",
	"HK":    "Hong Kong",
	"NP":    "Nepal",
	"MQ":    "Martinique",
	"SE":    "Sweden",
	"GH":    "Ghana",
	"CH":    "Switzerland",
	"AW":    "Aruba",
	"PK":    "Pakistan",
	"TW":    "Taiwan",
	"CA":    "Canada",
	"PR":    "Puerto Rico",
	"FR":    "France",
	"LV":    "Latvia",
	"MX":    "Mexico",
	"RO":    "Romania",
	"CW":    "Curaçao",
	"TZ":    "Tanzania",
	"MO":    "Macao",
	"PH":    "Philippines",
	"ZW":    "Zimbabwe",
	"IM":    "Isle of Man",
	"ID":    "Indonesia",
	"PA":    "Panama",
	"LB":    "Lebanon",
	"ZA":    "South Africa",
	"PT":    "Portugal",
	"UM":    "United States Minor Outlying Islands",
	"NG":    "Nigeria",
	"IL":    "Israel",
	"BG":    "Bulgaria",
	"JO":    "Jordan",
	"IT":    "Italy",
	"GB":    "United Kingdom",
	"HR":    "Croatia",
	"GR":    "Greece",
}

type Module string

// Enum values
const (
	Discovery         Module = "DISCOVERY"
	ProfilePage       Module = "PROFILE_PAGE"
	Leaderboard       Module = "LEADERBOARD"
	AccountTracking   Module = "ACCOUNT_TRACKING"
	CampaignReport    Module = "CAMPAIGN_REPORT"
	TopicResearch     Module = "TOPIC_RESEARCH"
	ProfileCollection Module = "PROFILE_COLLECTION"
)

type PlanType string

const (
	FreePlan PlanType = "FREE"
	SaasPlan PlanType = "SAAS"
	PaidPlan PlanType = "PAID"
)

type CollectionSource string

const (
	SaasCollection        CollectionSource = "SAAS"
	SaasAtCollection      CollectionSource = "SAAS-AT"
	GccCampaignCollection CollectionSource = "GCC_CAMPAIGN"
	CampaignCollection    CollectionSource = "CAMPAIGN"
)

type ProfileLinkingSource string

const (
	GccLinkedProfile  ProfileLinkingSource = "GCC"
	SaasLinkedProfile ProfileLinkingSource = "SAAS"
)

var TrackedModule = map[string]bool{
	"DISCOVERY":       true,
	"CAMPAIGN_REPORT": true,
	"TOPIC_RESEARCH":  true,
}

var PartnerLimitModules = map[string]bool{
	"DISCOVERY":        true,
	"CAMPAIGN_REPORT":  true,
	"PROFILE_PAGE":     true,
	"ACCOUNT_TRACKING": true,
}

type Activity string

// Enum values
const (
	DiscoveryActivity              Activity = "DISCOVERY"
	ProfilePageActivity            Activity = "PROFILE_PAGE"
	LeaderboardActivity            Activity = "LEADERBOARD"
	AccountTrackingActivity        Activity = "ACCOUNT_TRACKING"
	CampaignReportActivity         Activity = "CAMPAIGN_REPORT"
	TopicResearchActivity          Activity = "TOPIC_RESEARCH"
	ProfileCollectionActivity      Activity = "PROFILE_COLLECTION"
	AddProfileToCollectionActivity Activity = "ADD_PROFILE_TO_COLLECTION"
	AddPostToCollectionActivity    Activity = "ADD_POST_TO_COLLECTION"
	ProfilePageTimeSeriesActivity  Activity = "PROFILE_PAGE_TIMESERIES"
)

type SentimentType string

const (
	PositiveSentiment SentimentType = "positive"
	NegativeSentiment SentimentType = "negative"
	NeutralSentiment  SentimentType = "neutral"
)

type Platform string

// Enum values
const (
	InstagramPlatform Activity = "INSTAGRAM"
	YoutubePlatform   Activity = "YOUTUBE"
	GCCPlatform       Activity = "GCC"
)
