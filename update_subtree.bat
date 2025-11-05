@echo off
REM ============================================
REM 功能：更新一个已添加的 subtree 子项目
REM 用法：update-subtree.bat remote-name [prefix]
REM 示例：update-subtree.bat my-lib src/libs/my-lib
REM ============================================

set REMOTE_NAME=%1
set PREFIX=%2

if "%REMOTE_NAME%"=="" (
    echo [ERROR] 缺少参数！
    echo 用法: update-subtree.bat remote-name [prefix]
    exit /b 1
)

if "%PREFIX%"=="" (
    set PREFIX=%REMOTE_NAME%
)

echo [INFO] Fetching remote: %REMOTE_NAME%
git fetch %REMOTE_NAME%
if errorlevel 1 (
    echo [ERROR] git fetch 失败！
    exit /b 1
)

echo [INFO] Pulling subtree into: %PREFIX%
git subtree pull --prefix=%PREFIX% %REMOTE_NAME% main --squash
if errorlevel 1 (
    echo [ERROR] git subtree pull 失败！
    exit /b 1
)

echo [SUCCESS] 成功更新子项目 %REMOTE_NAME% 到子目录 %PREFIX%！
