from fastapi import FastAPI
from routes.predict import router as predict_router
from routes.extract_keywords import router as extract_keywords_router

app = FastAPI()

# ルータをインクルード
app.include_router(predict_router)
app.include_router(extract_keywords_router)
