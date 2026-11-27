# EventHub — Sistem Pengelolaan Event Kampus

Backend untuk membuat event kampus, mengelola pendaftaran, dan mengatur peserta kegiatan. Dibangun dengan **Go** dan **Clean Architecture** mengikuti modul praktikum, dengan tiga penyimpanan: **PostgreSQL**, **Redis**, dan **MongoDB**.

## Deskripsi

EventHub membantu panitia kegiatan dan mahasiswa mengelola event kampus (seminar, workshop, lomba) secara online.

Alurnya: Admin membuat event beserta kategori dan kuotanya. Peserta mencari event, lalu mengajukan pendaftaran dan menerima **kode tiket**. Panitia menyetujui atau menolak pendaftaran, dan setiap perubahan status mengirim **notifikasi otomatis** ke peserta secara asinkron. Pada hari pelaksanaan, panitia melakukan **check-in** peserta memakai kode tiket. Setelah event selesai, peserta yang hadir dapat memberi **ulasan dan rating**. Seluruh aksi penting tercatat pada **activity log**.

Sistem menjamin satu peserta hanya punya satu pendaftaran aktif per event, dan jumlah pendaftaran yang disetujui tidak pernah melebihi kuota event meskipun ada banyak pendaftar bersamaan.

## Role pada Sistem

| Role | Hak akses utama |
| --- | --- |
| **Admin** | Mengelola data event & kategori, melihat activity log |
| **Panitia** | Menyetujui / menolak pendaftaran, melihat daftar peserta, melakukan check-in kehadiran |
| **Peserta** | Mencari & melihat event, mendaftar, membatalkan pendaftaran, melihat riwayat & notifikasi, memberi ulasan |

## Fitur

Status: ✅ selesai · ⏳ direncanakan

### A. Autentikasi & Sesi

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 1 | Registrasi akun | Semua | ⏳ |
| 2 | Login (JWT *access token* + *refresh token*) | Semua | ⏳ |
| 3 | Perpanjang sesi lewat *refresh token* | Semua | ⏳ |
| 4 | Logout (token lama di-*blacklist*) | Semua | ⏳ |
| 5 | Lihat profil sendiri | Semua | ⏳ |
| 6 | Pembatasan percobaan login (*rate limit*) | Sistem | ⏳ |

### B. Katalog Event

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 7 | Kelola data event (tambah, lihat, ubah, hapus) | Admin | ✅ |
| 8 | Lihat daftar & detail event | Semua | ✅ |
| 9 | Pencarian, filter & pagination event | Semua | ⏳ |
| 10 | Kelola kategori event | Admin | ⏳ |

### C. Pendaftaran

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 11 | Ajukan pendaftaran ke event (dapat kode tiket) | Peserta | ⏳ |
| 12 | Setujui / tolak pendaftaran | Panitia | ⏳ |
| 13 | Batalkan pendaftaran | Peserta | ⏳ |
| 14 | Riwayat pendaftaran sendiri | Peserta | ⏳ |
| 15 | Lihat daftar peserta per event | Panitia | ⏳ |

### D. Kehadiran Peserta

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 16 | Check-in peserta memakai kode tiket | Panitia | ⏳ |
| 17 | Rekap kehadiran per event | Panitia | ⏳ |

### E. Notifikasi

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 18 | Notifikasi otomatis (asinkron) saat status pendaftaran berubah | Peserta | ⏳ |
| 19 | Lihat & tandai sudah dibaca | Semua | ⏳ |

### F. Ulasan Event

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 20 | Beri ulasan & rating event (khusus peserta yang hadir) | Peserta | ⏳ |
| 21 | Lihat ulasan & rata-rata rating event | Semua | ⏳ |

### G. Audit

| No | Fitur | Role | Status |
| --- | --- | --- | --- |
| 22 | Lihat activity log seluruh aksi penting | Admin | ⏳ |

### Di luar scope (pengembangan lanjutan)

Kelola user oleh Admin (ubah role, aktif/nonaktif), *waiting list* saat kuota penuh, export daftar peserta ke CSV, dashboard statistik, dan sertifikat peserta.

## Pembagian Penyimpanan Data

Tiga penyimpanan dipakai sesuai karakter datanya, bukan sekadar ditempelkan:

| Penyimpanan | Menyimpan | Alasan dipilih |
| --- | --- | --- |
| **PostgreSQL** | `users`, `event_categories`, `events`, `registrations`, `attendances`, `notifications` | Data relasional & transaksional. Butuh jaminan konsistensi: kuota tidak boleh terlampaui, satu pendaftaran aktif per user per event. |
| **Redis** | Cache event, *refresh token* & *blacklist* JWT, penghitung *rate limit*, *lock* kuota, antrean notifikasi | Akses sangat cepat, TTL otomatis (token & cache kedaluwarsa sendiri), operasi atomik untuk *lock*, sekaligus *broker* antrean. |
| **MongoDB** | `event_reviews`, `activity_logs` | Dokumen berskema fleksibel dan sifatnya *append-heavy*. Ulasan bisa punya struktur bersarang, dan activity log menyimpan cuplikan data sebelum/sesudah yang bentuknya berbeda tiap entitas. |

