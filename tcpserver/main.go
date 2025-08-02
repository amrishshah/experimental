package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func handleClient(nfd int) {
	defer unix.Close(nfd)
	buf := make([]byte, 1024)
	for {
		n, err := unix.Read(nfd, buf)

		if err != nil {
			log.Fatalf("Failed to read Data %v", err)
		}

		if n == 0 {
			log.Println("client disconnected")
			return
		}

		request := string(buf[:n])
		fmt.Println("---- Incoming Request ----")
		fmt.Println(request)
		//unix.Write(nfd, []byte("Echo: "))
		//unix.Write(nfd, buf)

		if strings.HasPrefix(request, "GET / ") {
			unix.Write(nfd, []byte("HTTP/1.1 200 ok \r\nContent-Length:1\r\nContent-Type:text/html\r\n\r\nT"))
		} else if strings.HasPrefix(request, "POST / ") {

			contentLength := 0
			var boundary string
			lines := strings.Split(request, "\r\n")
			for _, line := range lines {
				fmt.Println("Finding Length")
				if strings.HasPrefix(line, "Content-Length: ") {
					parts := strings.Split(line, ": ")
					contentLength, _ = strconv.Atoi(strings.Trim(parts[1], " "))
				}
				if strings.HasPrefix(line, "Content-Type: multipart/form-data;") {
					contentTypeSplit := strings.Split(line, "boundary=")
					if len(contentTypeSplit) == 2 {
						boundary = strings.Trim(contentTypeSplit[1], "-")
					}
					fmt.Println(boundary)
				}

			}
			body := ""
			if contentLength > 0 {
				bodySplit := strings.SplitN(request, "\r\n\r\n", 2)
				body = bodySplit[1]
				remaingData := contentLength - len(body)
				if remaingData > 0 {
					fmt.Println("Pending DATA")
					more := make([]byte, remaingData)
					unix.Read(nfd, more)
					body += string(more)
				}
			}

			if len(boundary) > 0 {
				parts := strings.Split(body, "--"+boundary)
				for _, part := range parts {
					partData := strings.SplitN(part, "\r\n\r\n", 2)
					if len(partData) != 2 {
						continue
					}

					headers := partData[0]
					content := partData[1]
					content = strings.TrimSuffix(content, "\r\n")
					content = strings.TrimSuffix(content, "------------------------")
					content = strings.TrimSpace(content)
					fmt.Println(headers)
					fmt.Println(content)

					var filename string
					for _, line := range strings.Split(headers, "\r\n") {
						if strings.Contains(line, "Content-Disposition") {
							for _, token := range strings.Split(line, ";") {
								if strings.Contains(token, "filename=") {
									filename = strings.Trim(strings.Split(token, "=")[1], `"`)
								}
							}
						}
					}
					if filename != "" {
						err := os.WriteFile(filename, []byte(content), 0644)
						if err != nil {
							log.Println("Error writing file:", err)
						} else {
							log.Println("Saved file:", filename)
						}
					}
				}
			}

			resp := fmt.Sprintf("HTTP/1.1 200 ok \r\nContent-Length:%d\r\nContent-Type:text/html\r\n\r\n%s", len(body), body)
			unix.Write(nfd, []byte(resp))
		} else {
			unix.Write(nfd, []byte("HTTP/1.1 404 Not Found\r\nContent-Length:0\r\n\r\n"))
		}
	}
}

func main() {

	//Create Socket
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("Failed to create Socket %v", err)
	}

	//bind Port
	addr := &unix.SockaddrInet4{Port: 8080}
	copy(addr.Addr[:], []byte{127, 0, 0, 1})

	if err := unix.Bind(fd, addr); err != nil {
		log.Fatalf("Failed to bind port %v", err)
	}

	//listen // 10 backlog connections
	if err := unix.Listen(fd, 10); err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	fmt.Print("listing")

	for {
		nfd, s, err := unix.Accept(fd)
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		fmt.Print(s)

		go handleClient(nfd)
	}
}
