package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"time"

	"github.com/fatih/color"
)

var request []byte
var APIRequest []byte
var WebRequest []byte
var target string
var loggedIn bool
var Profile map[string]string
var counter uint64
var _api = GetAPI()
var allow bool
var blockedEdit bool
var blockedEditInt uint64
var blockedSetInt uint64
var blockedSet bool
var bypass bool
var check bool
var blockedWeb int
var TAU string
var sleep int
var ClearConsole func()
var iter int
var claimed bool = false
var stop bool = false
var DiscURL = "WebHook URL"

var mx sync.Mutex
var wg sync.WaitGroup
var start = make(chan struct{})

var G = color.New(color.FgHiCyan, color.Bold)
var R = color.New(color.FgRed, color.Bold)
var Gr = color.New(color.FgGreen, color.Bold)
var Y = color.New(color.FgYellow, color.Bold)

func end(s int) {

	ClearConsole()

	fmt.Println()
	color.Red("▄▄▄█████▓ ██▓▄▄▄█████▓ ▄▄▄       ███▄    █ ")
	color.Red("▓  ██▒ ▓▒▓██▒▓  ██▒ ▓▒▒████▄     ██ ▀█   █ ")
	color.Red("▒ ▓██░ ▒░▒██▒▒ ▓██░ ▒░▒██  ▀█▄  ▓██  ▀█ ██▒")
	color.Red("░ ▓██▓ ░ ░██░░ ▓██▓ ░ ░██▄▄▄▄██ ▓██▒  ▐▌██▒")
	color.Red("  ▒██▒ ░ ░██░  ▒██▒ ░  ▓█   ▓██▒▒██░   ▓██░")
	color.Red("  ▒ ░░   ░▓    ▒ ░░    ▒▒   ▓▒█░░ ▒░   ▒ ▒ ")
	color.Red("    ░     ▒ ░    ░      ▒   ▒▒ ░░ ░░   ░ ▒░")
	color.Red("  ░       ▒ ░  ░        ░   ▒      ░   ░ ░ ")
	color.Red("          ░                 ░  ░         ░ ")

	fmt.Println()

	if s == 0 {
		color.Green("Successfully Claimed: " + target)
	} else if s == 1 {
		color.Red("Error ! or it closed by the Developer")
	} else if s == 3 {
		color.Red("Closed")
	}

	fmt.Println()
	color.HiBlue("By Hades, inst: @0xhades")
	fmt.Println()

	os.Exit(0)

}

func WebHook(url string, log bool) {

	if len(target) > 4 {
		if !log {
			return
		}
	}

	data := "{\"embeds\": [{\"description\": \"Swapped Successfully!\\nAttempts: " + fmt.Sprintf("%v", counter) + "\\nBy Faisal @3wv\", \"title\": \"@" + target + "\", \"image\": {\"url\": \"https://i.imgur.com/QyoizqY.jpg\"}}], \"username\": \"Titan\"}"

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Add("Content-Type", "application/json")
	transport := http.Transport{}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	client.Transport = &transport
	resp, err := client.Do(req)
	if err == nil {
		var res string
		var StatusCode int
		_ = StatusCode
		_ = res

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ = gzip.NewReader(resp.Body)
			defer reader.Close()
		default:
			reader = resp.Body
		}
		body, _ := ioutil.ReadAll(reader)
		res = string(body)

		if resp.StatusCode != 0 {
			StatusCode = resp.StatusCode
		}
		println(res)
	} else {
		panic(err)
	}

}

func getProcessOwner() string {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return strings.Replace(string(stdout), "\n", "", -1)
}

