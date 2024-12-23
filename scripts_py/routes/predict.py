from fastapi import APIRouter, HTTPException
from models import SentenceRequest
from pydantic import BaseModel
from transformers import BertTokenizer, BertForSequenceClassification
import torch
import logging

router = APIRouter()

# 感情ラベル
emotion_labels = ["Joy", "Sadness", "Anticipation", "Surprise", "Anger", "Fear", "Disgust", "Trust"]

# モデルとトークナイザーの読み込み
MODEL_PATH = "./ai_model/emotion_model"  # 保存したモデルのパス
tokenizer = BertTokenizer.from_pretrained(MODEL_PATH)
model = BertForSequenceClassification.from_pretrained(MODEL_PATH, use_safetensors=True)
model.eval()  # 推論モードに設定

@router.post("/predict")
async def predict_emotions(request: SentenceRequest):
    try:
        # 入力文をトークナイズ
        inputs = tokenizer(
            request.emotion, truncation=True, padding=True, max_length=128, return_tensors="pt"
        )
        
        # 推論
        with torch.no_grad():
            outputs = model(**inputs)
        
        # 感情スコアを取得
        scores = torch.softmax(outputs.logits, dim=1).numpy()[0]
        print("Scores (after softmax):", scores)
        # スコア最大値を持つインデックスを取得
        max_score_index = int(scores.argmax())

        return {"emotion": request.emotion, "score": max_score_index}
    except Exception as e:
        logging.error(f"Error occurred: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))