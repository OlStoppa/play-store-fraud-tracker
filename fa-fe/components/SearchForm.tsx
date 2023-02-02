import styles from "../styles/SearchForm.module.scss";

export default function searchForm({ fetchResults }: { fetchResults: () => void }) {
  return (
    <div className={styles.container}>
      <div>
        <span>Search Term:</span>
        <input type="text" />
      </div>
      <div>
        <span>App Title Must Include:</span>
        <input type="text" />
      </div>
      <button onClick={fetchResults} className={styles.searchBtn}>
        Search
      </button>
    </div>
  )
}