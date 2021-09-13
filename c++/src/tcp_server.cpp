#include "tcp_server.h"

void tcp_loop(int port)
{
  // 主线程：tcp
  EventLoop loop;
  // 笔记：InetAddress
  // 参数1(port) 监听端口号
  // 参数2(loopback only) 可选，为true则只使用回环地址(127.0.0.1)，为false则使用0.0.0.0
  // 参数3(ipv6) 是否使用IPv6协议
  InetAddress listenAddr(port);
  MyTcpServer server(&loop, listenAddr);
  
  server.start();
  // 这里printf没法输出内容到terminal，cout就行，为啥
  std::cout << "starting tcp server on port " << port << ", now no tcp service available" << std::endl;
  loop.loop();
}

void MyTcpServer::onConnection(const TcpConnectionPtr& conn)
{
  LOG_TRACE << conn->peerAddress().toIpPort() << " -> "
            << conn->localAddress().toIpPort() << " is "
            << (conn->connected() ? "UP" : "DOWN");
}

void MyTcpServer::onMessage(const TcpConnectionPtr& conn, Buffer* buf, Timestamp time)
{
  string msg(buf->retrieveAllAsString());
  LOG_TRACE << conn->name() << " recv " << msg.size() << " bytes at " << time.toString();
  conn->send(msg);
}