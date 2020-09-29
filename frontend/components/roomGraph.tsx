import CircularProgress from '@material-ui/core/CircularProgress'
import useSWR from 'swr'
import dayjs from 'dayjs'

type Log = {
  value: number
  updatedAt: Date
}

const fetchLogs: (url: string) => Promise<Log[]> = (url) => fetch(url).then(
  res => {
    if (!res.ok) {
      const error = new Error('An error occurred while fetching the data.')
      throw error
    }
    return res.json()
  }
)

const useRoomLog = ({ roomName, logType, start, end }: { roomName: string, logType: string, start: dayjs.Dayjs, end: dayjs.Dayjs }) => {
  const startISO = start.toISOString()
  const endISO = end.toISOString()
  const { data, error } = useSWR(`/api/rooms/${roomName}/logs/${logType}?start=${startISO}&end=${endISO}`, fetchLogs)

  return {
    logs: data,
    isLoading: !error && !data,
    isError: error
  }
}

export const RoomGraph = ({ roomName, logType, start, end }: { roomName: string, logType: string, start: dayjs.Dayjs, end: dayjs.Dayjs }) => {
  const { logs, isLoading, isError } = useRoomLog({ roomName, logType, start, end })
  if (isLoading) return <CircularProgress />
  if (isError) return <div>error</div>
  return <div>{logs.map((v) => <li>{v.value}:{v.updatedAt}</li>)}</div>
}