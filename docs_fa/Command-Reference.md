# مرجع دستورات

مرجع کامل برای تمام دستورات، پرچم‌ها و گزینه‌های tmq.

## خلاصه دستور

```bash
tmq [OPTIONS] [QUERY] [FILE...]
tmq [OPTIONS] --validate [FILE...]
tmq [OPTIONS] --compare FILE1 FILE2
```

## گزینه‌های عمومی

### گزینه‌های خروجی
- `-o, --output FORMAT`: قالب خروجی (`toml`, `json`, `yaml`)
  - پیش‌فرض: `toml`
  - مثال: `tmq '.data' config.toml -o json`

### گزینه‌های تغییر
- `-i, --inplace`: تغییر فایل‌ها در جای خود
  - باید همراه عملیات set/delete استفاده شود
  - مثال: `tmq '.version = "2.0"' -i config.toml`

### اجرای خشک
- `--dry-run`: پیش‌نمایش تغییرات بدون اعمال روی فایل
  - نمایش آنچه انجام می‌شود
  - کد خروج ۰ برای موفقیت، ۱ برای خطا
  - مثال: `tmq '.version = "2.0"' --dry-run config.toml`

### اعتبارسنجی
- `--validate`: اعتبارسنجی نحو TOML
  - کد خروج ۰ اگر معتبر، ۱ اگر نامعتبر
  - می‌تواند چندین فایل را پردازش کند
  - مثال: `tmq --validate config.toml`

### مقایسه
- `--compare FILE1 FILE2`: مقایسه دو فایل TOML
  - کد خروج ۰ اگر یکسان، ۱ اگر متفاوت
  - تفاوت‌های دقیق را نشان می‌دهد
  - مثال: `tmq --compare old.toml new.toml`

### اطلاعات
- `-v, --version`: نمایش اطلاعات نسخه
- `-h, --help`: نمایش متن راهنما

## نحو پرس‌وجو

### پرس‌وجوهای پایه
```bash
# Root key
.key

# Nested table
.table.key

# Deep nesting
.app.database.host

# Array element
.array[0]

# Array element with nesting
.servers[1].name
```

### عملیات set
```bash
# Simple assignment
.key = "value"

# Number assignment
.port = 8080

# Boolean assignment
.enabled = true

# Array assignment
.tags = ["web", "api"]

# Object assignment
.config = { host = "localhost", port = 5432 }
```

### عملیات حذف
```bash
# Delete root key
del(.key)

# Delete nested key
del(.table.key)

# Delete array element
del(.array[0])
```

## مثال‌ها

### مثال‌های پرس‌وجو
```bash
# Simple query
tmq '.version' config.toml

# Nested query
tmq '.database.host' config.toml

# Array query
tmq '.servers[0].name' config.toml

# Multiple files
tmq '.version' config/*.toml

# JSON output
tmq '.config' file.toml -o json

# YAML output
tmq '.config' file.toml -o yaml
```

### مثال‌های تغییر
```bash
# Set string value
tmq '.version = "2.0.0"' -i config.toml

# Set nested value
tmq '.database.host = "prod-db"' -i config.toml

# Set array
tmq '.ports = [8080, 8443]' -i config.toml

# Delete key
tmq 'del(.obsolete)' -i config.toml

# Dry run
tmq '.version = "test"' --dry-run config.toml
```

### مثال‌های اعتبارسنجی
```bash
# Single file
tmq --validate config.toml

# Multiple files
tmq --validate *.toml

# In script
if tmq --validate config.toml; then
    echo "Valid TOML"
fi
```

### مثال‌های مقایسه
```bash
# Compare files
tmq --compare config1.toml config2.toml

# In script
if tmq --compare expected.toml actual.toml; then
    echo "Files match"
fi
```

## کدهای خروج

| کد | معنی | توضیحات |
|------|---------|-------------|
| 0 | موفقیت | عملیات با موفقیت انجام شد |
| 1 | خطای تجزیه | تجزیه TOML یا خطای پرس‌وجو |
| 2 | خطای استفاده | آرگومان‌های نامعتبر خط فرمان |
| 3 | خطای امنیتی | عبور از مسیر یا نقض امنیت |
| 4 | خطای فایل | فایل یافت نشد، دسترسی رد شد و غیره |

## پیام‌های خطا

tmq پیام‌های خطای ساخت‌یافته ارائه می‌دهد:

```
ERROR: <error_type>
DETAILS: <detailed_description>
ACTION: <suggested_fix>
```

### خطاهای رایج

#### خطاهای تجزیه
```bash
tmq '.invalid..' config.toml
# ERROR: Invalid query path
# DETAILS: Query path cannot be empty
# ACTION: Check your query syntax
```

