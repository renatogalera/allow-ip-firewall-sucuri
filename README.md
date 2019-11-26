# allow-ip-firewall-sucuri

Allow automatically your dynamic IP in Firewall https://sucuri.net

Add API_KEY and API_SECRET in config.env

- Run on Docker Container

```
git clone https://github.com/renatogalera/allow-ip-firewall-sucuri

cd allow-ip-firewall-sucuri

#First edit/create config.env 

cp config.env.example config.env

vim config.env

docker build -t allow-ip-firewall-sucuri .

docker run -dit --restart unless-stopped --name allow-ip-firewall-sucuri allow-ip-firewall-sucuri

#Check logs

docker logs -f allow-ip-firewall-sucuri

```

- Run locally

```
git clone https://github.com/renatogalera/allow-ip-firewall-sucuri

cd allow-ip-firewall-sucuri

#First edit/create config.env 

cp config.env.example config.env

vim config.

#Before delete lines 111/121 for looping docker version

go build main.go

./main

#Add on linux crontab

crontab -l | { cat; echo "*/3 * * * * $(pwd)/main"; } | crontab -
```
