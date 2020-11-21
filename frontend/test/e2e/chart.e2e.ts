import dayjs from 'dayjs';
import puppeteer from 'puppeteer';

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
    if (request.url().startsWith(`${apiEndpoint}/rooms/hoge/logs/temperature`)) {
      request.respond({
        headers: { content: "application/json" },
        body: JSON.stringify(temperatureLogs)
      })
    } else if (request.url().startsWith(`${apiEndpoint}/rooms`)) {
      request.respond({
        headers: { content: "application/json" },
        body: JSON.stringify(["hoge", "fuga", "bar"])
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
  page.click("div#room-select");
  await page.waitForSelector("li");
  await page.waitFor(3000);
  const clicked = await page.screenshot();
  expect(clicked).toMatchImageSnapshot({
    comparisonMethod: 'ssim',
    failureThreshold: 0.01,
    failureThresholdType: 'percent'
  });
  page.click("li");
  await page.waitFor(3000);
  const selected = await page.screenshot();
  expect(selected).toMatchImageSnapshot({
    comparisonMethod: 'ssim',
    failureThreshold: 0.01,
    failureThresholdType: 'percent'
  });
})
