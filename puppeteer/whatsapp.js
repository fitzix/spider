const fs = require('fs');
const os = require('os');

const config = require('./package.json');

const puppeteer = require('puppeteer-core');

const chromePath = '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';

let phoneToCheckStr = fs.readFileSync('phone.txt').toString();

let phoneNumberArr = phoneToCheckStr.split(os.EOL).filter(el => el.trim().length > 0 && !isNaN(el));

let availables = [];

let av;

let run = async () => {
  const browser = await puppeteer.launch({
    headless: false,
    executablePath: config.whatsapp.chromePath,
    defaultViewport: null,
  });
  const page = await browser.newPage();

  page.on('dialog', async dialog => await dialog.accept());

  await page.goto('https://web.whatsapp.com/');

  await page.waitForNavigation({ timeout: 0 });

  await page.waitFor(3000);

  console.log('开始检测');

  for (const item of phoneNumberArr) {
    await page.goto(`https://web.whatsapp.com/send?phone=${item}&text=2333`);
    await page.waitForSelector('#app > div > span:nth-child(2) > div > span > div > div > div > div', { visible: true });
    await page.waitFor(500);
    let msg = await page.$eval('#app > div > span:nth-child(2) > div > span > div > div > div > div div > div._2Vo52', el => el.textContent).catch(() => '有效');
    if (msg === '透过网址分享的电话号码无效') {
      console.log(item, ' xxx无效xxx');
    } else {
      console.log(item, ' ---有效---');
      availables.push(item);
    }
  }

  await browser.close();
  fs.writeFile('result.csv', availables.join(os.EOL), { flag: 'w+' }, err => {});
};

run();
