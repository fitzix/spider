package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	successSymbol = "√√√√√√"
	errSymbol     = "××××××"
)

var (
	total, invalid, success, fail int
)

func Run() {
	f, err := os.Open("phone.txt")
	if err != nil {
		log.Fatal("phone.txt 文件不存在")
	}
	defer f.Close()

	resultFile, err := os.OpenFile(fmt.Sprintf("%s.csv", time.Now().Format("20060102031504")), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("创建结果文件失败--->%v", err)
	}
	defer resultFile.Close()

	allocator, allocatorCancel := chromedp.NewExecAllocator(
		context.Background(),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)
	defer allocatorCancel()

	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://web.whatsapp.com"),
		chromedp.WaitVisible("#pane-side", chromedp.ByID),
	); err != nil {
		log.Printf("登录失败--->%v", err)
		return
	}

	if _, err := resultFile.WriteString("\xEF\xBB\xBF"); err != nil {
		log.Fatalf("写入csv bom err: ---> %v", err)
	}

	w := csv.NewWriter(resultFile)

	_ = w.Write([]string{"phone", "status"})

	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		total++
		line = bytes.TrimSpace(line)
		if ok, err := regexp.Match(`(?m)[0-9]+`, line); err != nil || !ok {
			invalid++
			continue
		}
		if err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://web.whatsapp.com/send?phone=%s&text=2333", line)),
			chromedp.WaitVisible("#app > div > div > div:nth-child(4)", chromedp.ByQuery),
			chromedp.Sleep(time.Millisecond*500),
			chromedp.QueryAfter("#main > footer", func(ctx context.Context, node ...*cdp.Node) error {
				symbol := errSymbol
				if len(node) > 0 {
					success++
					symbol = successSymbol
				} else {
					fail++
				}

				log.Printf("%s ------> %s", line, symbol)

				if err := w.Write([]string{string(line), symbol}); err != nil {
					log.Printf("写入结果失败 ---> %s", err)
				}
				w.Flush()
				return nil
			}, chromedp.AtLeast(0)),
		); err != nil {
			fail++
			log.Printf("%s ------> %s", line, errSymbol)
		}
	}
	log.Printf("共检测: %d, 格式错误: %d, 有效: %d, 无效:%d, 有效率: %.2f%%", total, invalid, success, fail, float64(success*100)/float64(total))
}
