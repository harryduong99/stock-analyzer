package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"math"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/duongnam99/stock-analyzer/analyzer"
	"github.com/duongnam99/stock-analyzer/config"
)

func SendAnalyzeResult(results []analyzer.AnalyzeResult) {
	subject := "Daily stocks report | " + time.Now().Format("02-01-2006")
	rootPath, _ := os.Getwd()
	var dataToMail [][]string
	var data []string
	for _, stock := range results {
		data = []string{}
		data = append(data, stock.Code)
		data = append(data, strconv.FormatFloat(stock.OpenPrice, 'f', -1, 64))
		data = append(data, strconv.FormatFloat(stock.AdjustClosedPrice, 'f', -1, 64))
		data = append(data, strconv.FormatFloat(stock.YesterdayAdjustPrice, 'f', -1, 64))
		data = append(data, strconv.FormatFloat(stock.DayBeforeYesterdayAdjustPrice, 'f', -1, 64))
		data = append(data, stock.Change)
		data = append(data, strconv.FormatFloat(math.Round(stock.FluctuatedAmount*100), 'f', -1, 64)+"%")
		data = append(data, strconv.FormatFloat(math.Round(stock.FluctuatedPrice*100), 'f', -1, 64)+"%")
		dataToMail = append(dataToMail, data)
	}
	body := ParseTemplate(rootPath+"/httpcore/mail/template/analyze_result.html", map[string]interface{}{"stocks": dataToMail})

	Send(subject, body, config.EMAIL_TARGET, true)
}

func Send(subj string, body string, toMail string, isHtml bool) {
	from := mail.Address{config.EMAIL_SENDER_NAME, config.EMAIL_SENDER}
	to := mail.Address{"", toMail}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	if isHtml {
		headers["MIME-version"] = "1.0;"
		headers["Content-Type"] = "text/html; charset=\"UTF-8\";"
	}
	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	// if isHtml {
	// 	message += body
	// } else {
	message += "\r\n" + body
	// }

	// log.Fatalln(message)

	// Connect to the SMTP Server
	servername := config.EMAIL_HOST + ":" + config.EMAIL_PORT

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", config.EMAIL_SENDER, config.EMAIL_SENDER_PASSWORD, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}

func ParseTemplate(fileName string, data interface{}) string {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		log.Fatalln(err)
	}
	return buffer.String()
}
