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
  LOG4CPLUS_INFO(log, "TCP服务已在" << port << "启动");
  loop.loop();
}

void MyTcpServer::onConnection(const TcpConnectionPtr& conn)
{
  bool connected = conn->connected();
  std::string remote_addr = conn->peerAddress().toIpPort();
  LOG_TRACE << remote_addr << " -> "
            << conn->localAddress().toIpPort() << " is "
            << (connected ? "UP" : "DOWN");
  
  if (connected)
  {
    LOG4CPLUS_INFO(log, "TCP接收到来自" << remote_addr << "的远程连接");
  }
  else
  {
    LOG4CPLUS_INFO(log, "TCP已从" << remote_addr << "断开");
    //TODO 断开之后干什么
  }
}

void MyTcpServer::onMessage(const TcpConnectionPtr& conn, Buffer* buf, Timestamp time)
{
  string msg(buf->retrieveAllAsString());
  // LOG_TRACE << conn->name() << " recv " << msg.size() << " bytes at " << time.toString();
  LOG4CPLUS_TRACE(log, "接收到来自" << conn->peerAddress().toIpPort() << "的信息：[" << msg << "]");
  conn->send("this is a return from tcp server(c++)");
}