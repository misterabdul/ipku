package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type h map[string]interface{}

type flags_t struct {
	Port        *string
	BehindProxy *bool
}

type metainfo_t struct {
	Name        string
	Description string
	Version     string
	Author      string
	AuthorEmail string
	Repo        string
}

var (
	flags = &flags_t{
		Port:        new(string),
		BehindProxy: new(bool),
	}
	ingfo = metainfo_t{
		Name:        "IPKU",
		Description: "Get the public IP address of the client.",
		Version:     "1.0.0",
		Author:      "Abdul Pasaribu",
		AuthorEmail: "mail@misterabdul.moe",
		Repo:        "https://github.com/misterabdul/ipku",
	}
)

func main() {
	flags = flagsGet()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/about", handleAbout)

	fmt.Printf("IPKU server running on port: %s\n", *flags.Port)
	if err := http.ListenAndServe(":"+*flags.Port, nil); err != nil {
		log.Fatal(err)
	}
}

func flagsGet() (flags *flags_t) {
	flags = &flags_t{
		Port:        flag.String("port", "80", "set port for the server"),
		BehindProxy: flag.Bool("behind-proxy", false, "set wheter server is running behind a proxy or not"),
	}

	flag.CommandLine.Parse(os.Args[1:])

	return
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "" && req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	ip, err := getIp(req)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, "Error: %s.\n", err)
		return
	}

	if isCurl(req) {
		fmt.Fprintf(res, "%s\n", ip)
		return
	}

	if wantsJson(req) {
		if err := json.NewEncoder(res).Encode(h{"ipAddress": ip}); err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(res, "Error: %s\n", err.Error())
		}
		return
	}

	renderIpHtmlResponse(res, ip)
}

func handleAbout(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/about" {
		http.NotFound(res, req)
		return
	}

	if isCurl(req) {
		fmt.Fprintf(res, "Name        : %s\n"+
			"Description : %s\n"+
			"Version     : %s\n"+
			"Author      : %s <%s>\n"+
			"Repository  : %s\n",
			ingfo.Name,
			ingfo.Description,
			ingfo.Version,
			ingfo.Author,
			ingfo.AuthorEmail,
			ingfo.Repo,
		)
		return
	}

	if wantsJson(req) {
		if err := json.NewEncoder(res).Encode(h{
			"name":        ingfo.Name,
			"description": ingfo.Description,
			"version":     ingfo.Version,
			"author": h{
				"name":  ingfo.Author,
				"email": ingfo.AuthorEmail},
			"repository": ingfo.Repo,
		}); err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(res, "Error: %s\n", err.Error())
		}
		return
	}

	renderAboutHtmlResponse(res, ingfo)
}

func getIp(req *http.Request) (ip string, err error) {
	if forwardedIps := strings.Split(req.Header.Get("X-Forwarded-For"), ","); len(forwardedIps) > 0 {
		pattern, err := regexp.Compile(`[^a-f0-9\.:]+`)
		if err != nil {
			return "", err
		}

		if *flags.BehindProxy && len(forwardedIps) > 1 {
			ip = pattern.ReplaceAllString(forwardedIps[len(forwardedIps)-2], "")
		} else {
			ip = pattern.ReplaceAllString(forwardedIps[len(forwardedIps)-1], "")
		}
	}
	if ip == "" {
		if ip, _, err = net.SplitHostPort(req.RemoteAddr); err != nil {
			return "", err
		}
	}

	if netIP := net.ParseIP(ip); netIP != nil {
		if ip = netIP.String(); ip != "::1" {
			return ip, nil
		}
		return "127.0.0.1", nil
	}

	return "", errors.New("IP not found")
}

func isCurl(req *http.Request) bool {
	for name, headers := range req.Header {
		if name != "User-Agent" {
			continue
		}

		for _, header := range headers {
			if strings.Contains(header, "curl") {
				return true
			}
		}
	}

	return false
}

func wantsJson(req *http.Request) bool {
	for name, headers := range req.Header {
		if name != "Accept" {
			continue
		}

		for _, header := range headers {
			if strings.Contains(header, "application/json") {
				return true
			}
		}
	}

	return false
}

func renderIpHtmlResponse(res http.ResponseWriter, ip string) {
	fmt.Fprint(res, getHtmlDefault("IPKU", fmt.Sprintf(
		`<table>
			<tr><th>IP Address</th></tr>
			<tr><td>%s</td><tr>
		</table>`, ip)))
}

func renderAboutHtmlResponse(res http.ResponseWriter, metainfo metainfo_t) {
	fmt.Fprint(res, getHtmlDefault("About IPKU", fmt.Sprintf(
		`<table>
			<tr><th>Name</th><td>%s</td></tr>
			<tr><th>Description</th><td>%s</td></tr>
			<tr><th>Version</th><td>%s</td></tr>
			<tr><th>Author</th><td>%s &lt;%s&gt;</td></tr>
			<tr><th>Repository</th><td><a href="%s">%s</a></td></tr>
		</table>`,
		metainfo.Name,
		metainfo.Description,
		metainfo.Version,
		metainfo.Author, metainfo.AuthorEmail,
		metainfo.Repo, metainfo.Repo)))
}

func getHtmlDefault(title, content string) (htmlString string) {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>%s</title>
		<style>
			* {
				font-family: sans-serif;
				font-size: 1.2rem;
			}
			@media (prefers-color-scheme: dark) {
			    body {
			        background-color: #121212;
			        color: white;
			    }
				table, th, td {
					border: 5px solid white;
					border-collapse: collapse;
				}
			}
			@media (prefers-color-scheme: light) {
				table, th, td {
					border: 5px solid black;
					border-collapse: collapse;
				}
			}
			main {
				display: grid;
				place-items: center;
				min-height: 100vh;
			}
			th, td {
				padding: 0.5rem;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<main>
			<section>%s</section>
		</main>
	</body>
</html>`, title, content)
}
