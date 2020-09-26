import dayjs from 'dayjs'
import { FC } from 'react'
import useSWR from 'swr'

type Log = {
  value: number
  updatedAt: Date
}

const fetchLogs: (url: string) => Promise<Log[]> = (url) => fetch(url).then(res => res.json())

export const RoomGraph = ({ logType }: { logType: string }) => {
  const start = dayjs().add(-1, 'day').toISOString()
  const end = dayjs().toISOString()
  const { data, error } = useSWR(`/rooms/living/logs/${logType}?start=${start}&end=${end}`, fetchLogs)

  if (error) return <div>{error.message}</div>

  return (
    <div>
      {error && <div>{error.message}</div>}
      {!data && <div>loading.....</div>}
      {data && data.map((v) => <li>{v}</li>)}
    </div>
  )
}