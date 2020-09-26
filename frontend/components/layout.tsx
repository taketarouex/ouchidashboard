import Head from 'next/head'
import styles from './layout.module.css'
import Link from 'next/link'

export const siteTitle = 'Dashboard'

export default function Layout({ children, home }: { children: React.ReactNode, home?: boolean }) {
  return (
    <div className={styles.container}>
      <Head>
        <meta
          name="description"
          content="Dashboard for ouchi"
        />
        <meta name="og:title" content={siteTitle} />
      </Head>
      <header className={styles.header}>
      </header>
      <main>{children}</main>
      {!home && (
        <div className={styles.backToHome}>
          <Link href="/">
            <a>← Back to home</a>
          </Link>
        </div>
      )}
    </div>
  )
}
