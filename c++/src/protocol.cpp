#include "protocol.h"

bool putBool(const int key, const bool value);
bool putInt(const int key, const int value);
bool putInt64(const int key, const int64_t value);
bool putString(const int key, const std::string& value);
bool putProtocol(const int key, const Protocol& value);