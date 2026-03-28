# File Transfer Based on Fountain Codes

A file transfer tool using RaptorQ (fountain codes) for efficient data transmission.

## Project Structure

```
├── backend/           # Go backend
│   ├── cmd/          # Application entry point
│   └── internal/    # Internal packages
│       ├── encoder/  # RaptorQ encoding
│       ├── decoder/  # RaptorQ decoding
│       └── transfer/ # File transfer handling
└── frontend/         # Web frontend (to be added)
```

## Quick Start

```bash
cd backend
go run cmd/main.go
```

The server will start on http://localhost:8080

## API Endpoints

- `GET /health` - Health check
- `POST /api/upload` - Upload file
- `POST /api/encode` - Encode file using fountain codes
- `POST /api/decode` - Decode fountain code packets
- `GET /api/download/:fileId` - Download file
- `GET /api/status/:fileId` - Check file status
