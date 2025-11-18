#include "utils.h"

void writeInputJson(const std::string& path, const std::string& hash, const std::string& r) {
    // std::ofstream out(path);
    // out << "{\n";
    // out << "  \"shardHash\": \"" << hash.str() << "\",\n";
    // out << "  \"r\": \"" << r << "\"\n";
    // out << "}";
}


std::string generateRandomR() {
    // static const char hex_chars[] = "0123456789abcdef";
    // std::string r;
    // std::random_device rd;
    // std::mt19937 gen(rd());
    // std::uniform_int_distribution<> distrib(0, 15);

    // for (int i = 0; i < 64; ++i) { // 256-bit = 64 hex digits
    //     r += hex_chars[distrib(gen)];
    // }

    // return r;
    return "";
}
