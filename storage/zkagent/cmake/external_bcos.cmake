include(GNUInstallDirs)
include(ExternalProject)


# 编译库
if (NOT BCOSSDK_ROOT_DIR)
    set(INSTALL_DIR ${CMAKE_SOURCE_DIR}/deps)

    set(BCOSSDK_ROOT_DIR ${INSTALL_DIR})
    unset(INSTALL_DIR)
endif()


# 创建导入库
set(BCOSSDK_INCLUDE_DIR ${BCOSSDK_ROOT_DIR}/include)
set(BCOSSDK_LIBRARY_DIR ${BCOSSDK_ROOT_DIR}/lib)
file(MAKE_DIRECTORY ${BCOSSDK_INCLUDE_DIR})

add_library(Bcossdk SHARED IMPORTED)
set_target_properties(Bcossdk PROPERTIES
    IMPORTED_LOCATION ${BCOSSDK_LIBRARY_DIR}/libbcos-c-sdk.so
)

set_property(TARGET Bcossdk PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${LSSDK_INCLUDE_DIR})
