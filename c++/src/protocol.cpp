#include "protocol.h"

bool Protocol::putBool(const int key, const bool value)
{

}

bool Protocol::putInt(const int key, const int value);
bool Protocol::putInt64(const int key, const int64_t value);
bool Protocol::putString(const int key, const std::string& value);
bool Protocol::putProtocol(const int key, const Protocol& value);