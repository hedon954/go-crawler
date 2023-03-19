package tianyancha

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"

	"github.com/hedon954/go-crawler/fetcher"
)

const (
	OriginUrl = "https://www.tianyancha.com/"

	TaskNameTianYanCha = "tian_yan_cha"

	ruleNameIndustry      = "parse_industry_list"
	ruleNameCompanyList   = "parse_company_list"
	ruleNameCompanyDetail = "parse_company_detail"

	fieldIndustryName = "industry_name"
	fieldCompanyName  = "company_name"
	fieldCompanyId    = "company_id"
)

type Data struct {
	ID                uint   `json:"id" gorm:"primarykey;autoIncrement"`
	IndustryName      string `json:"industry_name" gorm:"column:industry_name"`
	CompanyName       string `json:"company_name" gorm:"column:company_name"`
	CompanyId         string `json:"company_id" gorm:"column:company_id"`
	CompanyType       string `json:"company_type" gorm:"column:company_type"`
	CreditCode        string `json:"credit_code" gorm:"column:credit_code"`
	Score             string `json:"score" gorm:"column:score"`
	KeyPeople         string `json:"key_people" gorm:"column:key_people"`
	Shareholder       string `json:"shareholder" gorm:"column:shareholder"`
	ForeignInvestment string `json:"foreign_investment" gorm:"column:foreign_investment"`
}

func (d Data) TableName() string {
	return "tian_yan_cha"
}

var TianYanChaTask = &fetcher.Task{
	Property: fetcher.Property{
		Name:     TaskNameTianYanCha,
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "xxx",
		Headers: map[string]string{
			"X-AUTH-TOKEN": "xxx",
			"X-TYCID":      "xxx",
		},
	},

	Rule: fetcher.RuleTree{
		Root: func() ([]*fetcher.Request, error) {
			roots := []*fetcher.Request{
				&fetcher.Request{
					Priority: 0,
					Url:      OriginUrl,
					Method:   "GET",
					RuleName: ruleNameIndustry,
				},
			}
			return roots, nil
		},
		Trunk: map[string]*fetcher.Rule{
			ruleNameIndustry:      &fetcher.Rule{ParseFunc: ParseHomeURL},
			ruleNameCompanyList:   &fetcher.Rule{ParseFunc: ParseCompanyList},
			ruleNameCompanyDetail: &fetcher.Rule{ParseFunc: ParseCompanyDetail},
		},
	},
}

// ParseHomeURL parses the homepage of the TianYanCha
func ParseHomeURL(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	// [1]: url
	// [2]: industry name
	regStr := `<a class="link-sub-hover-click index_item___BGg3 index_-right__1_Hlv" href="([^"]+)" rel="nofollow noreferrer" target="_blank">([^<]+)</a>`
	re := regexp.MustCompile(regStr)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := fetcher.ParseResult{}
	for _, m := range matches {
		// no vip can just get the first 5 pages
		for i := 1; i <= 5; i++ {
			req := &fetcher.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      string(m[1]) + "?pageNum=" + strconv.Itoa(i) + "&key=&sessionNo=1679119392.12280304",
				Depth:    ctx.Req.Depth + 1,
				RuleName: ruleNameCompanyList,
			}
			req.TempData = &fetcher.Temp{}
			_ = req.TempData.Set(fieldIndustryName, string(m[2]))
			result.Requests = append(result.Requests, req)
		}
	}
	return result, nil
}

// ParseCompanyList parses the company list
func ParseCompanyList(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	// [1]: url
	// [2]: company name
	regStr := `<a class="index_alink__zcia5 link-click" href="([^"]+)" target="_blank"><span>([^<]+)</span></a>`
	re := regexp.MustCompile(regStr)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := fetcher.ParseResult{}
	for _, m := range matches {
		req := &fetcher.Request{
			Priority: math.MaxInt32,
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: ruleNameCompanyDetail,
		}
		req.TempData = &fetcher.Temp{}
		_ = req.TempData.Set(fieldIndustryName, ctx.Req.TempData.Get(fieldIndustryName))
		_ = req.TempData.Set(fieldCompanyName, string(m[2]))
		_ = req.TempData.Set(fieldCompanyId, strings.TrimPrefix(string(m[1]), "https://www.tianyancha.com/company/"))
		result.Requests = append(result.Requests, req)
	}
	return result, nil
}

