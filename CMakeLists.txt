cmake_minimum_required(VERSION 3.12)
project(TopURL)

set(CMAKE_CXX_STANDARD 14)

add_executable(main main.cpp util.cpp)
add_executable(test test.cpp util.cpp)