## Rancangan Database

### PostgreSQL

| Tabel | Kolom utama | Relasi |
| --- | --- | --- |
| `users` | id, name, email (unik), password_hash, role, is_active, created_at | Satu user punya banyak pendaftaran |
| `event_categories` | id, name, slug (unik), created_at | Satu kategori punya banyak event |
| `events` | id, category_id, created_by, title, description, location, start_time, end_time, quota, is_active, created_at | Milik satu kategori, punya banyak pendaftaran |
| `registrations` | id, user_id, event_id, ticket_code (unik), status, created_at, updated_at | Milik satu user dan satu event |
| `attendances` | id, registration_id (unik), checked_in_by, checked_in_at | Satu pendaftaran maksimal satu kehadiran |
| `notifications` | id, user_id, type, title, body, is_read, created_at | Milik satu user |

**Status pendaftaran:** `pending`, `approved`, `rejected`, `cancelled`.

**Aturan bisnis yang dijaga di tingkat basis data dan aplikasi:**

- Indeks unik *partial* pada `registrations (user_id, event_id)` untuk status `pending` dan `approved` — mencegah pendaftaran ganda.
- Jumlah `registrations` berstatus `approved` tidak boleh melebihi `events.quota` — dijaga dengan *lock* Redis dan dicek ulang di dalam transaksi.
- Check-in hanya boleh untuk pendaftaran berstatus `approved`.
- Ulasan hanya boleh dari peserta yang tercatat hadir pada event yang sudah selesai.

### MongoDB

| Koleksi | Field utama | Keterangan |
| --- | --- | --- |
| `event_reviews` | `_id`, `event_id`, `user_id`, `user_name`, `rating` (1–5), `comment`, `created_at`, `updated_at` | Satu ulasan per peserta per event |
| `activity_logs` | `_id`, `actor_id`, `actor_role`, `action`, `entity` { type, id }, `before`, `after`, `note`, `request_id`, `ip`, `created_at` | Catatan audit semua aksi penting |

### Redis

| Pola kunci | Isi | TTL |
| --- | --- | --- |
| `event:list:<hash-query>` | Hasil daftar event per kombinasi filter | 60 detik |
| `event:detail:<id>` | Detail satu event | 5 menit |
| `refresh:<user_id>:<jti>` | *Refresh token* aktif milik user | sesuai masa refresh token |
| `blacklist:<jti>` | Penanda *access token* yang sudah logout | sisa masa access token |
| `ratelimit:login:<email>` | Penghitung percobaan login gagal | 15 menit |
| `lock:event:<id>:register` | *Lock* saat memproses persetujuan pendaftaran | 5 detik |
| antrean `notification:send` | *Job* pengiriman notifikasi (Asynq) | — |

### ERD

> _TODO: lampirkan gambar ERD (dbdiagram.io / draw.io) di `docs/erd.png`._

## Daftar Endpoint API

Semua endpoint diawali `/api/v1`. Tanda 🔒 berarti wajib *access token*.

| Method | Endpoint | Akses |
| --- | --- | --- |
| `GET` | `/health` | Publik |
| `POST` | `/auth/register` | Publik |
| `POST` | `/auth/login` | Publik |
| `POST` | `/auth/refresh` | Publik |
| `POST` | `/auth/logout` | 🔒 Semua |
| `GET` | `/auth/me` | 🔒 Semua |
| `GET` | `/categories` | Publik |
| `POST` `PUT` `DELETE` | `/categories` · `/categories/:id` | 🔒 Admin |
| `GET` | `/events` (`?q=&category=&from=&to=&page=&limit=`) | Publik |
| `GET` | `/events/:id` | Publik |
| `POST` `PUT` `DELETE` | `/events` · `/events/:id` | 🔒 Admin |
| `POST` | `/events/:id/registrations` | 🔒 Peserta |
| `GET` | `/events/:id/registrations` | 🔒 Panitia |
| `GET` | `/registrations/me` | 🔒 Peserta |
| `PATCH` | `/registrations/:id/approve` · `/reject` | 🔒 Panitia |
| `PATCH` | `/registrations/:id/cancel` | 🔒 Peserta |
| `POST` | `/events/:id/attendances` | 🔒 Panitia |
| `GET` | `/events/:id/attendances` | 🔒 Panitia |
| `GET` | `/notifications` | 🔒 Semua |
| `PATCH` | `/notifications/:id/read` | 🔒 Semua |
| `POST` | `/events/:id/reviews` | 🔒 Peserta |
| `GET` | `/events/:id/reviews` | Publik |
| `GET` | `/activity-logs` | 🔒 Admin |

## Struktur Proyek (Clean Architecture)

Mengikuti penamaan layer dari modul praktikum — **tanpa layer `controller` / `handler`**.

