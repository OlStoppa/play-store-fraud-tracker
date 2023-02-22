import styles from '../styles/AppBox.module.scss'
import { App } from '../types/index' 

export default function AppBox({ appData }: { appData: App }) {
  return (
    <div className={styles.container}>
      <a href={appData.link} rel="noopener noreferrer" target="_blank" className={styles.linkContainer}>
        <div className={styles.imgbox}>
          <img src={appData.imgSrc} />
        </div>
        <div className={styles.info}>
          <img src={appData.thumb} />
          <div className={styles.titles}>
            <h5>{appData.name}</h5>
            <span>{appData.author}</span>
            {appData.rating && <span>{appData.rating}✭</span>}
          </div>
        </div>
      </a>
    </div>
  )
}