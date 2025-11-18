#------------------------------------------------------------------------------
# @copyright (C) 2023 0shu.
#
# @brief 工具函数库
# @author Mingsheng Zhu
# @date 2023/06/16
#------------------------------------------------------------------------------
include(ProcessorCount)


# 判断参数是否为整数
# 参数:
#   - in: 待判断的参数
#   - out: 输出结果，TRUE：整数，FALSE：不是整数
function(is_integer in out)
    string(REGEX MATCH "^[0-9]+$" INT_VALUE ${in})
    if (NOT INT_VALUE)
        set(${out} FALSE PARENT_SCOPE)
    else()
        set(${out} TRUE PARENT_SCOPE)
    endif()
endfunction(is_integer)


# 获取用于编译外部库的线程数
# 参数:
#   - out: 输出线程数
function(get_compile_cores out)
    string(REGEX MATCH "^[0-9]+$" value "${${out}}")
    if (NOT value OR value EQUAL 0)
        ProcessorCount(CORES)
        math(EXPR COMPILE_CORES "${CORES} / 2")
        if (${COMPILE_CORES} LESS 1)
            set(COMPILE_CORES 1)
        endif()
        set(${out} ${COMPILE_CORES} PARENT_SCOPE)
    endif()
endfunction(get_compile_cores)


# 以黄色字体打印警告信息
# 参数:
#   - message: 指定打印的消息
function(print_warning message)
    execute_process(
        COMMAND echo -e "\\033[1;33mWarning: ${message}\\033[0m"
        OUTPUT_VARIABLE output
        RESULT_VARIABLE result
        ERROR_QUIET
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )
    if(NOT ${result} EQUAL 0)
        message(WARNING "Failed to execute echo command")
    endif()
    message("${output}")
endfunction()