```
eventhub/
├── main.go                            # Composition root: pilih role api | worker | all
├── worker.go                          # Composition root worker (Asynq consumer)
├── app/                               # Core: domain + application
│   ├── model/                         # Entity & DTO request/response
│   │   ├── user.go
│   │   ├── auth.go
│   │   ├── category.go
│   │   ├── event.go
│   │   ├── registration.go
│   │   ├── attendance.go
│   │   ├── notification.go
│   │   └── review.go
│   ├── repository/                    # Akses data
│   │   ├── common.go                  #   error sentinel bersama
│   │   ├── tx.go                      #   transaction manager
│   │   ├── user_repository.go         #   PostgreSQL
│   │   ├── token_repository.go        #   Redis (refresh token & blacklist)
│   │   ├── category_repository.go     #   PostgreSQL
│   │   ├── event_repository.go        #   PostgreSQL
│   │   ├── event_cache_repository.go  #   Redis (cache daftar & detail)
│   │   ├── registration_repository.go
│   │   ├── attendance_repository.go
│   │   ├── notification_repository.go
│   │   ├── review_repository.go       #   MongoDB
│   │   └── activity_log_repository.go #   MongoDB
│   └── service/                       # Use case + business rules
│       ├── auth_service.go
│       ├── category_service.go
│       ├── event_service.go
│       ├── registration_service.go
│       ├── registration_rules.go
│       ├── attendance_service.go
│       ├── notification_service.go
│       ├── notification_worker.go     #   consumer antrean notifikasi
│       ├── review_service.go
│       └── activity_log_service.go
├── route/                             # Endpoint + injeksi dependency
│   ├── route.go                       #   Setup + Dependencies
│   ├── auth_routes.go
│   ├── event_routes.go
│   ├── registration_routes.go
│   ├── attendance_routes.go
│   ├── notification_routes.go
│   └── review_routes.go
├── middleware/                        # HTTP middleware
│   ├── middleware.go                  #   recover, CORS, helmet, request logger
│   ├── auth.go                        #   RequireAuth (JWT)
│   ├── authorize.go                   #   RequireRole (RBAC)
│   └── ratelimit.go                   #   pembatas percobaan login
├── helper/                            # Lintas-layer
│   ├── response.go                    #   WebResponse / envelope
│   ├── errors.go                      #   AppError
│   ├── jwt.go                         #   JWTManager + CurrentUser
│   ├── validator.go
│   ├── pagination.go
│   └── ticket.go                      #   generator kode tiket
├── config/
│   ├── config.go                      #   env, NewApp (Fiber) + ErrorHandler
│   └── queue.go                       #   RedisOpt, NewWorkerServer (Asynq)
├── database/
│   ├── postgres.go
│   ├── redis.go
│   └── mongo.go
├── migrations/                        # skema SQL 001..00n
├── docs/                              # ERD & dokumentasi
├── Dockerfile
├── docker-compose.yml                 # app + postgres + redis + mongo
├── .dockerignore
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## Tech Stack

- **Bahasa:** Go
- **Web framework:** Fiber v2
- **Basis data:** PostgreSQL (pgx), Redis (go-redis), MongoDB (mongo-driver)
- **Antrean:** Asynq (broker Redis)
- **Autentikasi:** JWT (access + refresh token), bcrypt
- **Kontainerisasi:** Docker & Docker Compose
- **Arsitektur:** Clean Architecture (model → repository → service → route)

## Cara Menjalankan

### 1. Jalankan service pendukung

```bash
cp .env.example .env            # sesuaikan konfigurasi
docker compose up -d db redis mongo
```

### 2. Jalankan aplikasi

```bash
go mod tidy
go run . -role=all              # API + worker notifikasi dalam satu proses
```

Atau pisahkan prosesnya:

```bash
go run . -role=api              # hanya API
go run . -role=worker           # hanya worker notifikasi
```

API berjalan di `http://localhost:8080` (sesuaikan dengan `.env`). Cek dengan:

```bash
curl http://localhost:8080/health
```

### Menjalankan semuanya dengan Docker

```bash
docker compose up --build
```

> Catatan: migrasi SQL dijalankan otomatis oleh PostgreSQL saat volume data masih baru. Jika skema berubah, jalankan `docker compose down -v` lalu `up` lagi.

## Anggota Kelompok

Pembagian dibuat per fitur agar ketiganya bisa berjalan paralel. Anggota yang lebih berpengalaman juga bertugas membimbing dan mereview pekerjaan anggota lain.

| Anggota | Pengalaman | Fitur yang dipegang |
| --- | --- | --- |
| Ziyaa Danil Mubarok | Terbiasa Git & Go | **Git Master** · Autentikasi & Sesi (Redis) · Notifikasi & Antrean (Asynq) |
| Syihaabuddin Irfan Ahmadi | Pernah membuat REST API | Skema & Migrasi · Pendaftaran Event (lock kuota) · Kehadiran Peserta |
| Muhamad Ainul Yaqin | Masih belajar Go | Katalog Event: kategori, pencarian & cache (Redis) · Ulasan & Activity Log (MongoDB) · Dokumentasi |
