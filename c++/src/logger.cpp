#include "logger.h"

std::string Logger::logPrefix(const e_logger_level log_level)
{
    switch (log_level)
    {
        case LOGLV_TRACE:
            return "[TRACE]";
        case LOGLV_DEBUG:
            return "[DEBUG]";
        case LOGLV_INFO:
            return "[INFO]";
        case LOGLV_WARN:
            return "[WARN]";
        case LOGLV_ERROR:
            return "[ERROR]";
        case LOGLV_FATAL:
            return "[FATAL]";
        default:
            return "";
    }
}

bool Logger::writeLog(const std::string& log_detail)
{
    message_queue_.push(log_detail);
}

void Logger::threadRun()
{

}

void Logger::openStream(std::string& file_path)
{
    if (logger_file_path_.is_open())
    {
        std::cerr << "Logger::openStream error: logger is opening" << std::endl;
        return;
    }

    logger_file_path_.open(file_path, std::ofstream::app);
    if (!logger_file_path_)
    {
        std::cerr << "Logger::openStream error: open " << file_path << " failed." << std::endl;
        return;
    }

    logger_is_open_ = true;
}

void Logger::closeStream()
{
    if (!logger_file_path_.is_open())
    {
        std::cout << "Logger::closeStream error: logger is not opening" << std::endl;
        return;
    }

    logger_file_path_.close();
    logger_is_open_ = false;
}