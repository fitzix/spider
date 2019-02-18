const puppeteer = require('puppeteer')

const creds = require('./creds')

const selector = require('./selector')

async function run() {
    const browser = await puppeteer.launch()
    const page = await browser.newPage()
    await page.setViewport({ width: 1920, height: 1080 })
    await page.goto('https://github.com/login')
    await page.type(selector.USERNAME_SELECTOR, creds.username)
    await page.type(selector.PASSWORD_SELECTOR, creds.password)
    await page.click(selector.BUTTON_SELECTOR)
    await page.waitForNavigation()
    // await page.goto('https://github.com/search?q=john&type=Users&utf8=%E2%9C%93')
    let pageIndex = 1
    let emails = []

    while (true) {
        await page.goto(`https://github.com/search?p=${pageIndex++}&q=fitzi&type=Users&utf8=%E2%9C%93`)
        // await page.screenshot({ path: `screenshot/${await page.title()}-${pageIndex}.png`, fullPage: true })
        let elements = await page.$$('#user_search_results > div.user-list > div > div.d-flex.flex-auto > div')
        if (elements.length === 0) {
            break
        }
        for (const [index, element] of elements.entries()) {
            let usernameElement = await element.$('a > em')
            if (usernameElement === null) {
                continue
            }
            let emailElement = await element.$('ul > li:nth-child(2) > a')
            emails.push({
                username: await page.evaluate(usernameElement => usernameElement.textContent, usernameElement),
                email: emailElement === null ?  '' : await page.evaluate(emailElement => emailElement.textContent, emailElement)
            })
        }
    }
    await browser.close()
}

run()