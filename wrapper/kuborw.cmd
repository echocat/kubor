@ECHO OFF
SETLOCAL
REM ##############################################################################
REM ##                                                                          ##
REM ##  kubor bootstrap wrapper for Windows systems                             ##
REM ##                                                                          ##
REM ##############################################################################
REM ##  DO NOT EDIT!!                                                           ##
REM ##############################################################################

SET scriptName=%~nx0
SET dirName=%~dp0
SET os=windows
SET arch=386
SET ext=.exe
IF "%PROCESSOR_ARCHITECTURE%" == "AMD64" (
    SET arch=amd64
)

IF NOT EXIST "%dirName%\kuborw" (
    CALL :fatal This kubor wrapper was not initiated correctly. Try download kubor binary and run: kubor wrapper install
    EXIT /b 1
)
SET findCmd=FINDSTR /B "version=" "%dirName%\kuborw"
FOR /f %%i IN ('%findCmd%') DO SET versionLine=%%i
SET version=%versionLine:~9,-1%

SET binariesCacheDir=%LOCALAPPDATA%\kubor\binaries
IF NOT EXIST "%binariesCacheDir%" (
    md "%binariesCacheDir%"
)
IF NOT ERRORLEVEL 0 (
    CALL :fatal "Cannot create cache directory for storing binaries. See above."
)
SET binaryFileName=kubor-%os%-%arch%-%version%%ext%
SET binary=%binariesCacheDir%\%binaryFileName%

IF NOT EXIST "%binary%" (
    CALL :doDownload
) ELSE (
    "%binary%" version 2>&1 | find "%version%" > NUL
    IF NOT ERRORLEVEL 0 (
        CALL :doDownload
    )
)

IF "%ERRORLEVEL%" == "0" (
    "%binary%" %*
)
EXIT /b %ERRORLEVEL%
GOTO :eofSuccess

:doDownload
    SETLOCAL
    SET binaryDownloadUrl=https://github.com/levertonai/kubor/releases/download/%version%/kubor-%os%-%arch%%ext%
    CALL :info Downloading %binaryDownloadUrl%...

    SET tmpFile=%binary%.%RANDOM%.tmp
    PowerShell -Command "(New-Object Net.WebClient).DownloadFile('%binaryDownloadUrl%','%tmpFile%')"
    IF "%ERRORLEVEL%" == "0" (
        MOVE /Y "%tmpFile%" "%binary%" > NUL
        IF NOT "%ERRORLEVEL%" == "0" (
            CALL :fatal Was not able to move %tmpFile% to %binary%. See above.
        )
    ) ELSE (
        CALL :fatal Was not able to download binary from %binaryDownloadUrl%. See above.
    )
    ENDLOCAL
    EXIT /b %ERRORLEVEL%

:fatal
    ECHO.FATAL: %*
    GOTO :eofError
    EXIT /b 1

:info
    ECHO.INFO: %*
    EXIT /b 0

:eofError
EXIT /b 1
GOTO :eof

:eofSuccess
EXIT /b 0
