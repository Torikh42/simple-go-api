# Golang Backend Golden Rules (AGENTS.md)

File ini memuat aturan emas (Golden Rules), standar kode (Clean Code), dan arsitektur yang HARUS diikuti oleh semua AI Agent (termasuk saya) atau developer yang berkontribusi pada proyek ini. Aturan ini sangat spesifik untuk ekosistem Go, membedakannya dari kebiasaan di ekosistem Node.js/TypeScript.

## 1. Arsitektur: Standard Go Project Layout
Kita menggunakan struktur folder standar industri untuk memisahkan *concern*. Dilarang keras menaruh semua logika di `main.go`.

- **`cmd/api/`**: Tempat file `main.go` berada. File ini harus sekecil mungkin. Hanya berisi inisialisasi aplikasi (koneksi DB, inisialisasi router, menyalakan server).
- **`internal/`**: Kode yang spesifik untuk proyek ini dan tidak boleh di-import oleh repositori lain.
  - **`internal/handlers/`**: Layer HTTP. Hanya bertugas membaca Request (JSON/Params), memanggil *Service*, dan memformat Response. DILARANG KERAS menaruh logika bisnis atau query database di sini.
  - **`internal/services/`**: Layer Logika Bisnis. Bertugas melakukan validasi logika, kalkulasi, dll. Service tidak boleh tahu apakah request datang dari HTTP, gRPC, atau CLI.
  - **`internal/repository/`**: Layer Database. Semua operasi SQL/NoSQL wajib berada di sini. Dilarang keras memanggil koneksi database di luar layer ini.
  - **`internal/models/`**: Struktur data, *Entity*, dan *Struct* yang dipakai lintas layer.
- **`pkg/`**: Kode utilitas yang bersifat umum dan bisa di-import oleh proyek lain (misalnya library format waktu, custom logger).

## 2. Aturan Clean Code (Idiomatic Go)

### A. Penamaan Variabel (Naming Conventions)
Berbeda dengan JS/Java yang menyukai nama panjang dan eksplisit, Go menyukai nama pendek untuk *scope* yang kecil.
- Gunakan `ctx` untuk `context.Context` (BUKAN `context`).
- Gunakan `req` untuk HTTP Request, `w` untuk HTTP ResponseWriter.
- Gunakan `err` untuk error.
- Gunakan satu huruf untuk metode receiver, misalnya `func (s *Service) GetUser()`.
- Hindari penamaan yang berulang. Gunakan `user.Repository` bukan `user.UserRepository`.

### B. Error Handling (Wajib Eksplisit)
- **Jangan pernah mengabaikan error**. Jangan gunakan `_` untuk menampung error yang mengembalikan nilai kecuali Anda benar-benar yakin 100%.
- Hindari penggunaan `panic()`. Kembalikan error menggunakan `return err`.
- Bungkus (*wrap*) error untuk memberi konteks menggunakan `fmt.Errorf("gagal mengambil user: %w", err)` agar *stack trace* logikal terbentuk.

### C. Dependency Injection
- Hindari *Global Variables* (seperti `var DB *sql.DB`).
- Lakukan injeksi dependensi melalui pembuat (*constructor*).
  ```go
  // BENAR:
  type UserService struct {
      repo UserRepository
  }
  func NewUserService(repo UserRepository) *UserService {
      return &UserService{repo: repo}
  }
  ```

## 3. Concurrency (Goroutines & Context)
- Jangan pernah menjalankan `go func()` tanpa tahu kapan ia akan berhenti. Ini akan menyebabkan *Goroutine Leak*.
- Selalu terima `ctx context.Context` sebagai parameter pertama di setiap fungsi yang melakukan proses I/O (Database, HTTP Call) agar eksekusi bisa dibatalkan jika terjadi *timeout*.

## 4. Database (Tanpa ORM Berat)
- Hindari penggunaan ORM berat jika query mulai kompleks. 
- Gunakan raw SQL, query builder yang ringan, atau `sqlc`.
- Jika menggunakan PostgreSQL, *driver* standar yang disarankan adalah `pgx`.

## 5. Middleware & Routing
- Gunakan `net/http.ServeMux` (Go 1.22+) atau `go-chi/chi` karena idiomatik.
- Semua *route* harus dilindungi oleh middleware dasar: `Recoverer` (mencegah server mati karena *panic*), `Logger`, dan `Timeout`.
