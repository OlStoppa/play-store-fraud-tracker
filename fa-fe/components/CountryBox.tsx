import { ReactElement } from 'react'
import styles from '../styles/CountryBox.module.scss'
import { locales } from '@/constants'

export default function SelectedApp ({ children, country }: { children: ReactElement[], country: keyof typeof locales }) {
  return (
    <div className={styles.container}>
      <div className={styles.title}>
        <p>{ locales[country] }</p>
      </div>
      <div className={styles.apps}>
        { children }
      </div>
    </div>
  )
}