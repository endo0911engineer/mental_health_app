'use client'

import { useRouter } from 'next/navigation';
import { getEmotions, saveEmotion, deleteEmotion } from '../../api'; 
import { useEffect, useState } from 'react';
import CalendarComponent from '../../components/CalendarComponent';
import EmotionGraph from '../../components/EmotionGraph';
import styles from '../../DashboardPage.module.css';

const DashboardPage = () => {
    const [selectedDate, setSelectedDate] = useState<string | null>(null);
    const [highlightedDates, setHighlightedDates] = useState<{  [key: string]: string}>({});
    const [emotions, setEmotions] = useState<{ date: string, score: number, text: string }[]>([]);
    const [currentEmotion, setCurrentEmotion] = useState<{ text: string; score: number } | null>(null);
    const [isEditing, setIsEditing] = useState<boolean>(false);
    const [inputText, setInputText] = useState<string>('');
    const router = useRouter();

    const handleLogout = () => {
        localStorage.removeItem('token');
        router.push('/login');
    }

    // Fetch emotions when the page loads
    useEffect(() => {
        const fetchEmotions = async () => {
            const userId = parseInt(localStorage.getItem('userId') || '0', 10);
            if (!userId) return;

            try {
                const emotions = await getEmotions(userId);
                console.log('Fetched emotions:', emotions); // デバッグ用ログ

                // 日付ごとにデータを整理
                const organizedEmotions = emotions.map((emotion: { date: string; score: number; emotion: string }) => {
                    const date = new Date(emotion.date); // UTCの日付を取得
                    const jstDate = new Date(date.getTime() + 9 * 60 * 60 * 1000); // JSTに変換
                    return {
                        date: jstDate.toISOString().split('T')[0],
                        score: emotion.score,
                        text: emotion.emotion,
                    }
                });

                setEmotions(organizedEmotions);

                // カレンダーのハイライトデータを更新
                const highlights: { [key: string]: string } = {};
                organizedEmotions.forEach((emotion: { date: string; score: number; text: string }) => {
                const { date, score } = emotion;
                if (score >= 5) highlights[date] = 'green';
                else if (score >= 3) highlights[date] = 'yellow';
                else highlights[date] = 'red';
            });

                setHighlightedDates(highlights);
            } catch (error) {
                console.error('Error fetching emotions:', error)
            }
        };

        fetchEmotions();
    }, []);

    const handleDateClick = (date: string) => {
        console.log('Received date from CalendarComponent:', date); // デバッグログ
        setSelectedDate(date);

        const emotion = emotions.find((e) => e.date === date) || null;
        console.log('Selected emotion:', emotion); // デバッグ用ログ
        setCurrentEmotion(emotion ? { text: emotion.text, score: emotion.score } : null);
        setInputText(emotion?.text || '');
        setIsEditing(!!emotion);
    };

    const handleSaveEmotion = async () => {
        // ユーザーIDを取得
        const userId = parseInt(localStorage.getItem('userId') || '', 10);
        if (!selectedDate || !userId) return;

        try {
            const newEmotion = await saveEmotion(selectedDate, inputText, userId);
            setEmotions((prev) => 
                prev.filter((e) => e.date !== selectedDate).concat(newEmotion)
            );
            setHighlightedDates((prev) => ({
                ...prev,
                [selectedDate]: newEmotion.score >= 5 ? 'green' : newEmotion.score >= 3 ? 'yellow' : 'red',
            }));
            setCurrentEmotion({ text: inputText, score: newEmotion.score });
        } catch (error) {
            console.error('Error saving emotion:', error);
        }
    };

    const handleDeleteEmotion = async () => {
        const userId = parseInt(localStorage.getItem('userId') || '', 10);
        if (! selectedDate || !userId) return;

        try {
            await deleteEmotion(selectedDate, userId);
            setEmotions((prev) => prev.filter((e) => e.date !== selectedDate));
            setHighlightedDates((prev) => {
                const { [selectedDate]: _, ...rest } = prev;
                return rest;
            });
            setCurrentEmotion(null);
            setInputText('');
            setIsEditing(false);
        } catch (error) {
            console.error('Error deleting emotion:', error);
        }
    };

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>Dashboard</h1>
                <button className={styles.logoutButton} onClick={handleLogout}>
                    Logout
                </button>
            </div>

            <div className={styles.calendarGraphWrapper}>
                <div className={styles.calendarWrapper}>
                    <CalendarComponent onDateClick={handleDateClick} highlightedDates={highlightedDates} />
                </div>
                <div className={styles.graphWrapper}>
                    <EmotionGraph data={emotions}/>
                </div>
            </div>

            {selectedDate && (
                <div className={styles.emotionPanel}>
                    <h3>{selectedDate}</h3>
                    {isEditing ? (
                        <>
                        <p>Current: Emotion: {currentEmotion?.text}</p>
                        <textarea
                        value={inputText}
                        onChange={(e) => setInputText(e.target.value)}
                        placeholder="Edit your emotion..."
                        />
                        <button onClick={handleSaveEmotion}>Save</button>
                        <button onClick={handleDeleteEmotion}>Delete</button>
                        </>
                    ) : (
                        <>
                        <textarea
                        value={inputText}
                        onChange={(e) => setInputText(e.target.value)}
                        placeholder="Write about your emotion..."
                        />
                        <button onClick={handleSaveEmotion}>Add Emotion</button>
                        </>
                    )}
                </div>
            )}
        </div>
    );
};

export default DashboardPage;