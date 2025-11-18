include(GNUInstallDirs)
include(ExternalProject)


# 编译库
if (NOT RAPIDJSON_ROOT_DIR)
    ExternalProject_Add(rapidjson
        PREFIX ${CMAKE_SOURCE_DIR}/deps
        INSTALL_DIR ${CMAKE_SOURCE_DIR}/deps

        DOWNLOAD_NAME rapidjson-1.1.0.tar.gz
        DOWNLOAD_NO_PROGRESS TRUE
        DOWNLOAD_EXTRACT_TIMESTAMP FALSE
        URL https://codeload.github.com/Tencent/rapidjson/tar.gz/refs/tags/v1.1.0
        URL_HASH SHA256=bf7ced29704a1e696fbccf2a2b4ea068e7774fa37f6d7dd4039d0787f8bed98e

        CMAKE_ARGS -DCMAKE_INSTALL_PREFIX=<INSTALL_DIR>
            -DCMAKE_POSITION_INDEPENDENT_CODE=${BUILD_SHARED_LIBS}
            -DCMAKE_BUILD_TYPE=Release
            -DCMAKE_C_COMPILER=${CMAKE_C_COMPILER}
            -DCMAKE_CXX_COMPILER=${CMAKE_CXX_COMPILER}
            -DRAPIDJSON_BUILD_DOC=OFF
            -DRAPIDJSON_BUILD_EXAMPLES=OFF
            -DRAPIDJSON_BUILD_TESTS=OFF

        LOG_CONFIGURE TRUE
        LOG_BUILD TRUE
        LOG_INSTALL TRUE
    )

    ExternalProject_Get_Property(rapidjson INSTALL_DIR)
    set(RAPIDJSON_ROOT_DIR ${INSTALL_DIR})
    unset(INSTALL_DIR)
endif()


# 创建导入库
set(RAPIDJSON_INCLUDE_DIR ${RAPIDJSON_ROOT_DIR}/include)
set(RAPIDJSON_LIBRARY_DIR ${RAPIDJSON_ROOT_DIR}/lib)
file(MAKE_DIRECTORY ${RAPIDJSON_INCLUDE_DIR})

add_library(RapidJson STATIC IMPORTED)
set_property(TARGET RapidJson PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${RAPIDJSON_INCLUDE_DIR})
add_dependencies(RapidJson rapidjson)
