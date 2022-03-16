package main

import ( 
  "io"
  "log" 
  "bytes"
  "fmt" 
  "os"
  "encoding/csv" 
  "net/smtp"
  "text/template"
  "github.com/joho/godotenv"
)
 
/*
  function for reading CSV file and return receivers email
*/
func getReceiversEmails() []string {
  // read CSV
  f, err := os.Open("static/emails.csv")
  if err != nil { 
      log.Fatal(err)
  }

  r := csv.NewReader(f)

  var emails []string

  for {
      record, err := r.Read()
      if err == io.EOF {
          break
      }
      if err != nil {
          log.Fatal(err)
      }

      for value := range record { 
          str := fmt.Sprint(record[value])
          emails = append(emails, str) 
      }
  }

  return emails
}
 
/*
  main function
*/
func main() {
  // loading env file
  envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("Could not load .env file")
		os.Exit(1)
	}
  
  // sender data
  from := os.Getenv("EMAIL")
  password := os.Getenv("PASSWORD")
 
  // receivers email addresss
  to := getReceiversEmails()

  // smtp server configuration
  smtpHost := "smtp.gmail.com" 
  smtpPort := "587"

  // authentication
  auth := smtp.PlainAuth("", from, password, smtpHost)

  // email body
  t, _ := template.ParseFiles("static/template.html")

  var body bytes.Buffer

  mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
  body.Write([]byte(fmt.Sprintf("Subject: " + os.Getenv("SUBJECT") + " \n%s\n\n", mimeHeaders)))

  t.Execute(&body, struct {}{})

  // sending email
  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
  if err != nil {
    fmt.Println(err)
    return
  }
  
  fmt.Println("Email Sent!") 
}
