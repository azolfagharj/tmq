# عیب‌یابی

مسائل رایج و راه‌حل‌ها هنگام استفاده از tmq.

## مشکلات نصب

### دستور یافت نشد
```bash
tmq --version
# bash: tmq: command not found
```

**راه‌حل‌ها:**
1. **بررسی PATH:**
   ```bash
   echo $PATH
   which tmq
   ```

2. **استفاده از مسیر کامل:**
   ```bash
   /usr/local/bin/tmq --version
   ./tmq --version
   ```

3. **نصب مجدد:**
   ```bash
   # Download again
   wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
   chmod +x tmq-linux-amd64
   sudo mv tmq-linux-amd64 /usr/local/bin/tmq
   ```

### دسترسی رد شد
```bash
./tmq --version
# bash: ./tmq: Permission denied
```

**راه‌حل:**
```bash
chmod +x tmq
./tmq --version
```

### پلتفرم پشتیبانی‌نشده
```bash
# Trying to run ARM binary on x86 system
./tmq-linux-arm64 --version
# exec format error
```

**راه‌حل:** باینری مناسب پلتفرم خود را دانلود کنید:
- لینوکس x86_64: `tmq-linux-amd64`
- لینوکس ARM64: `tmq-linux-arm64`
- macOS اینتل: `tmq-darwin-amd64`
- macOS Apple Silicon: `tmq-darwin-arm64`
- ویندوز: `tmq-windows-amd64.exe`

## مشکلات پرس‌وجو

### مسیر پرس‌وجوی نامعتبر
```bash
tmq '.invalid..path' config.toml
# ERROR: Invalid query path
```

**اشتباهات رایج:**
1. **بخش‌های مسیر خالی:**
   ```bash
   # Wrong
   tmq '.database..' config.toml

   # Right
   tmq '.database' config.toml
   ```

2. **کاراکترهای نامعتبر:**
   ```bash
   # Wrong (spaces in path)
   tmq '.my key' config.toml

   # Right (use underscores or different structure)
   tmq '.my_key' config.toml
   ```

3. **پرانتز جفت‌نشده:**
   ```bash
   # Wrong
   tmq '.array[0' config.toml

   # Right
   tmq '.array[0]' config.toml
   ```

### کلید یافت نشد
```bash
tmq '.nonexistent' config.toml
# ERROR: key 'nonexistent' not found
```

**راه‌حل‌ها:**
1. **بررسی وجود کلید:**
   ```bash
   # List all keys
   tmq '.' config.toml

   # Check specific section
   tmq '.database' config.toml
   ```

2. **حساسیت به حروف:**
   ```toml
   [Database]  # Capital D
   host = "localhost"
   ```
   ```bash
   # Wrong
   tmq '.database.host' config.toml

   # Right
   tmq '.Database.host' config.toml
   ```

3. **آرایه در مقابل آبجکت:**
   ```toml
   # This is an array of tables
   [[servers]]
   name = "server1"

   # This is a table
   [database]
   host = "localhost"
   ```

### شاخص آرایه خارج از محدوده
```bash
tmq '.servers[10]' config.toml
# ERROR: array index 10 out of bounds
```

**راه‌حل:** ابتدا طول آرایه را بررسی کنید:
```bash
# See the array
tmq '.servers' config.toml

# Get array length (in scripts)
length=$(tmq '.servers | length' config.toml)
```

## مشکلات تغییر

### عدم امکان تغییر فایل
```bash
tmq '.version = "2.0"' -i config.toml
# ERROR: File operation error
# DETAILS: Cannot write to file 'config.toml'
```

**راه‌حل‌ها:**
1. **بررسی دسترسی‌ها:**
   ```bash
   ls -la config.toml
   chmod 644 config.toml
   ```

2. **بررسی فضای دیسک:**
   ```bash
   df -h
   ```

3. **فایل قفل شده:**
   ```bash
   # Check if file is open by another process
   lsof config.toml
   ```

