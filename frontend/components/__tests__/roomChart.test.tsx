import React from 'react';
import { render, waitForElementToBeRemoved } from '@testing-library/react'
import { RoomChart } from '../roomChart'
import '@testing-library/jest-dom'
import dayjs from 'dayjs'

import { enableFetchMocks } from 'jest-fetch-mock'

test('RoomChart', async () => {
  enableFetchMocks()
  fetchMock.mockResponseOnce(
    JSON.stringify([
      { value: 0, updatedAt: dayjs(new Date(2020, 7, 31, 0, 0, 0, 0)) }
    ]))
  const start = dayjs(new Date(2020, 7, 31, 0, 0, 0, 0))
  const end = dayjs(new Date(2020, 7, 31, 10, 0, 0, 0))

  const { container, getByRole } = render(<RoomChart roomName={"test"} logType={"test"} start={start} end={end} />)
  expect(fetchMock).toBeCalledWith(
    `/api/rooms/test/logs/test?start=${start.toISOString()}&end=${end.toISOString()}`)
  await waitForElementToBeRemoved(getByRole("progressbar"))
  expect(container).toMatchSnapshot()
})
