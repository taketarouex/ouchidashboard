import React from 'react';
import { render, waitForElementToBeRemoved, fireEvent, within } from '@testing-library/react'
import { RoomInfo } from '../RoomInfo'
import '@testing-library/jest-dom'
import dayjs from 'dayjs'

import { enableFetchMocks } from 'jest-fetch-mock'

test('RoomInfo', async () => {
  enableFetchMocks()
  fetchMock.mockResponses(
    JSON.stringify(["hoge", "fuga", "bar"]),
    JSON.stringify([
      { value: 0, updatedAt: dayjs(new Date(2020, 7, 31, 0, 0, 0, 0)) },
      { value: 1, updatedAt: dayjs(new Date(2020, 7, 31, 1, 0, 0, 0)) },
      { value: 2, updatedAt: dayjs(new Date(2020, 7, 31, 2, 0, 0, 0)) }
    ])
  )
  const start = dayjs(new Date(2020, 7, 31, 0, 0, 0, 0))
  const end = dayjs(new Date(2020, 7, 31, 10, 0, 0, 0))

  const { container, getByRole, getByText } = render(<RoomInfo logType="test" start={start} end={end} />)
  await waitForElementToBeRemoved(getByRole("progressbar"))
  expect(fetchMock).toBeCalledWith("/api/rooms")
  expect(getByText("select rooms")).toBeInTheDocument()
  fireEvent.mouseDown(getByRole('button'))
  const listbox = within(getByRole('listbox'))
  fireEvent.click(listbox.getByText("hoge"))

  await waitForElementToBeRemoved(getByRole("progressbar"))

  expect(fetchMock).toBeCalledWith(
    `/api/rooms/hoge/logs/test?start=${start.toISOString()}&end=${end.toISOString()}`)
  expect(container).toMatchSnapshot()
})