### عبارت set نامعتبر
```bash
tmq '.version = ' -i config.toml
# ERROR: Invalid set expression
```

**خطاهای نحو رایج:**
1. **مقدار گمشده:**
   ```bash
   # Wrong
   tmq '.version =' -i config.toml

   # Right
   tmq '.version = "2.0"' -i config.toml
   ```

2. **رشته‌های بدون نقل‌قول:**
   ```bash
   # Wrong (if value contains spaces or special chars)
   tmq '.path = /usr/local/bin' -i config.toml

   # Right
   tmq '.path = "/usr/local/bin"' -i config.toml
   ```

3. **مقادیر TOML نامعتبر:**
   ```bash
   # Wrong
   tmq '.value = [unclosed array' -i config.toml

   # Right
   tmq '.value = ["item1", "item2"]' -i config.toml
   ```

## مشکلات اعتبارسنجی

### خطای تجزیه TOML
```bash
tmq --validate config.toml
# ERROR: TOML parsing failed
# DETAILS: Expected newline at line 5, column 10
```

**خطاهای نحو TOML رایج:**
1. **نقل‌قول‌های گمشده:**
   ```toml
   # Wrong
   title = My App

   # Right
   title = "My App"
   ```

2. **نحو آرایه نامعتبر:**
   ```toml
   # Wrong
   items = [1, 2, 3

   # Right
   items = [1, 2, 3]
   ```

3. **تعریف جدول نامعتبر:**
   ```toml
   # Wrong
   [section
   key = "value"

   # Right
   [section]
   key = "value"
   ```

### مشکلات رمزگذاری فایل
```bash
tmq --validate config.toml
# ERROR: Invalid UTF-8 encoding
```

**راه‌حل:** فایل را به UTF-8 تبدیل کنید:
```bash
# Check current encoding
file config.toml

# Convert to UTF-8
iconv -f latin1 -t utf8 config.toml > config_utf8.toml
mv config_utf8.toml config.toml
```

## مشکلات مقایسه

### فایل‌ها متفاوت‌اند ولی باید یکسان باشند
```bash
tmq --compare file1.toml file2.toml
# Files are different
```

**علل محتمل:**
1. **قالب‌بندی متفاوت:**
   ```toml
   # file1.toml
   host="localhost"

   # file2.toml
   host = "localhost"
   ```
   از نظر معنایی یکسان‌اند ولی قالب‌بندی متفاوت دارند.

2. **توضیحات نادیده گرفته می‌شوند:**
   توضیحات TOML بخشی از ساختار داده نیستند.

3. **تفاوت در ترتیب:**
   ```toml
   # Different order, same content
   [table]
   b = 2
   a = 1

   [table]
   a = 1
   b = 2
   ```

## مشکلات عملکرد

### عملیات کند روی فایل‌های بزرگ
```bash
time tmq '.' large.toml > /dev/null
# Takes several seconds
```

**راه‌حل‌ها:**
1. **استفاده از پرس‌وجوهای مشخص:**
   ```bash
   # Instead of full file
   tmq '.' large.toml

   # Use specific path
   tmq '.database.host' large.toml
   ```

2. **خروجی به فایل به‌جای stdout:**
   ```bash
   tmq '.' large.toml -o json > output.json
   ```

3. **بررسی اندازه فایل:**
   ```bash
   ls -lh large.toml
   # If > 50MB, consider splitting
   ```

### استفاده زیاد از حافظه
```bash
# Monitor memory
/usr/bin/time -v tmq '.' large.toml
```

**راه‌حل‌ها:**
1. **پردازش به صورت بخش‌بخش:**
   ```bash
   # Instead of whole file, process sections
   tmq '.section1' large.toml
   tmq '.section2' large.toml
   ```

2. **استفاده از streaming برای خروجی‌های بزرگ:**
   ```bash
   tmq '.' large.toml | head -100
   ```

## مشکلات عملیات گروهی

### تعداد زیاد فایل
```bash
tmq '.version' config/*.toml
# Argument list too long
```

