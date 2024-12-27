from fastapi import APIRouter, HTTPException
import matplotlib.pyplot as plt
from matplotlib import rcParams
from models import KeywordRequest
import utils
import logging
from wordcloud import WordCloud
import networkx as nx
import io
from PIL import Image
from matplotlib import font_manager

# ロガーの設定
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

router = APIRouter()

@router.post("/extract_keywords")
async def extract_keywords(request: KeywordRequest):
    try:
        logger.info("Starting keyword extraction...")
        keywords = utils.extract_keywords(request.text)
        logger.info(f"Extracted keywords: {keywords}")

        # 関連性分析
        try:
            logger.info("Creating keyword graph...")
            G = utils.create_keyword_graph(keywords)
        except Exception as e:
            logger.error(f"Error in create_keyword_graph: {e}")
            raise HTTPException(status_code=500, detail=f"Graph creation failed: {e}")

        # グラフ画像生成
        try:
            logger.info("Generating graph image...")
            # 日本語フォントの設定（Windowsの場合）
            font_path = "C:\\Windows\\Fonts\\meiryo.ttc"  # 適切なフォントパスを指定（Windows）
            font_prop = font_manager.FontProperties(fname=font_path)
            plt.rcParams['font.family'] = font_prop.get_name() 
            nx.draw_networkx(G, with_labels=True, node_color='lightblue', edge_color='gray', node_size=500, font_size=10, font_family=font_prop.get_name())

            buffer = io.BytesIO()
            plt.savefig(buffer, format='png')
            buffer.seek(0)

            # PIL Imageを生成
            img = Image.open(buffer)

            graph_image = utils.generate_base64_image(img)
            logger.info("Graph image generated successfully.")
        except Exception as e:
            logger.error(f"Error in generating graph image: {e}")
            raise HTTPException(status_code=500, detail=f"Graph image generation failed: {e}")

        # ワードクラウド画像生成
        try:
            logger.info("Generating wordcloud image...")
            wordcloud_image_pil = utils.generate_wordcloud(keywords)  # すでにPIL.Imageが返される想定

            # デバッグ用にローカルに保存
            wordcloud_image_pil.save("debug_wordcloud.png")  # 保存
            logger.info("Wordcloud image saved to debug_wordcloud.png.")

            # Base64エンコードされた画像を生成
            wordcloud_image = utils.generate_base64_image(wordcloud_image_pil)

            logger.info("Wordcloud image generated successfully.")

        except Exception as e:
            logger.error(f"Error in generating wordcloud image: {e}")
            raise HTTPException(status_code=500, detail=f"Wordcloud generation failed: {e}")
        
        return {
            "keywords": keywords,
            "graph_image": graph_image,
            "wordcloud_image": wordcloud_image,
        }

    except Exception as e:
        logger.error(f"Unhandled exception: {e}")
        raise HTTPException(status_code=500, detail=f"Unhandled error: {e}")
