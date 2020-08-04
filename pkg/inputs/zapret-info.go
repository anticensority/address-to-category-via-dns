package zapretInfo

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"io"
	"bufio"
	"os"
	"strings"
	"runtime"
	"strconv"
	"regexp"

	"net"
	"bytes"
	"encoding/binary"

	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/net/idna"

  "github.com/anticensority/address-to-category-via-dns/pkg/types"
)

type blockProvider struct {
	urls []string
	rssUrl string
}

var blockProviders = []blockProvider{
	blockProvider {
		urls: []string{
			"https://sourceforge.net/p/z-i/code-0/HEAD/tree/dump.csv?format=raw",
			"https://svn.code.sf.net/p/z-i/code-0/dump.csv",
		},
		rssUrl: "https://sourceforge.net/p/z-i/code-0/feed",
	},
	blockProvider {
		urls: []string{
			"https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv",
		},
		rssUrl: "https://github.com/zapret-info/z-i/commits/master.atom",
	},
	//blockProvider {
	//	urls: []string{
	//		"https://www.assembla.com/spaces/z-i/git/source/master/dump.csv?_format=raw",
	//	},
	//	rssUrl: "https://app.assembla.com/spaces/z-i/stream.rss",
	//},
}

var get = func (url string) (*http.Response, error) {

	fmt.Println("GETting " + url)
	response, err := http.Get(url)
	fmt.Println("Got")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Negative status code: " + strconv.Itoa(response.StatusCode) + ". For url: " + url)
	}
	return response, nil
}
var getOrDie = func (url string) *http.Response {

	response, err := get(url)
	if err != nil {
		panic(err)
	}
	return response
}

func Pull() *types.AddressToIntCat {

	var (
		response *http.Response
		err error
	)
	lastUpdateMessage := ""

	updatedRegexp := regexp.MustCompile(`Updated: \d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d [+-]0000`)

	var bestProvider *blockProvider = nil
	for _, provider := range blockProviders {
		response, err := get(provider.rssUrl)
		if err != nil {
			fmt.Println("Skipping provider because of:", err)
			continue
		}
		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			match := updatedRegexp.FindString(scanner.Text())
			if match != "" {
				if lastUpdateMessage < match {
					lastUpdateMessage = match
					bestProvider = &provider
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		response.Body.Close()
		if bestProvider != nil {
			break
		}
	}
	if bestProvider == nil {
		fmt.Println("No newer dump.csv published yet!")
		os.Exit(0)
	}
	urls := bestProvider.urls
	fmt.Println("Best provider urls are:", urls)

	response = getOrDie("https://bitbucket.org/ValdikSS/antizapret/raw/master/ignorehosts.txt")
	fmt.Println("Downloaded ingoredhosts.")

	ignoredHostnames := make(map[string]bool)
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		ignoredHostnames[scanner.Text()] = true
	}
	response.Body.Close()
	fmt.Println("Parsed ingoredhosts.txt.")

	response = getOrDie("https://raw.githubusercontent.com/zapret-info/z-i/master/nxdomain.txt")
	fmt.Println("Downloaded nxdomians.")

	nxdomains := make(map[string]bool)
	scanner = bufio.NewScanner(response.Body)
	for scanner.Scan() {
		nxdomains[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	response.Body.Close()
	fmt.Println("Parsed nxdomians.")

	var lastError error
	for _, url := range urls {
		response, err = get(url)
		if err == nil {
			break
		}
		lastError = err
		response = nil
	}
	if response == nil {
		panic(lastError)
	}
	csvIn := bufio.NewReader(response.Body)
	fmt.Println("Downloaded csv.")

	_, err = csvIn.ReadString('\n')
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(transform.NewReader(csvIn, charmap.Windows1251.NewDecoder()))
	reader.Comma = ';'
	reader.FieldsPerRecord = 6
	idna := idna.New()
	hostnames := map[string]int{
		// Extremism:
		"pravdabeslana.ru": 1,
		// WordPress:
		"putinism.wordpress.com": 1,
		"6090m01.wordpress.com": 1,
		// Custom hosts
		"archive.org": 1,
		"bitcoin.org": 1,
		// LinkedIn
		"licdn.com": 1,
		"linkedin.com": 1,
		// Based on users complaints:
		"koshara.net": 1,
		"koshara.co": 1,
		"new-team.org": 1,
		"fast-torrent.ru": 1,
		"pornreactor.cc": 1,
		"joyreactor.cc": 1,
		"nnm-club.name": 1,
		"rutor.info": 1,
		"free-rutor.org": 1,
		// Rutracker complaints:
		"static.t-ru.org": 1,
		"rutrk.org": 1,

		"nnm-club.ws": 1,
		"lostfilm.tv": 1,
		"e-hentai.org": 1,
		"deviantart.net": 1, // https://groups.google.com/forum/#!topic/anticensority/uXFsOS1lQ2
		"kaztorka.org": 1, // https://groups.google.com/forum/#!msg/anticensority/vweNToREQ1o/3EbhCDjfAgAJ
	}
	ipv4        := make(map[string]int)
	ipv4Subnets := make(map[string]int)
	ipv6        := make(map[string]int)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		ips := strings.Split(record[0], " | ")
		for _, ip := range ips {
			ip = strings.Trim(ip, " \t")
			ifIpV6 := strings.ContainsAny(ip, ":")
			if ifIpV6 {
				ipv6[ip] = 1
				continue
			}
			ifSubnet := strings.ContainsAny(ip, "/")
			if ifSubnet {
				ipv4Subnets[ip] = 1
				continue
			}
			ipv4[ip] = 1
		}
		hostnamesSlice := strings.Split(record[1], " | ")
		for _, hostname := range hostnamesSlice {
			hostname = strings.Trim(hostname, " \t")
			if hostname != "" {
				hostname, err := idna.ToASCII(hostname)
				if err != nil {
					panic(err)
				}
				if strings.HasPrefix(hostname, "*.") {
					hostname = hostname[2:]
				}
				if nxdomains[hostname] || ignoredHostnames[hostname] {
					continue
				}
				if strings.HasPrefix(hostname, "www.") {
					hostname = hostname[4:]
				}
				hostnames[hostname] = 1
			}
		}
	}
	response.Body.Close()
	response = nil
	fmt.Println("Parsed csv.")
	runtime.GC()

	// Converts IP mask to 16 bit unsigned integer.
	addrToInt := func (in []byte) int {

		//var i uint16
		var i int32
		buf := bytes.NewReader(in)
		err := binary.Read(buf, binary.BigEndian, &i)
		if err != nil {
			panic(err)
		}
		return int(i)
	}
	getSubnets := func (m map[string]int) [][]int {

		keys := make([][]int, len(m))
		i := 0
		for maskedNet := range m {
			_, mask, err := net.ParseCIDR(maskedNet)
			if err != nil {
				panic(err)
			}
			keys[i] = []int{ addrToInt([]byte(mask.IP)), addrToInt([]byte(mask.Mask)) }
			i++
		}
		return keys
	}
	ipv4SubnetKeys := getSubnets(ipv4Subnets)
	results := &types.AddressToIntCat{
		Hostnames: hostnames,
		Ipv4: ipv4,
		Ipv6: ipv6,
		Ipv4Subnets: ipv4SubnetKeys,
	}
	fmt.Println(result)

	ipv4Subnets = nil

	fmt.Println("Done.")
	return results
}
