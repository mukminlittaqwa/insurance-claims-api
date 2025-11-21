# Insurance Claims API

Backend API untuk sistem pengajuan dan persetujuan klaim asuransi, dibangun menggunakan **Go (Golang)** dengan framework **Gin** dan database **MongoDB**.

## Fitur Utama

- Autentikasi & otorisasi berbasis JWT
- 3 role user:
  - **User** → mengajukan klaim
  - **Verificator** → memverifikasi dokumen & data klaim
  - **Approval** → menyetujui atau menolak klaim
- CRUD klaim dengan status workflow (Draft → Submitted → Verified → Approved/Rejected)
- Upload dokumen pendukung (PDF/Image)
- History & audit log perubahan status

## Tech Stack

- **Language**: Go (Gyuk, 1.22+)
- **Framework**: Gin Gonic
- **Database**: MongoDB (dengan official driver `go.mongodb.org/mongo-driver`)
- **Authentication**: JWT
- **Validation**: `github.com/go-playground/validator/v10`
- **Env Config**: `github.com/joho/godotenv`

## Struktur Folder

├── cmd/
│ └── api/ # entry point (main.go)
├── internal/
│ ├── handlers/ # HTTP handlers (Gin routes)
│ ├── middleware/ # Auth, role checker, dll
│ ├── models/ # Struct MongoDB & request/response
│ ├── repository/ # Interaksi langsung dengan MongoDB
│ ├── services/ # Business logic
│ └── utils/ # Helper functions
├── .env.example
├── go.mod
└── README.md

## Prerequisites

- Go 1.22 atau lebih baru
- MongoDB (lokal atau cloud seperti MongoDB Atlas)
- Git

##running

1. add env (.env)
2. go mod tidy
3. go run cmd/api/main.go #untuk running
