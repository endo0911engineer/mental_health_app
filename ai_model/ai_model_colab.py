import pandas as pd
import numpy as np
import torch
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
from transformers import BertTokenizer, BertForSequenceClassification
from transformers import Trainer, TrainingArguments
import os
os.environ["WANDB_DISABLED"] = "true"
# データの読み込み（CSVファイルなどをアップロード）
from google.colab import files
uploaded = files.upload()

# データセットの読み込み
data = pd.read_csv(next(iter(uploaded.keys())), sep="\t")

# データ確認
print(data.head())

data = data[[
    "Sentence",
    "Avg. Readers_Joy", "Avg. Readers_Sadness", "Avg. Readers_Anticipation",
    "Avg. Readers_Surprise", "Avg. Readers_Anger", "Avg. Readers_Fear",
    "Avg. Readers_Disgust", "Avg. Readers_Trust"
]]


# 入力データとラベルを分ける
X = data["Sentence"]
y = data.iloc[:, 1:].values  # 感情のスコアを配列として取得

# データをトレーニングと検証用に分割
train_texts, val_texts, train_labels, val_labels = train_test_split(
    X, y, test_size=0.2, random_state=42
)

# データセットの作成
class EmotionDataset(torch.utils.data.Dataset):
    def __init__(self, encodings, labels):
        self.encodings = encodings
        self.labels = labels

    def __len__(self):
        return len(self.labels)

    def __getitem__(self, idx):
        item = {key: torch.tensor(val[idx]) for key, val in self.encodings.items()}
        item["labels"] = torch.tensor(self.labels[idx], dtype=torch.float)  # ラベルをfloatに変換
        return item

train_dataset = EmotionDataset(train_encodings, train_labels)
val_dataset = EmotionDataset(val_encodings, val_labels)

# モデルの準備
model = BertForSequenceClassification.from_pretrained(
    "cl-tohoku/bert-base-japanese",
    num_labels=8  # 8つの感情スコアを予測
)

# トレーニング設定
training_args = TrainingArguments(
    output_dir="./results",          # 出力ディレクトリ
    evaluation_strategy="epoch",    # 各エポックごとに評価
    learning_rate=2e-5,             # 学習率
    per_device_train_batch_size=16, # バッチサイズ
    per_device_eval_batch_size=16,  # 評価バッチサイズ
    num_train_epochs=3,             # エポック数
    weight_decay=0.01,              # 重み減衰
    logging_dir="./logs",           # ログ保存先
    save_total_limit=2,             # 保存するモデルの数
)

# Trainerのセットアップ
trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=train_dataset,
    eval_dataset=val_dataset,
)

# モデルのトレーニング
trainer.train()

import torch
import numpy as np

# 感情ラベル
emotion_labels = ["Joy", "Sadness", "Anticipation", "Surprise", "Anger", "Fear", "Disgust", "Trust"]

# テストデータ
sample_sentence = ["今日は仕事がとても忙しかった。ちょっとイライラするときもあったが最後にはしっかり終わらせることができた。"]

# 入力データをトークナイズしてデバイスに移動
sample_encodings = tokenizer(
    sample_sentence, truncation=True, padding=True, max_length=128, return_tensors="pt"
).to(device)

# 推論を実行
with torch.no_grad():
    outputs = model(**sample_encodings)

# 各感情スコアの予測結果
predicted_scores = torch.softmax(outputs.logits, dim=1).cpu().numpy()[0]  # 1文のスコアのみ取得

# スコアを見やすく整形
predicted_scores_rounded = np.round(predicted_scores, 4)

# 感情名とスコアを対応付けて表示
result = dict(zip(emotion_labels, predicted_scores_rounded))
print("Predicted emotion scores:", result)

model.save_pretrained("./emotion_model")
tokenizer.save_pretrained("./emotion_model")

from google.colab import drive
import shutil
drive.mount('/content/drive')

# モデルをDriveに保存
shutil.copytree('./emotion_model', '/content/drive/MyDrive/emotion_model')