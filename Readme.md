# SSCM - Self Signed Certificate Management

SSCM manages your self-signed certificate, which usually used within internal networks.

SSCM recommendation when managing self-signed certificate is to have single Self-Signed Root Certificate Authority that can be distributed within the organization to be locally trusted. Then use those Self Signed RootCA to create domain/subdomain/IP specific certificates.

## Planned Endpoints

### [GET] /root - List of root certificates
### [POST] /root - Create a root certificate
### [GET] /root/id - Download specific root certificate (without the key)
### [GET] /cert?root=xx - List all child certificates
### [POST] /cert - Create a child certificate
### [GET] /cert/id - Download specific child certificate and the key
### [GET] /search?q=xx - Search all certificates by query
