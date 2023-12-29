package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	mockRemoteAddr := "127.0.0.1:80"
	reqs := []func() *http.Request{
		func() (req *http.Request) {
			req = httptest.NewRequest("GET", "https://ipku.misterabdul.moe/entah", nil)
			req.RemoteAddr = mockRemoteAddr
			return
		}, func() (req *http.Request) {
			req = httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
			req.RemoteAddr = mockRemoteAddr
			return
		}, func() (req *http.Request) {
			req = httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
			req.RemoteAddr = mockRemoteAddr
			req.Header.Add("User-Agent", "curl/8.2.1")
			return
		}, func() (req *http.Request) {
			req = httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
			req.RemoteAddr = mockRemoteAddr
			req.Header.Add("Accept", "application/json")
			return
		}}
	testResponses := []func(body []byte) int{
		func(body []byte) int {
			if string(body) != "404 page not found\n" {
				t.Errorf("Wrong html response in index handler")
			}
			return http.StatusNotFound
		}, func(body []byte) int {
			if string(body) == "" {
				t.Errorf("Empty html response in index handler")
			}
			return http.StatusOK
		}, func(body []byte) int {
			if string(body) != "127.0.0.1\n" {
				t.Errorf("Wrong curl response in index handler")
			}
			return http.StatusOK
		}, func(body []byte) int {
			bodyJson := make(map[string]interface{})
			if err := json.Unmarshal(body, &bodyJson); err != nil {
				t.Errorf("Error decoding json response in index handler: %s\n", err.Error())
			}
			if bodyJson["ipAddress"] != "127.0.0.1" {
				t.Errorf("Wrong json response in index handler")
			}

			return http.StatusOK
		}}

	for i, req := range reqs {
		res := httptest.NewRecorder()
		handleIndex(res, req())
		body, err := io.ReadAll(res.Result().Body)
		if err != nil {
			t.Errorf("Error reading response in index handler: %s\n", err.Error())
		}
		if res.Result().StatusCode != testResponses[i](body) {
			t.Errorf("Wrong status code in index handler.\n")
		}
	}
}

func TestHandleAbout(t *testing.T) {
	reqs := []func() *http.Request{
		func() (req *http.Request) {
			return httptest.NewRequest("GET", "https://ipku.misterabdul.moe/about/entah", nil)
		}, func() (req *http.Request) {
			return httptest.NewRequest("GET", "https://ipku.misterabdul.moe/about", nil)
		}, func() (req *http.Request) {
			req = httptest.NewRequest("GET", "https://ipku.misterabdul.moe/about", nil)
			req.Header.Add("User-Agent", "curl/8.2.1")
			return
		}, func() (req *http.Request) {
			req = httptest.NewRequest("GET", "https://ipku.misterabdul.moe/about", nil)
			req.Header.Add("Accept", "application/json")
			return
		}}
	testResponses := []func(body []byte) int{
		func(body []byte) int {
			if string(body) != "404 page not found\n" {
				t.Errorf("Wrong html response in about handler")
			}
			return http.StatusNotFound
		}, func(body []byte) int {
			if string(body) == "" {
				t.Errorf("Empty html response in about handler")
			}
			return http.StatusOK
		}, func(body []byte) int {
			expected := "Name        : IPKU\n" +
				"Description : Get the public IP address of the client.\n" +
				"Version     : 1.0.0\n" +
				"Author      : Abdul Pasaribu <mail@misterabdul.moe>\n" +
				"Repository  : https://github.com/misterabdul/ipku\n"
			if string(body) != expected {
				t.Errorf("Wrong curl response in about handler")
			}
			return http.StatusOK
		}, func(body []byte) int {
			bodyJson := make(map[string]interface{})
			if err := json.Unmarshal(body, &bodyJson); err != nil {
				t.Errorf("Error decoding json response in about handler: %s\n", err.Error())
			}
			authorData, ok := bodyJson["author"].(map[string]interface{})
			if !ok {
				t.Errorf("Error decoding author json data in about handler")
			}
			jsonAssert := bodyJson["name"] == ingfo.Name && bodyJson["description"] == ingfo.Description &&
				bodyJson["version"] == ingfo.Version && bodyJson["repository"] == ingfo.Repo &&
				authorData["name"] == ingfo.Author && authorData["email"] == ingfo.AuthorEmail
			if !jsonAssert {
				t.Errorf("Wrong json response in about handler")
			}

			return http.StatusOK
		}}

	for i, req := range reqs {
		res := httptest.NewRecorder()
		handleAbout(res, req())
		body, err := io.ReadAll(res.Result().Body)
		if err != nil {
			t.Errorf("Error reading response in about handler: %s\n", err.Error())
		}
		if res.Result().StatusCode != testResponses[i](body) {
			t.Errorf("Wrong status code in about handler.\n")
		}
	}
}

