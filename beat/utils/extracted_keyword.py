import re
import yake


def remove_numeric_and_emojis(text):
    # Remove numeric characters
    text = re.sub(r'\d+', '', text)
    
    # Remove emojis
    emoji_pattern = re.compile("["
                               u"\U0001F600-\U0001F64F"  # emoticons
                               u"\U0001F300-\U0001F5FF"  # symbols & pictographs
                               u"\U0001F680-\U0001F6FF"  # transport & map symbols
                               u"\U0001F1E0-\U0001F1FF"  # flags (iOS)
                               "]+", flags=re.UNICODE)
    text = emoji_pattern.sub(r'', text)
    
    # Remove non-alpha characters
    text = re.sub(r'[^a-zA-Z\s]', '', text)
    
    return text.strip()

    
def get_extracted_keywords(s: str):
    # s = remove_numeric_and_emojis(s)
    kw_extractor = yake.KeywordExtractor(n=1, top=5)
    extracted_keywords = kw_extractor.extract_keywords(s)
    extracted_keywords = [item[0] for item in extracted_keywords]    
    return str(extracted_keywords)