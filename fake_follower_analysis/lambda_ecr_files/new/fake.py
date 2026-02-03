import json
import pandas as pd
import re
from unidecode import unidecode
from itertools import permutations
import time
from rapidfuzz import fuzz as fuzzz
import numpy as np
import boto3
import csv
from indictrans import Transliterator

session = boto3.Session(
    aws_access_key_id='AKIAVIODXZPRP45NADH2',
    aws_secret_access_key='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',
)

def handler(event, context):
    # arr = np.random.randint(0, 10, (3, 3))
    response = model(event)
    print(response)
    sqs = session.resource('sqs', region_name='eu-north-1') 
    # sqs = boto3.client('sqs', region_name='eu-north-1')
    queue = sqs.get_queue_by_name(QueueName='output_queue')
    # queue_url = 'https://sqs.eu-north-1.amazonaws.com/361724627938/output_queue'

    # The message you want to send
    message = 'Your message here'
    response = queue.send_message(MessageBody=json.dumps(response))
    print(response)
    return {
        "statusCode": 200,
        "body": {
            "message": response
        }
    }

# df = pd.read_csv('../creators/Result_10_2.csv')
# df = df[:1000]
# follower_data = df.to_json(indent=4, orient='records')
# follower_data = json.loads(follower_data)


data={
    "hin":[
        "़","०","१","२","३","४","५","६","७","८","९","ॐ","ं","ँ","ः","अ","आ","इ","ई","उ","ऊ","ऋ","ऌ","ऍ","ए","ऐ","ऑ","ओ","औ","क","ख","ग","घ","ङ","च","छ","ज","झ","ञ","ट","ठ","ड","ढ","ण","त","थ","द","ध","न","प","फ","ब","भ","म","य","र","ल","ळ","व","श","ष","स","ह","ऽ","ा","ि","ी","ु","ू","ृ","ॄ","ॅ","े","ै","ॉ","ो","ौ","्"
    ]
    ,"pan":[
        "੦","੧","੨","੩","੪","੫","੬","੭","੮","੯","ੴ","ੳ","ਉ","ਊ","ਓ","ਅ","ਆ","ਐ","ਔ","ੲ","ਇ","ਈ","ਏ","ਸ","ਸ਼","ਹ","ਕ","ਖ","ਖ਼","ਗ","ਗ਼","ਘ","ਙ","ਚ","ਛ","ਜ","ਜ਼","ਝ","ਞ","ਟ","ਠ","ਡ","ਢ","ਣ","ਤ","ਥ","ਦ","ਧ","ਨ","ਪ","ਫ","ਫ਼","ਬ","ਭ","ਮ","ਯ","ਰ","ਲ","ਲ਼","ਵ","ੜ"
    ]

    ,"guj":[
        "઼","ૐ","ં","ઁ","ઃ","અ","અં","અઃ","આ","ઇ","ઈ","ઉ","ઊ","ઋ","ૠ","ઍ","એ","ઐ","ઑ","ઓ","ઔ","ક","ક્ષ","ખ","ગ","ઘ","ઙ","ચ","છ","જ","જ્ઞ","ઝ","ઞ","ટ","ઠ","ડ","ઢ","ણ","ત","ત્ર","થ","દ","ધ","ન","પ","ફ","બ","ભ","મ","ય","ર","લ","વ","શ","ષ","સ","હ","ળ","ઽ","ા","િ","ી","ુ","ૂ","ૃ","ૄ","ૅ","ે","ૈ","ૉ","ો","ૌ","્"
    ]

    ,"ben":[
        "়","৺","অ","আ","ই","ঈ","উ","ঊ","ঋ","ৠ","ঌ","ৡ","এ","ঐ","ও","ঔ","ং","ঃ","ঁ","ক","ক্ষ","খ","গ","ঘ","ঙ","চ","ছ","জ","ঝ","ঞ","ট","ঠ","ড","ড়","ঢ","ঢ়","ণ","ত","ৎ","থ","দ","ধ","ন","প","ফ","ব","ভ","ম","য","য়","র","ল","শ","ষ","স","হ","ঽ","া","ি","ী","ু","ূ","ৃ","ৄ","ৢ","ৣ","ে","ৈ","ো","ৌ","্","ৗ"
    ]

    ,'urd':[
        "ا", "آ", "ب", "پ", "ت", "ٹ", "ث", "ج", "چ", "ح", "خ", "د", "ڈ", "ذ", "ر", "ڑ", "ز", "ژ", "س", "ش", "ص", "ض", "ط", "ظ", "ع", "غ", "ف", "ق", "ک", "گ", "ل", "م", "ن", "ں", "ھ", "و", "ؤ", "ہ", "ھ", "ء", "ی", "ئ"
    ]

    ,'tam':[
        "அ","ஆ","இ","ஈ","உ","ஊ","எ","ஏ","ஐ","ஒ","ஓ","ஔ","ஂ","ஃ","க்","க","ங்","ங","ச்","ச","ஞ்","ஞ","ட்","ட","ண்","ண","த்","த","ந்","ந","ப்","ப","ம்","ம","ய்","ய","ர்","ர","ல்","ல","வ்","வ","ழ்","ழ","ள்","ள","ற்","ற","ன்","ன","ஜ்","ஜ","ஶ்","ஷ்","ஷ","ஸ்","ஸ","ஹ்","ஹ","க்ஷ்","க்ஷ","ா","ி","ீ","ு","ூ","ெ","ே","ை","ொ","ோ","ௌ","்"
    ]
    ,'mal':[
        "അ", "ആ", "ഇ", "ഈ", "ഉ", "ഊ", "ഋ", "ൠ", "എ", "ഏ",
        "ഐ", "ഒ", "ഓ", "ഔ", "അം", "അഃ", "ക", "ഖ", "ഗ", "ഘ",
        "ങ", "ച", "ഛ", "ജ", "ഝ", "ഞ", "ട", "ഠ", "ഡ", "ഢ",
        "ണ", "ത", "ഥ", "ദ", "ധ", "ന", "പ", "ഫ", "ബ", "ഭ",
        "മ", "യ", "ര", "ല", "വ", "ശ", "ഷ", "സ", "ഹ", "ള",
        "ഴ", "റ"
    ],
    "kan":[
        "಼","೦","೧","೨","೩","೪","೫","೬","೭","೮","೯","ಅ","ಆ","ಇ","ಈ","ಉ","ಊ","ಋ","ೠ","ಌ","ೡ","ಎ","ಏ","ಐ","ಒ","ಓ","ಔ","ಂ","ಃ","ೱ","ೲ","ಕ","ಖ","ಗ","ಘ","ಙ","ಚ","ಛ","ಜ","ಝ","ಞ","ಟ","ಠ","ಡ","ಢ","ಣ","ತ","ಥ","ದ","ಧ","ನ","ಪ","ಫ","ಬ","ಭ","ಮ","ಯ","ರ","ಱ","ಲ","ವ","ಶ","ಷ","ಸ","ಹ","ಳ","ೞ","ಽ","ಾ","ಿ","ೀ","ು","ೂ","ೃ","ೄ","ೆ","ೇ","ೈ","ೊ","ೋ","ೌ","್","ೕ","ೖ"
    ],
    "ori": [
        "଼","ଅ","ଆ","ଇ","ଈ","ଉ","ଊ","ଋ","ଏ","ଐ","ଓ","ଔ","ଁ","ଂ","ଃ","କ","ଖ","ଗ","ଘ","ଙ","ଚ","ଛ","ଜ","ଝ","ଞ","ଟ","ଠ","ଡ","ଡ଼","ଢ","ଢ଼","ଣ","ତ","ଥ","ଦ","ଧ","ନ","ପ","ଫ","ବ","ଭ","ମ","ଯ","ୟ","ର","ଲ","ଳ","ଵ","ୱ","ଶ","ଷ","ସ","ହ","କ୍ଷ","ା","ି","ୀ","ୁ","ୂ","ୃ","େ","ୈ","ୋ","ୌ","୍"
    ]
    ,
    "tel": [
        "అ","ఆ","ఇ","ఈ","ఉ","ఊ","ఋ","ౠ","ఌ","ౡ","ఎ","ఏ","ఐ","ఒ","ఓ","ఔ","ఁ","ం","ః","క","ఖ","గ","ఘ","ఙ","చ","ఛ","జ","ఝ","ఞ","ట","ఠ","డ","ఢ","ణ","త","థ","ద","ధ","న","ప","ఫ","బ","భ","మ","య","ర","ఱ","ల","వ","శ","ష","స","హ","ళ","ా","ి","ీ","ు","ూ","ృ","ౄ","ె","ే","ై","ొ","ో","ౌ","్","ౕ","ౖ"
    ]
}
char_to_lang = {}

