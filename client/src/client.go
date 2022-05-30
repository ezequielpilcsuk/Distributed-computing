package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	ConnHost = "localhost"
	ConnPort = ":8080"
	ConnType = "tcp4"
)

func main() {
	// Estabelecendo conexão
	tcpAddr, err := net.ResolveTCPAddr(ConnType, ConnHost+ConnPort)
	checkErr(err)

	conn, err := net.DialTCP(ConnType, nil, tcpAddr)
	checkErr(err)
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Minute * 3))

	request := make([]byte, 1024)

	for {
		// Lendo entrada do usuário
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(
			"---- INICIANDO REQUISIÇÃO ----\n" +
				"Selecione um entre tres arquivos usando o comando GET#nome_arquivo.\n" +
				"Os arquivos disponíveis são 1, 2, e 3. Exemplo de comando: GET#1\n" +
				"Ou cancele a conexão com o comando QUIT#")

		text, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro lendo entrada: %s", err.Error())
			return
		}

		// Enviando requisição para servidor
		_, err = conn.Write([]byte(text))
		checkErr(err)

		// Lendo resposta
		readLen, err := conn.Read(request)
		checkErr(err)
		result := string(request[:readLen])

		fmt.Printf("Resposta obtida:\n%v\n"+
			"---- FINALIZANDO REQUISIÇÃO ----\n", result)

		if strings.Contains(result, "QUIT#") {
			fmt.Println("------- FINALIZANDO CONEXÃO -------.")
			return
		}

	}

}

func checkErr(err error) {
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			fmt.Fprintf(os.Stderr, "Tempo de conexão excedida: %s", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Erro fatal: %s", err.Error())
		}
		os.Exit(1)
	}
}
