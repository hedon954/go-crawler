package douban

import (
	"github.com/hedon954/go-crawler/fetcher"
	"regexp"
	"strconv"
	"time"
)

const (
	bookUrl = "https://book.douban.com"

	fieldBookName = "book_name"

	regexpStrBookList  = `<a href="([^"]+)" class="tag">([^<]+)</a>`
	regexpStrBookIntro = `<a.*?href="([^"]+)" title="([^"]+)"`

	ruleNameBookTag    = "parse_book_tag"
	ruleNameBookIntro  = "parse_book_intro"
	ruleNameBookDetail = "parse_book_detail"

	TaskNameDoubanBook = "douban_book_list"
)

var DoubanBookTask = &fetcher.Task{
	Name:     TaskNameDoubanBook,
	WaitTime: 1 * time.Second,
	MaxDepth: 5,
	Cookie:   "Xxx",

	Rule: fetcher.RuleTree{
		Root: func() []*fetcher.Request {
			roots := []*fetcher.Request{
				&fetcher.Request{
					Priority: 1,
					Url:      bookUrl,
					Method:   "GET",
					RuleName: ruleNameBookTag,
				},
			}
			return roots
		},
		Trunk: map[string]*fetcher.Rule{
			ruleNameBookTag:   &fetcher.Rule{ParseFunc: ParseTag},
			ruleNameBookIntro: &fetcher.Rule{ParseFunc: ParseBookList},
			ruleNameBookDetail: &fetcher.Rule{
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

// ParseTag parses tags to get book list
func ParseTag(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	re := regexp.MustCompile(regexpStrBookList)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := fetcher.ParseResult{}

	for _, m := range matches {
		result.Requests = append(result.Requests, &fetcher.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      bookUrl + string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: ruleNameBookTag,
		})
	}
	return result, nil
}

// ParseBookList parse books list to get books' intro
func ParseBookList(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	re := regexp.MustCompile(regexpStrBookIntro)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := fetcher.ParseResult{}

	for _, m := range matches {
		req := &fetcher.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: ruleNameBookIntro,
		}
		req.TempData = &fetcher.Temp{}
		_ = req.TempData.Set(fieldBookName, string(m[2]))
		result.Requests = append(result.Requests, req)
	}

	return result, nil
}

// ParseBookDetail parses book's detail info
func ParseBookDetail(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	bookName := ctx.Req.TempData.Get(fieldBookName)
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))
	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  ExtraString(ctx.Body, autoRe),
		"页数":  page,
		"出版社": ExtraString(ctx.Body, public),
		"得分":  ExtraString(ctx.Body, scoreRe),
		"价格":  ExtraString(ctx.Body, priceRe),
		"简介":  ExtraString(ctx.Body, intoRe),
	}
	data := ctx.Output(book)
	result := fetcher.ParseResult{
		Items: []interface{}{data},
	}
	return result, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) > 2 {
		return string(match[1])
	}
	return ""
}
