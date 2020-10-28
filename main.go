package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type config struct {
	HostsFile  string
	DomainList string
}

const resolvFileD = "/etc/resolv.conf.dns-lite"

var (
	conf    config
	disable bool
)

func init() {
	flag.StringVar(&conf.HostsFile, "hosts", "/etc/hosts", "Hosts file override")
	flag.StringVar(&conf.DomainList, "domains", "domains.txt", "Newline delimited list of domains")
	flag.BoolVar(&disable, "disable", false, "Disable/Undo DNS changes")
	flag.Parse()
}

func readFile(file string) string {
	hosts, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(hosts), "\n")
}

func resolveDomains() string {
	domains := strings.Split(readFile(conf.DomainList), "\n")
	var entries []string
	for _, v := range domains {
		if v != "" { // Protect against empty lines
			ip, err := net.ResolveIPAddr("ip4", v) // IPv4 Only
			if err != nil {
				log.Println("Skipping lookup:", err)
				continue
			}
			entries = append(entries, net.Addr(ip).String()+"\t"+v+" www."+v)
		}
	}
	return strings.Join(entries, "\n")
}

func writeHostsFile() {
	hostRegexp := regexp.MustCompile(`\n(?s)## DNS-LITE BEGIN ##.*## DNS-LITE END ##`) // Strip out our entries
	orgHosts := hostRegexp.ReplaceAllString(readFile(conf.HostsFile), "")              // Get original hosts file contents
	err := ioutil.WriteFile(conf.HostsFile, []byte(orgHosts), 0644)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(conf.HostsFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if !disable {
		if _, err := f.WriteString("\n## DNS-LITE BEGIN ##\n" +
			resolveDomains() + "\n## DNS-LITE END ##\n"); err != nil {
			log.Println(err)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	var err error
	// Kind of hacky way to "disable" DNS by removing the resolv.conf file and only use the hosts file
	// We keep a copy of it to use the next time this program runs so we can re-resolve the domains
	if fileExists(resolvFileD) {
		if err = os.Rename(resolvFileD, "/etc/resolv.conf"); err != nil {
			log.Fatal(err)
		}
	}
	writeHostsFile()
	if !disable {
		if err = os.Rename("/etc/resolv.conf", resolvFileD); err != nil {
			log.Fatal(err)
		}
	}
}
