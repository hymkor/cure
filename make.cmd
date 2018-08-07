@setlocal
@set PROMPT=$G
call :"%1"
@endlocal
@exit /b

:""
:"build"
    go fmt
    for %%I in (386 amd64) do call :build %%I
    exit /b
:build
    setlocal
    set "GOARCH=%1"
    if not exist cmd mkdir cmd
    if not exist cmd\%1 mkdir cmd\%1
    go build -o cmd\%1\cure.exe
    endlocal
:"get"
    go get ./...
    exit /b
:"get2"
    go get -u ./...
    exit /b
:"clean"
    if exist cure.exe del cure.exe
    exit /b
:"snapshot"
    zip -9 cure-%DATE:/=%.zip readme.md cure.exe
    exit /b
