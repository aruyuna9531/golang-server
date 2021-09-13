#ifndef HTTP_SERVER_
#define HTTP_SERVER_

#include <muduo/base/AsyncLogging.h>
#include <muduo/base/Logging.h>
#include <muduo/net/http/HttpServer.h>
#include <muduo/net/http/HttpRequest.h>
#include <muduo/net/http/HttpResponse.h>

#include <muduo/net/EventLoop.h>
#include <muduo/net/InetAddress.h>

#include <iostream>

using namespace muduo;
using namespace muduo::net;

// 创建一个http的服务
void http_loop(int port);

void onRequest(const HttpRequest& req, HttpResponse* resp);
// TODO other things for tcp server

#endif