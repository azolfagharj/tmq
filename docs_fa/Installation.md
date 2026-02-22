# نصب

tmq به‌صورت باینری از پیش کامپایل‌شده برای چند پلتفرم عرضه می‌شود. هیچ وابستگی خارجی ندارد.

## دانلود باینری

آخرین نسخه را از [GitHub Releases](https://github.com/azolfagharj/tmq/releases) دانلود کنید.

### باینری‌های موجود

| پلتفرم | معماری | نام فایل |
|--------|--------|----------|
| لینوکس | AMD64 | `tmq-linux-amd64` |
| لینوکس | ARM64 | `tmq-linux-arm64` |
| مک‌اواس | اینتل | `tmq-darwin-amd64` |
| مک‌اواس | اپل سیلیکان | `tmq-darwin-arm64` |
| ویندوز | AMD64 | `tmq-windows-amd64.exe` |

## راه‌اندازی سریع

### لینوکس (AMD64)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
chmod +x tmq-linux-amd64
sudo mv tmq-linux-amd64 /usr/local/bin/tmq
```

### لینوکس (ARM64)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-arm64
chmod +x tmq-linux-arm64
sudo mv tmq-linux-arm64 /usr/local/bin/tmq
```

### مک‌اواس (اپل سیلیکان)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-darwin-arm64
chmod +x tmq-darwin-arm64
sudo mv tmq-darwin-arm64 /usr/local/bin/tmq
```

### مک‌اواس (اینتل)
```bash
wget https://github.com/azolfagharj/tmq/releases/latest/download/tmq-darwin-amd64
chmod +x tmq-darwin-amd64
sudo mv tmq-darwin-amd64 /usr/local/bin/tmq
```

### ویندوز (AMD64)
1. دانلود: https://github.com/azolfagharj/tmq/releases/latest/download/tmq-windows-amd64.exe
2. تغییر نام به `tmq.exe`
3. اضافه کردن به PATH

## نصب دستی

1. **دانلود** باینری مناسب برای سیستم
2. **قابل اجرا کردن** (لینوکس/مک):
   ```bash
   chmod +x tmq-*
   ```
3. **تغییر نام** به `tmq` (یا `tmq.exe` در ویندوز):
   ```bash
   mv tmq-linux-amd64 tmq
   ```
4. **انتقال به PATH** (اختیاری ولی توصیه‌شده):
   ```bash
   sudo mv tmq /usr/local/bin/
   ```

## کامپایل از سورس

در صورت تمایل به کامپایل از سورس:

### پیش‌نیازها
- Go 1.23 یا بالاتر

### مراحل کامپایل
```bash
git clone https://github.com/azolfagharj/tmq.git
cd tmq
go build -o bin/tmq ./cmd/tmq
```

## تأیید نصب

پس از نصب، عملکرد tmq را بررسی کنید:

```bash
tmq --version
# باید نشان دهد: tmq version 1.0.1

tmq --help
# باید متن راهنما را نشان دهد
```

## نیازمندی‌های سیستم

- **سیستم‌عامل**: لینوکس، مک‌اواس، ویندوز
- **معماری**: AMD64، ARM64
- **حافظه**: حداقلی (با کمتر از ۱۰ مگابایت رم کار می‌کند)
- **فضای ذخیره**: حدود ۵ مگابایت برای باینری
- **بدون وابستگی خارجی** — کاملاً مستقل

## به‌روزرسانی

برای به‌روزرسانی به آخرین نسخه:

1. باینری جدید را از [releases](https://github.com/azolfagharj/tmq/releases) دانلود کنید
2. باینری قبلی را جایگزین کنید
3. اطمینان حاصل کنید که قابل اجراست (`chmod +x tmq`)

## عیب‌یابی

### دسترسی رد شد
اگر هنگام اجرای tmq خطای «permission denied» می‌گیرید:
```bash
chmod +x /path/to/tmq
```

### دستور پیدا نشد
اگر دستور `tmq` پیدا نشود:
- مطمئن شوید در PATH است: `echo $PATH`
- یا از مسیر کامل استفاده کنید: `/usr/local/bin/tmq`

### مشکل دانلود
اگر wget کار نکرد، از curl استفاده کنید:
```bash
curl -L -o tmq https://github.com/azolfagharj/tmq/releases/latest/download/tmq-linux-amd64
```
