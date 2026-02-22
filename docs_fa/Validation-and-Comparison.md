# اعتبارسنجی و مقایسه

tmq ابزارهایی برای اعتبارسنجی فایل‌های TOML و مقایسهٔ اختلاف‌های بین آن‌ها فراهم می‌کند.

## اعتبارسنجی

### اعتبارسنجی نحو پایه
```bash
# بررسی معتبر بودن فایل TOML
tmq --validate config.toml

# کد خروجی 0 اگر معتبر، 1 اگر نامعتبر
echo $?  # بررسی کد خروجی
```

### اعتبارسنجی گروهی
```bash
# اعتبارسنجی چند فایل
tmq --validate config/*.toml

# استفاده با find
find . -name "*.toml" -exec tmq --validate {} \;
```

### اعتبارسنجی در اسکریپت‌ها
```bash
#!/bin/bash
validate_toml() {
    local file="$1"
    if tmq --validate "$file" >/dev/null 2>&1; then
        echo "✓ $file معتبر است"
        return 0
    else
        echo "✗ $file نامعتبر است"
        return 1
    fi
}

# اعتبارسنجی تمام فایل‌های TOML در پوشه
for file in *.toml; do
    validate_toml "$file" || exit 1
done
```

## مقایسه

### مقایسهٔ پایهٔ فایل‌ها
```bash
# مقایسهٔ دو فایل TOML
tmq --compare config1.toml config2.toml

# کدهای خروجی:
# 0 = فایل‌ها یکسان‌اند
# 1 = فایل‌ها متفاوت‌اند
```

### خروجی مقایسهٔ تفصیلی
```bash
# مشاهدهٔ اختلاف‌های تفصیلی
tmq --compare old-config.toml new-config.toml
```

### مقایسه در CI/CD
```bash
# اگر پیکربندی به‌طور غیرمنتظره تغییر کرده، build ناموفق شود
if ! tmq --compare expected.toml actual.toml >/dev/null; then
    echo "Configuration mismatch!"
    tmq --compare expected.toml actual.toml
    exit 1
fi
```

## استفادهٔ پیشرفته

### اعتبارسنجی با جزئیات خطا
```bash
# tmq پیام‌های خطای تفصیلی نشان می‌دهد
tmq --validate invalid.toml
# Error: parse error: expected newline at line 5, column 10
# DETAILS: Check your TOML syntax
# ACTION: Fix the syntax error and try again
```

### مقایسه با خروجی
```bash
# هدایت خروجی مقایسه به فایل
tmq --compare file1.toml file2.toml > differences.txt

# استفاده در اسکریپت
if tmq --compare "$EXPECTED" "$ACTUAL" > diff.log; then
    echo "Files match"
else
    echo "Differences found:"
    cat diff.log
fi
```

## مثال‌های یکپارچه‌سازی

### Pre-commit Hook
```bash
#!/bin/bash
# .git/hooks/pre-commit

# اعتبارسنجی تمام فایل‌های TOML
echo "Validating TOML files..."
if ! find . -name "*.toml" -exec tmq --validate {} \;; then
    echo "TOML validation failed"
    exit 1
fi

echo "All TOML files are valid"
```

### تشخیص انحراف پیکربندی
```bash
#!/bin/bash
# بررسی مطابقت پیکربندی production با پیکربندی مورد انتظار

PROD_CONFIG="prod-config.toml"
EXPECTED_CONFIG="expected-config.toml"

if tmq --compare "$EXPECTED_CONFIG" "$PROD_CONFIG" >/dev/null; then
    echo "✓ پیکربندی production با پیکربندی مورد انتظار مطابقت دارد"
    exit 0
else
    echo "✗ انحراف پیکربندی تشخیص داده شد!"
    echo "اختلاف‌ها:"
    tmq --compare "$EXPECTED_CONFIG" "$PROD_CONFIG"
    exit 1
fi
```

### یکپارچه‌سازی با CI Pipeline
```yaml
# .github/workflows/validate.yml
name: Validate Configuration

on:
  pull_request:
    paths:
      - 'config/*.toml'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: checkout@v4

      - name: Setup tmq
        run: |
          wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
          chmod +x tmq-linux-amd64
          sudo mv tmq-linux-amd64 /usr/local/bin/tmq

      - name: Validate TOML files
        run: |
          for file in config/*.toml; do
            echo "Validating $file..."
            tmq --validate "$file"
          done

      - name: Check configuration consistency
        run: |
          tmq --compare config/base.toml config/production.toml
```

## پیام‌های خطا

