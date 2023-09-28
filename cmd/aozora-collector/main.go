package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strings"

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

func findAuthorAndZIP(siteURL string) (string, string) {
	log.Println("query", siteURL)
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return "", ""
	}

	author := doc.Find("table[summary=作家データ] tr:nth-child(1) td:nth-child(2)").Text()

	zipURL := ""
	doc.Find("table.download a").Each(func(n int, elem *goquery.Selection) {
		href := elem.AttrOr("href", "")
		if strings.HasSuffix(href, ".zip") {
			zipURL = href
		}
	})

	if zipURL == "" {
		return author, ""
	}
	if strings.HasPrefix(zipURL, "http://") || strings.HasPrefix(zipURL, "https://") {
		return author, zipURL
	}

	u, err := url.Parse(siteURL)
	if err != nil {
		return author, ""
	}
	u.Path = path.Join(path.Dir(u.Path), zipURL)
	return author, u.String()
}

func findEntries(siteURL string) ([]Entry, error) {
	// URLからDOMオブジェクトを作成
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}
	// 処理
	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)

	entries := []Entry{}
	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		title := elem.Text()
		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/card%s.html", token[1], token[2])
		author, zipURL := findAuthorAndZIP(pageURL) // 作者とzipファイルのURLを取得
		if zipURL != "" {
			entries = append(entries, Entry{
				AuthorID: token[1],
				Author: author,
				TitleId: token[2],
				Title: title,
				InfoURL: siteURL,
				ZipURL: zipURL,
			})
		}
	})
	return entries, nil
}

func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person1346.html"

	entries, err := findEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Println(entry.Title, entry.ZipURL)
	}
}