cmake_minimum_required(VERSION 3.14)
project(zkagent VERSION "1.0.0")

# 版本后缀，可根据需求附加在版本中，例如："-rc1"
set(VERSION_SUFFIX "")
# 安装项目时跳过所有的依赖项安装
set(CMAKE_SKIP_INSTALL_ALL_DEPENDENCY TRUE)
# 设置可执行程序输出目录
set(EXECUTABLE_OUTPUT_PATH ${PROJECT_BINARY_DIR}/bin)
file(MAKE_DIRECTORY ${EXECUTABLE_OUTPUT_PATH})
# 添加 cmake 模块搜索目录
set(DATA_TRANS_CMAKE_DIR ${CMAKE_SOURCE_DIR}/cmake)
list(APPEND CMAKE_MODULE_PATH ${DATA_TRANS_CMAKE_DIR})
# 设置运行时路径
set(CMAKE_BUILD_RPATH ".")

add_definitions(-DOPENSSL_API_1_1)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# 限制 CPU 指令集，避免不兼容问题
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -march=x86-64 -mtune=generic -O2")

# 静态链接 libstdc++ 和 libgcc，避免 `glibc` 兼容性问题
set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -static-libstdc++ -static-libgcc")

# 定义编译选项（默认 OFF）
option(USE_BCOS "Build with bcos sdk" OFF)
option(USE_LSSDK "Build with lssdk" OFF)

# 防止两个同时打开
if (USE_BCOS AND USE_LSSDK)
    message(FATAL_ERROR "You can only enable one of USE_BCOS or USE_LSSDK!")
endif()


# 设置 Go 构建的输出目录
set(GO_OUTPUT_DIR ${CMAKE_SOURCE_DIR}/deps/go)

if(NOT EXISTS ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/go.mod)
    add_custom_command(
        OUTPUT ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/go.mod
        COMMAND go mod init chainmaker
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/
        COMMENT "Initializing go.mod"
    )

    # 定义一个虚拟目标，依赖 go.mod，这样 make golib 会触发生成
    add_custom_target(init_chainmaker_go_mod ALL
        DEPENDS ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/go.mod
    )
endif()


# 1. 定义生成 go 动态库的命令
add_custom_command(
    OUTPUT ${GO_OUTPUT_DIR}/libchainmaker_sdk.so ${GO_OUTPUT_DIR}/libchainmaker_sdk.h
    COMMAND ${CMAKE_COMMAND} -E make_directory ${GO_OUTPUT_DIR}
    COMMAND go mod tidy
    COMMAND go build -o ${GO_OUTPUT_DIR}/libchainmaker_sdk.so -buildmode=c-shared ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/main.go
    # 复制 so 和 h 到指定目录
    COMMAND ${CMAKE_COMMAND} -E copy ${GO_OUTPUT_DIR}/libchainmaker_sdk.so ${CMAKE_SOURCE_DIR}/deps/lib/libchainmaker_sdk.so
    COMMAND ${CMAKE_COMMAND} -E copy ${GO_OUTPUT_DIR}/libchainmaker_sdk.h  ${CMAKE_SOURCE_DIR}/deps/include/libchainmaker_sdk.h
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/
    DEPENDS ${CMAKE_SOURCE_DIR}/src/blockchain/chainmaker/main.go
    COMMENT "Building Go shared library"
)

# 2. 定义一个目标表示 go 库
add_custom_target(chainmaker_golib ALL
    DEPENDS ${GO_OUTPUT_DIR}/libchainmaker_sdk.so
)

if(NOT EXISTS ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/go.mod)
    add_custom_command(
        OUTPUT ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/go.mod
        COMMAND go mod init fabric
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/
        COMMENT "Initializing go.mod"
    )

    # 定义一个虚拟目标，依赖 go.mod，这样 make golib 会触发生成
    add_custom_target(init_fabric_go_mod ALL
        DEPENDS ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/go.mod
    )
endif()