for lang, chars in data.items():
    for char in chars:
        char_to_lang[char] = lang

# # print(char_to_lang)
name_data = pd.read_csv('baby_names.csv')
namess = name_data['Baby Names'].str.lower()

def symbol_name_convert(name):
    original = [
        "🅐🅑🅒🅓🅔🅕🅖🅗🅘🅙🅚🅛🅜🅝🅞🅟🅠🅡🅢🅣🅤🅥🅦🅧🅨🅩🅐🅑🅒🅓🅔🅕🅖🅗🅘🅙🅚🅛🅜🅝🅞🅟🅠🅡🅢🅣🅤🅥🅦🅧🅨🅩⓿➊➋➌➍➎➏➐➑➒",
        '🅰🅱🅲🅳🅴🅵🅶🅷🅸🅹🅺🅻🅼🅽🅾🅿🆀🆁🆂🆃🆄🆅🆆🆇🆈🆉',
        "🄰🄱🄲🄳🄴🄵🄶🄷🄸🄹🄺🄻🄼🄽🄾🄿🅀🅁🅂🅃🅄🅅🅆🅇🅈🅉",
        "ⒶⒷⒸⒹⒺⒻⒼⒽⒾⒿⓀⓁⓂⓃⓄⓅⓆⓇⓈⓉⓊⓋⓌⓍⓎⓏⓐⓑⓒⓓⓔⓕⓖⓗⓘⓙⓚⓛⓜⓝⓞⓟⓠⓡⓢⓣⓤⓥⓦⓧⓨⓩ⓪①②③④⑤⑥⑦⑧⑨",
        "𝐀𝐁𝐂𝐃𝐄𝐅𝐆𝐇𝐈𝐉𝐊𝐋𝐌𝐍𝐎𝐏𝐐𝐑𝐒𝐓𝐔𝐕𝐖𝐗𝐘𝐙𝐚𝐛𝐜𝐝𝐞𝐟𝐠𝐡𝐢𝐣𝐤𝐥𝐦𝐧𝐨𝐩𝐪𝐫𝐬𝐭𝐮𝐯𝐰𝐱𝐲𝐳𝟎𝟏𝟐𝟑𝟒𝟓𝟔𝟕𝟖𝟗",
        "𝗔𝗕𝗖𝗗𝗘𝗙𝗚𝗛𝗜𝗝𝗞𝗟𝗠𝗡𝗢𝗣𝗤𝗥𝗦𝗧𝗨𝗩𝗪𝗫𝗬𝗭𝗮𝗯𝗰𝗱𝗲𝗳𝗴𝗵𝗶𝗷𝗸𝗹𝗺𝗻𝗼𝗽𝗾𝗿𝘀𝘁𝘂𝘃𝘄𝘅𝘆𝘇𝟬𝟭𝟮𝟯𝟰𝟱𝟲𝟳𝟴𝟵",
        "𝘈𝘉𝘊𝘋𝘌𝘍𝘎𝘏𝘐𝘑𝘒𝘓𝘔𝘕𝘖𝘗𝘘𝘙𝘚𝘛𝘜𝘝𝘞𝘟𝘠𝘡𝘢𝘣𝘤𝘥𝘦𝘧𝘨𝘩𝘪𝘫𝘬𝘭𝘮𝘯𝘰𝘱𝘲𝘳𝘴𝘵𝘶𝘷𝘸𝘹𝘺𝘻0123456789",
        "𝘼𝘽𝘾𝘿𝙀𝙁𝙂𝙃𝙄𝙅𝙆𝙇𝙈𝙉𝙊𝙋𝙌𝙍𝙎𝙏𝙐𝙑𝙒𝙓𝙔𝙕𝙖𝙗𝙘𝙙𝙚𝙛𝙜𝙝𝙞𝙟𝙠𝙡𝙢𝙣𝙤𝙥𝙦𝙧𝙨𝙩𝙪𝙫𝙬𝙭𝙮𝙯0123456789",
        "𝙰𝙱𝙲𝙳𝙴𝙵𝙶𝙷𝙸𝙹𝙺𝙻𝙼𝙽𝙾𝙿𝚀𝚁𝚂𝚃𝚄𝚅𝚆𝚇𝚈𝚉𝚊𝚋𝚌𝚍𝚎𝚏𝚐𝚑𝚒𝚓𝚔𝚕𝚖𝚗𝚘𝚙𝚚𝚛𝚜𝚝𝚞𝚟𝚠𝚡𝚢𝚣𝟶𝟷𝟸𝟹𝟺𝟻𝟼𝟽𝟾𝟿",
        "𝔸𝔹ℂ𝔻𝔼𝔽𝔾ℍ𝕀𝕁𝕂𝕃𝕄ℕ𝕆ℙℚℝ𝕊𝕋𝕌𝕍𝕎𝕏𝕐ℤ𝕒𝕓𝕔𝕕𝕖𝕗𝕘𝕙𝕚𝕛𝕜𝕝𝕞𝕟𝕠𝕡𝕢𝕣𝕤𝕥𝕦𝕧𝕨𝕩𝕪𝕫𝟘𝟙𝟚𝟛𝟜𝟝𝟞𝟟𝟠𝟡",
        "𝕬𝕭𝕮𝕯𝕰𝕱𝕲𝕳𝕴𝕵𝕶𝕷𝕸𝕹𝕺𝕻𝕼𝕽𝕾𝕿𝖀𝖁𝖂𝖃𝖄𝖅𝖆𝖇𝖈𝖉𝖊𝖋𝖌𝖍𝖎𝖏𝖐𝖑𝖒𝖓𝖔𝖕𝖖𝖗𝖘𝖙𝖚𝖛𝖜𝖝𝖞𝖟0123456789",
        "𝓐𝓑𝓒𝓓𝓔𝓕𝓖𝓗𝓘𝓙𝓚𝓛𝓜𝓝𝓞𝓟𝓠𝓡𝓢𝓣𝓤𝓥𝓦𝓧𝓨𝓩𝓪𝓫𝓬𝓭𝓮𝓯𝓰𝓱𝓲𝓳𝓴𝓵𝓶𝓷𝓸𝓹𝓺𝓻𝓼𝓽𝓾𝓿𝔀𝔁𝔂𝔃0123456789",
        "ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ０１２３４５６７８９",
        "𝓐𝓑𝓒𝓓𝓔𝓕𝓖𝓗𝓘𝓙𝓚𝓛𝓜𝓝𝓞𝓟𝓠𝓡𝓢𝓣𝓤𝓥𝓦𝓧𝓨𝓩𝓪𝓫𝓬𝓭𝓮𝓯𝓰𝓱𝓲𝓳𝓴𝓵𝓶𝓷𝓸𝓹𝓺𝓻𝓼𝓽𝓾𝓿𝔀𝔁𝔂𝔃0123456789",
        "𝘈𝘉𝘊𝘋𝘌𝘍𝘎𝘏𝘐𝘑𝘒𝘓𝘔𝘕𝘖𝘗𝘘𝘙𝘚𝘛𝘜𝘝𝘞𝘟𝘠𝘡𝘢𝘣𝘤𝘥𝘦𝘧𝘨𝘩𝘪𝘫𝘬𝘭𝘮𝘯𝘰𝘱𝘲𝘳𝘴𝘵𝘶𝘷𝘸𝘹𝘺𝘻𝟢𝟣𝟤𝟥𝟦𝟧𝟨𝟩𝟪𝟫",
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    ]
    replaceAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    originalMap = {}
    for alphabet in original:
        originalMap.update(dict(zip(alphabet, replaceAlphabet)))

    result = "".join(originalMap.get(char, char) for char in name)
    return result
        
