const http = require('http');

http.createServer((request, response)=>{
    var args = new URL(request.url, "http://0.0.0.0:8000");
    // console.log(args);
    response.write('hello! you are trying to access path: ' + args.pathname + ", your params is: " + args.searchParams);
}).listen(8000, '0.0.0.0', ()=>{
    console.log("server started at 0.0.0.0:8000");
})