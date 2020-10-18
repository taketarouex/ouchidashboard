import Head from 'next/head'
import Layout, { siteTitle } from '../components/layout'
import { RoomChart } from '../components/roomChart'
import dayjs from 'dayjs'

export default function Home() {
  const start = dayjs("2020-01-23T00:00:00Z")
  const end = dayjs("2020-01-23T04:00:00Z")
  return (
    <Layout home>
      <Head>
        <title>{siteTitle}</title>
      </Head>
      <section>
        <h1>{siteTitle}</h1>
        <RoomChart roomName={"testRoom"} logType={"temperature"} start={start} end={end} />
      </section>
    </Layout>
  )
}
