package goweb

// use radix implement the route
// 1.match http method
// 2.match static route
// 3.match batch prefix
// - user/info -user/login
// 4.match dynamic
//
// GET /home.html HTTP/1.1
//Host: developer.mozilla.org
//User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:50.0) Gecko/20100101 Firefox/50.0
//Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
//Accept-Language: en-US,en;q=0.5
//Accept-Encoding: gzip, deflate, br
//Referer: https://developer.mozilla.org/testpage.html
// GET    /home.html  HTTP/1.1
// Method Request-URI HTTP-Version
// route design:how to the framework user use the route
