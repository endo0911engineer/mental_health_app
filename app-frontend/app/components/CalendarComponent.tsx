'use client';

import Calendar from 'react-calendar';
import 'react-calendar/dist/Calendar.css'

interface CalendarProps {
    onDateClick: (date: string) => void;
    highlightedDates: { [key: string]: string };
}

const CalendarComponent = ({ onDateClick, highlightedDates }: CalendarProps) => {
    const tileContent = ({ date }: { date: Date }) => {
        const jstDate = new Date(date.getTime() + 9 * 60 * 60 * 1000);
        const key = jstDate.toISOString().split('T')[0];
        const color = highlightedDates[key];
        return color ? <div style={{ backgroundColor: color, borderRadius: '50%', width: '10px', height: '10px' }} /> :null;
    };

    // `onClickDay` に渡すラッパー関数
    const handleClickDay = (date: Date) => {
        const jstDate = new Date(date.getTime() + 9 * 60 * 60 * 1000); // JSTオフセット（+9時間）
        const formattedDate = jstDate.toISOString().split('T')[0]; // Date → string に変換
        onDateClick(formattedDate); // 親コンポーネントの関数を呼び出す
    };

    return (
        <Calendar
        onClickDay={handleClickDay}
        tileContent={tileContent}
        />
    );
};

export default CalendarComponent;
