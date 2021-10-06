#ifndef PROTOCOL_H_
#define PROTOCOL_H_

#include <unordered_map>
#include <unordered_set>

class Protocol
{
public:
    bool putBool(const int key, const bool value);
    bool putInt(const int key, const int value);
    bool putInt64(const int key, const int64_t value);
    bool putString(const int key, const std::string& value);
    bool putProtocol(const int key, const Protocol& value);
private:
    std::unordered_set<int> existing_key_;
    std::unordered_map<int, bool> type_bool_data;
    std::unordered_map<int, int> type_int_data_;
    std::unordered_map<int, int64_t> type_int64_data_;
    std::unordered_map<int, std::string> type_string_data_;
    std::unordered_map<int, Protocol> type_protocol_data_;          //？？？
};

#endif