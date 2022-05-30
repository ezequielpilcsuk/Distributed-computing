package main

import (
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
	i := 1
	fmt.Println("Iniciando servidor. Conexões são encerradas após 3 minutos")
	tcpAddr, err := net.ResolveTCPAddr(ConnType, ConnHost+ConnPort)
	checkErr(err)
	listener, err := net.ListenTCP(ConnType, tcpAddr)
	checkErr(err)
	for {
		fmt.Println("Esperando pedido . . .")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, i)
		i++
	}
}

func handleClient(conn net.Conn, clientNumber int) {
	i := 1
	fmt.Printf("------- ACEITANDO CONEXÃO COM O CLIENTE %v -------\n", clientNumber)
	request := make([]byte, 1024)
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Minute * 3))

	for {
		fmt.Printf("---- TRATANDO REQUISIÇÃO %v DO CLIENTE %v ----\n", i, clientNumber)
		// Lendo requisicao
		readLen, err := conn.Read(request)
		checkErr(err)
		msg := string(request[:readLen])

		//fmt.Println("A mensagem obtida foi " + msg)

		splitMsg := strings.Split(msg, "#")
		if len(splitMsg) < 2 {
			msg := fmt.Sprintf("Requisicao invalida do cliente %v", clientNumber)
			_, err := conn.Write([]byte(msg))
			checkErr(err)
		}

		// Tratando requisicao do tipo GET
		if splitMsg[0] == "GET" {
			// Tentando ler conteudo do arquivo escolhido
			filePath := "server/files/" + strings.TrimSuffix(splitMsg[1], "\n") + ".txt"
			response, err := os.ReadFile(filePath)
			if err != nil {
				_, err := conn.Write([]byte("O arquivo selecionado nao existe"))
				checkErr(err)
			}
			result := string(response)
			_, err = conn.Write([]byte(result))
			checkErr(err)

			// Tratando requisicao do tipo QUIT
		} else if splitMsg[0] == "QUIT" {
			fmt.Printf("Encerrando conexao do cliente %v\n", clientNumber)
			msg := fmt.Sprintf("QUIT#%v\n", clientNumber)

			_, err := conn.Write([]byte(msg))
			checkErr(err)
			return

			// Tratando requisicao inválida
		} else {
			msg := fmt.Sprintf("Requisicao invalida do cliente %v", clientNumber)
			_, err := conn.Write([]byte(msg))
			checkErr(err)
		}
		fmt.Printf("---- FINALIZANDO REQUISIÇÃO %v DO CLIENTE %v ----\n", i, clientNumber)
		i++
	}

	fmt.Printf("------- FECHANDO CONEXÃO COM O CLIENTE %v -------\n", clientNumber)

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
