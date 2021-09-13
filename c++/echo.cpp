#include <muduo/net/TcpServer.h>

#include <functional>
#include <utility>

#include <cstdio>
#include <iostream>
#include <unistd.h>

#include "echo_server.h"

using namespace muduo;
using namespace muduo::net;

int kRollSize = 500*1000*1000;

std::unique_ptr<muduo::AsyncLogging> g_asyncLog;

void asyncOutput(const char* msg, int len)
{
  g_asyncLog->append(msg, len);
}

void setLogging(const char* argv0)
{
  muduo::Logger::setOutput(asyncOutput);
  char name[256];
  strncpy(name, argv0, 256);
  g_asyncLog.reset(new muduo::AsyncLogging(::basename(name), kRollSize));
  g_asyncLog->start();
}

static void onRequest(const HttpRequest& req, HttpResponse* resp)
{
  std::string query = req.query();
  std::string body = req.body();
  std::string path = req.path();

  cout << "path = " << path << ", query = " << query << ", body = " << body;
}

void http_loop()
{
  EventLoop loop;
  InetAddress http_listenAddr(static_cast<short>(8001));
  HttpServer http_server(&loop, http_listenAddr, "http");
  http_server.setHttpCallback(onRequest);
  http_server.setThreadNum(1);
  loop.loop();
}

int main(int argc, char* argv[])
{
  setLogging(argv[0]);

  LOG_INFO << "pid = " << getpid() << ", tid = " << CurrentThread::tid();

  // 子线程：http
  muduo::Thread http_thread((const muduo::Thread::ThreadFunc)boost::bind(&http_loop));
  http_thread.start();

  // 主线程：tcp
  EventLoop loop;
  // 笔记：InetAddress
  // 参数1(port) 监听端口号
  // 参数2(loopback only) 可选，为true则只使用回环地址(127.0.0.1)，为false则使用0.0.0.0
  // 参数3(ipv6) 是否使用IPv6协议
  InetAddress listenAddr(8000, false);
  EchoServer server(&loop, listenAddr);
  
  server.start();
  // 这里printf没法输出内容到terminal，cout就行，为啥
  std::cout << "starting server on port 8000" << std::endl;

  loop.loop();
  // loop()之后的代码不会被执行，因此如果服务器要兼顾多项服务的话，子线程在前面就要创建好。因此return 0在这里也没有意义。
}

