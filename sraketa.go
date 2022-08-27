package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

var loadedPage playwright.Page
var playwrightError error
var playwrightInstance *playwright.Playwright
var playwrightBrowser playwright.Browser

func sraketa_init() bool {
	playwrightError = playwright.Install()
	if playwrightError != nil {
		log.Printf("Could not install playwright")
		return false
	}

	playwrightInstance, playwrightError = playwright.Run()
	if playwrightError != nil {
		return false
	}

	playwrightBrowser, playwrightError = playwrightInstance.Firefox.Launch()
	if playwrightError != nil {
		return false
	}

	loadedPage, playwrightError = playwrightBrowser.NewPage()

	if playwrightError != nil {
		return false
	}

	if _, playwrightError = loadedPage.Goto("https://alerts.in.ua", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); playwrightError != nil {
		return false
	}
	return true

}

func sraketa() bool {
	if playwrightError == nil {

		if _, err := loadedPage.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("sraketa_current.png"),
		}); err != nil {
			log.Printf("could not create screenshot: %v", err)
			return false
		}
		return true
	} else {
		return false
	}

}

func sraketa_shutdown() {
	if err := playwrightBrowser.Close(); err != nil {
		log.Printf("could not close browser: %v", err)
	}
	if err := playwrightInstance.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}

}