def check_lang_other_than_indic(symbolic_name):
    if not symbolic_name:
        print(symbolic_name)
    return 1 if re.search(r'[Α-Ωα-ωԱ-Ֆა-ჰ一-鿿가-힣]+', symbolic_name, re.UNICODE) else 0
        

def load_dict(filename):
    with open(filename, 'r') as f:
        return dict(csv.reader(f))
try:
    vowels = load_dict('svar.csv')
    consonants = load_dict('vyanjan.csv')
except:
    print("failed")
    


def process_word(word):
    str1 = ""
    i = 0
    while i < len(word):
        if (i+1<len(word) and word[i+1].strip()==' ़'.strip()):
            c = word[i]+word[i+1]
            i += 2
        else:
            c = word[i]
            i += 1
        if (c in vowels):
            str1 += vowels[c]
        elif (c in consonants):
            if(i<len(word) and word[i] in consonants):
                if ((c=='झ' and i!=0) or (i!=0 and i+1<len(word) and word[i+1] in vowels)): 
                    str1 += consonants[c]
                else:
                    str1 += consonants[c]+'a'
            else:
                str1 += consonants[c]
        elif c in ['\n','\t',' ','!',',','।','-',':','\\','_','?'] or c.isalnum():
            str1 += c.replace('।','.')
    return str1

