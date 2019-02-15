const puppeteer = require('puppeteer')

const creds = require('./creds')

const selector = require('./selector')

async function run() {
    const browser = await puppeteer.launch()
    const page = await browser.newPage()

    await page.goto('https://github.com/login')
    await page.type(selector.USERNAME_SELECTOR, creds.username)
    await page.type(selector.PASSWORD_SELECTOR, creds.password)
    await page.click(selector.BUTTON_SELECTOR)
    await page.waitForNavigation()
    // await page.goto('https://github.com/search?q=john&type=Users&utf8=%E2%9C%93')
    let pageIndex = 4
    let emails = []

    while (true) {
        await page.goto(`https://github.com/search?p=${pageIndex++}&q=fitzi&type=Users&utf8=%E2%9C%93`)
        let elements = await page.$$('#user_search_results > div.user-list > div')
        if (elements.length === 0) {
            break
        }
        let pageData = elements.map(element => {
            let username = await element.$eval('div.d-flex.flex-auto > div > a > em', usernameElement => usernameElement.textContent)
            let email = await element.$eval('div.d-flex.flex-auto > div > ul > li:nth-child(2) > a', emailsElement => emailsElement ? emailsElement.textContent : '')
            return {
                username,
                email
            }
        })
        emails = emails.concat(pageData)
    }

    console.log(emails)
    // let email = await page.evaluate(emailElement => emailElement.textContent, emailElement)
    await browser.close()
}

run()