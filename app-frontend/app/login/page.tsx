'use client';

import { useRouter } from "next/navigation";
import { useState } from "react";
import { logIn, setAuthToken } from '../api';
import styles from '../LoginPage.module.css';

interface UserData {
    email: string,
    password: string
}

const LoginPage = () => {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [emailError, setEmailError] = useState<string>('');
    const [passwordError, setPasswordError] = useState<string>('');
    const router = useRouter();

    const handleLogIn = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
        e.preventDefault();
        
        // バリデーション
        let valid = true;
        setEmailError('');
        setPasswordError('');
        setError('');
  
        if (!email) {
            setEmailError('Email is required');
            valid = false;
        } else if (!/\S+@\S+\.\S+/.test(email)) {
            setEmailError('Invalid email format');
            valid = false;
        }
  
        if (!password) {
            setPasswordError('Password is required');
            valid = false;
        } else if (password.length < 6) {
            setPasswordError('Password must be at least 6 characters');
            valid = false;
        }
  
        if (!valid) return;
        

        try {
            const response = await logIn({ email, password } as UserData);
            const token = response.token;
            const userId = response.userId; //ユーザーIDを取得

            console.log("Response status:", response.status); // ステータスコード
            console.log("Response data:", response.data);     // データ全体

            setAuthToken(token);
            localStorage.setItem('token', token);
            localStorage.setItem('userId', userId);

            router.push(`/dashboard/${userId}`);
        } catch (error: any) {
            console.error('Login error:', error); 
            setError(error.response?.data?.message || 'Login failed');
        }
    };

    return (
        <div className={styles.container}>
            <h1 className={styles.title}>Log In</h1>
            <form onSubmit={handleLogIn}>
                <input
                className={styles.input}
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                />
                {emailError && <p className={styles.error}>{emailError}</p>}

                <input
                className={styles.input}
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                />
                {passwordError && <p className={styles.error}>{passwordError}</p>}

                <button className={styles.button} type="submit">Log In</button>
                {error && <p className={styles.error}>{error}</p>}
            </form>
        </div>
    );
};

export default LoginPage;