import { useEffect, useState } from 'react';
import styles from '../Modal.module.css';

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (text: string) => void;
    selectedDate: string | null;
    emotionText?: string;
}

const Modal = ({ isOpen, onClose, onSave, selectedDate, emotionText }: ModalProps) => {
    const [text, setText] = useState(emotionText || '');

    useEffect(() => {
        if (isOpen) {
            setText(emotionText || ''); // 編集時にテキストを設定
        }
    }, [isOpen, emotionText]);

    if (!isOpen) return null;

    const handleSave = () => {
        onSave(text);
        setText('');
    };

    return (
        <div className={styles.modal}>
            <div className={styles.modalContent}>
                <h2>{selectedDate ? `Edit Emotion for ${selectedDate}` : 'Write Your Emotion'}</h2>
                <textarea
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="How are you feeling today?"
                />
                <button onClick={handleSave}>Save</button>
                <button onClick={onClose}>Cancel</button>
            </div>
        </div>
    );
};

export default Modal;