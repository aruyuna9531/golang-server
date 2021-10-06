//日志（问自己：和Log4cplus还有哪些区别）
#ifndef LOGGER_H_
#define LOGGER_H_

#include <iostream>
#include <fstream>
#include <sstream>
#include <queue>
#include "time_util.h"

enum e_logger_level
{
    LOGLV_TRACE = 1,
    LOGLV_DEBUG = 2,
    LOGLV_INFO = 3,
    LOGLV_WARN = 4,
    LOGLV_ERROR = 5,
    LOGLV_FATAL = 6,
    LOGLV_LOWEST = LOGLV_TRACE,
    LOGLV_HIGHEST = LOGLV_FATAL,
};



class Logger
{
public:
    static Logger& Instance()
    {
        static Logger inst;
        return inst;
    }

    void openStream(std::string& file_path);
    void closeStream();

    static std::string logPrefix(const e_logger_level log_level);
    bool writeLog(const std::string& log_detail);
private:
    void threadRun();
private:
    ~Logger()
    {
        closeStream();
    }

private:
    bool logger_is_open_;                           //是否开启
    int file_count_size_;                           //日志文件最大数量
    int64_t file_max_size_;                         //单个日志文件最大size(字节，第一次超过后不再写入)
    std::ofstream logger_file_path_;                //写入输出流
    std::queue<std::string> message_queue_;         //待写入消息列表（异步队列）
};

#define WRITE_LOG(lv, x) \
    std::ostringstream oss; \
    oss << "[" << getTime() << "]" << Logger::logPrefix(lv) << x << " - at " << __FILE__ << ":" << __LINE__ << std::endl; \
    Logger::Instance()->writeLog(oss.str());

#endif