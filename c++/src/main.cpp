#include <muduo/base/AsyncLogging.h>
#include <muduo/base/Logging.h>
#include <muduo/base/Thread.h>

#include <functional>
#include <utility>

#include <cstdio>
#include <iostream>
#include <unistd.h>

#include "tcp_server.h"
#include "http_server.h"
#include "common.h"

using namespace boost;
using namespace muduo;
using namespace muduo::net;

int kRollSize = 500*1000*1000;

std::unique_ptr<AsyncLogging> g_asyncLog;

void asyncOutput(const char* msg, int len)
{
  g_asyncLog->append(msg, len);
}

void setLogging(const char* argv0)
{
  Logger::setOutput(asyncOutput);
  char name[256];
  strncpy(name, argv0, 256);
  g_asyncLog.reset(new muduo::AsyncLogging(::basename(name), kRollSize));
  g_asyncLog->start();
}

int main(int argc, char* argv[])
{
  if (argc < 2)
  {
    std::cerr << "usage: " << argv[0] << " [log4cplus_conf]" << std::endl;
    return -1;
  }

  setLogging(argv[0]);

  // log4cplus init start
  log4cplus::initialize();

  log4cplus::BasicConfigurator config;
  config.configure();

  log4cplus::PropertyConfigurator::doConfigure(LOG4CPLUS_TEXT(argv[1]));
  // log4cplus init end

  LOG_INFO << "pid = " << getpid() << ", tid = " << CurrentThread::tid();

  // 子线程：http
  muduo::Thread http_thread((const muduo::Thread::ThreadFunc)std::bind(&http_loop, 8002));
  http_thread.start();

  // 主线程：tcp
  tcp_loop(8003);
  // loop()之后的代码不会被执行，因此如果服务器要兼顾多项服务的话，子线程在前面就要创建好。因此return 0在这里也没有意义。
  // main必须要挂起，否则结束了程序就没了，这里用tcp_loop()内的loop()使进程挂起。
  // 注意：http_thread线程创建的loop()不会使得进程挂起。因此tcp_loop不能挂在Thread内。（否则可能连创建线程的cout都没有）
  // 如果强迫症要代码整齐，必须使用其他方式在最后使进程挂起。
}