**راه‌حل:** از `find` استفاده کنید یا به صورت دسته‌ای پردازش کنید:
```bash
# Use find
find config/ -name "*.toml" -exec tmq '.version' {} \;

# Process in batches
for file in config/*.toml; do
    tmq '.version' "$file"
done
```

### نتایج نامتناقض
```bash
tmq '.version' config/*.toml
# Some files show different output format
```

**علت:** فایل‌های مختلف ممکن است ساختارهای متفاوت داشته باشند.

**راه‌حل:** کلیدهای گمشده را مدیریت کنید:
```bash
# Use default values
tmq '.version // "unknown"' config/*.toml
```

## مشکلات اسکریپت‌نویسی

### مدیریت کد خروج
```bash
#!/bin/bash
tmq '.nonexistent' config.toml
echo "Exit code: $?"
# Always prints exit code, even if command fails
```

**راه‌حل:** از مدیریت خطای مناسب استفاده کنید:
```bash
#!/bin/bash
set -e  # Exit on first error

if ! tmq '.nonexistent' config.toml >/dev/null 2>&1; then
    echo "Key not found"
    exit 1
fi
```

### مشکلات تجزیه خروجی
```bash
# Trying to parse multi-line output
version=$(tmq '.' config.toml | grep version)
# Unreliable
```

**راه‌حل:** از پرس‌وجوهای مشخص استفاده کنید:
```bash
version=$(tmq '.version' config.toml)
```

### جایگزینی متغیر
```bash
# Wrong
tmq '.version = $VERSION' -i config.toml

# Right
tmq ".version = \"$VERSION\"" -i config.toml
```

## مشکلات خاص پلتفرم

### مشکلات مسیر ویندوز
```cmd
REM Wrong (backslashes in paths)
tmq ".path = C:\Program Files\app" -i config.toml

REM Right (forward slashes or quoted)
tmq ".path = \"C:/Program Files/app\"" -i config.toml
```

### مشکلات دسترسی macOS
```bash
# When installed in /usr/local/bin
sudo tmq --version
# May require sudo if installed system-wide
```

## مشکلات شبکه و دانلود

### مشکلات پروکسی
```bash
# If behind corporate proxy
export https_proxy=http://proxy.company.com:8080
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```

### مشکلات گواهینامه
```bash
# Skip SSL verification (not recommended for production)
wget --no-check-certificate https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```

### مسدود شدن دانلود توسط فایروال
```bash
# Check connectivity
curl -I https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64

# Use alternative download method
curl -L -o tmq https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```

## دریافت کمک

### اطلاعات اشکال‌زدایی
```bash
# Collect system info
uname -a
tmq --version

# Test with simple file
echo 'test = "value"' > test.toml
tmq '.test' test.toml
```

### گزارش باگ
هنگام گزارش باگ، این موارد را شامل کنید:
1. **نسخه tmq:** `tmq --version`
2. **پلتفرم:** `uname -a`
3. **دستوری که ناموفق است**
4. **خروجی کامل خطا**
5. **فایل نمونه TOML** (در صورت امکان)
6. **رفتار مورد انتظار در مقابل رفتار واقعی**

### پشتیبانی جامعه
- **GitHub Issues:** https://github.com/azolfagharj/tmq/issues
- **Discussions:** https://github.com/azolfagharj/tmq/discussions

## رفع سریع

| مشکل | رفع سریع |
|-------|-----------|
| دستور یافت نشد | `export PATH=$PATH:/usr/local/bin` |
| دسترسی رد شد | `chmod +x tmq` |
| پرس‌وجوی نامعتبر | بررسی نحو: `tmq --help` |
| فایل یافت نشد | `ls -la config.toml` |
| خطای تجزیه | اعتبارسنجی: `tmq --validate config.toml` |
| مشکلات حافظه | استفاده از پرس‌وجوهای مشخص به‌جای کل فایل |
| عملیات گروهی کند | پردازش در دسته‌های کوچک‌تر |
