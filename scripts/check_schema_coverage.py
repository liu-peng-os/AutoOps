from pathlib import Path
import re
import sys


ROOT = Path(__file__).resolve().parents[1]
MYSQL_SQL = ROOT / "docker" / "mysql" / "devops001.sql"
MIGRATIONS_DIR = ROOT / "api" / "migrations"


def read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8", errors="replace")


def main() -> int:
    old_sql = read_text(MYSQL_SQL)
    old_tables = sorted(set(re.findall(r"^CREATE TABLE `([^`]+)`", old_sql, re.M)))

    migration_sql = "\n".join(
        read_text(path) for path in sorted(MIGRATIONS_DIR.glob("*.sql"))
    )
    migration_tables = sorted(
        set(re.findall(r"CREATE TABLE IF NOT EXISTS\s+([a-zA-Z0-9_]+)", migration_sql))
    )

    missing = [table for table in old_tables if table not in migration_tables]
    extra = [table for table in migration_tables if table not in old_tables]

    print(f"old_tables={len(old_tables)}")
    print(f"migration_tables={len(migration_tables)}")
    print(f"missing={len(missing)}")
    if missing:
        print("missing_tables=" + ",".join(missing))
    print(f"extra={len(extra)}")
    if extra:
        print("extra_tables=" + ",".join(extra))

    return 1 if missing or extra else 0


if __name__ == "__main__":
    sys.exit(main())