func main() {

	if runtime.GOOS == "windows" {

		ClearConsole = func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}

	} else {

		if getProcessOwner() != "root" {
			R.Println("You need to be root!")
			os.Exit(0)
		}

		maxingFdsLimit()

		ClearConsole = func() {
			print("\033[H\033[2J")
		}

	}

	ClearConsole()
	var choice string

	fmt.Println()
	color.Red("▄▄▄█████▓ ██▓▄▄▄█████▓ ▄▄▄       ███▄    █ ")
	color.Red("▓  ██▒ ▓▒▓██▒▓  ██▒ ▓▒▒████▄     ██ ▀█   █ ")
	color.Red("▒ ▓██░ ▒░▒██▒▒ ▓██░ ▒░▒██  ▀█▄  ▓██  ▀█ ██▒")
	color.Red("░ ▓██▓ ░ ░██░░ ▓██▓ ░ ░██▄▄▄▄██ ▓██▒  ▐▌██▒")
	color.Red("  ▒██▒ ░ ░██░  ▒██▒ ░  ▓█   ▓██▒▒██░   ▓██░")
	color.Red("  ▒ ░░   ░▓    ▒ ░░    ▒▒   ▓▒█░░ ▒░   ▒ ▒ ")
	color.Red("    ░     ▒ ░    ░      ▒   ▒▒ ░░ ░░   ░ ▒░")
	color.Red("  ░       ▒ ░  ░        ░   ▒      ░   ░ ░ ")
	color.Red("          ░                 ░  ░         ░ ")

	fmt.Println()
	color.HiBlue("By Hades, inst: @0xhades")
	fmt.Println()

	var receiverCookiesMap = make(map[string]string)
	var sessionid string
	var sessioned bool

	for {
		G.Print("Session ID[S] / Login [L]: ")
		fmt.Scanln(&choice)
		if strings.ToLower(choice) == "s" {
			G.Print("Enter the API SessionID: ")
			fmt.Scanln(&sessionid)
			var res HttpResponse
			Profile, res = GetProfile(sessionid)
			if strings.Contains(res.Body, "consent_required") {
				updateBTHRes := updateBTH(sessionid)
				if updateBTHRes.ResStatus != 200 {
					println(updateBTHRes.Body)
					color.Red("Error Updating Day of birth")
					time.Sleep(time.Second * 2)
					G.Print("Do you wanna try again? [y/n]: ")
					fmt.Scanln(&choice)
					if strings.ToLower(choice) != "y" {
						end(2)
					} else {
						continue
					}
				}
			}
			Profile, res = GetProfile(sessionid)
			if Profile["username"] != "" {
				for i := 0; i < len(res.Cookies); i++ {
					receiverCookiesMap[res.Cookies[i].Name] = res.Cookies[i].Value
				}
				receiverCookiesMap["sessionid"] = sessionid
				TAU = Profile["username"]
				color.Green("Logged In @" + TAU + " Successfully")
				loggedIn = true
				//time.Sleep(time.Second * 2)
				sessioned = true
				break
			} else {
				println(res.Body)
				color.Red("Error Getting Profile 2")
				time.Sleep(time.Second * 2)
				G.Print("Do you wanna try again? [y/n]: ")
				fmt.Scanln(&choice)
				if strings.ToLower(choice) != "y" {
					end(2)
				} else {
					continue
				}
			}
		} else {
			break
		}
	}

	for {
		if sessioned {
			receiverCookiesMap["sessionid"] = sessionid
			break
		}

		G.Print("Enter the username: ")
		fmt.Scanln(&TAU)
		var TAP string
		G.Print("Enter the password: ")
		fmt.Scanln(&TAP)
		var res HttpResponse

		res = login(TAU, TAP, 60*1000)

		for i := 0; i < len(res.Cookies); i++ {
			if res.Cookies[i].Name == "sessionid" {
				loggedIn = true
				println()
				color.Green("Logged In Successfully")
				color.Green("Session ID: " + res.Cookies[i].Value)
				sessionid = res.Cookies[i].Value
				_Res := HttpResponse{}
				Profile, _Res = GetProfile(sessionid)
				if strings.Contains(_Res.Body, "consent_required") || _Res.Res.StatusCode != 200 {
					updateBTHRes := updateBTH(sessionid)
					if updateBTHRes.ResStatus != 200 {
						println(updateBTHRes.Body)
						color.Red("Error Updating Day of birth")
						time.Sleep(time.Second * 2)
						G.Print("Do you wanna try again? [y/n]: ")
						fmt.Scanln(&choice)
						if strings.ToLower(choice) != "y" {
							end(2)
						} else {
							continue
						}
					}
					Profile, _Res = GetProfile(sessionid)
					if Profile["username"] == "" {
						println(_Res.Body)
						color.Red("Error Getting Profile ")
						time.Sleep(time.Second * 2)
						G.Print("Do you wanna try again? [y/n]: ")
						fmt.Scanln(&choice)
						if strings.ToLower(choice) != "y" {
							end(2)
						} else {
							continue
						}
					}
				}
			}
			receiverCookiesMap[res.Cookies[i].Name] = res.Cookies[i].Value
		}

		if strings.Contains(res.Body, "ogged_in") && loggedIn && Profile["username"] != "" {
			break
		} else {
			if strings.Contains(res.Body, "challenge_required") {

				urlRegex := regexp.MustCompile("\"api_path\": \"(.*?)\"").FindStringSubmatch(res.Body)
				var url string

				if urlRegex == nil {
					println(res.Body)
					color.Red("Getting API Path Error")
					time.Sleep(time.Second * 2)
					G.Print("Do you wanna try again? [y/n]: ")
					fmt.Scanln(&choice)
					if strings.ToLower(choice) != "y" {
						return
					} else {
						continue
					}
				}

				url = urlRegex[1]

				_headers := make(map[string]string)
				loginCookies := res.Headers.Get("set-cookie")

				if loginCookies == "" {
					color.Red("Login's set-cookie is empty")
					time.Sleep(time.Second * 2)
					G.Print("Do you wanna try again? [y/n]: ")
					fmt.Scanln(&choice)
					if strings.ToLower(choice) != "y" {
						return
					} else {
						continue
					}
				}

				CSRFRegex := regexp.MustCompile("csrftoken=(.*?);").FindStringSubmatch(loginCookies)
				//MidRegex := regexp.MustCompile("mid=(.*?);").FindStringSubmatch(loginCookies)
				var csrftoken string

				if CSRFRegex == nil {
					println(loginCookies)
					color.Red("CSRF is empty")
					time.Sleep(time.Second * 2)
					G.Print("Do you wanna try again? [y/n]: ")
					fmt.Scanln(&choice)
					if strings.ToLower(choice) != "y" {
						return
					} else {
						continue
					}
				}

				csrftoken = CSRFRegex[1]

				_headers["X-CSRFToken"] = csrftoken

				SecureResult := instRequest(url, nil, "", _headers, GetAPI(), "", res.Cookies, true, 60*1000)

				em := false
				ph := false

				var Pass bool
				var email string
				var phone string
				var emailRegex []string
				var phoneRegex []string

				if strings.Contains(SecureResult.Body, "select_verify_method") {
					if strings.Contains(SecureResult.Body, "email") {
						emailRegex = regexp.MustCompile("\"email\": \"(.*?)\"").FindStringSubmatch(SecureResult.Body)
					}
					if strings.Contains(SecureResult.Body, "phone_number") {
						phoneRegex = regexp.MustCompile("\"phone_number\": \"(.*?)\"").FindStringSubmatch(SecureResult.Body)
					}
				} else {
					choice = "0"
					Pass = true
				}

				var contactPoint string

				if !Pass {
					if phoneRegex == nil && emailRegex == nil {
						println(SecureResult.Body)
						color.Red("No Verify Methods Found")
						time.Sleep(time.Second * 2)
						G.Print("Do you wanna try again? [y/n]: ")
						fmt.Scanln(&choice)
						if strings.ToLower(choice) != "y" {
							return
						} else {
							continue
						}
					}

					if phoneRegex != nil {
						phone = phoneRegex[1]
						ph = true
					}
					if emailRegex != nil {
						email = emailRegex[1]
						em = true
					}

					if em {
						G.Println("1) email [" + email + "]: ")
					}
					if ph {
						G.Println("0) phone number [" + phone + "]: ")
					}

					G.Print("Select Method: ")
					fmt.Scanln(&choice)

					if choice == "0" {
						contactPoint = phone
					}

					if choice == "1" {
						contactPoint = email
					}

					if choice != "1" && choice != "0" {
						println(SecureResult.Body)
						color.Red("Choose a correct verify method")
						time.Sleep(time.Second * 2)
						G.Print("Do you wanna try again? [y/n]: ")
						fmt.Scanln(&choice)
						if strings.ToLower(choice) != "y" {
							return
						} else {
							continue
						}
					}

				}

				SecureResult = instRequest(url, nil, "choice="+choice, nil, GetAPI(), "", res.Cookies, true, 60*1000)

				if strings.Contains(strings.ToLower(SecureResult.Body), "contact_point") {

					G.Println("A code has been sent to " + contactPoint)

					G.Print("Security Code: ")
					fmt.Scanln(&choice)
					choice = strings.Replace(choice, " ", "", -1)

					SecureResult = instRequest(url, nil, "security_code="+choice, nil, GetAPI(), "", res.Cookies, true, 60*1000)

					if strings.Contains(strings.ToLower(SecureResult.Body), "ok") || SecureResult.Res.StatusCode == 200 {

						for i := 0; i < len(SecureResult.Cookies); i++ {
							if SecureResult.Cookies[i].Name == "sessionid" {
								sessioned = true
								loggedIn = true
								println()
								color.Green("Logged In Successfully")
								color.Green("Session ID: " + SecureResult.Cookies[i].Value)
								sessionid = SecureResult.Cookies[i].Value
								_Res := HttpResponse{}
								Profile, _Res = GetProfile(sessionid)
								if strings.Contains(_Res.Body, "consent_required") || _Res.Res.StatusCode != 200 {
									updateBTHRes := updateBTH(sessionid)
									if updateBTHRes.ResStatus != 200 {
										println(updateBTHRes.Body)
										color.Red("Error Updating Day of birth")
										time.Sleep(time.Second * 2)
										G.Print("Do you wanna try again? [y/n]: ")
										fmt.Scanln(&choice)
										if strings.ToLower(choice) != "y" {
											end(2)
										} else {
											continue
										}
									}
									Profile, _Res = GetProfile(sessionid)
									if Profile["username"] == "" {
										println(_Res.Body)
										color.Red("Error Getting Profile ")
										time.Sleep(time.Second * 2)
										G.Print("Do you wanna try again? [y/n]: ")
										fmt.Scanln(&choice)
										if strings.ToLower(choice) != "y" {
											end(2)
										} else {
											continue
										}
									}
								}

							}
							receiverCookiesMap[SecureResult.Cookies[i].Name] = SecureResult.Cookies[i].Value

						}

					} else {
						println(SecureResult.Body)
						println("Code: " + choice)
						color.Red("Sending Activation Code Error")
						time.Sleep(time.Second * 2)
						G.Print("Do you wanna try again? [y/n]: ")
						fmt.Scanln(&choice)
						if strings.ToLower(choice) != "y" {
							return
						} else {
							continue
						}
					}

				} else if SecureResult.Res.StatusCode == 200 {

					for i := 0; i < len(SecureResult.Cookies); i++ {
						if SecureResult.Cookies[i].Name == "sessionid" {
							sessioned = true
							loggedIn = true
							println()
							color.Green("Logged In Successfully")
							color.Green("Session ID: " + SecureResult.Cookies[i].Value)
							sessionid = SecureResult.Cookies[i].Value
							_Res := HttpResponse{}
							Profile, _Res = GetProfile(sessionid)
							if strings.Contains(_Res.Body, "consent_required") || _Res.Res.StatusCode != 200 {
								updateBTHRes := updateBTH(sessionid)
								if updateBTHRes.ResStatus != 200 {
									println(updateBTHRes.Body)
									color.Red("Error Updating Day of birth")
									time.Sleep(time.Second * 2)
									G.Print("Do you wanna try again? [y/n]: ")
									fmt.Scanln(&choice)
									if strings.ToLower(choice) != "y" {
										end(2)
									} else {
										continue
									}
								}
								Profile, _Res = GetProfile(sessionid)
								if Profile["username"] == "" {
									println(_Res.Body)
									color.Red("Error Getting Profile ")
									time.Sleep(time.Second * 2)
									G.Print("Do you wanna try again? [y/n]: ")
									fmt.Scanln(&choice)
									if strings.ToLower(choice) != "y" {
										end(2)
									} else {
										continue
									}
								}
							}
						}
						receiverCookiesMap[SecureResult.Cookies[i].Name] = SecureResult.Cookies[i].Value
					}

				} else {
					println(SecureResult.Body)
					println(SecureResult.Res.Status)
					color.Red("Error choosing verify method")
					time.Sleep(time.Second * 2)
					G.Print("Do you wanna try again? [y/n]: ")
					fmt.Scanln(&choice)
					if strings.ToLower(choice) != "y" {
						return
					} else {
						continue
					}
				}

			}

			if sessioned || sessionid != "" {
				break
			}

			println()
			color.Red("Error Logging into the account")
			println(res.Body)
			time.Sleep(time.Second * 2)
			G.Print("Do you wanna try again? [y/n]: ")
			fmt.Scanln(&choice)
			if strings.ToLower(choice) != "y" {
				end(2)
			} else {
				continue
			}

		}

	}

	ThreadsPerMoment := 50

	println()

	var PM string
	G.Print("Do you want to log this swapping session? (Y/N): ")
	fmt.Scanln(&PM)

	if strings.ToLower(PM) == "y" {
		bypass = true
	} else {
		bypass = false
	}

	for {
		var TPM string
		G.Print("Enter Threads (Best=30, Ultimate=100): ")
		fmt.Scanln(&TPM)

		if _, err := strconv.Atoi(TPM); err == nil && TPM != "0" && !strings.Contains(TPM, "-") {
			_int64, _ := strconv.ParseInt(TPM, 0, 64)
			ThreadsPerMoment = int(_int64)
			break
		} else {
			R.Print("Enter a correct number")
			time.Sleep(time.Second * 2)
		}
	}

	for {
		var TPM string
		G.Print("Enter Sleep (Milliseconds) (Best=1000, Ultimate=800): ")
		fmt.Scanln(&TPM)

		if _, err := strconv.Atoi(TPM); err == nil && !strings.Contains(TPM, "-") {
			_int64, _ := strconv.ParseInt(TPM, 0, 64)
			sleep = int(_int64)
			break
		} else {
			R.Print("Enter a correct number")
			time.Sleep(time.Second * 2)
		}
	}

	for {
		var TPM string
		G.Print("Enter loops (Best=2, Ultimate=10): ")
		fmt.Scanln(&TPM)

		if _, err := strconv.Atoi(TPM); err == nil && TPM != "0" && !strings.Contains(TPM, "-") {
			_int64, _ := strconv.ParseInt(TPM, 0, 64)
			iter = int(_int64)
			break
		} else {
			R.Print("Enter a correct number")
			time.Sleep(time.Second * 2)
		}
	}

	G.Print("Enter Target: ")
	fmt.Scanln(&target)

	ClearConsole()

	max := ThreadsPerMoment*5 + 2
	runtime.GOMAXPROCS(max)

	var headers = make(map[string]string)
	var payload = make(map[string]string)

	headers["Host"] = "i.instagram.com"
	headers["User-Agent"] = GetAPI().USERAGENT
	headers["Accept"] = "*/*"
	headers["Cookie2"] = "$Version=1"
	headers["X-IG-Capabilities"] = GetAPI().CAPABILITIES
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["X-IG-Connection-Type"] = "WIFI"
	headers["Accept-Language"] = "en-US"
	headers["X-FB-HTTP-Engine"] = "Liger"
	headers["Accept-Encoding"] = "gzip, deflate"
	headers["Content-Encoding"] = "gzip"
	headers["Connection"] = "Keep-Alive"

	post := make(map[string]string)
	post["username"] = target
	JSONBytes, _ := json.Marshal(post)
	payload["signed_body"] = fmt.Sprintf("SIGNATURE.%s", string(JSONBytes))

	request = parseRequest("i.instagram.com", "/api/v1/accounts/set_username/", "POST", headers, payload, "", receiverCookiesMap)

	params := url.Values{}
	params.Set("email", Profile["email"])
	params.Set("username", target)
	params.Set("biography", "By 0xhades, @0xhades")
	params.Set("external_url", "instagram.com/0xhades")
	params.Set("gender", Profile["gender"])
	if Profile["phone_number"] != "" {
		params.Set("phone_number", Profile["phone_number"])
	}

	APIRequest = parseRequest("i.instagram.com", "/api/v1/accounts/edit_profile/", "POST", headers, nil, params.Encode(), receiverCookiesMap)

	params = url.Values{}
	params.Add("email", Profile["email"])
	params.Add("username", target)
	params.Add("biography", "By 0xhades, @0xhades")
	params.Add("external_url", "instagram.com/0xhades")
	params.Add("chaining_enabled", "on")
	params.Add("phone_number", Profile["phone_number"])

	var WebHeaders = make(map[string]string)
	WebHeaders["Host"] = "i.instagram.com"
	WebHeaders["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:81.0) Gecko/20100101 Firefox/81.0"
	WebHeaders["Accept"] = "*/*"
	WebHeaders["Accept-Language"] = "ar,en-US;q=0.7,en;q=0.3"
	WebHeaders["Accept-Encoding"] = "gzip, deflate, br"
	WebHeaders["X-CSRFToken"] = "5zs9b2moqkd5ndYddCjQI1pzZuUBatGn"
	WebHeaders["X-Instagram-AJAX"] = "8895c2bab672"
	WebHeaders["X-IG-App-ID"] = "936619743392459"
	WebHeaders["X-IG-WWW-Claim"] = "hmac.AR0gbMXlxalK1rnxdtE9GBbOFaQVz8vCfl7E4EgrVs7T_RIx"
	WebHeaders["Content-Type"] = "application/x-www-form-urlencoded"
	WebHeaders["X-Requested-With"] = "XMLHttpRequest"
	WebHeaders["Origin"] = "https://i.instagram.com"
	WebHeaders["Connection"] = "Keep-Alive"
	WebHeaders["Referer"] = "https://i.instagram.com/accounts/edit/"
	WebHeaders["Cookie"] = "webSession"

	WebRequest = parseRequest("i.instagram.com", "/accounts/edit/", "POST", WebHeaders, nil, params.Encode(), nil)

	for i := 0; i < ThreadsPerMoment; i++ {
		wg.Add(1)
		go sender()
	}

	fmt.Println()
	color.Red("▄▄▄█████▓ ██▓▄▄▄█████▓ ▄▄▄       ███▄    █ ")
	color.Red("▓  ██▒ ▓▒▓██▒▓  ██▒ ▓▒▒████▄     ██ ▀█   █ ")
	color.Red("▒ ▓██░ ▒░▒██▒▒ ▓██░ ▒░▒██  ▀█▄  ▓██  ▀█ ██▒")
	color.Red("░ ▓██▓ ░ ░██░░ ▓██▓ ░ ░██▄▄▄▄██ ▓██▒  ▐▌██▒")
	color.Red("  ▒██▒ ░ ░██░  ▒██▒ ░  ▓█   ▓██▒▒██░   ▓██░")
	color.Red("  ▒ ░░   ░▓    ▒ ░░    ▒▒   ▓▒█░░ ▒░   ▒ ▒ ")
	color.Red("    ░     ▒ ░    ░      ▒   ▒▒ ░░ ░░   ░ ▒░")
	color.Red("  ░       ▒ ░  ░        ░   ▒      ░   ░ ░ ")
	color.Red("          ░                 ░  ░         ░ ")

	fmt.Println()
	color.HiBlue("By Hades, inst: @0xhades")
	fmt.Println()

	if runtime.GOOS == "windows" {
		MessageBoxPlain("TitanSwap", "Ready?")
	} else {
		color.Yellow("Click any key to start ...")
		fmt.Scanln()
	}

	close(start)
	go superVisior(&counter, ThreadsPerMoment)

	wg.Wait()

	color.Yellow("\nClick any key to exit ...")
	fmt.Scanln()

}