#### خطاهای فایل
```bash
tmq '.' nonexistent.toml
# ERROR: File operation error
# DETAILS: Cannot read file 'nonexistent.toml'
# ACTION: Ensure the file exists and is readable
```

#### خطاهای اعتبارسنجی
```bash
tmq --validate malformed.toml
# ERROR: TOML parsing failed
# DETAILS: Expected newline at line 5, column 10
# ACTION: Fix the TOML syntax error
```

## متغیرهای محیطی

tmq از متغیرهای محیطی برای پیکربندی استفاده نمی‌کند. تمام گزینه‌ها از طریق پرچم‌های خط فرمان مشخص می‌شوند.

## قالب‌های فایل

### ورودی TOML
tmq قالب استاندارد TOML 1.0.0 را می‌پذیرد:

```toml
# Comments are preserved in output
title = "Example"

[database]
host = "localhost"
port = 5432

[[servers]]
name = "web1"
ip = "192.168.1.1"
```

### خروجی JSON
```json
{
  "title": "Example",
  "database": {
    "host": "localhost",
    "port": 5432
  },
  "servers": [
    {
      "name": "web1",
      "ip": "192.168.1.1"
    }
  ]
}
```

### خروجی YAML
```yaml
title: Example
database:
  host: localhost
  port: 5432
servers:
- name: web1
  ip: 192.168.1.1
```

## کارایی

### معیارها
- عملیات پرس‌وجو: O(1) - زمان ثابت
- تجزیه فایل: O(n) که n اندازه فایل است
- استفاده از حافظه: کمتر از ۱۰ مگابایت برای فایل‌های معمولی
- فایل‌های بزرگ: به‌طور خطی با اندازه مقیاس می‌پذیرد

### نکات بهینه‌سازی
- به‌جای خروجی کامل فایل از پرس‌وجوهای مشخص استفاده کنید
- برای استفاده برنامه‌نویسی خروجی JSON/YAML را ترجیح دهید
- برای چندین فایل از عملیات گروهی استفاده کنید

## پشتیبانی پلتفرم

### پلتفرم‌های پشتیبانی‌شده
- **لینوکس**: amd64, arm64
- **macOS**: amd64 (اینتل), arm64 (Apple Silicon)
- **ویندوز**: amd64

### نام فایل‌های باینری
- `tmq-linux-amd64` (لینوکس اینتل/AMD)
- `tmq-linux-arm64` (لینوکس ARM)
- `tmq-darwin-amd64` (macOS اینتل)
- `tmq-darwin-arm64` (macOS Apple Silicon)
- `tmq-windows-amd64.exe` (ویندوز)

## محدودیت‌ها

### محدودیت‌های فعلی
- تبدیل JSON/YAML به TOML وجود ندارد (برنامه‌ریزی شده)
- حفظ توضیحات در تغییرات وجود ندارد (برنامه‌ریزی شده)
- سیستم افزونه وجود ندارد (برنامه‌ریزی شده)
- API کتابخانه وجود ندارد (برنامه‌ریزی شده)

### محدودیت‌های اندازه فایل
- حداکثر اندازه فایل: ۱۰۰ مگابایت
- توصیه شده: کمتر از ۱۰ مگابایت برای عملکرد بهینه

### محدودیت‌های طول مسیر
- حداکثر طول مسیر: ۱۰۲۴ کاراکتر
- حداکثر عمق تودرتویی: ۱۰۰ سطح

## عیب‌یابی

### مسائل رایج

#### "command not found"
```bash
# Check if tmq is in PATH
which tmq

# Or use full path
/path/to/tmq --version
```

#### "permission denied"
```bash
# Make executable
chmod +x tmq

# Or run with full permissions
sudo ./tmq --version
```

#### "directory not found"
```bash
# Check file exists
ls -la config.toml

# Check current directory
pwd

# Use absolute path
tmq '.' /full/path/to/config.toml
```

#### Invalid query
```bash
# Check syntax
tmq --help

# Test with simple query
tmq '.' config.toml
```

### حالت اشکال‌زدایی
```bash
# Use verbose output (when available)
tmq --help

# Check file contents
cat config.toml

# Validate file
tmq --validate config.toml
```

### دریافت کمک
```bash
# Show help
tmq --help

# Show version
tmq --version

# Report issues: https://github.com/azolfagharj/tmq/issues
```

## تاریخچه نسخه‌ها

### v1.0.1
- انتشار نخست
- عملیات پرس‌وجو، set و حذف
- اعتبارسنجی و مقایسه
- عملیات گروهی
- خروجی JSON/YAML
- باینری‌های چندپلتفرمی
