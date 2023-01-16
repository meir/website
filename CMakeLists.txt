cmake_minimum_required(VERSION 3.16)
project(static_builder)
set(CMAKE_C_STANDARD 11)

set(SOURCE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/src)
set(SOURCE_FILES ${SOURCE_DIR}/main.c)

add_executable(static_builder ${SOURCE_FILES})