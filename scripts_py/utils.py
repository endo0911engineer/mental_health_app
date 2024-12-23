import spacy
import networkx as nx
from wordcloud import WordCloud
import io
import base64

# spaCyの英語モデルをロード
nlp = spacy.load("ja_core_web_sm")

# キーワードの抽出と関連性分析の補助関数
def extract_keywords(text: str):
    doc = nlp(text)
    keywords = [token.text for token in doc if token.pos_ in ['NOUN', 'VERB']]
    return keywords

def create_keyword_graph(keywords):
    G = nx.Graph()
    for i, keyword1 in enumerate(keywords):
        for keyword2 in keywords[i+1:]:
            G.add_edge(keyword1, keyword2)
    return G

def generate_wordcloud(keywords):
    wordcloud = WordCloud(width=800, height=400).generate(" ".join(keywords))
    img = wordcloud.to_image()
    return img

def generate_base64_image(img):
    byte_arr = io.BytesIO()
    img.save(byte_arr, format='PNG')
    byte_arr.seek(0)
    return base64.b64encode(byte_arr.read()).decode('utf-8')