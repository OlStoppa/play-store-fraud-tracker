import { ReactElement, useState } from "react"
import DashLayout from "@/components/DashLayout"
import SearchForm from "@/components/SearchForm"
import AppBox from "@/components/AppBox";
import CountryBox from "@/components/CountryBox";
import styles from "../../styles/Search.module.scss"
import { App, SearchItem } from '../../types/index'
import { locales } from '../../constants/index';

export default function Page() {

  const [results, setResults] = useState<SearchItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedLocales, setLocales] = useState(() => {
    return Object.keys(locales);
  })

  const fetchResults = ({ searchTerm, keyword }: { searchTerm: string, keyword: string}) => {
    setLoading(true);
    const url = `/api/search?searchTerm=${searchTerm}&keyword=${keyword}`;
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({locales: selectedLocales.filter(item => item !== 'ALL')})
    })
      .then(res => res.json())
      .then((ret) => {
        setResults(ret as SearchItem[]);
        setLoading(false);
      })
  }

  const renderMain = () => {
    if (loading) return <div className={styles.spinner}/>;
    if (results.length) return results.map(res => {
      return (
        <CountryBox country={res.locale}>
          { res.apps.map(appData => <AppBox appData={appData} />) }
        </CountryBox>
      )
    });
    return
  }
  return (
    <>
      <SearchForm
        fetchResults={fetchResults}
        selectedLocales={selectedLocales}
        updateLocales={setLocales}
        isLoading={loading}
      />
      <div className={styles.container}>{renderMain()}</div>
    </>
  )
}

Page.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashLayout>
      {page}
    </DashLayout>
  )
}