import { ReactElement } from ".pnpm/@types+react@18.0.27/node_modules/@types/react"
import DashLayout from "@/components/DashLayout"

export default function Page() {
  return (
    <>
      <h1>I am a saved page</h1>
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