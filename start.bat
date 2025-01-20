@echo off

@REM call .\docker-utils\cmd\init-network.bat

@REM call .\docker-utils\cmd\add-client.bat -n frontend1 -i 10.0.11.3 -p 3001:3000 -e "NEXT_PUBLIC_API_URL=http://10.0.10.2:8080"

call .\docker-utils\cmd\add-server.bat -m mongo-server-1 -i 10.0.10.3 -s scrapper-server-1 -j 10.0.10.2 -p 8080:8080