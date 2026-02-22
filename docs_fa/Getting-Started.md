# شروع سریع

این راهنما مقدمات کار با tmq برای فایل‌های TOML را پوشش می‌دهد.

## مفاهیم پایه

### ساختار TOML
فایل‌های TOML از یک فرمت ساده و خوانا برای انسان استفاده می‌کنند:

```toml
# config.toml
title = "My Application"

[database]
host = "localhost"
port = 5432
enabled = true

[server]
host = "0.0.0.0"
port = 8080

[[users]]
name = "Alice"
role = "admin"

[[users]]
name = "Bob"
role = "user"
```

### نحو کوئری
tmq از نحو نقطه‌ای برای دسترسی به مقادیر TOML استفاده می‌کند:
- `title` → کلیدهای سطح ریشه
- `database.host` → مقادیر جدول تودرتو
- `users[0].name` → دسترسی به عنصر آرایه

## اولین قدم‌ها

### بررسی نصب
```bash
tmq --version
tmq --help
```

### کوئری پایه
یک فایل TOML آزمایشی بسازید:
```bash
cat > config.toml << 'EOF'
title = "My App"
version = "1.0.0"

[database]
host = "localhost"
port = 5432
EOF
```

کوئری مقادیر:
```bash
# Get title
tmq '.title' config.toml
# Output: "My App"

# Get version
tmq '.version' config.toml
# Output: "1.0.0"

# Get database host
tmq '.database.host' config.toml
# Output: "localhost"

# Get database port
tmq '.database.port' config.toml
# Output: 5432
```

### نمایش کل داده‌ها
```bash
# Show entire file
tmq '.' config.toml

# Format as JSON
tmq '.' config.toml -o json

# Format as YAML
tmq '.' config.toml -o yaml
```

## الگوهای رایج

### ورودی از stdin
```bash
# Read from stdin
cat config.toml | tmq '.version'

# Use with other tools
echo 'version = "2.0.0"' | tmq '.version'
```

### مدیریت خطا در اسکریپت‌ها
```bash
#!/bin/bash
VERSION=$(tmq '.version' config.toml)
if [ $? -ne 0 ]; then
    echo "Error reading version from config.toml"
    exit 1
fi
echo "Version: $VERSION"
```

### عملیات فایل
```bash
# Check if file exists and is valid TOML
if tmq '.' config.toml > /dev/null 2>&1; then
    echo "config.toml is valid"
else
    echo "config.toml is invalid or missing"
fi
```

## قدم‌های بعدی

اکنون که مقدمات را آموختید:

1. **عملیات کوئری**: کوئری پیشرفته در [عملیات کوئری](Query-Operations.md)
2. **تغییرات**: تغییر فایل‌های TOML در [عملیات تغییر](Modification-Operations.md)
3. **اعتبارسنجی**: بررسی اعتبار فایل در [اعتبارسنجی و مقایسه](Validation-and-Comparison.md)
4. **مثال‌ها**: نمونه‌های جامع در [مثال‌ها](Examples.md)

## مرجع سریع

| دستور | توضیح |
|-------|-------|
| `tmq '.key' file.toml` | کوئری یک مقدار |
| `tmq '.' file.toml` | نمایش کل فایل |
| `tmq '.' file.toml -o json` | خروجی به صورت JSON |
| `cat file.toml \| tmq '.key'` | خواندن از stdin |
| `tmq --validate file.toml` | اعتبارسنجی نحو TOML |
| `tmq --help` | نمایش راهنما |

## کدهای خروجی

tmq از کدهای خروجی استاندارد استفاده می‌کند:
- `0`: موفقیت
- `1`: خطای parse/اجرا
- `2`: خطای استفاده
- `3`: خطای امنیتی
- `4`: خطای فایل
