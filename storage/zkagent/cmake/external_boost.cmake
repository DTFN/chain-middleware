cmake_minimum_required(VERSION 3.15)

include(GNUInstallDirs)
include(ExternalProject)
include(utilities)


# 编译库
if (NOT BOOST_ROOT_DIR)
    get_compile_cores(COMPILE_CORES)

    set(TOOLSET)
    if (${CMAKE_CXX_COMPILER_ID} MATCHES "Clang")
        set(TOOLSET "toolset=clang")
    endif()

    set(CXXFLAGS "cxxflags=-march=x86-64 -fPIC -mtune=generic -fvisibility=hidden -fvisibility-inlines-hidden")
    if (ARCH_NATIVE)
        set(CXXFLAGS "cxxflags=-march=native -fPIC -mtune=generic -fvisibility=hidden -fvisibility-inlines-hidden")
    endif()

    if(WIN32)
        set(CXXFLAGS "define=_WIN32_WINNT=0x0602")
        set(BOOST_BOOTSTRAP_COMMAND bootstrap.bat)
    elseif (APPLE)
        set(SED_CMMAND sed -i .bkp)
        set(BOOST_BOOTSTRAP_COMMAND ./bootstrap.sh COMMAND ${SED_CMMAND} "s/-fcoalesce-templates//g" ${CMAKE_SOURCE_DIR}/deps/src/boost/tools/build/src/tools/darwin.jam)
    else()
        set(SED_CMMAND sed -i)
        set(BOOST_BOOTSTRAP_COMMAND ./bootstrap.sh)
    endif()

    set(BOOST_LIB_PREFIX ${CMAKE_SOURCE_DIR}/deps/src/boost/stage/lib/libboost_)
    set(BOOST_LIB_SUFFIX ${CMAKE_STATIC_LIBRARY_SUFFIX})
    if (MSVC)
        set(BOOST_LIB_SUFFIX "-vc${MSVC_TOOLSET_VERSION}-mt-x64-1_86${CMAKE_STATIC_LIBRARY_SUFFIX}")
    endif()
    set(BOOST_BUILD_FILES ${BOOST_LIB_PREFIX}system${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}chrono${BOOST_LIB_SUFFIX}
            ${BOOST_LIB_PREFIX}date_time${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}regex${BOOST_LIB_SUFFIX}
            ${BOOST_LIB_PREFIX}filesystem${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}random${BOOST_LIB_SUFFIX}
            ${BOOST_LIB_PREFIX}thread${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}serialization${BOOST_LIB_SUFFIX}
            ${BOOST_LIB_PREFIX}timer${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}iostreams${BOOST_LIB_SUFFIX}
            ${BOOST_LIB_PREFIX}log${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}program_options${BOOST_LIB_SUFFIX}
            ${BOOST_LIB_PREFIX}unit_test_framework${BOOST_LIB_SUFFIX} ${BOOST_LIB_PREFIX}json${BOOST_LIB_SUFFIX})

    ExternalProject_Add(boost
        PREFIX ${CMAKE_SOURCE_DIR}/deps
        INSTALL_DIR ${CMAKE_SOURCE_DIR}/deps

        DOWNLOAD_NO_PROGRESS TRUE
        DOWNLOAD_EXTRACT_TIMESTAMP FALSE
        DOWNLOAD_NAME boost_1_86_0.tar.bz2
        URL https://boostorg.jfrog.io/artifactory/main/release/1.86.0/source/boost_1_86_0.tar.bz2
        URL_HASH SHA256=1bed88e40401b2cb7a1f76d4bab499e352fa4d0c5f31c0dbae64e24d34d7513b

        BUILD_IN_SOURCE TRUE
        CONFIGURE_COMMAND ${BOOST_BOOTSTRAP_COMMAND}
        BUILD_COMMAND ./b2 stage threading=multi link=static variant=release address-model=64 ${TOOLSET} ${CXXFLAGS}
            --with-system
            --with-chrono
            --with-date_time
            --with-regex
            --with-filesystem
            --with-random
            --with-thread
            --with-serialization
            --with-timer
            --with-iostreams
            --with-log
            --with-program_options
            --with-test
            --with-json
            -s NO_BZIP2=1 -s NO_LZMA=1 -s NO_ZSTD=1
            -j${COMPILE_CORES}
        INSTALL_COMMAND ""
        #INSTALL_COMMAND "./b2 install --prefix=${CMAKE_SOURCE_DIR}/deps"
        BUILD_BYPRODUCTS ${BOOST_BUILD_FILES}

        LOG_CONFIGURE TRUE
        LOG_BUILD TRUE
        LOG_INSTALL TRUE
    )

    ExternalProject_Get_Property(boost SOURCE_DIR)
    set(BOOST_ROOT_DIR ${SOURCE_DIR})
    unset(SOURCE_DIR)
endif()


# 创建导入库
set(BOOST_INCLUDE_DIR ${BOOST_ROOT_DIR})
set(BOOST_LIBRARY_DIR ${BOOST_ROOT_DIR}/stage/lib)

add_library(Boost::System STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::System PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_system${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::System PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::System boost)

add_library(Boost::Chrono STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Chrono PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_chrono${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Chrono PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::Chrono boost)

add_library(Boost::DataTime STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::DataTime PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_date_time${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::DataTime PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
set_property(TARGET Boost::DataTime PROPERTY INTERFACE_LINK_LIBRARIES Boost::System)
add_dependencies(Boost::DataTime boost)

add_library(Boost::Regex STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Regex PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_regex${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Regex PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
set_property(TARGET Boost::Regex PROPERTY INTERFACE_LINK_LIBRARIES Boost::System)
add_dependencies(Boost::Regex boost)

add_library(Boost::Filesystem STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Filesystem PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_filesystem${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Filesystem PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
set_property(TARGET Boost::Filesystem PROPERTY INTERFACE_LINK_LIBRARIES Boost::System)
add_dependencies(Boost::Filesystem boost)

add_library(Boost::Random STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Random PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_random${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Random PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::Random boost)

add_library(Boost::Thread STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Thread PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_thread${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Thread PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
set_property(TARGET Boost::Thread PROPERTY INTERFACE_LINK_LIBRARIES Boost::Chrono Boost::DataTime Boost::Regex)
add_dependencies(Boost::Thread boost)

add_library(Boost::Serialization STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Serialization PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_serialization${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Serialization PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::Serialization boost)

add_library(Boost::timer STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::timer PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_timer${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::timer PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
set_property(TARGET Boost::timer PROPERTY INTERFACE_LINK_LIBRARIES Boost::Chrono)
add_dependencies(Boost::timer boost)

add_library(Boost::iostreams STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::iostreams PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_iostreams${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::iostreams PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::iostreams boost)

add_library(Boost::Log STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Log PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_log${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Log PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
set_property(TARGET Boost::Log PROPERTY INTERFACE_LINK_LIBRARIES Boost::Filesystem Boost::Thread)
add_dependencies(Boost::Log boost)

add_library(Boost::program_options STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::program_options PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_program_options${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::program_options PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::program_options boost)

add_library(Boost::UnitTestFramework STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::UnitTestFramework PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_unit_test_framework${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::UnitTestFramework PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::UnitTestFramework boost)

# 添加 Boost.JSON（新的）
add_library(Boost::Json STATIC IMPORTED GLOBAL)
set_property(TARGET Boost::Json PROPERTY IMPORTED_LOCATION ${BOOST_LIBRARY_DIR}/libboost_json${BOOST_LIB_SUFFIX})
set_property(TARGET Boost::Json PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${BOOST_INCLUDE_DIR})
add_dependencies(Boost::Json boost)