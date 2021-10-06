#ifndef TCP_SERVER_
#define TCP_SERVER_

#include <muduo/base/AsyncLogging.h>
#include <muduo/base/Logging.h>
#include <muduo/net/TcpServer.h>

#include <muduo/net/EventLoop.h>
#include <muduo/net/InetAddress.h>

#include <iostream>
#include "common.h"

using namespace muduo;
using namespace muduo::net;
class MyTcpServer
{
 public:
  MyTcpServer(EventLoop* loop, const InetAddress& listenAddr)
    : loop_(loop),
      server_(loop, listenAddr, "MyTcpServer")
  {
    server_.setConnectionCallback(
        std::bind(&MyTcpServer::onConnection, this, _1));
    server_.setMessageCallback(
        std::bind(&MyTcpServer::onMessage, this, _1, _2, _3));
  }

  void start()
  {
    server_.start();
  }

 private:
  void onConnection(const TcpConnectionPtr& conn);

  void onMessage(const TcpConnectionPtr& conn, Buffer* buf, Timestamp time);

  EventLoop* loop_;
  TcpServer server_;
};

void tcp_loop(int port);

#endif