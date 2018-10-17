import json
import os
import time
from typing import Dict

from openpyxl import Workbook
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.webdriver.support.wait import WebDriverWait


print('开始创建浏览器')
chrome_options = webdriver.chrome.options.Options()
chrome_options.add_argument('headless')

chrome = webdriver.Chrome(options=chrome_options)

path = os.path.abspath(os.path.dirname(os.path.dirname(__file__)))

def check_cookie():
    print('检测是否有cookie')
    if os.path.exists('./data/cookie.json'):
        print('开始读取cookie')
        chrome.get('https://www.tianyancha.com/login')
        with open(os.path.join(path, 'data/cookie.json'), 'r') as fs:
            cookies = json.loads(fs.read())
            for cookie in cookies:
                chrome.add_cookie(cookie)

def login():
    chrome.find_element_by_xpath("//input[@onfocus=\"clearMsg('phone')\"]").send_keys(int(input('请输入手机号:')))
    chrome.find_element_by_css_selector("input[placeholder='请输入密码']").send_keys(input('请输入密码:'))
    chrome.find_element_by_css_selector("div[tyc-event-ch='Login.Login']").click()
    WebDriverWait(chrome, 10).until(EC.presence_of_element_located((By.CSS_SELECTOR, 'a[event-name=导航-用户中心]')))
    print('登陆成功,开始获取cookie')
    jsonStr = json.dumps(chrome.get_cookies())
    with open(os.path.join(path, 'data/cookie.json'), 'w') as f:
        f.write(jsonStr)

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
    name = chrome.find_element_by_css_selector('div.content div.header a').text
    status = chrome.find_element_by_css_selector('div.content div.header div').text

    legal_person = chrome.find_element_by_css_selector('div.info div:nth-child(1) a').text
    registered_capital = chrome.find_element_by_css_selector('div.info div:nth-child(2) span').text
    reg_date = chrome.find_element_by_css_selector('div.info div:nth-child(3) span').text

    phone = chrome.find_element_by_css_selector('div.contact div:nth-child(1) span.link-hover-click').text
    return name, status, legal_person, registered_capital, reg_date, phone


check_cookie()

chrome.get('https://www.tianyancha.com/login')


try:
    print('检测是否登录成功')
    chrome.find_element_by_css_selector('a[event-name=导航-用户中心]')
    print('登录成功')
except NoSuchElementException:
    print('cookie登录失败, 开始登录')
    login()
    

# 读取查询列表
with open(os.path.join(path, 'data/search.txt'), 'r') as f:
    companies = f.readlines()
    wb = Workbook()
    sheet = wb.active
    sheet.title = '天眼查结果'
    sheet.append(['名称', '状态', '法人代表', '注册资本', '注册时间', '联系电话'])
    for company in companies:
        result = find_result(company)
        if result:
            sheet.append(result)
    wb.save(os.path.join(path, 'data/result.xlsx'))

chrome.close()
chrome.quit()
