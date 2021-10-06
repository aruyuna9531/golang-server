#ifndef TIME_UTIL_H_
#define TIME_UTIL_H_

#include <string>
#include <sstream>
#include <time.h>
#include "common.h"

std::string getTime()
{
    struct timeval time;
 
    gettimeofday(&time, NULL);
    time.tv_sec time.tv_usec

	struct tm tm;
	char s[100] = {0};
 
    time_t stampTime = time(NULL);
	tm = *localtime(&stampTime);
	strftime(s, sizeof(s), "%Y-%m-%d %H:%M:%S", &tm);

    std::ostringstream oss;
    oss << std::string(s) << "." << time.tv_usec;
    return oss.str();
}

template<typename T>
std::string toString(const T& value)
{
    std::ostringstream oss;
    oss << value;
    return oss.str();
}
#endif