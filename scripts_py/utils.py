import spacy
import networkx as nx
from wordcloud import WordCloud
import io
import base64
import matplotlib.pyplot as plt
from matplotlib.figure import Figure
from PIL import Image

# spaCyの日本語モデルをロード
nlp = spacy.load("ja_core_news_sm")

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
    wordcloud = WordCloud(
        font_path='C:/Windows/Fonts/meiryo.ttc',
        width=800, 
        height=400,
        max_words=100,
        background_color='white',
        colormap='coolwarm',
    ).generate(" ".join(keywords))

    img = wordcloud.to_image()
    return img

def generate_base64_image(img):
    print(f"Image type: {type(img)}")
    byte_arr = io.BytesIO()
    
    # 型をチェックして適切な処理を実行
    if isinstance(img, Figure):
        img.savefig(byte_arr, format='PNG')  # Matplotlib Figureの場合
    elif isinstance(img, Image.Image):
        img.save(byte_arr, format='PNG')  # Pillow Imageの場合
    else:
        raise ValueError("Unsupported image type: img must be a Matplotlib Figure or Pillow Image")
    
    byte_arr.seek(0)
    return base64.b64encode(byte_arr.read()).decode('utf-8')