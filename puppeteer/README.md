## whatsapp

- 配置 npm taobao 源 npm config set registry https://registry.npm.taobao.org
- 进入 puppeteer 目录, 执行 npm install
- 在配置文件 package.josn whatapp/chromePath 中配置本机 chrome.exe(windows)/Google Chrome(Macos)的绝对路径. 对于 window 一般在`C:\Users\username\AppData\Local\Google\Chrome SxS\Application\chrome.exe`
- 在 phone.txt 中配置需要检测的手机号,每个一行. 号码需带有国际区号且去掉+号
- 执行 npm run whatsapp
- 扫码登录,等待脚本执行

## FAQ

- A: 如何查看 chrome 执行文件路径 位置  
  Q: 在 chrome 浏览器地址栏输入 chrome://version, 在 "可执行文件路径" 下可找到 chrome 的执行文件路径.  
  **注意: 对于 windows 系统, 由于转义符问题 你需要把路径中的反斜线"\\"替换成正斜线"/"才可以配置在 package.json 中**
