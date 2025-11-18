#pragma once

#include <condition_variable>
#include <functional>
#include <iostream>
#include <mutex>
#include <queue>
#include <string>
#include <thread>

struct FileStruct {
    std::string user;
    std::string raw_file_id;
    std::string pk;
    std::string file_id;
    std::string file_name;
    std::string full_fule_path;
};

class TaskQueue {
public:
    static TaskQueue& Instance() {
        static TaskQueue task_queue;
        return task_queue;
    }
    ~TaskQueue() {}

    void Push(const FileStruct& file_info) {
        {
            std::lock_guard<std::mutex> lock(mutex_);
            tasks_.push(file_info);
        }
        cv_.notify_one();
    }

    void RunWorker(std::function<void(const std::string&, const std::string&, const std::string&, const std::string&,
            const std::string&, const std::string&)>
            callback) {
        std::thread([this, callback]() {
            while (true) {
                FileStruct file_info;
                {
                    std::unique_lock<std::mutex> lock(mutex_);
                    cv_.wait(lock, [this]() { return !tasks_.empty() || stop_; });

                    if (stop_ && tasks_.empty()) return;

                    file_info = tasks_.front();
                    tasks_.pop();
                }
                callback(file_info.user, file_info.raw_file_id, file_info.pk, file_info.file_id, file_info.file_name,
                    file_info.full_fule_path);
            }
        }).detach();
    }

    void Shutdown() {
        {
            std::lock_guard<std::mutex> lock(mutex_);
            stop_ = true;
        }
        cv_.notify_all();
    }

private:
    TaskQueue() {}
    std::queue<FileStruct> tasks_;
    std::mutex mutex_;
    std::condition_variable cv_;
    bool stop_ = false;
};
