# 📖 Penjelasan Lengkap Kode Go API

Dokumen ini menjelaskan **seluruh kode** di proyek Anda secara mendetail, termasuk fungsi-fungsi Go yang dipakai dan perbandingannya dengan Node.js/Express.

---

## Struktur Folder Proyek

```
go-api/
├── cmd/api/main.go                 # Titik masuk aplikasi (seperti server.js)
├── internal/
│   ├── models/product.go           # Definisi data Product (seperti TypeScript interface)
│   ├── repository/
│   │   ├── product.go              # Interface + In-Memory DB
│   │   └── postgres_product.go     # Implementasi PostgreSQL
│   ├── services/product.go         # Logika bisnis (validasi, aturan)
│   ├── handlers/
│   │   ├── product.go              # HTTP handler (seperti controller di Express)
│   │   └── health.go               # Health check endpoint
│   ├── routes/routes.go            # Daftar semua endpoint
│   └── middleware/logger.go        # Logger middleware
├── migrations/
│   └── 001_create_products.sql     # SQL untuk membuat tabel
├── .env                            # Konfigurasi (DB_HOST, DB_PORT, dll)
├── .gitignore                      # File yang diabaikan Git
├── .air.toml                       # Konfigurasi hot-reload (seperti nodemon)
├── go.mod                          # Daftar dependensi (seperti package.json)
└── go.sum                          # Hash dependensi (seperti package-lock.json)
```

---

## Alur Data Request (Gambaran Besar)

Saat user mengirim request `GET /products/1`, berikut perjalanannya:

```
Browser/Postman
     │
     ▼
┌──────────────┐
│  Middleware   │  middleware/logger.go → Catat waktu mulai
│   (Logger)    │
└──────┬───────┘
       ▼
┌──────────────┐
│   Routes     │  routes/routes.go → Cocokkan URL dengan handler
└──────┬───────┘
       ▼
┌──────────────┐
│   Handler    │  handlers/product.go → Baca URL param, panggil Service
└──────┬───────┘
       ▼
┌──────────────┐
│   Service    │  services/product.go → Validasi logika bisnis
└──────┬───────┘
       ▼
┌──────────────┐
│  Repository  │  repository/postgres_product.go → Jalankan query SQL
└──────┬───────┘
       ▼
┌──────────────┐
│  PostgreSQL  │  Database sungguhan
└──────────────┘
```

> [!TIP]
> Data mengalir ke bawah (request), lalu kembali ke atas (response). Setiap layer hanya tahu layer di bawahnya.

---

## File 1: `models/product.go` — Cetak Biru Data

```go
package models

type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price int64  `json:"price"`
    Stock int    `json:"stock"`
}
```

### Penjelasan Baris per Baris

| Kode | Penjelasan |
|------|------------|
| `package models` | Semua file di folder yang sama harus punya nama package yang sama. Ini seperti folder organisasi. |
| `type Product struct` | **Struct** adalah cetak biru data di Go. Mirip `interface` di TypeScript atau `class` di Java. Bedanya, struct di Go **tidak bisa punya inheritance** (turunan). |
| `int`, `string`, `int64` | Go adalah bahasa **strongly typed** — setiap variabel harus punya tipe data yang jelas. `int` = angka biasa, `int64` = angka 64-bit (untuk harga agar tidak overflow). |
| `` `json:"id"` `` | Ini disebut **Struct Tag**. Fungsinya memberi tahu Go: *"Saat kamu mengubah struct ini ke JSON, pakai nama `id` (huruf kecil), bukan `ID` (huruf besar)."* |

**Perbandingan dengan TypeScript:**
```typescript
// TypeScript
interface Product {
    id: number;
    name: string;
    price: number;
    stock: number;
}
```

---

## File 2: `cmd/api/main.go` — Jantung Aplikasi

Ini adalah file pertama yang dijalankan Go. Seperti `server.js` di proyek Express.

### Bagian 1: Import

