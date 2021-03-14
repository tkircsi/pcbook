echo "Clean *.pem files"
rm *.pem

# 1. Generate CA's private key and self-signed certificate
# openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem
echo ""
echo "Generate CA's private key and certification"
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=HU/ST=Hungary/L=Budapest/O=Tibor Kircsi/OU=R&D/CN=*.tkircsi.net/emailAddress=tkircsi@gmail.com" -nodes
# -nodes : does not ask for passphrase

echo ""
echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
echo ""
echo "Generate server private key and CSR"
openssl req -newkey rsa:4096 -keyout server-key.pem -out server-req.pem -subj "/C=HU/ST=Hungary/L=Budapest/O=PC BOOK/OU=Computer/CN=*.pcbook.net/emailAddress=tkircsi@gmail.com" -nodes
# -nodes : does not ask for passphrase

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
echo "Genereate the server's signed certificate"
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -days 60 -extfile server-ext.cnf

echo "Server signed certificate"
openssl x509 -in server-cert.pem -noout -text

# 4. Verify certificate
echo ""
echo "Verify server certificate"
openssl verify -CAfile ca-cert.pem server-cert.pem

# 5. Generate client's private key and certificate signing request (CSR)
echo ""
echo "Generate client private key and CSR"
openssl req -newkey rsa:4096 -keyout client-key.pem -out client-req.pem -subj "/C=HU/ST=Hungary/L=Budapest/O=PC CLIENT/OU=Computer/CN=*.pcclient.net/emailAddress=tkircsi@gmail.com" -nodes
# -nodes : does not ask for passphrase

# 6. Use CA's private key to sign web server's CSR and get back the signed certificate
echo "Genereate the clients's signed certificate"
openssl x509 -req -in client-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -days 60 -extfile client-ext.cnf

echo "Client signed certificate"
openssl x509 -in client-cert.pem -noout -text

# 7. Verify certificate
echo ""
echo "Verify client certificate"
openssl verify -CAfile ca-cert.pem client-cert.pem

