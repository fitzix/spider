const fs = require('fs');
const os = require('os');

const puppeteer = require('puppeteer');

let phoneToCheckStr = fs.readFileSync('phone.txt').toString();

let phoneNumberArr = phoneToCheckStr.split(os.EOL).filter((el) => el.trim().length > 0 && !isNaN(el));

let availables = [];

let run = async () => {
  const browser = await puppeteer.launch({
    headless: false,
    defaultViewport: null,
  });
  const page = await browser.newPage();

  page.on('dialog', async (dialog) => await dialog.accept());

  await page.goto('https://web.whatsapp.com/');

  await page.waitForNavigation({ timeout: 0 });

  // await page.waitFor(3000);

  console.log('开始检测');

  for (const item of phoneNumberArr) {
    process.stdout.write(`号码: ${item}  ------- `);
    await page.goto(`https://web.whatsapp.com/send?phone=${item}&text=2333`);
    let find = await page.waitForSelector('#app > div > div > div:nth-child(4)', { visible: true }).catch(() => 'timeout');

    if (find === 'timeout') {
      process.stdout.write('超时' + '\n');
      continue;
    }

    await page.waitFor(1000);
    if ((await page.$('#main > footer')) !== null) {
      process.stdout.write('√√√√√√' + '\n');
    } else {
      process.stdout.write('××××××' + '\n');
      availables.push(item);
    }
  }

  await browser.close();
  fs.writeFile('result.csv', availables.join(os.EOL), { flag: 'w+' }, (err) => {});
};

run();
