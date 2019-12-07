const puppeteer = require('puppeteer-core');

(async () => {
  const browser = await puppeteer.launch({
    headless: false,
    executablePath: '/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary',
  });
  const page = await browser.newPage();
  await page.goto('https://keepa.com');
  await page.setViewport({ width: 1366, height: 768 });
  await page.waitForSelector('#userMenuPanel');
  //   page.waitFor()
  await page.click('#userMenuPanel');

  await page.focus('#username');
  await page.type('#username', 'zengweigang');
  await page.focus('#password');
  await page.type('#password', 'gang12345');

  await page.click('#submitLogin');

  await page.waitForFunction('document.getElementById("panelUsername").innerText == "zengweigang"');

  await page.goto('https://keepa.com/#!categorytree');

  await page.waitForSelector('#grid-wrapper-category');

  // await browser.close();
})();