```go
import (
    "context"
    "fmt"
    "net/http"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"

    // Import internal packages
    "go-api/internal/handlers"
    "go-api/internal/middleware"
    "go-api/internal/repository"
    "go-api/internal/routes"
    "go-api/internal/services"
)
```

| Package | Fungsi | Padanan Express |
|---------|--------|-----------------|
| `context` | Mengontrol pembatalan & timeout operasi | Tidak ada padanan langsung |
| `fmt` | Format & cetak teks (seperti `console.log`) | `console.log` |
| `net/http` | Web server bawaan Go | `express` (tapi bawaan, tanpa install) |
| `os` | Akses environment variables & sistem operasi | `process.env` |
| `pgxpool` | Driver PostgreSQL dengan Connection Pool | `pg` (node-postgres) |
| `godotenv` | Membaca file `.env` | `dotenv` |

### Bagian 2: Load `.env`

```go
if err := godotenv.Load(); err != nil {
    fmt.Println("Peringatan: File .env tidak ditemukan...")
}
```

**Apa yang terjadi:**
1. `godotenv.Load()` membaca file `.env` di root proyek dan memasukkan isinya ke environment sistem.
2. `if err := ...` — Ini pola khas Go. Kita **mendeklarasikan** variabel `err` dan **langsung mengeceknya** dalam satu baris. Jika `err` bukan `nil` (bukan kosong), berarti ada masalah.
3. `fmt.Println()` — Mencetak teks ke terminal (seperti `console.log()`). Kita hanya memberi peringatan, bukan menghentikan server, karena mungkin environment sudah di-set secara global.

**Di Express:**
```javascript
require('dotenv').config(); // Kalau gagal, diam saja
```

### Bagian 3: `fmt.Sprintf()` — Merangkai String

```go
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
)
```

**Ini adalah bintang utama yang Anda tanyakan!**

`fmt.Sprintf()` adalah fungsi untuk **merangkai string dengan memasukkan nilai-nilai** ke dalam template. Huruf `S` di depan `printf` artinya **S**tring — hasilnya dikembalikan sebagai string, bukan langsung dicetak ke terminal.

| Fungsi | Apa yang dilakukan |
|--------|-------------------|
| `fmt.Println("halo")` | Cetak teks ke terminal |
| `fmt.Printf("halo %s", nama)` | Cetak teks **dengan format** ke terminal |
| `fmt.Sprintf("halo %s", nama)` | **Kembalikan** teks berformat sebagai string (tidak dicetak) |
| `fmt.Fprint(w, "halo")` | Tulis teks ke "penulis" (misalnya HTTP response) |

**Format specifier `%s`** artinya: *"taruh string di sini"*.

Jadi baris di atas menghasilkan string seperti ini:
```
host=localhost port=5432 user=postgres password=postgres dbname=go_api_db sslmode=disable
```

String ini disebut **DSN (Data Source Name)** — alamat lengkap menuju database Anda.

**`os.Getenv("DB_HOST")`** membaca nilai dari environment variable. Setelah `godotenv.Load()` dijalankan, variabel `DB_HOST` sudah terisi `localhost` (dari file `.env`).

**Di Node.js:**
```javascript
// Template literal (backtick)
const dsn = `host=${process.env.DB_HOST} port=${process.env.DB_PORT} ...`;
```

### Bagian 4: Koneksi Database (Connection Pool)

```go
dbPool, err := pgxpool.New(context.Background(), dsn)
if err != nil {
    fmt.Printf("Gagal terhubung ke database: %s\n", err)
    os.Exit(1)
}
defer dbPool.Close()
```

**Penjelasan per baris:**

**`pgxpool.New(context.Background(), dsn)`**
- Membuat **Connection Pool** ke PostgreSQL. Pool artinya: bukan 1 koneksi, tapi **kumpulan koneksi** yang bisa dipakai bergantian. Bayangkan antrian kasir di supermarket — ada 5 kasir (koneksi), jadi 5 customer (request) bisa dilayani bersamaan.
- `context.Background()` artinya: *"Jalankan ini tanpa batas waktu."* Context adalah cara Go mengontrol **timeout** dan **pembatalan**. Kita pakai `Background()` karena ini di tahap inisialisasi, belum ada request.

