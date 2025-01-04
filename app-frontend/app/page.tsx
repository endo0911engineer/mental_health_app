'use client';

import styles from './HomePage.module.css';
import Link from 'next/link';

export default function HomePage() {
  return (
    <main className={styles.main}>
      {/* ヘッダー */}
      <header className={styles.header}>
        <h1 className={styles.logo}>AI Journal</h1>
        <nav className={styles.navWrapper}>
          <Link href="/login" className={`${styles.navLink}`}>
            ログイン
          </Link>
          <Link href="/signup" className={`${styles.navLink} ${styles.navSignUp}`}>
            サインアップ
          </Link>
        </nav>
      </header>

      {/* メインコンテンツ */}
      <section>
        <h2 className={styles.heading}>
          AIで毎日の日記を <span className={styles.highlight}>より深く理解</span>
        </h2>
        <p className={styles.description}>
          あなたの気持ちや考えを簡単に記録し、AIが分析してサポート。
          自分の心の状態を見つめ直すための新しい方法です。
        </p>
        <Link href="/signup">
          <button className={styles.ctaButton}>今すぐ始める</button>
        </Link>
      </section>

      {/* 特徴 */}
      <section className={styles.features}>
        <div className={styles.featureCard}>
          <h3>AI解析の力</h3>
          <p>日記の中のポジティブ・ネガティブな感情を可視化し、トレンドを提供します。</p>
        </div>
        <div className={styles.featureCard}>
          <h3>簡単操作</h3>
          <p>シンプルなインターフェースで、毎日数分で記録可能です。</p>
        </div>
        <div className={styles.featureCard}>
          <h3>完全プライバシー保護</h3>
          <p>データは暗号化され、安心してご利用いただけます。</p>
        </div>
      </section>

      {/* フッター */}
      <footer className={styles.footer}>
        © 2025 AI Journal. All rights reserved.
      </footer>
    </main>
  );
}