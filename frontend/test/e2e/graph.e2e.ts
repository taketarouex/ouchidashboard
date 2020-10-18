import dayjs from 'dayjs';

const apiEndpoint = "http://localhost:3000/api";
const uiEndpoint = "http://localhost:3000"

test('graph', async () => {
  const temperatureLogs = [
    {
      value: 0,
      updatedAt: dayjs(new Date(2020, 7, 31, 0, 0, 0, 0))
    },
    {
      value: 1,
      updatedAt: dayjs(new Date(2020, 7, 31, 1, 0, 0, 0))
    },
    {
      value: 2,
      updatedAt: dayjs(new Date(2020, 7, 31, 2, 0, 0, 0))
    },
  ]
  page.setDefaultNavigationTimeout(0);
  await page.setRequestInterception(true)
  page.on("request", (request) => {
    if (request.url().startsWith(`${apiEndpoint}/rooms`)) {
      request.respond({
        headers: { content: "application/json" },
        body: JSON.stringify(temperatureLogs)
      })
    } else {
      request.continue()
    }
  })
  await page.goto(`${uiEndpoint}`, { waitUntil: ['networkidle2'] })
  const rendered = await page.screenshot();
  expect(rendered).toMatchImageSnapshot({
    comparisonMethod: 'ssim',
    failureThreshold: 0.01,
    failureThresholdType: 'percent'
  });
})
