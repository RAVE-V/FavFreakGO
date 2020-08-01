package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spaolacci/murmur3"
)

func banner() {
	color.Cyan(`

	8888888888                8888888888                       888       .d8888b.   .d88888b.  
	888                       888                              888      d88P  Y88b d88P" "Y88b 
	888                       888                              888      888    888 888     888 
	8888888  8888b.  888  888 8888888 888d888 .d88b.   8888b.  888  888 888        888     888 
	888         "88b 888  888 888     888P"  d8P  Y8b     "88b 888 .88P 888  88888 888     888 
	888     .d888888 Y88  88P 888     888    88888888 .d888888 888888K  888    888 888     888 
	888     888  888  Y8bd8P  888     888    Y8b.     888  888 888 "88b Y88b  d88P Y88b. .d88P 
	888     "Y888888   Y88P   888     888     "Y8888  "Y888888 888  888  "Y8888P88  "Y88888P"   `)
	color.HiGreen("\n\t\t\t\t\t\t\t    FavFreakGo: FavFreak.py ported to GO")
	color.HiYellow("\t\t\t\t\t\t\t    github.com/RAVE-V/FavFreakGO\n\n")
}

func getDomains(domain []string) []string {
	if len(os.Args) <= 1 {
		color.HiRed("No Domains Specified")
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		domain = []string{os.Args[1]}
		fmt.Println(os.Args)
	} else {
		for i, val := range os.Args {
			if i != 0 {
				if strings.HasPrefix(val, "http://") == false && strings.HasPrefix(val, "https://") == false {
					val = "https://" + val
				}
				if val[len(val)-1] == '/' {
					val = val + "favicon.ico"
				} else {
					val = val + "/favicon.ico"
				}
				domain = append(domain, val)
			}
		}
	}
	return domain
}

func StandBase64(braw []byte) []byte {
	//https://github.com/Becivells/iconhash/blob/0eabbf376e7812050809dd07847bf2e7224b580b/config.go#L219
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}

func Mmh3Hash32(raw []byte) string {
	//https://github.com/Becivells/iconhash/blob/0eabbf376e7812050809dd07847bf2e7224b580b/config.go#L210
	var h32 hash.Hash32 = murmur3.New32()
	h32.Write([]byte(raw))
	/*if objx.Value().IsUint32(h32) {
		return fmt.Sprintf("%d", h32.Sum32())
	}*/
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

func calHash(resp *http.Response) string {
	content, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	favbase := StandBase64(content)
	hash := Mmh3Hash32(favbase)
	return hash
}

func downloadFavicon(domain []string) map[string][]string {
	hashMap := make(map[string][]string)
	for _, url := range domain {
		resp, err := http.Get(url)
		if err != nil {
			color.HiRed("[ERR] Not Fetched %s", url)
			continue
		} else {
			color.HiGreen("[INFO] Fetched %s", url)
		}
		hash := calHash(resp)
		if _, found := hashMap[hash]; found {
			hashMap[hash] = append(hashMap[hash], url)
		} else {
			hashMap[hash] = []string{url}
		}
		//fmt.Println("\n", hashMap)
		defer resp.Body.Close()
	}
	return hashMap

}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func printHashes(hashMap map[string][]string) {
	color.Blue("\n----------- Hashes Found -------------\n")
	for key := range hashMap {
		color.Cyan("------- Favion Hash : %s--------\n", key)
		for _, value := range hashMap[key] {
			color.White(value + "\n")
		}
	}
}

func main() {
	var domain []string
	banner()
	domain = getDomains(domain)
	//fmt.Println("domain ", domain)
	hashMap := downloadFavicon(domain)
	printHashes(hashMap)

}
