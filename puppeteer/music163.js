const puppeteer = require('puppeteer');
const fs = require('fs');

const neteasyMusicBaseURL = 'https://music.163.com/#';

async function run() {
  const start = new Date();

  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto('https://music.163.com/#/discover/toplist');

  let results = [];
  // 获取auto id
  const contentFrame = page.frames().find(frame => frame.name() === 'contentFrame');
  let top100Elements = await contentFrame.$$('#song-list-pre-cache > div > div:nth-child(1) > table > tbody > tr');
  for (const [index, item] of top100Elements.entries()) {
    let name = await item.$eval(`${index > 2 ? ':nth-child(2)' : 'td.rank'} > div > div > div > span > a > b`, el => el.getAttribute('title'));
    let author = await item.$eval('td:nth-child(4) > div > span', el => el.getAttribute('title'));
    let songUrl = await item.$eval('td:nth-child(2) > div > div > div > span > a', el => el.getAttribute('href'));
    results.push({
      name,
      author,
      url: neteasyMusicBaseURL + songUrl,
      comments: [],
    });
  }

  for (const result of results) {
    await page.goto(result.url);
    await page.waitFor(1000);
    let curContentFrame = page.frames().find(frame => frame.name() === 'contentFrame');
    let commentElements = await curContentFrame.$$('#comment-box > div > div.m-cmmt > div.cmmts.j-flag > div > div.cntwrap > div:nth-child(1)');
    for (const commentElement of commentElements) {
      let comment = await commentElement.$eval('div', el => el.textContent);
      let author = await commentElement.$eval('div > a', el => el.textContent);
      result.comments.push({
        comment,
        author,
      });
    }
    console.log(result.comments);
  }

  const end = new Date();
  fs.writeFile('results.json', JSON.stringify(results), err => {
    console.error(err);
  });

  console.log('use', (end.getTime() - start.getTime()) / 1000);
  // name
  // #song-list-pre-cache > div > div:nth-child(1) > table > tbody > tr:nth-child(1) > td:nth-child(4) > div > span

  // await page.waitForNavigation()
  // await page.goto('https://github.com/search?q=john&type=Users&utf8=%E2%9C%93')

  // while (true) {
  //     await page.goto(`https://github.com/search?p=${pageIndex++}&q=fitzi&type=Users&utf8=%E2%9C%93`)
  //     // await page.screenshot({ path: `screenshot/${await page.title()}-${pageIndex}.png`, fullPage: true })
  //     let elements = await page.$$('#user_search_results > div.user-list > div > div.d-flex.flex-auto > div')
  //     if (elements.length === 0) {
  //         break
  //     }
  //     for (const [index, element] of elements.entries()) {
  //         let usernameElement = await element.$('a > em')
  //         if (usernameElement === null) {
  //             continue
  //         }
  //         let emailElement = await element.$('ul > li:nth-child(2) > a')
  //         emails.push({
  //             username: await page.evaluate(usernameElement => usernameElement.textContent, usernameElement),
  //             email: emailElement === null ?  '' : await page.evaluate(emailElement => emailElement.textContent, emailElement)
  //         })
  //     }
  // }
  await browser.close();
}

run();
