# SSCM - Self Signed Certificate Management

SSCM manages your self-signed certificate, which usually used within internal networks.

SSCM recommendation when managing self-signed certificate is to have single Self-Signed Root Certificate Authority that can be distributed within the organization to be locally trusted. Then use those Self Signed RootCA to create domain/subdomain/IP specific certificates.

## Planned Endpoints

### [GET] /cert?root=xx - List all child certificates (or root if ?root=0)
### [POST] /cert?parent=xx - Create a child certificate (or root if ?parent=0)
### [GET] /cert/id - Download specific certificate (root without the key, child with key)
### [GET] /search?q=xx - Search all certificates by query

## Development

### FrontEnd

To develop the frontend part, enter web directory and run develpoment mode:
```bash
cd web
npm run dev
```

### Backend

## Building / Compiling

Building this app requires 2 step. First the Frontend and then the Backend.
The Frontend must be built first, because the static files will be bundled into the Backend.

```bash
cd web
npm run build
cd ..
go build --trimpath .
```

After compilation, you only need the exe/binary file.