func TestGetIp(t *testing.T) {
	mockRemotes := []string{"172.17.0.4:80", "172.17.0.4:80", "[2001:0db8::1001]:80", "[2001:0db8::1001]:80"}
	mockForwardedIps := []string{"127.0.0.1", "192.168.0.1", "::1", "2001:db8:3333:4444:5555:6666:7777:8888"}

	for _, mockRemote := range mockRemotes {
		req := httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
		req.RemoteAddr = mockRemote

		ip, err := getIp(req)
		if err != nil {
			t.Errorf("Error in get ip: %s\n", err.Error())
		}
		expectedIp, _, err := net.SplitHostPort(mockRemote)
		if err != nil {
			t.Errorf("Error in split host port mock remote ip: %s\n", err.Error())
		}
		if ip != expectedIp {
			t.Errorf("Get IP is not working as expected.\n")
		}
	}

	for i, mockForwardedIp := range mockForwardedIps {
		req := httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
		req.RemoteAddr = mockRemotes[i]
		req.Header.Add("X-Forwarded-For", mockRemotes[i]+", "+mockForwardedIp)

		ip, err := getIp(req)
		if err != nil {
			t.Errorf("Error in get ip: %s\n", err.Error())
		}
		if ip != mockForwardedIp {
			t.Errorf("Get IP is not working as expected.\n")
		}
	}
}

func TestIsCurl(t *testing.T) {
	isCurlReq := httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
	isCurlReq.Header.Add("User-Agent", "curl/8.2.1")
	if !isCurl(isCurlReq) {
		t.Errorf("Is curl not working as expected.\n")
	}

	isntCurlReq := httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
	isntCurlReq.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")
	if isCurl(isntCurlReq) {
		t.Errorf("Is curl not working as expected.\n")
	}
}

func TestWantsJson(t *testing.T) {
	wantsJsonReq := httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
	wantsJsonReq.Header.Add("Accept", "application/json")
	if !wantsJson(wantsJsonReq) {
		t.Errorf("Wants JSON not working as expected.\n")
	}

	notWantJsonReq := httptest.NewRequest("GET", "https://ipku.misterabdul.moe", nil)
	notWantJsonReq.Header.Add("Accept", "*/*")
	if wantsJson(notWantJsonReq) {
		t.Errorf("Wants JSON not working as expected.\n")
	}
}

func TestRenderIpHtmlResponse(t *testing.T) {
	mockIp := "127.0.0.1"
	res := httptest.NewRecorder()
	renderIpHtmlResponse(res, mockIp)
	body, err := io.ReadAll(res.Result().Body)
	if err != nil {
		t.Errorf("Error in render ip html response: %s\n", err.Error())
	}
	if string(body) != getHtmlDefault("IPKU", fmt.Sprintf(
		`<table>
			<tr><th>IP Address</th></tr>
			<tr><td>%s</td><tr>
		</table>`, mockIp,
	)) {
		t.Errorf("Wrong content in render ip html response.\n")
	}
}

func TestRenderAboutHtmlResponse(t *testing.T) {
	mockIngfo := metainfo_t{
		Name:        "IPKU",
		Description: "Get the public IP address of the client.",
		Version:     "1.0.0",
		Author:      "Abdul Pasaribu",
		AuthorEmail: "mail@misterabdul.moe",
		Repo:        "https://github.com/misterabdul/ipku",
	}

	res := httptest.NewRecorder()
	renderAboutHtmlResponse(res, mockIngfo)
	body, err := io.ReadAll(res.Result().Body)
	if err != nil {
		t.Errorf("Error in render about html response: %s\n", err.Error())
	}
	if string(body) != getHtmlDefault("About IPKU", fmt.Sprintf(
		`<table>
			<tr><th>Name</th><td>%s</td></tr>
			<tr><th>Description</th><td>%s</td></tr>
			<tr><th>Version</th><td>%s</td></tr>
			<tr><th>Author</th><td>%s &lt;%s&gt;</td></tr>
			<tr><th>Repository</th><td><a href="%s">%s</a></td></tr>
		</table>`,
		mockIngfo.Name,
		mockIngfo.Description,
		mockIngfo.Version,
		mockIngfo.Author, mockIngfo.AuthorEmail,
		mockIngfo.Repo, mockIngfo.Repo,
	)) {
		t.Errorf("Wrong content in render about html response.\n")
	}
}

func TestGetHtmlDefault(t *testing.T) {
	defaultHtmlEmpty := getHtmlDefault("", "")
	if defaultHtmlEmpty == "" {
		t.Errorf("Default html is empty.\n")
	}

	defaultHtmlTitle := getHtmlDefault("Entah", "")
	if defaultHtmlTitle == "" || defaultHtmlTitle == defaultHtmlEmpty {
		t.Errorf("Default html with title is not working.\n")
	}

	defaultHtmlTitleContent := getHtmlDefault("Entah", "Entah")
	if defaultHtmlTitleContent == "" || defaultHtmlTitleContent == defaultHtmlTitle {
		t.Errorf("Default html with title and content is not working.\n")
	}
}
