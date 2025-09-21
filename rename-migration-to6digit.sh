#!/bin/bash

# rename-migrations-fix-gap.sh
# Purpose: Fix migration numbering gaps and remove version 0
# Assumes all files are 6-digit: 000001_... up to 000010_...

MIGRATIONS_DIR="./internal/db/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "❌ Error: Directory '$MIGRATIONS_DIR' not found."
  exit 1
fi

cd "$MIGRATIONS_DIR" || exit 1

echo "🔍 Cleaning up migration files..."

# Remove any 000000 files
rm -f 000000_*.up.sql 000000_*.down.sql

# Find all .up.sql and .down.sql files, sort numerically
for file in $(ls *.up.sql *.down.sql | sort -n); do
  [[ ! -f "$file" ]] && continue
  if [[ $file =~ ^([0-9]{6})_(.+)\.(up|down)\.sql$ ]]; then
    num="${BASH_REMATCH[1]}"
    rest="${BASH_REMATCH[2]}"
    ext="${BASH_REMATCH[3]}"

    # Convert numeric value to integer
    n=$(printf "%d" "$num")

    # If number >= 9, shift it down by 1 (to fill gap at 000008)
    if [ "$n" -ge 9 ]; then
      new_n=$((n - 1))
      new_num=$(printf "%06d" "$new_n")
      new_name="${new_num}_${rest}.${ext}.sql"
      echo "🔄 Renaming: $file → $new_name"
      mv "$file" "$new_name"
    fi
  else
    echo "⚠️ Skipping (unexpected format): $file"
  fi
done

echo "✅ Done! All migrations renumbered to remove gap and eliminate 000000."