package tianyancha

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hedon954/go-crawler/fetcher"
)

var (
	TaskNameTianYanCha2 = "tian_yan_cha_2"

	fieldShareholder = "shareholder"
	fieldScore       = "score"
	fieldCreditCode  = "credit_code"
	fieldCompanyType = "company_type"

	baseUrlFormat   = "https://www.tianyancha.com/company/%d"
	shareUrlFormat  = `https://capi.tianyancha.com/cloud-company-background/company/getShareHolderStructure?gid=%s&rootGid=%s&operateType=1`
	investUrlFormat = `https://capi.tianyancha.com/cloud-company-background/company/getShareHolderStructure?gid=%s&rootGid=%s&operateType=2`

	ruleNameParseBase   = "parse_base_info"
	ruleNameParseShare  = "parse_shareholder"
	ruleNameParseInvest = "parse_investment"
)

type ShareOrInvest struct {
	Data struct {
		Name     string               `json:"name"`
		Children []ShareOrInvestChild `json:"children"`
	} `json:"data"`
}

type ShareOrInvestChild struct {
	Gid     string `json:"gid"`
	Name    string `json:"name"`
	Amount  string `json:"amount"`
	Percent string `json:"percent"`
}

type TYCData struct {
	ID          uint   `json:"id" gorm:"primarykey;autoIncrement"`
	CompanyName string `json:"company_name" gorm:"column:company_name"`
	CompanyId   string `json:"company_id" gorm:"column:company_id"`
	CompanyType string `json:"company_type" gorm:"column:company_type"`
	CreditCode  string `json:"credit_code" gorm:"column:credit_code"`
	Score       string `json:"score" gorm:"column:score"`
	Shareholder string `json:"shareholder" gorm:"column:shareholder"`
	Investment  string `json:"investment" gorm:"column:investment"`
}

func (sh TYCData) TableName() string {
	return "tian_yan_cha_2"
}

var TianYanCha2Task = &fetcher.Task{
	Property: fetcher.Property{
		Name:     TaskNameTianYanCha2,
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "jsid=SEO-BAIDU-ALL-SY-000001; TYCID=1af08280c53d11ed81b4bdbf95e432f8; ssuid=8176940507; _ga=GA1.2.292669649.1679110386; _gid=GA1.2.1479135818.1679110386; RTYCID=79f96e3cf3344cac957774235738d5c9; bannerFlag=true; HWWAFSESID=53f4f2084618297f90c; HWWAFSESTIME=1679209205819; csrfToken=HejZGLuJo14EAR9Ld69Md9x3; bdHomeCount=2; Hm_lvt_e92c8d65d92d534b0fc290df538b4758=1679195712,1679198477,1679205683,1679209207; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22284632286%22%2C%22first_id%22%3A%22186f2c3fbd1926-076f09e7232fd44-1f525634-1296000-186f2c3fbd2b8e%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTg2ZjJjM2ZiZDE5MjYtMDc2ZjA5ZTcyMzJmZDQ0LTFmNTI1NjM0LTEyOTYwMDAtMTg2ZjJjM2ZiZDJiOGUiLCIkaWRlbnRpdHlfbG9naW5faWQiOiIyODQ2MzIyODYifQ%3D%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%24identity_login_id%22%2C%22value%22%3A%22284632286%22%7D%2C%22%24device_id%22%3A%22186f2c3fbd1926-076f09e7232fd44-1f525634-1296000-186f2c3fbd2b8e%22%7D; searchSessionId=1679210092.14816511; cid=8660845; ss_cidf=1; tyc-user-info={%22state%22:%225%22%2C%22vipManager%22:%220%22%2C%22mobile%22:%2217144837089%22%2C%22isExpired%22:%220%22}; tyc-user-info-save-time=1679210971673; auth_token=eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxNzE0NDgzNzA4OSIsImlhdCI6MTY3OTIxMDk3MSwiZXhwIjoxNjgxODAyOTcxfQ.L7toTBppSeV4FvN_tTZu2RmXN9aV2WUtras0nvw-6V5blip-k9c-Q-6TjnysrD5qbjVLGgohjjlSEQbVLeBHrA; tyc-user-phone=%255B%252217144837089%2522%252C%2522156%25202320%25205156%2522%255D; cloud_token=e31b86ad17644638b09ed03aab9fda34; Hm_lpvt_e92c8d65d92d534b0fc290df538b4758=1679212380",
		Headers: map[string]string{
			"X-AUTH-TOKEN": "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxNzE0NDgzNzA4OSIsImlhdCI6MTY3OTIxMDk3MSwiZXhwIjoxNjgxODAyOTcxfQ.L7toTBppSeV4FvN_tTZu2RmXN9aV2WUtras0nvw-6V5blip-k9c-Q-6TjnysrD5qbjVLGgohjjlSEQbVLeBHrA",
			"X-TYCID":      "1af08280c53d11ed81b4bdbf95e432f8",
			"version":      "TYC-Web",
		},
	},

	Rule: fetcher.RuleTree{
		Root: func() ([]*fetcher.Request, error) {
			start := 1062000
			end := 5000000
			roots := make([]*fetcher.Request, end-start)
			for i := start; i < end; i++ {
				req := &fetcher.Request{
					Priority: 0,
					Url:      fmt.Sprintf(baseUrlFormat, i),
					RuleName: ruleNameParseBase,
					TempData: &fetcher.Temp{},
				}
				_ = req.TempData.Set(fieldCompanyId, strconv.Itoa(i))
				roots[i-start] = req
			}
			return roots, nil
		},
		Trunk: map[string]*fetcher.Rule{
			ruleNameParseBase:   &fetcher.Rule{ParseFunc: ParseBaseInfo},
			ruleNameParseShare:  &fetcher.Rule{ParseFunc: ParseShareInfo},
			ruleNameParseInvest: &fetcher.Rule{ParseFunc: ParseInvestInfo},
		},
	},
}

