# Panduan Optimasi Database & Pencarian (Database Optimization Guide)

Dokumen ini menjelaskan strategi untuk meningkatkan performa database PostgreSQL, khususnya terkait pencarian (Full-Text Search) dan indexing.

---

## 1. Memahami Indexing

Index adalah struktur data yang meningkatkan kecepatan operasi pengambilan data pada tabel database dengan mengorbankan penulisan tambahan dan ruang penyimpanan.

### Jenis Index yang Kami Gunakan:

1.  **B-Tree Index (Default)**
    *   **Kegunaan**: Digunakan untuk perbandingan kesetaraan (`=`) dan rentang (`<`, `>`, `BETWEEN`).
    *   **Otomatis**: GORM otomatis membuat ini untuk primary key (`id`), foreign key, dan kolom dengan tag `uniqueIndex` (seperti `email`).
    *   **Contoh**: Mencari user berdasarkan `email` atau rentang `created_at`.

2.  **GIN (Generalized Inverted Index)**
    *   **Kegunaan**: Sangat efisien untuk **Full-Text Search** (pencarian teks) dan tipe data JSONB/Array.
    *   **Penting**: Tanpa index GIN, pencarian teks (`to_tsvector`) akan melakukan *sequential scan* yang sangat lambat pada data besar.

---

## 2. Implementasi Full-Text Search (Pencarian Cerdas)

Fitur "Smart Search" di `UserRepository` menggunakan kemampuan Full-Text Search PostgreSQL.

### Query Dasar
```sql
to_tsvector('indonesian', name || ' ' || email) @@ to_tsquery('indonesian', 'query_prefix:*')
```
*   `to_tsvector`: Mengubah teks menjadi vektor kata-kata yang dapat di-index. Kami menggabungkan `name` dan `email`.
*   `to_tsquery`: Mengubah input user menjadi query pencarian.

### Menambahkan Index untuk Performa (Wajib!)

Agar query di atas cepat, kita **harus** membuat index yang sesuai dengan ekspresi `to_tsvector` tersebut. Ini dilakukan di `migrations/user.go`.

#### Raw SQL Migration:
```go
// Pastikan index ini ada agar pencarian ngebut!
db.Exec("CREATE INDEX IF NOT EXISTS idx_users_fulltext ON users USING GIN (to_tsvector('indonesian', name || ' ' || email));")
```

**Kenapa Manual?**
GORM belum mendukung pembuatan index GIN berbasis ekspresi (function-based index) secara native via struct tags dengan mudah.

---

## 3. Menganalisa Performa Query (Debugging)

Jika API terasa lambat, gunakan teknik ini untuk mencari tahu penyebabnya.

### Menggunakan `EXPLAIN ANALYZE`

Jalankan query SQL langsung di database client (DBeaver, pgAdmin, atau terminal) dengan awalan `EXPLAIN ANALYZE`.

**Contoh Query Lambat (Tanpa Index):**
```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE to_tsvector('indonesian', name) @@ to_tsquery('budi');
```
*Output mungkin mengandung*: `Seq Scan on users ... (actual time=0.050..15.420 rows=5 loops=1)` -> **BURUK** (Scan semua baris).

**Contoh Query Cepat (Dengan Index GIN):**
```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE to_tsvector('indonesian', name) @@ to_tsquery('budi');
```
*Output*: `Bitmap Heap Scan on users ... (actual time=0.015..0.080 rows=5 loops=1)` -> **BAGUS** (Menggunakan Index).

### Logging GORM

Di development (`make dev`), GORM dikonfigurasi untuk log query yang lambat (>200ms). Perhatikan terminal Anda:

```text
[120.34ms] [rows:5] SELECT * FROM "users" ...
```

Jika waktu eksekusi tinggi padahal datanya sedikit, cek apakah index sudah digunakan.

---

## 4. Best Practices

1.  **Gunakan Pagination**: Jangan pernah me-return semua data (`Select *`). Selalu gunakan `Limit` dan `Offset` (sudah dihandle oleh `utils/pagination.go`).
2.  **Index Foreign Keys**: Pastikan kolom foreign key (seperti `role_id`, `user_id`) di-index. GORM biasanya melakukan ini otomatis.
3.  **Hindari `LIKE '%...%'`**: Pencarian `LIKE` dengan wildcard di depan (`%budi`) **TIDAK BISA** menggunakan index B-Tree biasa. Gunakan Full-Text Search sebagai gantinya.
4.  **Vacuum**: PostgreSQL perlu maintenance rutin (`VACUUM ANALYZE`) untuk memperbarui statistik query planner, terutama setelah banyak delete/update.

---

## 5. Menambahkan Index Baru

Jika Anda menambahkan fitur baru yang butuh pencarian cepat:

1.  Tentukan kolom yang sering dicari.
2.  Jika pencarian teks kompleks: Buat **GIN Index** via raw SQL di `migrations/migrate.go`.
3.  Jika pencarian exact match/range: Tambahkan tag `` `gorm:"index"` `` di struct Entity.
