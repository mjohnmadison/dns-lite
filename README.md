# DNS-LITE
A simple script for my kids raspberry pi. Instead of blocking all internet access, I wanted a lightweight way to allow only certain sites.

TL;DR this will:
    * resolve your domains in your domain.txt (newline delimited) file
    * add the domains, with a current IP to /etc/hosts
    * temporarily rename /etc/resolv.conf so no other domains will be resolved

Until my kid figures out how DNS works, this should work. Please refrain from installing in a production environment ;)

## Installing
If you are on a raspberry pi with ARMv7:
```
curl -o /home/${USER}/dns-lite blahblahlbah
```

Create a file named domains.txt (default name, check usage below to chagne name):
```

```

I added this script to rc.local to run on boot (sleep to ensure DHCP has assigned IP). You can do this with:
````
sudo sed -i "$ i\sleep 10 && /home/${USER}/dns-lite -domains /home/${USER}/domains.txt" /etc/rc.local
```

Add a cron to run every 5 minutes:
````
echo "*/5 * * * * /home/${USER}/dns-lite -domains /home/${USER}/domains.txt >/dev/null 2>&1" | sudo tee /etc/cron.d/dns-lite
```

## Usage
```
Usage of ./dns-lite:
  -disable
        Disable/Undo DNS changes
  -domains string
        Newline delimited list of domains (default "domains.txt")
  -hosts string
        Hosts file override (default "/etc/hosts")
```

### Releases
There is a binary compiled for ARMv7 in the releases for easy download/use
