# Palindromic Fuel Cost Calculator

Because sometimes you just need to know when filling up your tank creates a mathematically satisfying receipt. 🚗✨

## 🤔 Wait, What?

Ever notice when your fuel total comes to £32.23 and it reads the same backwards? This completely unnecessary but oddly satisfying tool finds ALL of those moments.

**Example:** At 128.9p/litre, pump exactly 25 litres → pay £32.23. Palindromic perfection! 🎯

## 🎪 Why Does This Exist?

I have an irrational need to see satisfying numbers on receipts. When filling up my partner's car, I _may or may not_ spend an embarrassing amount of time trying to hit palindromic totals. 

After one too many attempts at the perfect fill-up, I thought "there has to be a better way" and built this. Now it's a blazingly fast Go program that tells me exactly when to stop pumping. Because if you're going to be weird about receipts, might as well be _efficiently_ weird.

## 📦 Installation

```bash
# Download the binary or build it yourself
go build -o palindromic-fuel main.go
./palindromic-fuel -price=128.9 -max=100
```

## 💡 Usage

### Find all palindromes
```bash
./palindromic-fuel -price=128.9 -max=100
```

Output:
```
25 litres = £32.23 (whole number litres)
38.83 litres = £50.05 (palindromic decimal litres)  # DOUBLE PALINDROME!
42.24 litres = £54.45 (palindromic decimal litres)
50 litres = £64.46 (whole number litres)
```

### Export to CSV (for fake legitimacy)
```bash
./palindromic-fuel -price=128.9 -max=100 -csv=results.csv
```

### "I want to spend £50..."
```bash
./palindromic-fuel -price=128.9 -reverse-price=50.00 -radius=500
```

### Check multiple prices (you're in deep now)
```bash
./palindromic-fuel -batch=128.9,135.7,142.3 -max=1000
```

## 🔧 All The Flags

| Flag | What It Does |
|------|--------------|
| `-price` | Fuel price in pence |
| `-max` | How many litres to check (default: 10000) |
| `-batch` | Multiple prices, comma-separated |
| `-reverse-litres` | Find nearest palindrome to X litres |
| `-reverse-price` | Find nearest palindrome to £X |
| `-radius` | Search radius (default: 100) |
| `-csv` | Export to CSV |

## 🧮 The Clever Bit

Instead of checking every litre amount, we:
1. Generate palindromic pence values (3223, 5005, 6446...)
2. Check if they stay palindromic as pounds (£32.23 ✓, £50.05 ✓)
3. Calculate how many litres that is

**Result:** ~4.75x faster than the Node.js version. We only check ~2,266 values instead of 10,000.

## 📊 Performance

| Range | Time |
|-------|------|
| 100 litres | ~0.2ms |
| 10,000 litres | ~1.7ms |

Stupidly fast for something stupidly pointless. 🚀

## 🎯 "Legitimate" Use Cases

- **The Receipt Perfectionist** (the real reason this exists)
- **Filling partner's car** (they don't need to know you spent 5 minutes on this)
- **Party trick** (still impresses exactly nobody)
- **Personal challenge** (only ever pay palindromic amounts)
- **The smug satisfaction** when the pump hits exactly £50.05

## 🙋 FAQ

**Q: Why?**  
A: Because seeing £32.23 on a receipt is objectively more satisfying than £32.14.

**Q: Is this useful?**  
A: To my partner? No. To my inexplicable receipt number obsession? Absolutely.

**Q: Do you actually do this at petrol stations?**  
A: I refuse to answer on the grounds that it may incriminate me.

**Q: Most satisfying palindrome?**  
A: £50.05 for 38.83 litres. Double palindrome. *Chef's kiss*.

**Q: Did you really optimize this?**  
A: Yes. Because if you're going to be weird about receipts, be _efficiently_ weird.

## 🔗 My Other Projects

- [fuelaround.me](https://fuelaround.me) - UK fuel price tracker (actually useful!)
- [throwaway.cloud](https://throwaway.cloud) - Disposable email detection
- [abusedb.cloud](https://abusedb.cloud) - IP abuse contact identification

## 📄 License

Do whatever you want with this. If you make money from palindromic fuel costs, I want to hear that story.

---

Built with ❤️ and an unhealthy obsession with palindromes

*P.S. - If you actually pump exactly 38.83 litres for that £50.05 receipt, send me a photo. You're my hero.*
