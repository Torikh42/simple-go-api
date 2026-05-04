# 🚀 Kurikulum: Dari Node.js ke Golang Backend Engineer

Kurikulum ini dirancang khusus untuk Anda yang sudah mahir dengan ekosistem Node.js (Express, Hono, Prisma, TS) dan ingin bertransisi ke Go (Golang) untuk membangun backend berkinerja tinggi.

---

## Modul 1: Kejutan Budaya (Node.js vs Go)
Sebelum menyentuh *framework* web, kita harus meluruskan pola pikir (mental model).
- **Go tidak punya Class**: Memahami `Struct`, `Methods`, dan `Interfaces` (pengganti OOP di TypeScript).
- **Pointer & Value Semantics**: Kapan menggunakan `*` dan `&` (Referencing vs Copying), sesuatu yang diabstraksi habis-habisan di JS.
- **Error Handling**: Selamat tinggal `try-catch`. Selamat datang `if err != nil`.
- **Tidak ada `package.json`**: Memahami `go.mod`, `go.sum`, dan perintah `go mod tidy` / `go get`.

## Modul 2: Arsitektur & Struktur Proyek
Di Node.js Anda bebas menaruh file di mana saja. Go sangat opiniated tentang ini.
- Membedah **Standard Go Project Layout** (`cmd/`, `internal/`, `pkg/`).
- Memisahkan kode agar tidak terekspos menggunakan direktori `internal/`.
- Memahami alur: `Router -> Handler/Controller -> Service/Usecase -> Repository`.

## Modul 3: Web Server & Routing
- Mengenal kekuatan `net/http` di Go 1.22 (Routing bawaan yang sudah secanggih Express).
- Jika perlu, kita menggunakan `go-chi` untuk kemudahan middleware.
- Membaca request (JSON Body, URL Params, Query Params).
- Mengirim response JSON yang konsisten.
- Mengganti Zod dengan Struct Tags (`json:"name" validate:"required"`) dan `validator`.

## Modul 4: Database & Repository (Meninggalkan ORM Berat)
- Mengapa ekosistem Go lebih suka SQL murni daripada ORM berat seperti Prisma/Drizzle.
- Menggunakan `pgx` untuk koneksi PostgreSQL tercepat.
- Menggunakan **`sqlc`** (Alat sakti untuk men-generate type-safe Go code dari raw SQL).
- Mengelola *database migrations* di Go.

## Modul 5: Konteks & Middleware (Go Context)
- Rahasia terbesar Go: `context.Context`.
- Bagaimana Go membatalkan (*cancel*) proses DB jika *user* menutup *browser* (*Timeout & Cancellation*).
- Membuat *Middleware* (Logger, Authentication, Recover dari Crash).

## Modul 6: Concurrency (Superpower Go)
Ini adalah pengganti mutlak dari `async/await` dan *Event Loop* Node.js.
- **Goroutines**: Membuat fungsi berjalan di *background* hanya dengan kata `go` (sangat ringan, bisa ratusan ribu goroutine).
- **Channels**: Cara Goroutine berkomunikasi satu sama lain secara aman tanpa *race condition*.
- **Sync Package**: Menggunakan `sync.WaitGroup` dan `sync.Mutex` untuk kasus tingkat lanjut.

## Modul 7: Clean Code & Dependency Injection
- *Idiomatic Go*: Penamaan variabel pendek (`ctx`, `req`, `err`), bukan *camelCase* panjang seperti di Java/JS.
- Menggunakan *Interface* untuk *Dependency Injection*.
- Mengapa kita tidak butuh *framework* DI di Go (cukup mem-passing dependencies melalui konstruktor/fungsi pembuat).
- Menulis Unit Test (`go test`) tanpa framework tambahan.