# 1. 定义生成 go 动态库的命令
add_custom_command(
    OUTPUT ${GO_OUTPUT_DIR}/libfabric_sdk.so ${GO_OUTPUT_DIR}/libfabric_sdk.h
    COMMAND ${CMAKE_COMMAND} -E make_directory ${GO_OUTPUT_DIR}
    COMMAND go mod tidy
    COMMAND go build -o ${GO_OUTPUT_DIR}/libfabric_sdk.so -buildmode=c-shared ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/main.go
    # 复制 so 和 h 到指定目录
    COMMAND ${CMAKE_COMMAND} -E copy ${GO_OUTPUT_DIR}/libfabric_sdk.so ${CMAKE_SOURCE_DIR}/deps/lib/libfabric_sdk.so
    COMMAND ${CMAKE_COMMAND} -E copy ${GO_OUTPUT_DIR}/libfabric_sdk.h  ${CMAKE_SOURCE_DIR}/deps/include/libfabric_sdk.h
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/
    DEPENDS ${CMAKE_SOURCE_DIR}/src/blockchain/fabric/main.go
    COMMENT "Building Go shared library"
)

# 2. 定义一个目标表示 go 库
add_custom_target(fabric_golib ALL
    DEPENDS ${GO_OUTPUT_DIR}/libfabric_sdk.so
)

include(external_boost)
include(external_yamlcpp)
include(external_rapidjson)
# 根据选项 include
if (USE_LSSDK)
  message(STATUS ">>> building with lssdk")
  add_definitions(-DUSE_LSSDK)
  include(external_lssdk)
endif()
include(external_chainmaker)
include(external_fabric)
if (USE_BCOS)
  message(STATUS ">>> building with bcos")
  add_definitions(-DUSE_BCOS)
  include(external_bcos)
endif()

include_directories(
  "${CMAKE_SOURCE_DIR}"
  "${CMAKE_SOURCE_DIR}/src")

file(GLOB_RECURSE SRC_LIST "src/*.cpp" "src/*.c")
file(GLOB_RECURSE HEADERS "src/*.h")
add_executable(zkagent src/main.cpp ${SRC_LIST} ${HEADERS})
target_include_directories(zkagent PRIVATE ..)
install(TARGETS zkagent DESTINATION bin)
if (USE_LSSDK)
  message(STATUS ">>> building with lssdk")
  target_link_libraries(zkagent PUBLIC YamlCpp Lssdk Boost::Log Boost::System Boost::Filesystem Boost::Json dl m z pthread curl Chainmakersdk Fabricsdk ) 
elseif (USE_BCOS)
  message(STATUS ">>> building with bcos")
  target_link_libraries(zkagent PUBLIC YamlCpp Boost::Log Boost::System Boost::Filesystem Boost::Json Bcossdk ssl crypto dl m z pthread curl Chainmakersdk Fabricsdk ) 
else(!USE_BCOS AND !USE_LSSDK)
  message(STATUS ">>> building with other")
  target_link_libraries(zkagent PUBLIC YamlCpp Boost::Log Boost::System Boost::Filesystem Boost::Json dl m z pthread curl Chainmakersdk Fabricsdk ) 
endif()

#add_custom_command(
#    TARGET DataTrans POST_BUILD
#    COMMAND ${CMAKE_COMMAND} -E copy
#        ${CMAKE_BINARY_DIR}/bin/DataTrans
#        ${CMAKE_BINARY_DIR}/bin/server/DataTrans
#    COMMAND ${CMAKE_COMMAND} -E copy
#        ${CMAKE_BINARY_DIR}/bin/DataTrans
#        ${CMAKE_BINARY_DIR}/bin/client/DataTrans
#)

# add_custom_command(
#     TARGET zkagent POST_BUILD
#     COMMAND ${CMAKE_COMMAND} -E copy
#             "${CMAKE_SOURCE_DIR}/script/generate_poseidon_hash.js"
#             "${CMAKE_BINARY_DIR}/bin/script/generate_poseidon_hash.js"
# )

# add_custom_command(
#     TARGET zkagent POST_BUILD
#     COMMAND ${CMAKE_COMMAND} -E copy
#             "${CMAKE_SOURCE_DIR}/contracts/Challenge.sol"
#             "${CMAKE_BINARY_DIR}/bin/contracts/solidity/Challenge.sol"
# )

# add_custom_command(
#     TARGET zkagent POST_BUILD
#     COMMAND ${CMAKE_COMMAND} -E copy
#             "${CMAKE_SOURCE_DIR}/contracts/FileStore.sol"
#             "${CMAKE_BINARY_DIR}/bin/contracts/solidity/FileStore.sol"
# )

# add_custom_command(
#     TARGET zkagent POST_BUILD
#     COMMAND ${CMAKE_COMMAND} -E copy
#             "${CMAKE_SOURCE_DIR}/contracts/Groth16Verifier.sol"
#             "${CMAKE_BINARY_DIR}/bin/contracts/solidity/Groth16Verifier.sol"
# )