func ParseBaseInfo(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	creditCode := extraString(ctx.Body, regexCode)
	if creditCode == "" {
		creditCode = extraString(ctx.Body, regexTax)
	}
	score := getCompanyScore(ctx)
	companyType := extraString(ctx.Body, regexCompanyType)

	companyId := ctx.Req.TempData.Get(fieldCompanyId).(string)
	result := fetcher.ParseResult{}
	req := &fetcher.Request{
		Priority: math.MaxInt32,
		Method:   "GET",
		Task:     ctx.Req.Task,
		Url:      fmt.Sprintf(shareUrlFormat, companyId, companyId),
		Depth:    ctx.Req.Depth + 1,
		RuleName: ruleNameParseShare,
	}
	req.TempData = ctx.Req.TempData.Copy()
	_ = req.TempData.Set(fieldScore, score)
	_ = req.TempData.Set(fieldCreditCode, creditCode)
	_ = req.TempData.Set(fieldCompanyType, companyType)
	result.Requests = append(result.Requests, req)
	return result, nil
}

func ParseShareInfo(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	var res ShareOrInvest
	_ = json.Unmarshal(ctx.Body, &res)
	if res.Data.Name == "" {
		return fetcher.ParseResult{}, nil
	}
	bs, _ := json.Marshal(res.Data.Children)

	companyId := ctx.Req.TempData.Get(fieldCompanyId).(string)
	result := fetcher.ParseResult{}
	req := &fetcher.Request{
		Priority: math.MaxInt32,
		Method:   "GET",
		Task:     ctx.Req.Task,
		Url:      fmt.Sprintf(investUrlFormat, companyId, companyId),
		Depth:    ctx.Req.Depth + 1,
		RuleName: ruleNameParseInvest,
	}
	req.TempData = ctx.Req.TempData.Copy()
	_ = req.TempData.Set(fieldCompanyName, res.Data.Name)
	_ = req.TempData.Set(fieldShareholder, string(bs))
	result.Requests = append(result.Requests, req)
	return result, nil
}

func ParseInvestInfo(ctx *fetcher.Context) (fetcher.ParseResult, error) {
	result := fetcher.ParseResult{}
	var res ShareOrInvest
	_ = json.Unmarshal(ctx.Body, &res)
	if res.Data.Name == "" {
		return result, nil
	}
	bs, _ := json.Marshal(res.Data.Children)
	companyId := ctx.Req.TempData.Get(fieldCompanyId).(string)
	companyName := ctx.Req.TempData.Get(fieldCompanyName).(string)
	companyType := ctx.Req.TempData.Get(fieldCompanyType).(string)
	score := ctx.Req.TempData.Get(fieldScore).(string)
	creditCode := ctx.Req.TempData.Get(fieldCreditCode).(string)
	sh := ctx.Req.TempData.Get(fieldShareholder).(string)
	cData := ctx.OutputStruct(TYCData{
		CompanyId:   companyId,
		CompanyName: companyName,
		CompanyType: companyType,
		Score:       score,
		CreditCode:  creditCode,
		Shareholder: sh,
		Investment:  string(bs),
	})
	result.Items = append(result.Items, cData)
	return result, nil
}
