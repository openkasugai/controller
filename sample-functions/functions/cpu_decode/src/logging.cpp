/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <chrono>
#include <ctime>
#include <iomanip>
#include "logging.hpp"

static LogLevel level = LogLevel::INFO;
static bool is_stdout = false;
static std::ofstream ofs;

namespace Logging {
	// Log start time prefix
	static std::ostringstream time_prefix() {
		const std::chrono::system_clock::time_point now = std::chrono::system_clock::now();
		const std::time_t t = std::chrono::system_clock::to_time_t(now);
		std::tm *lt = std::localtime(&t);
		std::ostringstream oss;
		oss << lt->tm_year + 1900;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_mon + 1;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_mday;
		oss << "-";
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_hour;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_min;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_sec;
		return oss;
	}

	// Log File Start
	static int32_t open() {
		std::ostringstream tp = time_prefix();
		std::string logfile = "app_" + tp.str() + ".log";
		ofs.open(logfile);
		if (!ofs) {
			return -1;
		}
		ofs << "log start..." << tp.str() << std::endl;
		ofs << "loglevel: " << std::to_string(static_cast<uint32_t>(level)) << std::endl;
		if (is_stdout) {
			std::cout << "log start..." << tp.str() << std::endl;
			std::cout << "loglevel: " << std::to_string(static_cast<uint32_t>(level)) << std::endl;
		}

		return 0;
	}

	// Log Level Settings
	void setlevel(LogLevel lv) {
		level = lv;
	}

	// Log Level Capture
	LogLevel getlevel() {
		return level;
	}

	// Standard output setting
	void set_stdoutmode(bool i) {
		is_stdout = i;
	}

	// standard output setting acquisition
	bool get_stdoutmode() {
		return is_stdout;
	}

	// Log Set
	int32_t set(LogLevel lv, const std::string &s) {
		if (lv > level) {
			return 0;
		}

		if (!ofs.is_open()) {
			if (open() < 0) {
				std::cerr << "Logging.set(): logfile open failed." << std::endl;
				return -1;
			}
		}
		std::ostringstream oss;

		const std::chrono::system_clock::time_point now = std::chrono::system_clock::now();
		const std::time_t t = std::chrono::system_clock::to_time_t(now);
		std::tm *lt = std::localtime(&t);
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_hour << ":";
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_min << ":";
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_sec << " ";

		switch (lv) {
			case LogLevel::PRINT:
				break;
			case LogLevel::ERROR:
				oss << "[error] ";
				break;
			case LogLevel::WARN:
				oss << "[warn ] ";
				break;
			case LogLevel::INFO:
				oss << "[info ] ";
				break;
			case LogLevel::DEBUG:
				oss << "[debug] ";
				break;
			default:
				oss << "[?????] ";
				break;
		}

		ofs << oss.str() << s << std::endl;
		if (is_stdout) {
			std::cout << oss.str() << s << std::endl;
		}
		oss.str("");

		return 0;
	}

	// End Log File
	void close() {
		if (ofs.is_open()) {
			ofs.close();
		}
	}
}