def detect_language(word):
    for char in word:
        lang = char_to_lang.get(char)
        if lang is not None:
            if lang == 'hin':
                return process_word(word)
            else:
                trn = Transliterator(source=lang, target='eng', build_lookup=True)
                return trn.transform(word)
        # if lang is not None:
        #     out = xlit_engine.translit_sentence(word, lang_code=lang)
        #     return out
    return word




def uni_decode(row):
    return unidecode(row, errors='preserve')    





def process(follower_handle, cleaned_handle, cleaned_name):
    SPECIAL_CHARS = ('_', '-', '.')
    
    if any(char in follower_handle for char in SPECIAL_CHARS):
        if not ' ' in cleaned_name:
            if generate_similarity_score(cleaned_handle, cleaned_name)>80:
                return 0
            else:
                return 1
        else:
            return 0
    else:
        return 2
    

def clean_handle(handle):
    if not handle or isinstance(handle, float):
        return ''
    cleaned_handle = re.sub(r'[_\-.]', ' ', handle)
    cleaned_handle = re.sub(r'[^\w\s]', '', cleaned_handle)
    cleaned_handle = re.sub(r'\d', '', cleaned_handle).lower()
    cleaned_handle = re.sub(r'[^a-zA-Z\s]', '', cleaned_handle).strip()
    return cleaned_handle

