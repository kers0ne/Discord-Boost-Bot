# Local Domain Setup Guide

## For Windows:
1. Open Notepad as Administrator
2. Open file: `C:\Windows\System32\drivers\etc\hosts`
3. Add this line at the bottom:
```
127.0.0.1    kers0neserverclone.com
```
4. Save the file
5. Start a local server in your project folder:
```
python -m http.server 8000
```
6. Visit: http://kers0neserverclone.com:8000

## For Mac/Linux:
1. Open Terminal
2. Edit hosts file:
```bash
sudo nano /etc/hosts
```
3. Add this line:
```
127.0.0.1    kers0neserverclone.com
```
4. Save and exit
5. Start local server:
```bash
python3 -m http.server 8000
```
6. Visit: http://kers0neserverclone.com:8000

## For Live Domain:
To get a real domain, you need to:
1. Buy `kers0neserverclone.com` from Namecheap ($10-15/year)
2. Host on Netlify (free)
3. Connect domain to hosting

Would you like me to help with deployment setup instead?