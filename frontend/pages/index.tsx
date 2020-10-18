import Head from 'next/head'
import Layout, { siteTitle } from '../components/layout'
import { RoomGraph } from '../components/roomGraph'
import dayjs from 'dayjs'

export default function Home() {
  const start = dayjs(new Date(2020, 10, 18, 0, 0, 0, 0))
  const end = dayjs(new Date(2020, 10, 19, 0, 0, 0, 0))
  return (
    <Layout home>
      <Head>
        <title>{siteTitle}</title>
      </Head>
      <section>
        <h1>{siteTitle}</h1>
        <RoomGraph roomName={"fuwOX1K6757LpAo3y05j"} logType={"temperature"} start={start} end={end} />
      </section>
    </Layout>
  )
}
