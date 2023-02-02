import { ReactElement, useState } from "react"
import DashLayout from "@/components/DashLayout"
import SearchForm from "@/components/SearchForm"
import AppBox from "@/components/AppBox";

export default function Page() {

  const [results, setResults] = useState<{ name: string }[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchResults = () => {
    setLoading(true);
    fetch('http://localhost:9000/api')
      .then(res => res.json())
      .then(ret => {
        setResults(ret);
        setLoading(false);
      })
  }

  const renderMain = () => {
    if (loading) return <h1>Loading</h1>;
    else if (!results.length) return <h1>Try a new search!</h1>;
    else return results.map(res => <AppBox appData={res} />)
  }
  return (
    <>
      <SearchForm fetchResults={fetchResults} />
      <div>{renderMain()}</div>
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