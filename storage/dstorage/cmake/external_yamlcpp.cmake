include(GNUInstallDirs)
include(ExternalProject)
include(utilities)


# 编译库
if (NOT YAML_ROOT_DIR)
    get_compile_cores(COMPILE_CORES)

    ExternalProject_Add(yaml-cpp
        PREFIX ${CMAKE_SOURCE_DIR}/deps
        INSTALL_DIR ${CMAKE_SOURCE_DIR}/deps

        DOWNLOAD_NAME yaml-cpp-0.7.0.tar.gz
        DOWNLOAD_NO_PROGRESS TRUE
        DOWNLOAD_EXTRACT_TIMESTAMP FALSE
        URL https://github.com/jbeder/yaml-cpp/archive/refs/tags/yaml-cpp-0.7.0.tar.gz 
        URL_HASH SHA256=43e6a9fcb146ad871515f0d0873947e5d497a1c9c60c58cb102a97b47208b7c3

        CMAKE_ARGS -DCMAKE_INSTALL_PREFIX=<INSTALL_DIR>
            -DCMAKE_POSITION_INDEPENDENT_CODE=${BUILD_SHARED_LIBS}
            -DCMAKE_BUILD_TYPE=Release
            -DCMAKE_C_COMPILER=${CMAKE_C_COMPILER}
            -DCMAKE_CXX_COMPILER=${CMAKE_CXX_COMPILER}
            -DCMAKE_INSTALL_LIBDIR=lib
            -DYAML_CPP_BUILD_TESTS=OFF
        BUILD_COMMAND make -j${COMPILE_CORES}
        BUILD_BYPRODUCTS <INSTALL_DIR>/lib/libyaml-cpp.a

        LOG_CONFIGURE TRUE
        LOG_BUILD TRUE
        LOG_INSTALL TRUE
    )

    ExternalProject_Get_Property(yaml-cpp INSTALL_DIR)
    set(YAML_ROOT_DIR ${INSTALL_DIR})
    unset(INSTALL_DIR)
endif()


# 创建导入库
set(YAML_INCLUDE_DIR ${YAML_ROOT_DIR}/include)
set(YAML_LIBRARY_DIR ${YAML_ROOT_DIR}/lib)
file(MAKE_DIRECTORY ${YAML_INCLUDE_DIR})

add_library(YamlCpp STATIC IMPORTED)
set_property(TARGET YamlCpp PROPERTY IMPORTED_LOCATION ${YAML_LIBRARY_DIR}/${CMAKE_STATIC_LIBRARY_PREFIX}yaml-cpp${CMAKE_STATIC_LIBRARY_SUFFIX})
set_property(TARGET YamlCpp PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${YAML_INCLUDE_DIR})
add_dependencies(YamlCpp yaml-cpp)