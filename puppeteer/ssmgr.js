const puppeteer = require('puppeteer-core');

(async () => {
  const browser = await puppeteer.launch({
    headless: false,
    executablePath: '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
  });
  const page = await browser.newPage();
  await page.goto('https://2simple.dev/home/login');
  await page.waitForSelector('#input_0');

  await page.focus('#input_0');
  await page.type('#input_0', 'caojunkaiv@gmail.com');
  await page.focus('#input_1');
  await page.type('#input_1', '232323');
  await page.click('body > div > div.layout-fill.ng-scope > div > section > div > div.layout-row.flex > md-content > div > div.hide-gt-sm.ng-scope.layout-column > div:nth-child(2) > div.layout-align-space-around-stretch.layout-column.flex > form > div > div > div:nth-child(2) > button');

  // await browser.close();
})();
