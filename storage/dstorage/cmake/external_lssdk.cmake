include(GNUInstallDirs)
include(ExternalProject)


# 编译库
if (NOT LSSDK_ROOT_DIR)
    ExternalProject_Add(lssdk
        PREFIX  ${CMAKE_SOURCE_DIR}/deps
        INSTALL_DIR ${CMAKE_SOURCE_DIR}/deps

        GIT_REPOSITORY git@192.168.1.232:dtfn2.0/cpp-sdk.git
        GIT_TAG        master
        GIT_SUBMODULES_RECURSE TRUE
        UPDATE_COMMAND ""
        
        CMAKE_ARGS
            -DCMAKE_BUILD_TYPE=Release
            -DBUILD_SHARED_LIBS=on
            -DCMAKE_INSTALL_PREFIX=<INSTALL_DIR>
        
        BUILD_COMMAND make -j4
        INSTALL_COMMAND make install
        #${CMAKE_COMMAND} -E copy <BINARY_DIR>/lib/liblsc-sdk.so <INSTALL_DIR>/lib/liblsc-sdk.so &&
        #               ${CMAKE_COMMAND} -E copy_directory <BINARY_DIR>/cppsdk/include/ <INSTALL_DIR>/include/
        BUILD_IN_SOURCE 1
    )

    ExternalProject_Get_Property(lssdk INSTALL_DIR)
    set(LSSDK_ROOT_DIR ${INSTALL_DIR})
    unset(INSTALL_DIR)
endif()


# 创建导入库
set(LSSDK_INCLUDE_DIR ${LSSDK_ROOT_DIR}/include)
set(LSSDK_LIBRARY_DIR ${LSSDK_ROOT_DIR}/lib)
file(MAKE_DIRECTORY ${LSSDK_INCLUDE_DIR})

add_library(Lssdk SHARED IMPORTED)
set_target_properties(Lssdk PROPERTIES
    IMPORTED_LOCATION ${LSSDK_LIBRARY_DIR}/liblsc-sdk.so
)

set_property(TARGET Lssdk PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${LSSDK_INCLUDE_DIR})
add_dependencies(Lssdk lssdk)