**`os.Exit(1)`**
- Hentikan program sepenuhnya dengan kode error `1`. Jika database tidak bisa dikoneksikan, tidak ada gunanya server berjalan.

**`defer dbPool.Close()`** ⭐ **Konsep penting Go!**
- `defer` artinya: *"Jalankan perintah ini NANTI, saat fungsi `main()` selesai."*
- Ini memastikan koneksi database **pasti** ditutup saat server mati, mencegah "kebocoran koneksi" (connection leak).
- Bayangkan `defer` seperti catatan tempel: "Jangan lupa matikan lampu sebelum pulang." Anda menulisnya sekarang, tapi dikerjakan belakangan.

**Di Express:**
```javascript
const pool = new Pool({
    host: process.env.DB_HOST,
    // ...
});
// Di Node.js, kita harus ingat memanggil pool.end() secara manual
process.on('SIGTERM', () => pool.end());
```

### Bagian 5: Ping Database

```go
if err := dbPool.Ping(context.Background()); err != nil {
    fmt.Printf("Database tidak dapat dijangkau: %s\n", err)
    os.Exit(1)
}
fmt.Println("✅ Koneksi database berhasil!")
```

`Ping()` mengirim sinyal kecil ke database untuk memastikan koneksi benar-benar hidup. Seperti mengetuk pintu sebelum masuk rumah.

### Bagian 6: Dependency Injection (Wiring)

```go
productRepo := repository.NewPostgresProductRepository(dbPool)
productService := services.NewProductService(productRepo)
productHandler := handlers.NewProductHandler(productService)
```

Inilah momen **Dependency Injection** terjadi. Bayangkan ini seperti merakit Lego:
1. **Buat "Kaset Database"** — `productRepo` (berisi koneksi ke PostgreSQL)
2. **Pasang kaset ke "Otak Bisnis"** — `productService` (menerima repo)
3. **Pasang otak ke "Mesin HTTP"** — `productHandler` (menerima service)

Setiap layer **tidak membuat dependensinya sendiri**, tapi **menerima dari luar**. Inilah inti DI.

### Bagian 7: Setup & Start Server

```go
mux := http.NewServeMux()
routes.SetupRoutes(mux, productHandler)

port := os.Getenv("APP_PORT")
if port == "" {
    port = "8080"
}

handler := middleware.Logger(mux)
if err := http.ListenAndServe(":"+port, handler); err != nil {
    fmt.Printf("Gagal menyalakan server: %s\n", err)
}
```

| Kode | Penjelasan |
|------|------------|
| `http.NewServeMux()` | Membuat router baru (seperti `express.Router()`) |
| `routes.SetupRoutes(mux, ...)` | Mendaftarkan semua endpoint ke router |
| `os.Getenv("APP_PORT")` | Ambil port dari `.env`. Jika kosong, pakai default `8080` |
| `middleware.Logger(mux)` | Membungkus router dengan middleware Logger |
| `http.ListenAndServe(":"+port, handler)` | Nyalakan server di port tersebut (seperti `app.listen(8080)`) |

---

## File 3: `routes/routes.go` — Peta Jalan Endpoint

```go
func SetupRoutes(mux *http.ServeMux, productHandler *handlers.ProductHandler) {
    mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Selamat datang di Go Backend!")
    })
    mux.HandleFunc("GET /health", handlers.HealthHandler)

    mux.HandleFunc("GET /products", productHandler.GetProducts)
    mux.HandleFunc("POST /products", productHandler.CreateProduct)
    mux.HandleFunc("GET /products/{id}", productHandler.GetProductByID)
    mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
    mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)
}
```

**`mux.HandleFunc("GET /products/{id}", ...)`**
- Format `"METHOD /path"` adalah fitur baru Go 1.22. Sebelumnya, Go tidak bisa membedakan `GET` dan `POST` di level routing.
- `{id}` adalah **wildcard** — menangkap bagian URL yang dinamis (seperti `:id` di Express).