### خطاهای اعتبارسنجی
```bash
# خطاهای نحو
tmq --validate malformed.toml
# ERROR: TOML parsing failed
# DETAILS: Expected newline at line 5, column 10
# ACTION: Check your TOML syntax and fix any formatting issues

# فایل پیدا نشد
tmq --validate nonexistent.toml
# ERROR: File operation error
# DETAILS: Cannot read file 'nonexistent.toml'
# ACTION: Ensure the file exists and is readable
```

### خطاهای مقایسه
```bash
# فایل پیدا نشد
tmq --compare missing1.toml missing2.toml
# ERROR: File operation error
# DETAILS: Cannot read file 'missing1.toml'
# ACTION: Ensure both files exist and are readable

# خطای parse در یک فایل
tmq --compare valid.toml invalid.toml
# ERROR: TOML parsing failed
# DETAILS: Invalid TOML in 'invalid.toml'
# ACTION: Fix the TOML syntax in invalid.toml
```

## فرمت‌های خروجی

### خروجی اعتبارسنجی
```bash
# اعتبارسنجی موفق (بدون خروجی، کد خروجی 0)
tmq --validate valid.toml

# اعتبارسنجی ناموفق (پیام خطا، کد خروجی 1)
tmq --validate invalid.toml
```

### خروجی مقایسه
```bash
# فایل‌های یکسان (بدون خروجی، کد خروجی 0)
tmq --compare same1.toml same1.toml

# فایل‌های متفاوت (نمایش اختلاف‌ها، کد خروجی 1)
tmq --compare old.toml new.toml
# Files are different:
# - old.toml: version = "1.0.0"
# + new.toml: version = "2.0.0"
```

## بهترین روش‌ها

### اعتبارسنجی در توسعه
```bash
# افزودن به فرایند build
make validate-config:
    @for file in config/*.toml; do \
        echo "Validating $$file..."; \
        tmq --validate "$$file" || exit 1; \
    done
    @echo "All configurations valid"
```

### تست پیکربندی
```bash
# تست تغییرات پیکربندی
#!/bin/bash

CONFIG_FILE="app-config.toml"
BACKUP_FILE="${CONFIG_FILE}.backup"

# بکاپ پیکربندی فعلی
cp "$CONFIG_FILE" "$BACKUP_FILE"

# تغییرات تست
tmq '.debug = true' -i "$CONFIG_FILE"
tmq '.test_mode = true' -i "$CONFIG_FILE"

# اعتبارسنجی تغییرات
if tmq --validate "$CONFIG_FILE" >/dev/null; then
    echo "✓ تغییرات پیکربندی معتبرند"
    # اینجا اپلیکیشن را تست کنید
    # ./run-tests.sh
else
    echo "✗ پیکربندی پس از تغییرات نامعتبر است"
fi

# بازگردانی بکاپ
mv "$BACKUP_FILE" "$CONFIG_FILE"
```

### یکپارچه‌سازی با کنترل نسخه
```bash
# بررسی تغییرات پیکربندی در commitها
git diff --name-only HEAD~1 | grep '\.toml$' | while read -r file; do
    echo "Checking $file..."
    if ! tmq --validate "$file" >/dev/null; then
        echo "Invalid TOML in $file"
        exit 1
    fi
done
```

### نظارت بر تغییرات پیکربندی
```bash
#!/bin/bash
# نظارت بر پیکربندی برای تغییرات غیرمجاز

BASELINE="config-baseline.toml"
CURRENT="config.toml"

if [ ! -f "$BASELINE" ]; then
    echo "Creating baseline configuration..."
    cp "$CURRENT" "$BASELINE"
    exit 0
fi

if ! tmq --compare "$BASELINE" "$CURRENT" >/dev/null; then
    echo "WARNING: Configuration has changed!"
    echo "Differences:"
    tmq --compare "$BASELINE" "$CURRENT"

    read -p "Accept changes and update baseline? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cp "$CURRENT" "$BASELINE"
        echo "Baseline updated"
    fi
fi
```

## عملکرد

- اعتبارسنجی سریع است — معمولاً کمتر از ۱۰۰ میلی‌ثانیه برای اکثر فایل‌ها
- سرعت مقایسه به اندازهٔ فایل و میزان اختلاف بستگی دارد
- مصرف حافظه برای هر دو عملیات حداقلی است

## خلاصهٔ کدهای خروجی

| عملیات | موفقیت | خطا |
|--------|--------|-----|
| اعتبارسنجی | 0 | 1 (خطای parse) |
| مقایسه | 0 (یکسان) | 1 (متفاوت) |
| خطاهای فایل | - | 4 (فایل پیدا نشد/قابل خواندن نیست) |
