import styles from "../styles/SearchForm.module.scss";

interface Props { 
  searchParams: {
    searchTerm: string,
    keyword: string,
  }
  clearSearch: () => void
};

export default function resultBar({ searchParams, clearSearch }: Props) {

  return (
    <div className={styles.container}>
      <div className={styles.topbar}>
        <div>
          <p>Search Term:</p>
          <p>{searchParams.searchTerm}</p>
        </div>
        <div>
          <p>Keyword:</p>
          <p>{searchParams.keyword}</p>
        </div>
        <button
          onClick={clearSearch}
          className={styles.searchBtn}
        >
          Clear
        </button>
      </div>
    </div>
  )
}