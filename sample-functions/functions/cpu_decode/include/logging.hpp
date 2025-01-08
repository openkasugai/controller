/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef LOGGING_HPP
#define LOGGING_HPP

#include <string>

// Log level definition
enum class LogLevel : uint8_t {
	NOTHING = 0,
	PRINT   = 1,
	ERROR   = 2,
	WARN    = 3,
	INFO    = 4,
	DEBUG   = 5,
	ALL     = 6
};

namespace Logging {
	extern void setlevel(LogLevel lv);
	extern LogLevel getlevel();
	extern void set_stdoutmode(bool i);
	extern bool get_stdoutmode();
	extern int32_t set(LogLevel lv, const std::string &s);
	extern void close();
}

#endif //LOGGING_HPP
