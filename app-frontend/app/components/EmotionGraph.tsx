'use client';

import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface EmotionGraphProps {
    data: { date: string; score: number }[];
}

const EmotionGraph = ({ data }: EmotionGraphProps) => {
    return (
        <div style={{ width: '100%', height: 300 }}>
            <h2>Emotion Trends</h2>
            <ResponsiveContainer>
                <LineChart data={data}>
                    <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
                    <XAxis dataKey="date" />
                    <YAxis domain={[0, 7]} />
                    <Tooltip />
                    <Line type="monotone" dataKey="score" stroke="#8884d8" activeDot={{ r: 8 }} />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};

export default EmotionGraph