package douban

import (
	"time"

	"github.com/hedon954/go-crawler/fetcher"
)

const (
	TaskNameFindSunRoomJs = "js_find_douban_sun_room"
)

var DoubanJsTask = &fetcher.TaskModel{
	Property: fetcher.Property{
		Name:     TaskNameFindSunRoomJs,
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "xxx",
	},
	Root: `
		var arr = new Array();
		for (var i = 25; i <= 25; i+=25) {
		    var obj = {
		        Url: "https://www.douban.com/group/szsh/discussion?start=" + i,
		        Priority: 1,
		        RuleName: "解析网站URL",
		        Method: "GET",
		    };
		    arr.push(obj);
		};
		console.log(arr[0].Url);
		AddJsReq(arr);
	`,
	Rules: []fetcher.RuleModel{
		{
			Name: "解析网站URL",
			ParseFunc: `
				ctx.ParseJSReg("解析阳台房","(https://www.douban.com/group/topic/[0-9a-z]+/)\"[^>]*>([^<]+)</a>");
			`,
		},
		{
			Name: "解析阳台房",
			ParseFunc: `
			ctx.OutputJS("<div class=\"topic-content\">[\\s\\S]*?阳台[\\s\\S]*?<div class=\"aside\">");
			`,
		},
	},
}
