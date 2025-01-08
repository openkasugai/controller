/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef CONNECT_THREAD_QUEUE_HPP
#define CONNECT_THREAD_QUEUE_HPP

#include <queue>
#include <mutex>
#include <condition_variable>

// Queue for passing data between asynchronous threads
template <typename T>
class ConnectThreadQueue {
	private:
		std::queue<T> que_;
		mutable std::mutex mtx_;
	 	mutable	std::condition_variable  empty_wait_;

	public:
		ConnectThreadQueue(){}

		~ConnectThreadQueue(){}

		// blocking push: adding data to queue
		void push(const T& in) {
			std::unique_lock<std::mutex> lock(mtx_);
			const bool is_empty = que_.empty();
			que_.push(in);
			if (is_empty) {
				empty_wait_.notify_all();
			}
		}

		// blocking pop: retrieving data from queue
		T pop() {
			std::unique_lock<std::mutex> lock(mtx_);
			empty_wait_.wait(lock, [&]{return !que_.empty();});
			T out = que_.front();
			que_.pop();
			return out;
		}

		// Get Queue Size
		size_t size() {
			std::unique_lock<std::mutex> lock(mtx_);
			return que_.size();
		}
};

#endif //CONNECT_THREAD_QUEUE_HPP
