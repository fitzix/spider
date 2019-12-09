const fs = require('fs');
const os = require('os');

const puppeteer = require('puppeteer');

let phoneToCheckStr = fs.readFileSync('phone.txt').toString();

let phoneNumberArr = phoneToCheckStr.split(os.EOL).filter(el => el.trim().length > 0 && !isNaN(el));

let availables = [];

let run = async () => {
  const browser = await puppeteer.launch({
    headless: false,
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
    let find = await page.waitForSelector('#app > div > span:nth-child(2) > div > span > div > div > div > div', { visible: true }).catch(() => 'timeout');

    if (find === 'timeout') {
      console.log(item, ' ×××××× -- 超时');
      continue;
    }

    await page.waitFor(1000);
    let msg = await page.$eval('#app > div > span:nth-child(2) > div > span > div > div > div > div div > div._2Vo52', el => el.textContent).catch(() => '有效');
    if (msg === '透过网址分享的电话号码无效') {
      console.log(item, ' ××××××');
    } else {
      console.log(item, ' √√√√√√');
      availables.push(item);
    }
  }

  await browser.close();
  fs.writeFile('result.csv', availables.join(os.EOL), { flag: 'w+' }, err => {});
};

run();