**`fmt.Fprint(w, "Selamat datang...")`**
- `Fprint` menulis teks ke `w` (HTTP ResponseWriter). Huruf `F` artinya **F**ile/writer — dia menulis ke "tujuan" yang diberikan, bukan ke terminal.

**Di Express:**
```javascript
router.get('/products/:id', productHandler.getProductByID);
```

---

## File 4: `handlers/product.go` — Penerima Request

### Constructor & Struct

```go
type ProductHandler struct {
    service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
    return &ProductHandler{
        service: service,
    }
}
```

**`&ProductHandler{...}`** — Tanda `&` membuat **pointer** (alamat memori). Artinya: kita tidak mengirim salinan data, tapi **alamat** di mana data tersebut disimpan. Ini jauh lebih efisien untuk struct besar.

**`*ProductHandler`** — Tanda `*` artinya "ini adalah pointer ke `ProductHandler`". Jadi fungsi ini mengembalikan **alamat** ke struct, bukan struct-nya langsung.

> [!NOTE]
> **Analogi Pointer:** Bayangkan Anda ingin memberi tahu teman alamat rumah Anda. Anda bisa:
> - Membangun rumah baru yang persis sama untuk dia (tanpa pointer — boros memori)
> - Atau memberikan alamat rumah Anda di kertas (pointer — efisien)

### Method `GetProducts`

```go
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
    products, _ := h.service.GetAll()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}
```

| Kode | Penjelasan |
|------|------------|
| `(h *ProductHandler)` | **Method receiver** — ini yang membuat fungsi menjadi "milik" struct `ProductHandler`. `h` adalah konvensi untuk handler (seperti `this` di JS, tapi eksplisit). |
| `w http.ResponseWriter` | Alat untuk menulis response ke browser (seperti `res` di Express) |
| `r *http.Request` | Data request dari browser (seperti `req` di Express) |
| `products, _ := h.service.GetAll()` | Panggil service. Tanda `_` artinya: *"Saya tahu ada error yang dikembalikan, tapi saya sengaja mengabaikannya."* |
| `w.Header().Set(...)` | Menambahkan header HTTP (seperti `res.setHeader()`) |
| `json.NewEncoder(w).Encode(products)` | Mengubah struct Go menjadi JSON lalu menulis langsung ke response (seperti `res.json(products)`) |

### Method `CreateProduct`

```go
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var newProduct models.Product

    if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
        http.Error(w, "Data tidak valid", http.StatusBadRequest)
        return
    }

    if err := h.service.Create(&newProduct); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newProduct)
}
```

| Kode | Penjelasan |
|------|------------|
| `var newProduct models.Product` | Membuat variabel kosong bertipe `Product`. Seperti menyiapkan wadah kosong. |
| `json.NewDecoder(r.Body).Decode(&newProduct)` | Membaca JSON dari request body dan mengisi wadah `newProduct`. Tanda `&` penting! Artinya: *"Ini alamat wadahnya, silakan isi langsung."* Tanpa `&`, Go akan membuat salinan dan wadah aslinya tetap kosong. |
| `http.Error(w, "...", http.StatusBadRequest)` | Cara cepat mengirim error response dengan status code 400. |
| `w.WriteHeader(http.StatusCreated)` | Mengirim status code 201 Created (bukan default 200 OK). |
| `err.Error()` | Mengubah objek error menjadi string pesan error yang bisa dibaca. |

### Method `GetProductByID`

```go
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
    idString := r.PathValue("id")

    id, err := strconv.Atoi(idString)
    if err != nil {
        http.Error(w, "ID harus berupa angka", http.StatusBadRequest)
        return
    }

    product, err := h.service.GetByID(id)
    // ...
}
```

