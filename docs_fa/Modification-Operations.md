# عملیات تغییر

tmq تغییر در جای فایل‌های TOML را پشتیبانی می‌کند و اجازهٔ تنظیم مقادیر جدید یا حذف کلیدهای موجود را می‌دهد.

## تنظیم مقادیر

### تخصیص مقدار ساده
```bash
# تنظیم مقدار رشته ساده
tmq '.version = "2.0.0"' -i config.toml

# تنظیم عدد
tmq '.port = 8080' -i config.toml

# تنظیم بولین
tmq '.enabled = true' -i config.toml
```

### مقادیر جدول تودرتو
```toml
# قبل
[database]
host = "oldhost"
```

```bash
# به‌روزرسانی مقدار تودرتو
tmq '.database.host = "newhost"' -i config.toml
```

```toml
# بعد
[database]
host = "newhost"
```

### ایجاد کلیدهای جدید
```bash
# افزودن کلید سطح ریشه
tmq '.new_key = "new_value"' -i config.toml

# افزودن کلید تودرتو
tmq '.database.pool_size = 10' -i config.toml
```

### تودرتوی عمیق
```bash
# ایجاد ساختار تودرتوی عمیق
tmq '.app.cache.redis.ttl = 3600' -i config.toml

# این ایجاد می‌کند:
# [app.cache.redis]
# ttl = 3600
```

### مقادیر آرایه
```bash
# تنظیم آرایه رشته
tmq '.tags = ["web", "api", "prod"]' -i config.toml

# تنظیم آرایه عدد
tmq '.ports = [8080, 8443, 9000]' -i config.toml
```

### آبجکت‌های پیچیده
```bash
# تنظیم جدول درون‌خطی
tmq '.credentials = { username = "admin", password = "secret" }' -i config.toml

# تنظیم آبجکت تودرتو
tmq '.database = { host = "localhost", port = 5432 }' -i config.toml
```

## عملیات حذف

### حذف کلیدهای ریشه
```bash
# حذف کلید سطح بالا
tmq 'del(.obsolete_key)' -i config.toml
```

### حذف کلیدهای تودرتو
```bash
# حذف از جدول تودرتو
tmq 'del(.database.old_setting)' -i config.toml
```

### حذف عناصر آرایه
```bash
# حذف اندیس مشخص آرایه
tmq 'del(.servers[1])' -i config.toml
```

## حالت Dry Run

### پیش‌نمایش تغییرات
```bash
# مشاهدهٔ تغییرات بدون اصلاح فایل
tmq '.version = "3.0.0"' --dry-run config.toml

# پیش‌نمایش حذف
tmq 'del(.obsolete_key)' --dry-run config.toml
```

### تغییرات امن
```bash
# همیشه ابتدا با dry-run تست کنید
tmq '.database.host = "prod-db"' --dry-run config.toml

# اگر درست بود اعمال کنید
tmq '.database.host = "prod-db"' -i config.toml
```

## مثال‌های پیشرفته

### به‌روزرسانی پیکربندی
```toml
# config.toml قبل
[app]
version = "1.0.0"
debug = true

[database]
host = "dev-db"
port = 5432
```

```bash
# به‌روزرسانی برای استقرار production
tmq '.app.version = "1.1.0"' -i config.toml
tmq '.app.debug = false' -i config.toml
tmq '.database.host = "prod-db"' -i config.toml
```

```toml
# config.toml بعد
[app]
version = "1.1.0"
debug = false

[database]
host = "prod-db"
port = 5432
```

### پیکربندی محیط‌محور
```bash
# تنظیمات توسعه
tmq '.database.host = "localhost"' -i config.toml
tmq '.debug = true' -i config.toml

# تنظیمات production
tmq '.database.host = "prod.example.com"' -i config.toml
tmq '.debug = false' -i config.toml
```

### عملیات پاکسازی
```bash
# حذف تنظیمات منسوخ
tmq 'del(.legacy_feature)' -i config.toml
tmq 'del(.old_database_url)' -i config.toml

# حذف کاربران تست
tmq 'del(.test_users)' -i config.toml
```

## مدیریت خطا

### مسیرهای ناموجود
```bash
# تنظیم والد ناموجود ساختار را ایجاد می‌کند
tmq '.new.deep.key = "value"' -i config.toml
# ایجاد می‌کند: [new.deep]
#          key = "value"
```

### تعارض نوع
```bash
# بازنویسی انواع مختلف مجاز است
tmq '.value = "string"' -i config.toml  # قبلاً عدد بود
tmq '.value = 42' -i config.toml        # قبلاً رشته بود
```

### عملیات نامعتبر
```bash
# نام‌های کلید نامعتبر
tmq '.invalid key = "value"' -i config.toml
# خطا: عبارت set نامعتبر

# نبودن کوتیشن برای رشته
tmq '.name = John' -i config.toml
# خطا: عبارت set نامعتبر
```

## استراتژی بکاپ

### بکاپ دستی
```bash
# همیشه قبل از تغییر بکاپ بگیرید
cp config.toml config.toml.backup

# تغییرات
tmq '.version = "2.0.0"' -i config.toml

# تأیید
tmq '.version' config.toml
```

### بکاپ اسکریپتی
```bash
#!/bin/bash
CONFIG_FILE="config.toml"
BACKUP_FILE="${CONFIG_FILE}.backup.$(date +%Y%m%d_%H%M%S)"

cp "$CONFIG_FILE" "$BACKUP_FILE"
echo "Backup created: $BACKUP_FILE"

# تغییرات
tmq '.version = "2.0.0"' -i "$CONFIG_FILE"

# تأیید
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
# اعتبارسنجی پس از تغییرات
tmq '.' config.toml > /dev/null || echo "Invalid TOML after modification"
```

### عملیات اتمیک
```bash
# برای به‌روزرسانی‌های حیاتی از فایل موقت استفاده کنید
TEMP_FILE=$(mktemp)
cp config.toml "$TEMP_FILE"

tmq '.critical_setting = "new_value"' -i "$TEMP_FILE"

# اعتبارسنجی
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
# قبل و بعد از تغییرات commit کنید
git add config.toml
git commit -m "Update configuration via tmq"

# تغییرات
tmq '.version = "2.0.0"' -i config.toml

git add config.toml
git commit -m "Bump version to 2.0.0"
```

### یکپارچه‌سازی با اسکریپت
```bash
#!/bin/bash
set -e

# تابع برای به‌روزرسانی امن پیکربندی
update_config() {
    local key="$1"
    local value="$2"
    local file="$3"

    echo "Updating $key = $value in $file"

    # ابتدا dry run
    if tmq "$key = $value" --dry-run "$file" >/dev/null; then
        tmq "$key = $value" -i "$file"
        echo "✓ Updated successfully"
    else
        echo "✗ Update failed"
        return 1
    fi
}

# به‌روزرسانی چند تنظیم
update_config '.version' '"2.0.0"' config.toml
update_config '.debug' 'false' config.toml
update_config '.database.port' '5432' config.toml
```
