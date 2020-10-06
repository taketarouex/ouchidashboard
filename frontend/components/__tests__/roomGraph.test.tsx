import React from 'react';
import { render, waitForElementToBeRemoved, screen } from '@testing-library/react'
import { RoomGraph } from '../roomGraph'
import '@testing-library/jest-dom'
import dayjs from 'dayjs'

import { enableFetchMocks } from 'jest-fetch-mock'

test('RoomGraph', async () => {
  enableFetchMocks()
  fetchMock.mockResponseOnce(
    JSON.stringify([{ value: 0, updatedAt: dayjs(new Date(2020, 7, 31, 1, 0, 0, 0)) }]))
  const start = dayjs(new Date(2020, 7, 31, 0, 0, 0, 0))
  const end = dayjs(new Date(2020, 7, 31, 10, 0, 0, 0))

  const { getByRole } = render(<RoomGraph roomName={"test"} logType={"test"} start={start} end={end} />)
  await waitForElementToBeRemoved(getByRole("progressbar"))
  expect(screen.getByRole("listitem")).toBeInTheDocument()
  expect(fetchMock).toBeCalledWith(
    `/api/rooms/test/logs/test?start=${start.toISOString()}&end=${end.toISOString()}`)
})
