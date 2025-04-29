```sh
sudo mkdir /etc/vcenter-bot
```

```sh
sudo cp configs/config.yaml /etc/vcenter-bot/config.yaml
```

```sh
sudo chmod 777 -R /etc/vcenter-bot
```

```sh
sudo mkdir /var/lib/vcenter-bot
```

```sh
sudo chmod 777 -R /var/lib/vcenter-bot
```

```sh
sudo chmod +x /usr/bin/vcenter-bot
```

```sh
sudo touch /var/log/vcenter-bot.log
```

```sh
sudo chmod 777 /var/log/vcenter-bot.log
```

```sh
sudo cp init/vcenter-bot.service /etc/systemd/system/vcenter-bot.service
```

```sh
sudo systemctl enable vcenter-bot
```

```sh
sudo systemctl start vcenter-bot
```