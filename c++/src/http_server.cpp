#include "http_server.h"

void http_loop(int port)
{
  EventLoop loop;
  InetAddress http_listenAddr(port);
  HttpServer http_server(&loop, http_listenAddr, "http");
  http_server.setHttpCallback(onRequest);
  http_server.setThreadNum(1);
  http_server.start();
  std::cout << "listening http in port " << port << ", you can try an http request as \"http://127.0.0.1:" << port << "/hello\"" << std::endl;
  loop.loop();
}

void onRequest(const HttpRequest& req, HttpResponse* resp)
{
  std::cout << "http request accepted" << std::endl;

  std::string query = req.query();
  // std::string body = req.body();
  std::string path = req.path();

  std::cout << "path = " << path << ", query = " << query << std::endl;
}