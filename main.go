package main

import (
        "crypto/sha256"
        "encoding/base64"
        "flag"
        "log"
        "net/http"
        "net/http/httputil"
        "net/url"
)

const Username = "admin"
const Salt = "h5h8DAkx2fpRs9BhmjDr"

func main() {
        passHash := flag.String("passHash", "", "hash of the pwd you expect")
        remote := flag.String("remote", "", "remote to proxy to ex. http://localhost")
        certFile := flag.String("cert", "cert.pem", "path of the cert file")
        keyFile := flag.String("key", "key.pem", "path of the key file")
        host := flag.String("host", ":8090", "host to run this proxy on")

        flag.Parse()

        url, err := url.Parse(*remote)
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(url)
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                username, password, ok := r.BasicAuth()
                if ok {

                        passwordHash := sha256.Sum256([]byte(password + Salt))
                        if username == Username && base64.StdEncoding.EncodeToString(passwordHash[:]) == *passHash {
                                proxy.ServeHTTP(w, r)
                                return
                        }
                }
                w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
                http.Error(w, "Unauthorized", http.StatusUnauthorized)

        })

        log.Fatal(http.ListenAndServeTLS(*host, *certFile, *keyFile, nil))
}