# Function to count numerical digits in a string
def count_numerical_digits(text):
    if not isinstance(text, str):
        text = str(text)
    return sum(c.isdigit() for c in text)

# function to detect if numerical digit count is more than 4
def fake_real_more_than_4_digit(number):
    return 1 if number>4 else 0



# Function to generate a similarity score between cleaned handle and cleaned name
def generate_similarity_score(handle, name):
    start = time.time()
    name = name.split()
    if len(name)<=4:
        name_permutations = [' '.join(p) for p in permutations(name)]
    else :
        name_permutations = name
    similarity_score = -1
    cleaned_handle = handle.replace(' ', '')
    for name in name_permutations:
        cleaned_name = name.replace(' ', '')  # Remove spaces from the name
        partial_ratio = fuzzz.partial_ratio(cleaned_handle.lower(), cleaned_name.lower())
        token_sort_ratio = fuzzz.token_sort_ratio(cleaned_handle.lower(), cleaned_name.lower())
        token_set_ratio = fuzzz.token_set_ratio(cleaned_handle.lower(), cleaned_name.lower())

        # Calculate a weighted average of the scores
        similarity_score = max(similarity_score, (2 * partial_ratio + token_sort_ratio + token_set_ratio) / 4)
    end = time.time()
    print("similirity_score ---------:",(end-start) * 10**3, "ms")
    return similarity_score


def based_on_partial_ratio(similarity_score):
    if similarity_score>90:
        return 0
    return 1

def process1(fake_real_based_on_lang, number_more_than_4_handle, chhitij_logic):
    if fake_real_based_on_lang:
        return 1
    if number_more_than_4_handle:
        return 1
    if chhitij_logic==1:
        return 1
    elif chhitij_logic==2:
        return 2
    return 0

def final(fake_real_based_on_lang, similarity_score, number_more_than_4_handle, chhitij_logic):
    if fake_real_based_on_lang:
        return 1

    if 0 < similarity_score <= 40:
        return 0.33

    if number_more_than_4_handle:
        return 1
    if chhitij_logic==1:
        return 1
    elif chhitij_logic==2:
        return 0

    return 0

def only_numeric(name):
    if re.search(r'[a-zA-Z]', name):
        return 0
    else:
        return 1


