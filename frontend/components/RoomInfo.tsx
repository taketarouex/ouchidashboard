import React, { FC, useState } from 'react';
import CircularProgress from '@material-ui/core/CircularProgress'
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { LineChart, XAxis, YAxis, Line, Tooltip, ResponsiveContainer } from 'recharts'
import useSWR from 'swr'
import dayjs, { Dayjs } from 'dayjs'

type Log = {
  Value: number
  UpdatedAt: Date
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

const useRoomLog = ({ roomName, logType, start, end }: { roomName: string, logType: string, start: Dayjs, end: Dayjs }) => {
  const startISO = start.toISOString()
  const endISO = end.toISOString()
  const { data, error } = useSWR(`/api/rooms/${roomName}/logs/${logType}?start=${startISO}&end=${endISO}`, fetchLogs)

  return {
    data: data,
    isLoading: !error && !data,
    isError: error
  }
}

const RoomChart = ({ roomName, logType, start, end }: { roomName: string, logType: string, start: dayjs.Dayjs, end: dayjs.Dayjs }) => {
  if (roomName === "") return <div>select rooms</div>
  const { data, isLoading, isError } = useRoomLog({ roomName, logType, start, end })
  if (isLoading) return <CircularProgress />
  if (isError) return <div>error</div>
  return (
    <ResponsiveContainer width={'90%'} height={400}>
      <LineChart
        data={data}
      >
        <YAxis />
        <XAxis dataKey="UpdatedAt" domain={['dataMin', 'dataMax']} />
        <Line type="monotone" dataKey="Value" stroke="#8884d8" />
        <Tooltip />
      </LineChart>
    </ResponsiveContainer>
  )
}

const fetchRoomNames: (url: string) => Promise<string[]> = (url) => fetch(url).then(
  res => {
    if (!res.ok) {
      const error = new Error('An error occurred while fetching the data.')
      throw error
    }
    return res.json()
  }
)

const useRoomNames = () => {
  const { data, error } = useSWR(`/api/rooms`, fetchRoomNames)
  return {
    data: data,
    isLoading: !error && !data,
    isError: error
  }
}

const RoomSelect = ({ room, handleChange }: { room: string, handleChange: (event: React.ChangeEvent<{ value: unknown }>) => void }) => {
  const { data, isLoading, isError } = useRoomNames();
  if (isLoading) return <CircularProgress />
  if (isError) return <div>error</div>
  return (
    <Select
      labelId="room"
      id="room-select"
      value={room}
      onChange={handleChange}
    >
      {data.map(v =>
        <MenuItem value={v} key={v}>{v}</MenuItem>
      )}
    </Select>
  )
}

export const RoomInfo = ({ logType, start, end }: { logType: string, start: Dayjs, end: Dayjs }) => {
  const [room, setRoom] = useState<string>("")
  const handleChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setRoom(event.target.value as string)
  }
  return (
    <div>
      <RoomSelect room={room} handleChange={handleChange} />
      <RoomChart roomName={room} logType={logType} start={start} end={end} />
    </div>
  )
}