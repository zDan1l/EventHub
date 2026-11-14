# EventHub — Sistem Pengelolaan Event Kampus

Backend sederhana untuk membuat event kampus, mengelola pendaftaran, dan mengatur peserta kegiatan. Dibangun dengan **Go** dan **Clean Architecture** mengikuti modul praktikum.

## Deskripsi

EventHub membantu panitia kegiatan dan mahasiswa mengelola event kampus (seminar, workshop, lomba) secara online. Panitia membuat event dengan kuota tertentu, peserta mengajukan pendaftaran, lalu panitia menyetujui atau menolaknya. Sistem mencegah pendaftaran ganda (user yang sama pada event yang sama) dan mencegah pendaftaran yang disetujui melebihi kuota event.

## Role pada Sistem

| Role | Hak akses utama |
| --- | --- |
| **Admin** | Mengelola user dan data event (tambah, ubah, hapus, aktif/nonaktif) |
| **Panitia** | Menyetujui / menolak pendaftaran, melihat & mengelola daftar peserta event |
| **Peserta** | Melihat event, mendaftar, membatalkan pendaftaran, melihat riwayat |

## Fitur

| No | Fitur | Role |
| --- | --- | --- |
| 1 | Registrasi dan login | Semua |
| 2 | Kelola data event (tambah, lihat, ubah, hapus) | Admin |
| 3 | Lihat daftar & detail event | Semua |
| 4 | Daftar (registrasi) ke event | Peserta |
| 5 | Setujui / tolak pendaftaran | Panitia |
| 6 | Batalkan pendaftaran | Peserta |
| 7 | Riwayat pendaftaran | Peserta, Panitia |

## Rancangan Database

| Tabel | Kolom utama | Relasi |
| --- | --- | --- |
| `users` | id, name, email, password_hash, role, created_at | Satu user punya banyak pendaftaran |
| `events` | id, title, description, location, start_time, end_time, quota, is_active, created_at | Satu event punya banyak pendaftaran |
| `registrations` | id, user_id, event_id, status, created_at | Milik satu user dan satu event |
| `registration_logs` | id, registration_id, actor_id, from_status, to_status, note, created_at | Mencatat setiap perubahan status pendaftaran |

**Status pendaftaran:** `pending`, `approved`, `rejected`, `cancelled`.

**Aturan:** satu user hanya boleh punya satu pendaftaran aktif per event; jumlah pendaftaran berstatus `approved` tidak boleh melebihi `quota` event.

### ERD

> _TODO: lampirkan gambar ERD (dbdiagram.io / draw.io) di sini, mis. `docs/erd.png`._

## Struktur Proyek (Clean Architecture)

Mengikuti penamaan layer dari modul praktikum — **tanpa layer `controller` / `handler`**.

```
eventhub-backend/
├── main.go                 # Composition root: setup config, DB, route
├── app/                    # Core: domain + application
│   ├── model/              # Entity & DTO request/response
│   │   ├── user.go
│   │   ├── event.go
│   │   └── registration.go
│   ├── repository/         # Akses data (interface + impl PostgreSQL)
│   │   ├── common.go
│   │   ├── user_repository.go
│   │   ├── event_repository.go
│   │   └── registration_repository.go
│   └── service/            # Use case + business rules
│       ├── auth_service.go
│       ├── event_service.go
│       └── registration_service.go
├── route/route.go          # Endpoint + injeksi dependency
├── middleware/             # Auth (JWT) & authorization (role)
├── helper/                 # response, error, jwt, validator
├── config/                 # env & setup app
├── database/               # koneksi PostgreSQL
├── migrations/             # skema SQL
├── Dockerfile              # image untuk service backend (Go)
├── docker-compose.yml      # orkestrasi service: app + PostgreSQL
├── .dockerignore
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## Tech Stack

- **Bahasa:** Go
- **Database:** PostgreSQL
- **Kontainerisasi:** Docker & Docker Compose (mengelola service app + database)
- **Arsitektur:** Clean Architecture (model → repository → service → route)

## Cara Menjalankan

> _Detail dilengkapi pada tahap boilerplate (Langkah 3)._

### Dengan Docker (disarankan)

Menjalankan backend + PostgreSQL sekaligus, tanpa perlu install Go/DB di laptop.

```bash
cp .env.example .env     # sesuaikan konfigurasi
docker compose up --build
```

API berjalan di `http://localhost:8080` (sesuaikan dengan `.env`).

### Tanpa Docker (lokal)

Membutuhkan Go dan PostgreSQL terpasang di laptop.

```bash
cp .env.example .env     # sesuaikan konfigurasi database
go mod tidy
go run main.go
```

## Anggota Kelompok


| Anggota | Pengalaman | Peran |
| --- | --- | --- |
| Ziyaa Danil Mubarok | Terbiasa Git & Go | **Git Master** + logika bisnis (service) |
| Syihaabuddin Irfan Ahmadi| Pernah membuat REST API | Database + akses data (repository, migrations) |
| Muhamad Ainul Yaqin | Masih belajar Go | Entity (model) + routing API (route) + dokumentasi |
