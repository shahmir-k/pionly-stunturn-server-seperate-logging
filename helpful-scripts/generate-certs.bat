@echo off
echo Generating self-signed SSL certificates...

REM Create certs directory if it doesn't exist
if not exist "..\certs" mkdir "..\certs"

REM Generate private key
openssl genrsa -out ..\certs\privkey.pem 2048

REM Generate certificate signing request
openssl req -new -key ..\certs\privkey.pem -out ..\certs\cert.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

REM Generate self-signed certificate
openssl x509 -req -in ..\certs\cert.csr -signkey ..\certs\privkey.pem -out ..\certs\fullchain.pem -days 365

REM Clean up CSR file
del ..\certs\cert.csr

echo SSL certificates generated successfully!
echo ..\certs\fullchain.pem - Certificate file
echo ..\certs\privkey.pem - Private key file
echo.
echo Note: These are self-signed certificates for development only.
echo For production, use certificates from a trusted Certificate Authority.
pause 