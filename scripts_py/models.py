from pydantic import BaseModel

class SentenceRequest(BaseModel):
    emotion: str

class KeywordRequest(BaseModel):
    text: str