| Kode | Penjelasan |
|------|------------|
| `r.PathValue("id")` | Mengambil nilai `{id}` dari URL. Misalnya URL `/products/5`, maka hasilnya `"5"` (string!). Di Express: `req.params.id` |
| `strconv.Atoi(idString)` | **A**lphanumeric **to** **I**nteger. Mengubah string `"5"` menjadi angka `5`. Di JS ini otomatis, di Go harus eksplisit. |

---

## File 5: `services/product.go` — Otak Bisnis

```go
type productService struct {
    repo repository.ProductRepository
}

func (s *productService) Create(product *models.Product) error {
    if product.Price < 0 {
        return errors.New("harga tidak boleh negatif")
    }
    if product.Stock < 0 {
        return errors.New("stok tidak boleh negatif")
    }
    if product.Name == "" {
        return errors.New("nama produk wajib diisi")
    }

    return s.repo.Create(product)
}
```

**Layer ini TIDAK tahu soal HTTP** (tidak ada `w`, `r`, atau status code). Dia hanya peduli: *"Apakah data ini masuk akal secara bisnis?"*

| Kode | Penjelasan |
|------|------------|
| `errors.New("...")` | Membuat objek error baru dengan pesan tertentu. Di JS: `new Error("...")` |
| `return s.repo.Create(product)` | Jika semua validasi lolos, serahkan ke repository untuk disimpan ke database. |

---

## File 6: `repository/postgres_product.go` — Jembatan ke Database

### GetAll (SELECT)

```go
func (r *postgresProductRepo) GetAll() ([]models.Product, error) {
    rows, err := r.db.Query(context.Background(),
        "SELECT id, name, price, stock FROM products")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var p models.Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, nil
}
```

| Kode | Penjelasan |
|------|------------|
| `r.db.Query(ctx, sql)` | Menjalankan query SQL yang mengembalikan **banyak baris** (rows). Di Node: `pool.query('SELECT ...')` |
| `defer rows.Close()` | Pastikan "kursor database" ditutup setelah fungsi selesai. Jika tidak, koneksi bisa bocor. |
| `for rows.Next()` | Loop setiap baris hasil query (seperti `for...of` pada array di JS) |
| `rows.Scan(&p.ID, &p.Name, ...)` | Mengambil **setiap kolom** dari baris database dan memasukkannya ke variabel Go. Urutannya harus **sama persis** dengan urutan di `SELECT`. Tanda `&` wajib — karena `Scan` perlu **mengisi** variabel tersebut. |
| `append(products, p)` | Menambahkan elemen ke slice (seperti `array.push(p)` di JS) |

### GetByID (SELECT WHERE)

```go
err := r.db.QueryRow(ctx,
    "SELECT id, name, price, stock FROM products WHERE id = $1",
    id,
).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
```

| Kode | Penjelasan |
|------|------------|
| `QueryRow()` | Seperti `Query()`, tapi hanya mengembalikan **1 baris**. Tidak perlu loop. |
| `$1` | **Parameterized query** — placeholder untuk nilai yang akan diisi. Ini **WAJIB** untuk mencegah SQL Injection! Di Node `pg`: pakai `$1`, `$2`, dst. |
| `pgx.ErrNoRows` | Error khusus dari pgx yang muncul saat query tidak menemukan data apapun. |
| `errors.Is(err, pgx.ErrNoRows)` | Cara aman mengecek tipe error di Go. Lebih baik dari `err == pgx.ErrNoRows` karena bisa menangkap error yang "dibungkus" (wrapped). |

### Create (INSERT RETURNING)

```go
err := r.db.QueryRow(ctx,
    "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id",
    product.Name, product.Price, product.Stock,
).Scan(&product.ID)
```

| Kode | Penjelasan |
|------|------------|
| `RETURNING id` | Fitur PostgreSQL yang mengembalikan kolom tertentu setelah INSERT. Karena ID di-generate oleh database (SERIAL), kita perlu "menangkap" hasilnya. |
| `.Scan(&product.ID)` | Memasukkan ID yang baru di-generate langsung ke struct `product` yang kita punya. |

### Update & Delete (Exec + RowsAffected)

