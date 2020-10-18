import dayjs from 'dayjs';

const apiEndpoint = "http://localhost:3000/api";
const uiEndpoint = "http://localhost:3000"

test('chart', async () => {
  const temperatureLogs = [
    {
      Value: 0,
      UpdatedAt: dayjs("2020-07-31T00:00:00Z")
    },
    {
      Value: 1,
      UpdatedAt: dayjs("2020-07-31T01:00:00Z")
    },
    {
      Value: 2,
      UpdatedAt: dayjs("2020-07-31T02:00:00Z")
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
