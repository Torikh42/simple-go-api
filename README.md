# Simple Product & Auth API — Golang (Dockerized)

Proyek RESTful API berbasis Golang dan PostgreSQL yang diorkestrasi menggunakan Docker Compose dengan arsitektur Clean Code, multi-stage Docker build, dan volume persistensi.

> Proyek ini disusun untuk memenuhi **Tugas 1: Dockerisasi Backend dengan Docker Compose** — Study Club Backend Intermediate, KSM Veterantech UPNVJ 2026.

---

## 🛠 Tech Stack & Tools

- **Backend**: Go 1.22+ (`net/http`, Standard Go Layout)
- **Database**: PostgreSQL 17 (via Docker Container)
- **Database Driver**: [pgx/v5](https://github.com/jackc/pgx)
- **Query Generator**: [sqlc](https://sqlc.dev)
- **Authentication**: JWT (`golang-jwt/jwt/v5`) & `golang.org/x/crypto/bcrypt`
- **Containerization**: Docker & Docker Compose (Multi-stage build)
- **Database GUI**: JetBrains DataGrip / DBeaver / pgAdmin

## 🏛 Application Architecture & Project Structure

Mengikuti prinsip *Clean Architecture* dan *Standard Go Project Layout*:

```text
Handler (HTTP) → Service (Business Logic) → Repository (Database Layer) → PostgreSQL
```

```text
go-api/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point aplikasi (injeksi dependensi & router)
├── db/
│   ├── queries/
│   │   ├── product.sql             # SQL queries untuk produk (input sqlc)
│   │   └── users.sql               # SQL queries untuk users & auth (input sqlc)
│   └── schema/
│       ├── products.sql            # Skema DDL tabel products
│       └── users_and_sessions.sql  # Skema DDL tabel users & sessions
├── internal/
│   ├── db/                         # Kode type-safe yang digenerate oleh sqlc
│   ├── handlers/                   # Layer HTTP (Product & Auth Handler)
│   ├── middleware/                 # Middleware (Auth JWT, Logger, Recoverer)
│   ├── models/                     # Model data / entity
│   ├── repository/                 # Layer akses database
│   ├── routes/                     # Definisi rute URL & middleware binding
│   └── services/                   # Layer logika bisnis & validasi
├── migrations/                     # Skrip SQL untuk inisialisasi database di Docker
├── Dockerfile                      # Multi-stage Docker build
└── docker-compose.yml              # Orkestrasi multi-container backend & database
```

---

## 🏗 Arsitektur Container (Docker Compose)

Aplikasi dan database berjalan di dalam container terisolasi yang dihubungkan melalui satu bridge network (`app-network`):

```text
[Host Machine (Komputer Lokal)]
  │
  ├── Port 8080 ──────────────> [Container: simple-go-api-app (Go Backend)]
  │                                    │
  │                              (Bridge Network: app-network)
  │                                    │
  └── Port 5434 (DataGrip) ───> [Container: simple-go-api-db (PostgreSQL 17)]
                                       │
                                 (Named Volume: pgdata)
```

---

## 📋 Prasyarat (Prerequisites)

Sebelum menjalankan aplikasi, pastikan komputer Anda telah terinstal:
- [Docker Engine & Docker Compose](https://docs.docker.com/get-docker/) (Docker Desktop)
- [Git](https://git-scm.com/)
- *(Opsional)* [Go 1.22+](https://go.dev/dl/) jika ingin menjalankan secara lokal tanpa Docker
- *(Opsional)* [DataGrip](https://www.jetbrains.com/datagrip/) atau Database Client lainnya

---

## 🚀 Panduan Setup & Menjalankan Aplikasi

### 1. Clone the repository

```bash
git clone https://github.com/Torikh42/simple-go-api.git
cd simple-go-api
```

### 2. Setup Environment Variables
Salin template konfigurasi `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

### 3. Build & Jalankan Seluruh Layanan (Docker Compose)
Jalankan perintah berikut untuk mem-build image backend dan menyalakan seluruh service:
```bash
docker compose up --build -d
```
> **Catatan:**
> - Service `postgres` akan otomatis menginisialisasi skrip SQL pada folder `./migrations` saat pertama kali container dibuat.
> - Service `app` memiliki *healthcheck dependency*, sehingga backend hanya akan menyala setelah PostgreSQL siap menerima koneksi.

Periksa status container yang sedang berjalan:
```bash
docker compose ps
```

Untuk melihat log aplikasi secara langsung:
```bash
docker compose logs -f app
```

---

## 🔌 Panduan Koneksi Database ke DataGrip

Untuk melihat data di database secara visual melalui DataGrip:
1. Buka **DataGrip** $\rightarrow$ klik **+ (New)** $\rightarrow$ **Data Source** $\rightarrow$ **PostgreSQL**.
2. Masukkan parameter koneksi:
   - **Host**: `localhost` (atau `127.0.0.1`)
   - **Port**: `5434` *(port yang diexpose ke host)*
   - **User**: `postgres`
   - **Password**: `postgres`
   - **Database**: `go_api_db`
3. Klik **Test Connection**. Jika diminta unduh driver, klik *Download Driver Files*.
4. Setelah status berhasil (*Succeeded*), klik **Apply** $\rightarrow$ **OK**.

---

## 📡 Pengujian API Endpoints

Setelah container menyala, API dapat diakses di `http://localhost:8080`.

### Daftar Endpoint:

| Method | Endpoint | Keterangan | Akses |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Healthcheck API | Public |
| `POST` | `/register` | Registrasi user baru | Public |
| `POST` | `/login` | Login user & generate JWT | Public |
| `GET` | `/products` | Mengambil seluruh daftar produk | 🔒 Butuh JWT |
| `POST` | `/products` | Menambah produk baru | 🔒 Butuh JWT |
| `GET` | `/products/{id}` | Mengambil detail produk | 🔒 Butuh JWT |
| `PUT` | `/products/{id}` | Mengubah data produk | 🔒 Butuh JWT |
| `DELETE` | `/products/{id}` | Menghapus produk | 🔒 Butuh JWT |

---

### Contoh Pengujian Request:

#### 1. Registrasi User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "torikh",
    "email": "torikh@example.com",
    "password": "password123"
  }'
```

#### 2. Login untuk Mendapatkan Token
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "torikh@example.com",
    "password": "password123"
  }'
```
*Salin token dari field `token` pada response JSON.*

#### 3. Membuat Produk Baru (Write Data)
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_JWT_ANDA>" \
  -d '{
    "name": "Keyboard Mechanical",
    "price": 750000,
    "stock": 15
  }'
```

#### 4. Membaca Produk (Read Data)
```bash
curl -X GET http://localhost:8080/products \
  -H "Authorization: Bearer <TOKEN_JWT_ANDA>"
```

---

## 💾 Pembuktian Persistensi Data (Rubrik Tugas 1)

Untuk membuktikan data di database tetap tersimpan (*persistent*) meskipun container dimatikan atau dibuat ulang:

1. **Tambah data produk baru** melalui endpoint `POST /products` atau cek data di DataGrip.
2. **Hentikan dan hapus container** (tanpa menghapus named volume):
   ```bash
   docker compose down
   ```
3. **Jalankan kembali container**:
   ```bash
   docker compose up -d
   ```
4. **Cek kembali data** menggunakan endpoint `GET /products` atau *refresh* di DataGrip:
   ```bash
   curl -X GET http://localhost:8080/products -H "Authorization: Bearer <TOKEN_JWT_ANDA>"
   ```
5. **Hasil:** Data produk yang telah dibuat sebelumnya tetap ada dan tidak hilang, membuktikan volume `pgdata` berfungsi dengan baik.

---