```go
result, err := r.db.Exec(ctx,
    "UPDATE products SET name = $1 ... WHERE id = $4",
    product.Name, ..., product.ID,
)
if result.RowsAffected() == 0 {
    return errors.New("produk tidak ditemukan")
}
```

| Kode | Penjelasan |
|------|------------|
| `Exec()` | Menjalankan query yang **tidak mengembalikan data** (UPDATE, DELETE). Berbeda dari `Query/QueryRow` yang mengembalikan rows. |
| `result.RowsAffected()` | Berapa baris yang berhasil diubah/dihapus. Jika `0`, berarti ID yang dicari tidak ada. |

---

## File 7: `middleware/logger.go` — Pencatat Aktivitas

```go
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        fmt.Printf("[%s] %s %s %s\n",
            time.Now().Format("2006-01-02 15:04:05"),
            r.Method, r.URL.Path, time.Since(start),
        )
    })
}
```

| Kode | Penjelasan |
|------|------------|
| `func Logger(next http.Handler) http.Handler` | Menerima handler berikutnya, mengembalikan handler baru yang "membungkus"-nya. Pola **Decorator**. |
| `time.Now()` | Ambil waktu saat ini (seperti `Date.now()`) |
| `next.ServeHTTP(w, r)` | Serahkan request ke handler berikutnya (seperti `next()` di Express middleware) |
| `time.Now().Format("2006-01-02 15:04:05")` | Format waktu di Go unik — menggunakan **tanggal referensi** `Jan 2, 2006 3:04:05 PM` (bukan `YYYY-MM-DD`). Ini tanggal yang dipilih Go secara konvensi. |
| `time.Since(start)` | Hitung selisih waktu dari `start` sampai sekarang (durasi request) |

---

## Ringkasan Fungsi-Fungsi Penting Go

| Fungsi Go | Padanan JS/Express | Kegunaan |
|-----------|-------------------|----------|
| `fmt.Println()` | `console.log()` | Cetak teks ke terminal |
| `fmt.Printf()` | `console.log()` dengan template literal | Cetak teks berformat ke terminal |
| `fmt.Sprintf()` | Template literal (backtick) | Rangkai string berformat, **kembalikan** sebagai string |
| `fmt.Fprint(w, ...)` | `res.send()` | Tulis teks ke HTTP response |
| `os.Getenv()` | `process.env.KEY` | Baca environment variable |
| `os.Exit(1)` | `process.exit(1)` | Hentikan program |
| `json.NewEncoder(w).Encode()` | `res.json()` | Kirim JSON response |
| `json.NewDecoder(r.Body).Decode()` | `req.body` (dengan body-parser) | Baca JSON dari request body |
| `strconv.Atoi()` | `parseInt()` | Ubah string ke integer |
| `errors.New()` | `new Error()` | Buat error baru |
| `errors.Is()` | `err instanceof ErrorType` | Cek tipe error |
| `context.Background()` | — | Buat context kosong tanpa timeout |
| `defer` | — | Jadwalkan kode untuk dijalankan saat fungsi selesai |
| `r.PathValue("id")` | `req.params.id` | Ambil URL parameter |
| `http.Error(w, msg, code)` | `res.status(400).send(msg)` | Kirim error response |
| `w.WriteHeader(201)` | `res.status(201)` | Set status code |

---

## Simbol-Simbol Go yang Wajib Dipahami

| Simbol | Arti | Contoh |
|--------|------|--------|
| `&` | "Berikan alamat memori" (reference) | `&product` → alamat product |
| `*` | "Ikuti alamat ini ke datanya" (dereference) atau "tipe pointer" | `*models.Product` → pointer ke Product |
| `:=` | Deklarasi + assignment sekaligus (shorthand) | `name := "Kopi"` |
| `_` | "Saya sengaja mengabaikan nilai ini" | `products, _ := service.GetAll()` |
| `...` | Spread operator (membongkar slice) | `append(a[:i], a[i+1:]...)` |
