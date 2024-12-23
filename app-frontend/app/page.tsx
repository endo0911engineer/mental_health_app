'use client';

import { useState } from "react";
import { analyzeSentiment } from "./api";

interface SentimentResult {
  emotion: string;
  confidence: number;
}

export default function Home() {
  const [inputText, setInputText] =  useState<string>("");
  const [result, setResult] = useState<SentimentResult | null>(null);

  const handleAnalyze = async () => {
    if (inputText.trim() === "") {
      alert("Please enter some text!");
      return;
    }

    try {
      const data = await analyzeSentiment(inputText);
      setResult(data);
    } catch (error) {
      console.error("Error:", error);
      alert("A error occurred while analyzing the sentiment.");
    }
  };

  return (
    <div className="container">
      <h1>Sentiment Analysis</h1>
      <textarea
      value={inputText}
      onChange={(e) => setInputText(e.target.value)}
      placeholder="Enter your text here..."
      />
      <button onClick={handleAnalyze}>Analyze</button>

      {result && (
        <div id="resultContainer">
          <h2>Analysis Result</h2>
          <p>
            <strong>Label:</strong> {result.emotion}
          </p>
          <p>
            <strong>Confidence:</strong>{(result.confidence * 100).toFixed(2)}%
          </p>
        </div>
      )}
    </div>
  )
}