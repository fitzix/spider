from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import NoSuchElementException

import time
import json
import os

def create_brower():
    chrome_options = webdriver.chrome.options.Options()
    chrome_options.add_argument('headless')
    return webdriver.Chrome(options=chrome_options)

def check_cookie():
    if os.path.exists('./data/cookie.json'):
        chrome.get("https://www.tianyancha.com")
        with open('./data/cookie.json', 'r') as f:
            cookies = json.loads(f.read())
            for cookie in cookies:
                chrome.add_cookie(cookie)

chrome = create_brower()

check_cookie()

# 将刚刚复制的帖在这
chrome.get("https://www.tianyancha.com/search?key=百度")


try:
    chrome.find_element_by_xpath("//div[text()='登录']")
#  关闭无头浏览器
    chrome.close()
    chrome.quit()
# 
    chrome = webdriver.Chrome()

    chrome.find_element_by_xpath("//input[@onfocus=\"clearMsg('phone')\"]").send_keys(13699146887)
    chrome.find_element_by_xpath("//div[text()='短信验证码登录']").click()
    chrome.find_element_by_id('smsCodeBtn').click()
    time.sleep(1)
    chrome.find_element_by_xpath("//div[@class='pb10 over-hide position-rel']/input[@placeholder='请输入验证码']").send_keys(int(input('输入验证码:')))
    chrome.find_element_by_xpath("//div[@onclick='loginByMes()']").click()

    time.sleep(10)
    print('开始获取cookie')
    jsonStr = json.dumps(chrome.get_cookies())
    with open('./data/cookie.json', 'w') as f:
        f.write(jsonStr)
    
    chrome.close()
    chrome.quit()


    chrome = create_brower()
    check_cookie()

except NoSuchElementException:
    print('已登录')


#  获取详细信息
def find_result(name):
    print('正在查询---%s' % name, end='')
    chrome.get("https://www.tianyancha.com/search?key=%s" % name)
    try:
        results = chrome.find_elements_by_class_name('search-result-single')
    except NoSuchElementException:
        print('没有结果')
        return
    if len(results) == 0:
        return
    result = {}
    result['name'] = chrome.find_element_by_css_selector('div.content div.header a').text
    result['status'] = chrome.find_element_by_css_selector('div.content div.header div').text

    result['legalPerson'] = chrome.find_element_by_css_selector('div.info div:nth-child(1) a').text
    result['registeredCapital'] = chrome.find_element_by_css_selector('div.info div:nth-child(2) span').text
    result['regDate'] = chrome.find_element_by_css_selector('div.info div:nth-child(3) span').text

    result['phone'] = chrome.find_element_by_css_selector('div.contact div:nth-child(1) span.link-hover-click').text
    print(result)
    return result

# 读取查询列表

with open('./data/search.txt') as f:
    companies = f.readlines()
    with open('./data/result.txt', 'w') as r:
        for company in companies:
            r.write(str(find_result(company)) + '\n')

chrome.close()
chrome.quit()

