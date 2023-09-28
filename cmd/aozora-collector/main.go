package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Entry struct {
	AuthorID string
	Author string
	TitleId string
	Title string
	InfoURL string
	ZipURL string
}

func findEntries(siteURL string) ([]Entry, error) {
	// URLからDOMオブジェクトを作成
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}
	// 処理
	doc.Find("ol li a").Each(func(i int, elem *goquery.Selection) {
		println(elem.Text(), elem.AttrOr("href", ""))
	})
	return nil, nil
}

func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person338.html"

	entries, err := findEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Println(entry.Title, entry.ZipURL)
	}
}