#!/bin/bash

# Generate self-signed SSL certificates for development
echo "Generating self-signed SSL certificates..."

# Create certs directory if it doesn't exist
mkdir -p ../certs

# Generate private key
openssl genrsa -out ../certs/privkey.pem 2048

# Generate certificate signing request
openssl req -new -key ../certs/privkey.pem -out ../certs/cert.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

# Generate self-signed certificate
openssl x509 -req -in ../certs/cert.csr -signkey ../certs/privkey.pem -out ../certs/fullchain.pem -days 365

# Clean up CSR file
rm ../certs/cert.csr

echo "SSL certificates generated successfully!"
echo "../certs/fullchain.pem - Certificate file"
echo "../certs/privkey.pem - Private key file"
echo ""
echo "Note: These are self-signed certificates for development only."
echo "For production, use certificates from a trusted Certificate Authority." 