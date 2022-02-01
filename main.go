package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

const (
	googleSignin = "https://accounts.google.com"

	itviecBasePath = "https://itviec.com"
	itviecJobsPath = "/it-jobs"
	itviecSignin   = "/sign_in"

	vnwBasePath = "https://secure.vietnamworks.com"
	vnwSignin   = "/login/vi?client_id=3"
)

func newChromedp() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("start-fullscreen", true),

		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("remote-debugging-port", "9222"),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	googleTask(ctx)
	itviecWithGoogleTask(ctx)
	extractItviecTask(ctx)
	// itviecTask(ctx)

	// vietnamworksWithGoogleTask(ctx)
	// vietnamworksTask(ctx)

	return ctx, cancel
}

func googleTask(ctx context.Context) {
	email := "//*[@id='identifierId']"
	password := "//*[@id='password']/div[1]/div/div[1]/input"
	buttonEmailNext := "//*[@id='identifierNext']/div/button"
	buttonPasswordNext := "//*[@id='passwordNext']/div/button/span"

	task := chromedp.Tasks{
		chromedp.Navigate(googleSignin),
		chromedp.SendKeys(email, "email"),
		chromedp.Sleep(1 * time.Second),

		chromedp.Click(buttonEmailNext),
		chromedp.Sleep(1 * time.Second),

		chromedp.SendKeys(password, "pw"),
		chromedp.Sleep(1 * time.Second),

		chromedp.Click(buttonPasswordNext),
		chromedp.Sleep(2 * time.Second),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
}

func itviecTask(ctx context.Context) {
	url := fmt.Sprintf("%s%s", itviecBasePath, itviecSignin)
	email := "//*[@id='user_email']"
	password := "//*[@id='user_password']"
	label := "//*[@id='container']/div[2]/div/div[2]/form/div[4]/div/div/div/iframe"
	button := "//*[@id='container']/div[2]/div/div[2]/form/div[5]/div/button"
	task := chromedp.Tasks{
		chromedp.Navigate(url),

		chromedp.SendKeys(email, "email"),
		chromedp.Sleep(2 * time.Second),

		chromedp.SendKeys(password, "pw"),
		chromedp.Sleep(2 * time.Second),

		chromedp.Click(label),
		chromedp.Sleep(2 * time.Second),

		chromedp.Click(button),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
}

func vietnamworksTask(ctx context.Context) {
	url := fmt.Sprintf("%s%s", vnwBasePath, vnwSignin)
	email := "//*[@id='email']"
	password := "//*[@id='login__password']"
	button := "//*[@id='button-login']"
	task := chromedp.Tasks{
		chromedp.Navigate(url),

		chromedp.SendKeys(email, "email"),
		chromedp.Sleep(2 * time.Second),

		chromedp.SendKeys(password, "pw"),
		chromedp.Sleep(2 * time.Second),

		chromedp.Click(button),
		chromedp.Sleep(2 * time.Second),

		chromedp.Navigate("https://www.vietnamworks.com/tim-viec-lam/tat-ca-viec-lam"),
		chromedp.Sleep(3 * time.Second),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
}

func vietnamworksWithGoogleTask(ctx context.Context) {
	url := fmt.Sprintf("%s%s", vnwBasePath, vnwSignin)
	button := "/html/body/div[2]/div[1]/div/div/a[2]"

	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),

		chromedp.Click(button),
		chromedp.Sleep(2 * time.Second),

		chromedp.Navigate("https://www.vietnamworks.com/tim-viec-lam/tat-ca-viec-lam"),
		chromedp.Sleep(2 * time.Second),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
}

func itviecWithGoogleTask(ctx context.Context) {
	url := fmt.Sprintf("%s%s", itviecBasePath, itviecSignin)

	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
}

func extractItviecTask(ctx context.Context) error {
	task := chromedp.Tasks{
		chromedp.Navigate("https://itviec.com/it-jobs/java-developer-mysql-spring-oracle-cj-olivenetworks-vina-co-ltd-0324"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return err
			}

			doc.Find("div.job-details__overview div.svg-icon__text").Each(func(index int, info *goquery.Selection) {
				text := info.Text()
				fmt.Println(text)
			})

			return nil
		}),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
	return nil
}

func main() {
	_, _ = newChromedp()

	// close chrome
	// _, cancel := newChromedp()
	// defer cancel()
}