// ParseCompanyDetail parses the company detail info
func ParseCompanyDetail(ctx *fetcher.Context) (fetcher.ParseResult, error) {

	// [1] company type
	regStr := `<div class="index_company-tag__ZcJFV([^"]*)" style="color:#0084FF;background:#EBF5FF">([^<]+)`
	re := regexp.MustCompile(regStr)
	matches := re.FindAllSubmatch(ctx.Body, -1)

	isSmall := false
	isA := false

	for _, m := range matches {
		// [小微企业]
		if strings.Contains(string(m[2]), "小微企业") {
			isSmall = true
		}
		// [A 股]
		if strings.Contains(string(m[2]), "A股") {
			isA = true
		}
	}

	result := fetcher.ParseResult{
		Items: make([]interface{}, 0),
	}
	if isSmall {
		parseSmallCompanyDetail(ctx, &result)
	}
	if isA {
		parseACompanyDetail(ctx, &result)
	}
	return result, nil
}

var (
	regexCode        = regexp.MustCompile(`"creditCode":"([^"]+)"`)
	regexTax         = regexp.MustCompile(`"taxNumber":"([^"]+)"`)
	regexCompanyType = regexp.MustCompile(`"companyShowBizTypeName":"([^"]+)"`)
)

func parseSmallCompanyDetail(ctx *fetcher.Context, result *fetcher.ParseResult) {
	industryName := ctx.Req.TempData.Get(fieldIndustryName)
	companyName := ctx.Req.TempData.Get(fieldCompanyName)
	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	cData := Data{
		IndustryName: industryName.(string),
		CompanyName:  companyName.(string),
		CompanyId:    companyId.(string),
		CompanyType:  "小微企业",
		CreditCode:   extraString(ctx.Body, regexCode),
	}

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		cData.Score = getCompanyScore(ctx)
	}()
	go func() {
		defer wg.Done()
		cData.KeyPeople = getCompanyKeyPerson(ctx)
	}()
	go func() {
		defer wg.Done()
		cData.Shareholder = getSmallCompanyShareholder(ctx)
	}()
	go func() {
		defer wg.Done()
		cData.ForeignInvestment = getCompanyForeignInvestment(ctx)
	}()

	wg.Wait()
	data := ctx.OutputStruct(cData)
	result.Items = append(result.Items, data)
}

func parseACompanyDetail(ctx *fetcher.Context, result *fetcher.ParseResult) {
	industryName := ctx.Req.TempData.Get(fieldIndustryName)
	companyName := ctx.Req.TempData.Get(fieldCompanyName)
	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	cData := Data{
		IndustryName:      industryName.(string),
		CompanyName:       companyName.(string),
		CompanyId:         companyId.(string),
		CompanyType:       "A股上市",
		CreditCode:        extraString(ctx.Body, regexCode),
		Score:             getCompanyScore(ctx),
		KeyPeople:         "",
		Shareholder:       getACompanyShareholder(ctx),
		ForeignInvestment: getCompanyForeignInvestment(ctx),
	}
	if cData.CreditCode == "" {
		cData.CreditCode = extraString(ctx.Body, regexTax)
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		cData.Score = getCompanyScore(ctx)
	}()
	go func() {
		defer wg.Done()
		cData.Shareholder = getACompanyShareholder(ctx)
	}()
	go func() {
		defer wg.Done()
		cData.ForeignInvestment = getCompanyForeignInvestment(ctx)
	}()

	wg.Wait()

	data := ctx.OutputStruct(cData)
	result.Items = append(result.Items, data)
}

func extraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}

func getCompanyScore(ctx *fetcher.Context) string {
	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	if companyId == "" {
		return ""
	}

	var res struct {
		Data struct {
			BaseScore int `json:"baseScore"`
		} `json:"data"`
	}

	scoreUrlFormat := `https://capi.tianyancha.com/cloud-other-information/companyinfo/claim/score?_=%d&companyGid=%s`
	gorequest.New().Get(fmt.Sprintf(scoreUrlFormat, time.Now().UnixMilli(), companyId)).
		Timeout(5*time.Second).
		Retry(3, 1*time.Second).
		EndStruct(&res)

	return strconv.Itoa(res.Data.BaseScore)
}