func sender() {

	var streams []*tls.Conn

	for i := 0; i < 2; i++ {

		raddr, err := net.ResolveTCPAddr("tcp", "i.instagram.com:443")
		Conn, err := net.DialTCP("tcp", nil, raddr)
		Conn.SetNoDelay(true)
		Conn.SetKeepAlive(true)
		Conn.SetLinger(0)

		if err != nil {
			panic(err)
		}

		cf := &tls.Config{Rand: rand.Reader, InsecureSkipVerify: true}
		ssl := tls.Client(Conn, cf)

		go reader(ssl, i)

		streams = append(streams, ssl)

	}

	Requests := [][]byte{request, APIRequest}
	for in := 0; in < 2; in++ {
		go func(init int) {
			<-start
			for {
				for i := 0; i < iter; i++ {
					_, err := streams[init].Write(Requests[init])
					if err != nil {
						if !strings.Contains(err.Error(), "use of closed connection") {
							appendToFile("titan_log", "\nin: "+strconv.Itoa(init)+", Errors: "+err.Error())
							i--
							continue
						}
						if init == 0 {
							blockedSet = true
						} else {
							blockedEdit = true
						}
						break
					}
				}
				time.Sleep(time.Millisecond * time.Duration(sleep))
			}
		}(in)
	}

}

