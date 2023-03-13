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