# عملیات تغییر

tmq تغییر در جای فایل‌های TOML را پشتیبانی می‌کند و اجازهٔ تنظیم مقادیر جدید یا حذف کلیدهای موجود را می‌دهد.

## تنظیم مقادیر

### تخصیص مقدار ساده
```bash
# Set a simple string value
tmq '.version = "2.0.0"' -i config.toml

# Set a number
tmq '.port = 8080' -i config.toml

# Set a boolean
tmq '.enabled = true' -i config.toml
```

### مقادیر جدول تودرتو
```toml
# Before
[database]
host = "oldhost"
```

```bash
# Update nested value
tmq '.database.host = "newhost"' -i config.toml
```

```toml
# After
[database]
host = "newhost"
```

### ایجاد کلیدهای جدید
```bash
# Add new root-level key
tmq '.new_key = "new_value"' -i config.toml

# Add nested key
tmq '.database.pool_size = 10' -i config.toml
```

### تودرتوی عمیق
```bash
# Create deep nested structure
tmq '.app.cache.redis.ttl = 3600' -i config.toml

# This creates:
# [app.cache.redis]
# ttl = 3600
```

### مقادیر آرایه
```bash
# Set array of strings
tmq '.tags = ["web", "api", "prod"]' -i config.toml

# Set array of numbers
tmq '.ports = [8080, 8443, 9000]' -i config.toml
```

### آبجکت‌های پیچیده
```bash
# Set an inline table
tmq '.credentials = { username = "admin", password = "secret" }' -i config.toml

# Set nested object
tmq '.database = { host = "localhost", port = 5432 }' -i config.toml
```

## عملیات حذف

### حذف کلیدهای ریشه
```bash
# Delete a top-level key
tmq 'del(.obsolete_key)' -i config.toml
```

### حذف کلیدهای تودرتو
```bash
# Delete from nested table
tmq 'del(.database.old_setting)' -i config.toml
```

### حذف عناصر آرایه
```bash
# Delete specific array index
tmq 'del(.servers[1])' -i config.toml
```

## حالت Dry Run

### پیش‌نمایش تغییرات
```bash
# See what would be changed without modifying the file
tmq '.version = "3.0.0"' --dry-run config.toml

# Preview deletion
tmq 'del(.obsolete_key)' --dry-run config.toml
```

### تغییرات امن
```bash
# Always test with dry-run first
tmq '.database.host = "prod-db"' --dry-run config.toml

# Then apply if it looks correct
tmq '.database.host = "prod-db"' -i config.toml
```

## مثال‌های پیشرفته

### به‌روزرسانی پیکربندی
```toml
# config.toml before
[app]
version = "1.0.0"
debug = true

[database]
host = "dev-db"
port = 5432
```

```bash
# Update for production deployment
tmq '.app.version = "1.1.0"' -i config.toml
tmq '.app.debug = false' -i config.toml
tmq '.database.host = "prod-db"' -i config.toml
```

```toml
# config.toml after
[app]
version = "1.1.0"
debug = false

[database]
host = "prod-db"
port = 5432
```

### پیکربندی محیط‌محور
```bash
# Development settings
tmq '.database.host = "localhost"' -i config.toml
tmq '.debug = true' -i config.toml

# Production settings
tmq '.database.host = "prod.example.com"' -i config.toml
tmq '.debug = false' -i config.toml
```

### عملیات پاکسازی
```bash
# Remove deprecated settings
tmq 'del(.legacy_feature)' -i config.toml
tmq 'del(.old_database_url)' -i config.toml

# Remove test users
tmq 'del(.test_users)' -i config.toml
```

## مدیریت خطا

### مسیرهای ناموجود
```bash
# Setting non-existent parent creates the structure
tmq '.new.deep.key = "value"' -i config.toml
# Creates: [new.deep]
#          key = "value"
```

### تعارض نوع
```bash
# Overwriting different types is allowed
tmq '.value = "string"' -i config.toml  # was a number
tmq '.value = 42' -i config.toml        # was a string
```

### عملیات نامعتبر
```bash
# Invalid key names
tmq '.invalid key = "value"' -i config.toml
# Error: invalid set expression

# Missing quotes for strings
tmq '.name = John' -i config.toml
# Error: invalid set expression
```

## استراتژی بکاپ

### بکاپ دستی
```bash
# Always backup before modification
cp config.toml config.toml.backup

# Make changes
tmq '.version = "2.0.0"' -i config.toml

# Verify
tmq '.version' config.toml
```

### بکاپ اسکریپتی
```bash
#!/bin/bash
CONFIG_FILE="config.toml"
BACKUP_FILE="${CONFIG_FILE}.backup.$(date +%Y%m%d_%H%M%S)"

cp "$CONFIG_FILE" "$BACKUP_FILE"
echo "Backup created: $BACKUP_FILE"

# Make changes
tmq '.version = "2.0.0"' -i "$CONFIG_FILE"

# Verify
if tmq '.version' "$CONFIG_FILE" >/dev/null; then
    echo "Update successful"
else
    echo "Update failed, restoring backup"
    cp "$BACKUP_FILE" "$CONFIG_FILE"
fi
```

## ملاحظات عملکرد

- تغییرات در جای فایل کارآمد است — فقط بخش‌های تغییرکرده بازنویسی می‌شوند
- فایل‌های بزرگ ممکن است به‌خاطر بازنویسی کامل بیشتر طول بکشند
- برای فایل‌های بزرگ یا حیاتی از dry-run استفاده کنید

## بهترین روش‌ها

### اعتبارسنجی
```bash
# Validate after modifications
tmq '.' config.toml > /dev/null || echo "Invalid TOML after modification"
```

### عملیات اتمیک
```bash
# Use temporary files for critical updates
TEMP_FILE=$(mktemp)
cp config.toml "$TEMP_FILE"

tmq '.critical_setting = "new_value"' -i "$TEMP_FILE"

# Validate
if tmq '.' "$TEMP_FILE" >/dev/null; then
    mv "$TEMP_FILE" config.toml
    echo "Update successful"
else
    rm "$TEMP_FILE"
    echo "Update failed - invalid TOML"
fi
```

### کنترل نسخه
```bash
# Commit before and after modifications
git add config.toml
git commit -m "Update configuration via tmq"

# Make changes
tmq '.version = "2.0.0"' -i config.toml

git add config.toml
git commit -m "Bump version to 2.0.0"
```

### یکپارچه‌سازی با اسکریپت
```bash
#!/bin/bash
set -e

# Function to safely update config
update_config() {
    local key="$1"
    local value="$2"
    local file="$3"

    echo "Updating $key = $value in $file"

    # Dry run first
    if tmq "$key = $value" --dry-run "$file" >/dev/null; then
        tmq "$key = $value" -i "$file"
        echo "✓ Updated successfully"
    else
        echo "✗ Update failed"
        return 1
    fi
}

# Update multiple settings
update_config '.version' '"2.0.0"' config.toml
update_config '.debug' 'false' config.toml
update_config '.database.port' '5432' config.toml
```
