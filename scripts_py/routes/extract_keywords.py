from fastapi import APIRouter, HTTPException
import matplotlib.pyplot as plt
from models import KeywordRequest
import utils

router = APIRouter()

@router.post("/extract_keywords")
async def extract_keywords(request: KeywordRequest):
    try:
        keywords = utils.extract_keywords(request.text)

        # 関連性分析
        G = utils.create_keyword_graph(keywords)

        # グラフ画像生成
        graph_image = utils.generate_base64_image(plt.figure(figsize=(8, 6)))
        wordcloud_image = utils.generate_base64_image(utils.generate_wordcloud(keywords))

        return {
            "keywords": keywords,
            "graph_image": graph_image,
            "wordcloud_image": wordcloud_image
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))