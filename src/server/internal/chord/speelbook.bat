@echo off
REM Running protoc for chord.proto
protoc --proto_path=.\ --go_out=.\ --go-grpc_out=. chord.proto
if %errorlevel% neq 0 (
    echo Error running protoc.
    exit /b %errorlevel%
)
echo protoc executed successfully.
pause