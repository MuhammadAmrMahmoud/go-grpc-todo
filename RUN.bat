REM Run the server
echo Starting server...
start /B server/server.exe  REM Start server in background

REM Wait for a few seconds to ensure the server is up and running
timeout /t 2 /nobreak

REM Run the client
echo Starting client...
start client/client.exe