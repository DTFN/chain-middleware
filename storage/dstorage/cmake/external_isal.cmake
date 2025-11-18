include(GNUInstallDirs)
include(ExternalProject)


# 编译库
if (NOT ISAL_ROOT_DIR)
    ExternalProject_Add(isal
        PREFIX  ${CMAKE_SOURCE_DIR}/deps
        INSTALL_DIR ${CMAKE_SOURCE_DIR}/deps

        GIT_REPOSITORY https://github.com/intel/isa-l.git
        GIT_TAG        master
        UPDATE_COMMAND ""
        CONFIGURE_COMMAND ./autogen.sh && ./configure --prefix=<INSTALL_DIR>
        BUILD_COMMAND make -j4
        INSTALL_COMMAND make install
        BUILD_IN_SOURCE 1
    )

    ExternalProject_Get_Property(isal INSTALL_DIR)
    set(ISAL_ROOT_DIR ${INSTALL_DIR})
    unset(INSTALL_DIR)
endif()


# 创建导入库
set(ISAL_INCLUDE_DIR ${ISAL_ROOT_DIR}/include)
set(ISAL_LIBRARY_DIR ${ISAL_ROOT_DIR}/lib)
file(MAKE_DIRECTORY ${ISAL_INCLUDE_DIR})

add_library(Isal SHARED IMPORTED)
set_target_properties(Isal PROPERTIES
    IMPORTED_LOCATION ${ISAL_LIBRARY_DIR}/libisal.so
)
set_property(TARGET Isal PROPERTY INTERFACE_INCLUDE_DIRECTORIES ${ISAL_INCLUDE_DIR})
add_dependencies(Isal isal)
