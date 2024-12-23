'use client';

import { useRouter } from "next/navigation";
import { useState } from "react";
import { signUp } from '../api';
import styles from '../SignUpPage.module.css';

interface UserData {
    email: string,
    password: string
}

const SignUpPage = () => {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [emailError, setEmailError] = useState<string>('');
    const [passwordError, setPasswordError] = useState<string>('');
    const router = useRouter();

    const handleSignUp = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
        e.preventDefault();

        // バリデーションチェック
        let valid = true;
        setEmailError('');
        setPasswordError('');

        if (!email) {
            setEmailError('Email is required');
            valid = false;
        }else if (!/\S+@\S+\.\S+/.test(email)) {
            setEmailError('Invalid email format');
            valid = false;
        }

        if (!password) {
            setPasswordError('Password is required');
            valid = false;
        }else if (password.length < 6) {
            setPasswordError('Password must be at least 6 characters');
            valid = false;
        }

        if (!valid) return;

        try {
            await signUp({ email, password } as UserData);
            alert('Sign up successgull! Please login');
            router.push('/login')
        } catch (error: any) {
            setError(error.message || 'Sign up failed');
        }
    };

    return (
        <div className={styles.conrtainer}>
            <h1 className={styles.title}>Sign Up</h1>
            <form onSubmit={handleSignUp}>
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

                <button className={styles.button} type="submit">Sign Up</button>
                {error && <p className={styles.error}>{error}</p>}
            </form>
        </div>
    );
};

export default SignUpPage;