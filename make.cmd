@setlocal
@set PROMPT=$G
if not "%1" == "" goto %1

:build
    go build
    goto end
:get
    for %%I in (github.com/mattn/go-runewidth github.com/shiena/ansicolor github.com/zetamatta/nyagos/conio) do ( go get %%I & cd %GOPATH%\src\%%I & git pull origin master)
    goto end
:fmt
    go fmt
    goto end
:clean
    if exist cure.exe del cure.exe
    goto end
:snapshot
    zip -9 cure-%DATE:/=%.zip readme.md cure.exe
:end
