import Head from 'next/head'
import Layout, { siteTitle } from '../components/layout'
import { RoomGraph } from '../components/roomGraph'
import dayjs from 'dayjs'

export default function Home() {
  const start = dayjs(new Date(2020, 9, 28, 0, 0, 0, 0))
  const end = dayjs(new Date(2020, 9, 28, 0, 0, 0, 0))
  return (
    <Layout home>
      <Head>
        <title>{siteTitle}</title>
      </Head>
      <section>
        <h1>{siteTitle}</h1>
        <RoomGraph roomName={"living"} logType={"temperature"} start={start} end={end} />
      </section>
    </Layout>
  )
}
