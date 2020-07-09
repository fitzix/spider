package whatsapp

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Run() {
	f, err := os.Open("phone.txt")
	if err != nil {
		log.Fatal("phone.txt 文件不存在")
	}
	defer f.Close()

	resultFile, err := os.OpenFile(fmt.Sprintf("%s.csv", time.Now().Format("2006-01-02.03:15:04")), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("创建结果文件失败--->%v", err)
	}
	defer resultFile.Close()

	allocator, allocatorCancel := chromedp.NewExecAllocator(
		context.Background(),
		chromedp.Flag("headless", false),
		chromedp.NoFirstRun,
	)
	defer allocatorCancel()

	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://web.whatsapp.com"),
		chromedp.WaitVisible("#pane-side", chromedp.ByID),
	); err != nil {
		log.Fatalf("登录失败--->%v", err)
	}

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			page.HandleJavaScriptDialog(true)
		}
	})

	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if errors.Is(err, io.EOF) {
			return
		}
		line = bytes.TrimSpace(line)
		if ok, err := regexp.Match(`(?m)[0-9]+`, line); err != nil || !ok {
			continue
		}
		if err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://web.whatsapp.com/send?phone=%s&text=2333", line)),
			chromedp.WaitVisible("#app > div > div > div:nth-child(4)", chromedp.ByQuery),
			chromedp.QueryAfter("#main > footer", func(ctx context.Context, node ...*cdp.Node) error {
				if len(node) > 0 {
					log.Printf("%s ------> √√√√√√", line)
					if _, err := fmt.Fprintln(resultFile, string(line)); err != nil {
						log.Printf("写入结果[%s]失败, err: %v", line, err)
					}
					return nil
				}
				return errors.New("not found")
			}, chromedp.AtLeast(0)),
		); err != nil {
			log.Printf("%s ------> ××××××", line)
		}
	}
}