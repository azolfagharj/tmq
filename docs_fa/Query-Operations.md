# عملیات کوئری

tmq توانایی‌های قدرتمند کوئری برای استخراج داده از فایل‌های TOML با نحو نقطه‌ای ساده فراهم می‌کند.

## کوئری‌های پایه

### کلیدهای سطح ریشه
```toml
# config.toml
title = "My App"
version = "1.0.0"
enabled = true
```

```bash
tmq '.title' config.toml     # "My App"
tmq '.version' config.toml   # "1.0.0"
tmq '.enabled' config.toml   # true
```

### جدول‌های تودرتو
```toml
[database]
host = "localhost"
port = 5432
ssl = false

[server]
host = "0.0.0.0"
port = 8080
```

```bash
tmq '.database.host' config.toml    # "localhost"
tmq '.database.port' config.toml    # 5432
tmq '.server.host' config.toml      # "0.0.0.0"
```

### تودرتوی عمیق
```toml
[app]
[app.config]
[app.config.database]
host = "db.example.com"
port = 3306

[app.config.cache]
type = "redis"
ttl = 3600
```

```bash
tmq '.app.config.database.host' config.toml   # "db.example.com"
tmq '.app.config.cache.type' config.toml      # "redis"
tmq '.app.config.cache.ttl' config.toml       # 3600
```

## عملیات آرایه

### دسترسی به عناصر آرایه
```toml
[[servers]]
name = "web1"
ip = "192.168.1.1"

[[servers]]
name = "web2"
ip = "192.168.1.2"

[[servers]]
name = "db1"
ip = "192.168.1.10"
```

```bash
# اولین سرور
tmq '.servers[0].name' config.toml    # "web1"
tmq '.servers[0].ip' config.toml      # "192.168.1.1"

# دومین سرور
tmq '.servers[1].name' config.toml    # "web2"

# سرور دیتابیس
tmq '.servers[2].name' config.toml    # "db1"
```

### آرایه مقادیر
```toml
ports = [8080, 8443, 9000]
tags = ["web", "api", "admin"]
```

```bash
# دسترسی به کل آرایه‌ها
tmq '.ports' config.toml      # [8080, 8443, 9000]
tmq '.tags' config.toml       # ["web", "api", "admin"]

# دسترسی به عناصر آرایه
tmq '.ports[0]' config.toml   # 8080
tmq '.ports[1]' config.toml   # 8443
tmq '.tags[2]' config.toml    # "admin"
```

## فرمت‌های خروجی

### خروجی پیش‌فرض TOML
```bash
tmq '.database' config.toml
# خروجی: host = "localhost"
#         port = 5432
```

### خروجی JSON
```bash
tmq '.database' config.toml -o json
# خروجی: {"host":"localhost","port":5432}
```

### خروجی YAML
```bash
tmq '.database' config.toml -o yaml
# خروجی: host: localhost
#         port: 5432
```

## کوئری‌های پیشرفته

### ساختارهای پیچیده
```toml
[app]
name = "myapp"

[app.database]
host = "db.example.com"
credentials = { username = "admin", password = "secret" }

[app.features]
logging = true
metrics = { enabled = true, port = 9090 }
```

```bash
# دسترسی به آبجکت‌های تودرتو
tmq '.app.database.credentials' config.toml
# خروجی: username = "admin"
#         password = "secret"

tmq '.app.database.credentials.username' config.toml    # "admin"
tmq '.app.features.metrics' config.toml
# خروجی: enabled = true
#         port = 9090
```

### انواع داده مخلوط
```toml
# config.toml
version = "1.2.3"
debug = true
timeout = 30
pi = 3.14159

[metadata]
created = 2024-01-15T10:30:00Z
tags = ["production", "stable"]
```

```bash
tmq '.version' config.toml     # "1.2.3"
tmq '.debug' config.toml       # true
tmq '.timeout' config.toml     # 30
tmq '.pi' config.toml          # 3.14159
tmq '.metadata.tags' config.toml    # ["production", "stable"]
```

## کوئری ریشه

### دسترسی به همه چیز
```bash
# نمایش کل فایل
tmq '.' config.toml

# همان بالا (ریشه صریح)
tmq '. .' config.toml
```

### ریشه با فرمت‌های مختلف
```bash
# خروجی JSON کل فایل
tmq '.' config.toml -o json

# خروجی YAML کل فایل
tmq '.' config.toml -o yaml
```

## مدیریت خطا

### کلیدهای ناموجود
```bash
tmq '.nonexistent' config.toml
# خطا: کلید 'nonexistent' پیدا نشد
# کد خروج: 1
```

### مسیرهای نامعتبر
```bash
tmq '.invalid..path' config.toml
# خطا: مسیر کوئری نامعتبر
# کد خروج: 1
```

### عدم تطابق نوع
```bash
# تلاش برای دسترسی به اندیس آرایه روی غیر-آرایه
tmq '.title[0]' config.toml
# خطا: نمی‌توان به رشته اندیس زد
# کد خروج: 1
```

## نکات عملکرد

- کوئری‌ها در زمان ثابت O(1) ارزیابی می‌شوند
- مصرف حافظه با اندازهٔ دادهٔ کوئری‌شده مقیاس می‌یابد
- فایل‌های بزرگ (بیش از ۱۰۰ مگابایت) ممکن است نیاز به افزایش محدودیت حافظه داشته باشند

## بهترین روش‌ها

### اسکریپت‌نویسی
```bash
#!/bin/bash
# کوئری امن با مدیریت خطا
DB_HOST=$(tmq '.database.host' config.toml 2>/dev/null) || {
    echo "Error: Could not read database host from config"
    exit 1
}
echo "Database host: $DB_HOST"
```

### اعتبارسنجی قبل از کوئری
```bash
# بررسی معتبر بودن فایل قبل از کوئری
if tmq '.' config.toml >/dev/null 2>&1; then
    VERSION=$(tmq '.version' config.toml)
    echo "Version: $VERSION"
else
    echo "Invalid TOML file"
    exit 1
fi
```

### استفاده از فرمت خروجی مناسب
```bash
# برای اسکریپت‌ها، خروجی raw (پیش‌فرض)
HOST=$(tmq '.database.host' config.toml)

# برای پردازش داده، از JSON استفاده کنید
tmq '.servers' config.toml -o json | jq '.[0].name'
```
