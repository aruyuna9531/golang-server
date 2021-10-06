#include "http_server.h"

void http_loop(int port)
{
  EventLoop loop;
  InetAddress http_listenAddr(port);
  HttpServer http_server(&loop, http_listenAddr, "http");
  http_server.setHttpCallback(onRequest);
  http_server.setThreadNum(1);
  http_server.start();
  LOG4CPLUS_INFO(log, "HTTP服务器在" << port << "开启，可以尝试请求\"http://127.0.0.1:" << port << "/hello\"");
  loop.loop();
}

void onRequest(const HttpRequest& req, HttpResponse* resp)
{
  LOG4CPLUS_INFO(log, "接收到HTTP请求");

  std::string query = req.query();
  // std::string body = req.body();
  std::string path = req.path();

  LOG4CPLUS_INFO(log, "path = " << path << ", query = " << query);

  //回写响应
  resp->setStatusCode(HttpResponse::k200Ok);
  resp->setStatusMessage("OK");
  resp->setContentType("text/html");
  resp->addHeader("Content-Type", "text/html");
  resp->addHeader("Connection", "close");
  resp->setBody("this is a return from http server(c++)");
}