def score(i, first_name):
    i = i.lower()
    first_name = first_name.lower()
    ratio = fuzzz.ratio(i, first_name)
    token_sort_ratio = fuzzz.token_sort_ratio(i, first_name)
    token_set_ratio = fuzzz.token_set_ratio(i, first_name)
    return (2 * ratio + token_sort_ratio + token_set_ratio) / 4

def check_indian_names(name):
    if len(name)<2:
        return 1
    else:
        similarity_score = 0
        name = name.split()
        first_name = name[0]
        last_name = name[1] if len(name) >= 2 else None

        for i in namess:
            similarity_score = max(similarity_score, score(i, first_name))
            
        if last_name:
            if len(last_name)<2:
                similarity_score = 1
            else:
                for i in namess:
                    similarity_score = max(similarity_score, score(i, last_name))
        return similarity_score



def model(event):
    follower_data = [event]
    response = {}
    total_time = 0
    for index, temp in enumerate(follower_data, start=1):
        print(follower_data)
        print(index)
        start = time.time()
        symbolic_name = symbol_name_convert(temp['follower_full_name'])
        fake_real_based_on_lang = check_lang_other_than_indic(symbolic_name)
        transliterated_follower_name = detect_language(symbolic_name)
        decoded_name = uni_decode(transliterated_follower_name)
        end = time.time()
        print("The time of execution of decoded_name ---------:",(end-start) * 10**3, "ms")
        
        start1 = time.time()
        cleaned_handle = clean_handle(temp['follower_handle'])
        cleaned_name = clean_handle(decoded_name)
        chhitij_logic = process(temp['follower_handle'], cleaned_handle, cleaned_name)
        number_handle = count_numerical_digits(temp['follower_handle'])
        number_name = count_numerical_digits(temp['follower_full_name'])
        number_more_than_4_handle = fake_real_more_than_4_digit(number_handle)
        number_more_than_4_name = fake_real_more_than_4_digit(number_name)
        similarity_score = generate_similarity_score(cleaned_handle, cleaned_name)
        fake_real_based_on_fuzzy_score_90 = based_on_partial_ratio(similarity_score)
        process1_ = process1(fake_real_based_on_lang, number_more_than_4_handle, chhitij_logic)
        final_ = final(fake_real_based_on_lang, similarity_score, number_more_than_4_handle, chhitij_logic)
        numeric_handle = only_numeric(temp['follower_handle'])

        # grouped = temp_df.groupby(['handle', 'cleaned_name']).size().reset_index(name='count')
        # temp_df = temp_df.merge(grouped, on=['handle', 'cleaned_name'], how='left', suffixes=('', '_grouped'))
        # temp_df['count'].fillna(0, inplace=True)
        # temp_df['duplicate_more_than_3'] = (temp_df['count'] > 3).astype(int)
        # temp_df['duplicate_more_than_5'] = (temp_df['count'] > 5).astype(int)
        # temp_df['duplicate_more_than_10'] = (temp_df['count'] > 10).astype(int)
        indian_name_score = check_indian_names(cleaned_name)
        # temp_df['indian_name_score'] = temp_df['cleaned_name'].apply(check_indian_names)
        # filtered_score = df[df['indian_name_score']>80]
        score_80 = 1 if indian_name_score>80 else 0
        end1 = time.time()
        # temp_df['score>80'] = (temp_df['indian_name_score'] > 80).astype(int)
        end = time.time()
        total_time += (end-start) * 10**3
        
        print("The time of execution of final ---------:",(end-start) * 10**3, "ms")
        print(decoded_name)
        print(score_80)

        response = {
            "symbolic_name": symbolic_name,
            "fake_real_based_on_lang": fake_real_based_on_lang,
            "transliterated_follower_name": transliterated_follower_name,
            "decoded_name": decoded_name,
            "cleaned_handle": cleaned_handle,
            "cleaned_name": cleaned_name,
            "chhitij_logic": chhitij_logic,
            "number_handle": number_handle,
            "number_more_than_4_handle": number_more_than_4_handle,
            "similarity_score": similarity_score,
            "fake_real_based_on_fuzzy_score_90": fake_real_based_on_fuzzy_score_90,
            "process1_": process1_,
            "final_": final_,
            "numeric_handle": numeric_handle,
            "indian_name_score": indian_name_score,
            "score_80": score_80
        }
        

    print(total_time/len(follower_data))
    return response