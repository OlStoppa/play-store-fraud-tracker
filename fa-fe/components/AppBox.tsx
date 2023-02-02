import styles from '../styles/AppBox.module.scss'

export default function AppBox({ appData }) {
  return (
    <div className={styles.container}>
      <div className={styles.imgbox}>
        <img src={appData.imgSrc} />
      </div>
      <div className={styles.info}>
        <img src={appData.thumb} />
        <div className={styles.titles}>
          <h5>{appData.name}</h5>
          <span>{appData.author}</span>
          <span>{appData.rating}✭</span>
        </div>
      </div>
    </div>
  )
}