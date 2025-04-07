# DBox - DB Toolbox

**A simple migration manager written in Go.**

Originally built for my own Go backend project, so the file structure might feel a bit custom.

### Built-in support for MySQL, PostgreSQL, and SQLite
---

## ✨ Features

1. `./dbox make [name]` or `./dbox create [name]` — create a new migration
2. `./dbox up` or `./dbox migrate` — run all pending migrations
3. `./dbox migrate --pretend` or `./dbox rollback --pretend` — show SQL that would run, without touching the DB (`-p` for short)
4. `./dbox down` or `./dbox rollback` — roll back the last migration
5. `./dbox init` — initialize the database and `.env` file
6. `./dbox clean` — remove migration records in the DB with no matching folder
7. `./dbox refresh` — roll back everything and run all migrations from scratch
8. `./dbox status` - view the status of all migrations

---

## ⚙️ How it works

Instead of a single file per migration, DBox creates a **folder** for each one  
Each folder contains two files: `up.sql` and `down.sql`

- `up.sql` → raw SQL that runs when you migrate
- `down.sql` → SQL that runs when you roll back

No extra syntax or parsing — just clean SQL

---

## 🛠️ Setup

**Make sure you have [Go installed](https://go.dev/doc/install)**

Then:

```bash
git clone https://github.com/your-username/dbox.git
cd dbox
go build dbox.go
```
If anything's missing, install dependencies:
```
go mod tidy
```
That's it.

### Enjoy!
