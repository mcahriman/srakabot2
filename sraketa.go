package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func sraketa() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Firefox.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://alerts.in.ua", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("sraketa_current.png"),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}

}