func superVisior(c *uint64, Threads int) {
	for {
		if blockedEdit && blockedSet {
			R.Print("\nYou got blocked for spamming too many requests, reached: " + fmt.Sprintf("%v", *c))
			break
		}

		if stop {
			break
		}

		if claimed {
			R.Print("\n[+] @" + TAU + " Successfully Claimed @" + target + " - Attempts: " + fmt.Sprintf("%v", *c))
			break
		}

		Y.Print("@" + TAU + " Claiming [" + target + "] - Counter: " + fmt.Sprintf("%v", *c) + "\r")

		time.Sleep(time.Millisecond * 10)
	}

	for i := 0; i < Threads; i++ {
		wg.Done()
	}

}

func reader(ssl *tls.Conn, i int) {

	innerCounter := 0

	tp := textproto.NewReader(bufio.NewReader(ssl))

	<-start
	for {

		if claimed {
			return
		}

		line, _ := tp.ReadLine()

		if strings.Contains(line, "HTTP/1.1") {
			if strings.Contains(line, "HTTP/1.1 200") {
				mx.Lock()
				ssl.Close()
				claimed = true
				WebHook(DiscURL, bypass)
				mx.Unlock()
				return
			} else if strings.Contains(line, "HTTP/1.1 400") {
				atomic.AddUint64(&counter, 1)
				innerCounter++
			} else if strings.Contains(line, "HTTP/1.1 429") {
				if blockedEditInt > 2 && blockedSetInt > 2 {
					ssl.Close()
				}
				if i == 0 {
					atomic.AddUint64(&blockedSetInt, 1)
				}
				if i == 1 {
					atomic.AddUint64(&blockedEditInt, 1)
				}
				if blockedEdit && blockedSet {
					return
				}
			} else if strings.Contains(line, "HTTP/1.1 403") || strings.Contains(line, "HTTP/1.1 302") {
				mx.Lock()
				R.Print("Not logged in")
				stop = true
				mx.Unlock()
				ssl.Close()
				os.Exit(0)
			} else if strings.Contains(line, "HTTP/1.1 404") {
				mx.Lock()
				stop = true
				R.Print("API not found?")
				mx.Unlock()
				ssl.Close()
				os.Exit(0)
			} else {
				mx.Lock()
				appendToFile("titan_log", "Status Code: "+line+"\n")
				mx.Unlock()
			}
		}

		time.Sleep(time.Millisecond * 50)

	}
}

func appendToFile(filename string, data string) error {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}
