import { ReactElement, useState, useEffect } from "react";
import DashLayout from "@/components/DashLayout";
import SearchForm from "@/components/SearchForm";
import AppBox from "@/components/AppBox";
import CountryBox from "@/components/CountryBox";
import ResultsBar from "@/components/ResultsBar";
import styles from "@/styles/Search.module.scss";
import { App, SearchItem } from '@/types/index';
import { locales } from '@/constants/index';
import axios from 'axios';
import Router from 'next/router';
import { NextPageWithLayout } from './_app';

const Page: NextPageWithLayout<{ username: string }> = ({ username }: { username: string}) => {
  const [results, setResults] = useState<SearchItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedLocales, setLocales] = useState(() => {
    return Object.keys(locales);
  })
  const [searchParams, setSearchParams] = useState({ searchTerm: '', keyword: ''});
  
  useEffect(() => {
    if (!username) Router.replace('/login');
  }, [username])

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
        setSearchParams({ searchTerm, keyword })
        setLoading(false);
      })
  }
  
  const handleClearSearch = () => {
    setResults([]);
    setSearchParams({ searchTerm: '', keyword: '' })
  }

  const renderMain = () => {
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
      {
        results.length
        ? <ResultsBar
            searchParams={searchParams}
            clearSearch={handleClearSearch}
          />
        : <SearchForm
            fetchResults={fetchResults}
            selectedLocales={selectedLocales}
            updateLocales={setLocales}
            isLoading={loading}
          />
      }
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

Page.getInitialProps = async ({ req, res }) => {
  if (req && res) {
    try {
      const response = await axios.get<{username: string}>(
        'http://pssearch.bespoken.live/api/get-user',
        { headers: req.headers }
      );
      if (!response.data.username) {
       throw new Error()
      }
      return response.data;
    } catch (e) {
      res.writeHead(301, { Location: '/login' });
      res.end();
    }
  } else {
    try {
      const response = await axios.get<{username: string}>('/api/get-user')
      if (!response.data.username) {
        throw new Error();
      }
      return response.data;
    } catch (e) {
      Router.replace('/login')
    }
  }
  return { username: ''};
}

export default Page;