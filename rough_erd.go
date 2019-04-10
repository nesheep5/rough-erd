package rough_erd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"text/template"
)

const (
	DataSourceMySQL      = "mysql"
	DataSourcePostgreSQL = "postgres"
	mapper               = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
)

type Option struct {
	Database string
	User     string
	Password string
	Port     int
	Protocol string
	Name     string
}

func Run(option *Option) error {
	fmt.Printf("run!!! option: %+v\n", option)
	conn := &ConnectInfo{
		User:     option.User,
		Password: option.Password,
		Port:     option.Port,
		Protocol: option.Protocol,
		DBName:   option.Name,
	}
	db, err := CreateDatabase(option.Database, conn)
	if err != nil {
		return err
	}
	defer db.Close()
	tables, err := db.Tables(conn.DBName)
	if err != nil {
		return err
	}
	uml := makePruntUml(tables)
	fmt.Println("-----------------------------")
	fmt.Println(uml)
	fmt.Println("-----------------------------")

	encoded := encodeAsTextFormat([]byte(uml))

	fmt.Println("Open → http://www.plantuml.com/plantuml/uml/" + encoded)
	return nil
}

const umlTemplate = `
@startuml
title Sample ERDiagram
{{range . -}}
entity "{{.Name}}" {
{{range .IDColumns -}}
{{.}}
{{end -}}
}
{{end -}}
{{range . -}}
{{$t := .Name -}}
{{range .RelayedTables -}}
{{$t}} -- {{.TableName}} :{{.IDName}}
{{end -}}
{{end -}}
@enduml
`

func makePruntUml(tables []*Table) string {
	tmpl, err := template.New("uml").Parse(umlTemplate)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, tables); err != nil {
		panic(err)
	}
	s := doc.String()
	return s
}

func encodeAsTextFormat(raw []byte) string {
	compressed := deflate(raw)
	return base64_encode(compressed)
}

func deflate(input []byte) []byte {
	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write(input)
	w.Close()
	return b.Bytes()
}

func base64_encode(input []byte) string {
	var buffer bytes.Buffer
	inputLength := len(input)
	for i := 0; i < 3-inputLength%3; i++ {
		input = append(input, byte(0))
	}

	for i := 0; i < inputLength; i += 3 {
		b1, b2, b3, b4 := input[i], input[i+1], input[i+2], byte(0)

		b4 = b3 & 0x3f
		b3 = ((b2 & 0xf) << 2) | (b3 >> 6)
		b2 = ((b1 & 0x3) << 4) | (b2 >> 4)
		b1 = b1 >> 2

		for _, b := range []byte{b1, b2, b3, b4} {
			buffer.WriteByte(byte(mapper[b]))
		}
	}
	return string(buffer.Bytes())
}
