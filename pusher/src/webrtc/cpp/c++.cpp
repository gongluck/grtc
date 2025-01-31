/*
 * @Author: gongluck
 * @Date: 2025-01-30 00:47:13
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-31 19:53:20
 */

extern "C"
{
#include "c++.h"
}

#include <iostream>
#include <memory>
#include <string>

int CppFunction()
{
    std::shared_ptr<std::string> str = std::make_shared<std::string>("Hello from C++!");
    std::cout << *str << std::endl;
    return 0;
}