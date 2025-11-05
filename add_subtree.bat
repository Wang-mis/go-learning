@echo off
REM ============================================
REM 功能：添加一个本地项目作为 subtree
REM 用法：add-subtree.bat remote-name local-path [prefix]
REM 示例：add-subtree.bat my-lib D:\workspace\my-lib src/libs/my-lib
REM ============================================

set REMOTE_NAME=%1
set LOCAL_PATH=%2
set PREFIX=%3

REM 参数检查
if "%REMOTE_NAME%"=="" (
    echo [ERROR] 缺少 remote-name 参数！
    echo 用法: add-subtree.bat remote-name local-path [prefix]
    exit /b 1
)

if "%LOCAL_PATH%"=="" (
    echo [ERROR] 缺少 local-path 参数！
    echo 用法: add-subtree.bat remote-name local-path [prefix]
    exit /b 1
)

if "%PREFIX%"=="" (
    set PREFIX=%REMOTE_NAME%
)

REM 检查 remote 是否已存在
set REMOTE_EXISTS=
for /f "tokens=1" %%i in ('git remote') do (
    if "%%i"=="%REMOTE_NAME%" (
        set REMOTE_EXISTS=1
    )
)

REM 如果 remote 不存在，则添加
if not defined REMOTE_EXISTS (
    echo [INFO] 添加本地 remote: %REMOTE_NAME%  ->  %LOCAL_PATH%
    git remote add %REMOTE_NAME% %LOCAL_PATH%
    if errorlevel 1 (
        echo [ERROR] git remote add 失败！
        exit /b 1
    )
) else (
    echo [INFO] 远程 %REMOTE_NAME% 已存在，跳过添加
)

echo [INFO] Fetching remote: %REMOTE_NAME%
git fetch %REMOTE_NAME%
if errorlevel 1 (
    echo [ERROR] git fetch 失败！
    exit /b 1
)

echo [INFO] Adding subtree to: %PREFIX%
git subtree add --prefix=%PREFIX% %REMOTE_NAME% main --squash
if errorlevel 1 (
    echo [ERROR] git subtree add 失败！
    exit /b 1
)

echo [SUCCESS] 成功将本地项目 "%LOCAL_PATH%" 添加为子项目 "%PREFIX%"！
