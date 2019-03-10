# allow-ip-firewall-sucuri

Add API_KEY and API_SECRET in config.env

- Installation

```
git clone https://github.com/renatoguilhermini/allow-ip-firewall-sucuri

cd allow-ip-firewall-sucuri

go get github.com/joho/godotenv

go build allow-ip-sucuri.go

chmod +x allow-ip-sucuri
```

On linux, add task in Crontab - Obs: Edit address folder (dir/...). The task check run every three minutes.

```
crontab -l | { cat; echo "*/3 * * * * dir/allow-ip-firewall-sucuri/allow-ip-sucuri"; } | crontab -
```
