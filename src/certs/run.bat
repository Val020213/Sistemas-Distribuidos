@echo off

echo [SAN] > san.cnf
for /L %%i in (1,1,254) do (
    echo subjectAltName=IP:10.0.10.%%i >> san.cnf
)
echo subjectAltName=IP:10.0.10.0 >> san.cnf

openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=server" -addext "subjectAltName=IP:10.0.10.0" -config "C:\openssl.cnf"

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile san.cnf
