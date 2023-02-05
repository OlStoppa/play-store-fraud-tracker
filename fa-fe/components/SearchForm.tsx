import styles from "../styles/SearchForm.module.scss";
import { useState } from "react";
import { locales } from '../constants/index';
import { Dispatch, SetStateAction } from "react";

type fetchResultsFn = ({ searchTerm, keyword } : { searchTerm: string, keyword: string }) => void;
interface Props { 
  fetchResults: fetchResultsFn,
  selectedLocales: string[],
  updateLocales: Dispatch<SetStateAction<string[]>>
  isLoading: boolean
};

export default function searchForm({ fetchResults, selectedLocales, updateLocales, isLoading }: Props) {
  const [searchTerm, updateSearchTerm] = useState('');
  const [keyword, updateKeyword] = useState('');
  
  const handleLocale = (id: string) => {
    if (id === 'ALL') {
      const update = selectedLocales.includes(id) ? [] : Object.keys(locales).map(key => key);
      updateLocales(update);
      return;
    }
    
    const update = selectedLocales.includes(id)
      ? selectedLocales.filter(key => key !== id)
      : [...selectedLocales, id];

    updateLocales(update);
  }

  return (
    <div className={styles.container}>
      <div className={styles.topbar}>
        <div>
          <span>Search Term:</span>
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => updateSearchTerm(e.target.value)}
          />
        </div>
        <div>
          <span>App Title Must Include:</span>
          <input type="text"
            value={keyword}
            onChange={(e) => updateKeyword(e.target.value)}
          />
        </div>
        <button
          onClick={() => fetchResults({ searchTerm, keyword })}
          disabled={!searchTerm.length || !keyword.length || !selectedLocales.length || isLoading }
          className={styles.searchBtn}
        >
          Search
        </button>
      </div>
      <p>Search Countries: </p>
      <div className={styles.locales}>
        {
          Object.entries(locales).map(locale => (
            <div className={styles.locale}>
              <input
                type="checkbox" id={locale[0]}
                checked={selectedLocales.includes(locale[0])}
                onChange={() => handleLocale(locale[0])}
              />
              <label htmlFor={locale[0]}>{locale[1]}</label>
            </div>
          ))
        }
      </div>
    </div>
  )
}