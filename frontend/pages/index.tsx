import Head from 'next/head'
import Layout, { siteTitle } from '../components/layout'
import { RoomInfo } from '../components/RoomInfo'
import dayjs from 'dayjs'

export default function Home() {
  const start = dayjs().set('hour', 0).set('minute', 0).set('second', 0).set('millisecond', 0);
  const end = start.add(1, 'day');
  return (
    <Layout home>
      <Head>
        <title>{siteTitle}</title>
      </Head>
      <section>
        <h1>{siteTitle}</h1>
        <RoomInfo logType={"temperature"} start={start} end={end} />
      </section>
    </Layout>
  )
}