func getCompanyKeyPerson(ctx *fetcher.Context) string {
	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	if companyId == "" {
		return ""
	}
	var res struct {
		Data struct {
			Result []struct {
				TypeSore           string `json:"typeSore"`                     // 职位
				Percent            string `json:"percent"`                      // 持股
				FinalBenefitShares string `json:"finalBenefitShares,omitempty"` // 最终受益股份
				Name               string `json:"name"`                         // 姓名
			} `json:"result"`
		} `json:"data"`
	}

	urlFormat := `https://capi.tianyancha.com/cloud-company-background/company/dim/staff?_=%d&gid=%s&pageSize=20&pageNum=1`
	_, _, _ = gorequest.New().Get(fmt.Sprintf(urlFormat, time.Now().UnixMilli(), companyId)).
		Timeout(5*time.Second).
		Retry(3, 1*time.Second).
		EndStruct(&res)
	bs, _ := json.Marshal(res.Data.Result)
	return string(bs)
}

func getSmallCompanyShareholder(ctx *fetcher.Context) string {

	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	if companyId == "" {
		return ""
	}

	urlFormat := `https://capi.tianyancha.com/cloud-company-background/companyV2/dim/holderForWeb?_=%d`

	var req = struct {
		PageSize     int    `json:"pageSize"`
		PageNum      int    `json:"pageNum"`
		Gid          string `json:"gid"`
		PercentLevel int    `json:"percentLevel"`
		SortField    string `json:"sortField"`
		SortType     int    `json:"sortType"`
	}{
		PageSize:     20,
		PageNum:      1,
		Gid:          companyId.(string),
		PercentLevel: -100,
		SortType:     -100,
		SortField:    "capitalAmount",
	}

	var res struct {
		Data struct {
			Result []struct {
				Name    string `json:"name"` // 股东名称
				Capital []struct {
					Amomon  string `json:"amomon"`  // 认缴出金额
					Time    string `json:"time"`    // 认缴出资日期
					Percent string `json:"percent"` // 持股比例
				} `json:"capital"`
			} `json:"result"`
		} `json:"data"`
	}

	gorequest.New().Post(fmt.Sprintf(urlFormat, time.Now().UnixMilli())).
		Timeout(5*time.Second).
		Retry(3, 1*time.Second).
		SendStruct(&req).
		EndStruct(&res)

	bs, _ := json.Marshal(res.Data.Result)
	return string(bs)
}

func getACompanyShareholder(ctx *fetcher.Context) string {
	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	if companyId == "" {
		return ""
	}

	urlFormat := `https://capi.tianyancha.com/cloud-listed-company/listed/holder/topTen?_=%d&gid=%s&pageSize=20&pageNum=1&percentLevel=-100&type=1`

	var res struct {
		Data struct {
			HolderList []struct {
				Proportion string `json:"proportion"` // 占总股比例
				ShareType  string `json:"shareType"`  // 股份类型
				Name       string `json:"name"`       // 股东名称
			} `json:"holderList"`
		} `json:"data"`
	}

	gorequest.New().Get(fmt.Sprintf(urlFormat, time.Now().UnixMilli(), companyId)).
		Timeout(5*time.Second).
		Retry(3, 1*time.Second).
		EndStruct(&res)

	bs, _ := json.Marshal(res.Data.HolderList)
	return string(bs)
}

func getCompanyForeignInvestment(ctx *fetcher.Context) string {

	companyId := ctx.Req.TempData.Get(fieldCompanyId)
	if companyId == "" {
		return ""
	}

	var res struct {
		Data struct {
			Category []struct {
				Name string `json:"name"`
				Num  int    `json:"num"`
			} `json:"category"`
			Area []struct {
				Name string `json:"name"`
				Num  int    `json:"num"`
			} `json:"area"`
		} `json:"data"`
	}

	urlFormat := `https://capi.tianyancha.com/cloud-company-background/company/invest/statistics?_=%d&gid=%s`
	gorequest.New().Get(fmt.Sprintf(urlFormat, time.Now().UnixMilli(), companyId)).
		Timeout(5*time.Second).
		Retry(3, 1*time.Second).
		EndStruct(&res)

	bs, _ := json.Marshal(res.Data)
	return string(bs)
}
