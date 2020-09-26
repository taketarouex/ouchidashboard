import Head from 'next/head'
import Layout, { siteTitle } from '../components/layout'
import { RoomGraph } from '../components/roomGraph'

export default function Home() {
  return (
    <Layout home>
      <Head>
        <title>{siteTitle}</title>
      </Head>
      <section>
        <h1>{siteTitle}</h1>
        <RoomGraph logType={"temperature"} />
      </section>
    </Layout>
  )
}
