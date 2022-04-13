package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)
	var mensaje_retorno string

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				log.Println(clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		mensaje_retorno = obtenerBroadcast(clientRequest)

		// Responding to the client request
		if _, err = con.Write([]byte(mensaje_retorno + "\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}

func obtenerBroadcast(mensaje string) string {

	var mascara net.IPMask
	var ip net.IP
	var ip_decimal [4]int
	var ip_bytes [4]byte
	var partes_ip []string
	var direccion_broadcast net.IP
	var direccion_red net.IP
	var mensaje_partes []string
	var numero_mascara int

	const mensaje_error string = "Error en la dirección recibida"

	//Antes que nada, remover el salto de línea
	mensaje = strings.TrimSuffix(mensaje, "\n")

	//Paso 1, partir el mensaje

	mensaje_partes = strings.Split(mensaje, "/")

	//Paso 2, obtener los unos de la máscara
	numero_mascara, err := strconv.Atoi(mensaje_partes[1])

	if err != nil {
		fmt.Print("Error al obtener los unos de la máscara")
		return mensaje_error
	}

	//Paso 3, construir una máscara de 32 bits

	mascara = net.CIDRMask(numero_mascara, 32)

	//Paso 4, partir la parte del mensaje con la dirección IP

	partes_ip = strings.Split(mensaje_partes[0], ".")

	//Paso 5, convertir el fragmento a decimal

	for contador := 0; contador < 4; contador++ {

		ip_decimal[contador], err = strconv.Atoi(partes_ip[contador])

		if err != nil {
			fmt.Println("Error al convertir las partes de la ip")
			return mensaje_error
		}
	}

	//Paso 6, construir una dirección IP de verdad

	for contador := 0; contador < 4; contador++ {
		ip_bytes[contador] = byte(ip_decimal[contador])
	}

	ip = net.IP(ip_bytes[:])

	//Paso 7, construir la dirección broadcast

	direccion_broadcast = net.IP(make([]byte, 4))

	//Operación OR entre la IP y el inverso de la máscara

	for contador := range ip {
		direccion_broadcast[contador] = ip[contador] | ^mascara[contador]
	}

	//Paso 8, construir la dirección de red
	direccion_red = net.IP(make([]byte, 4))

	//Operación AND entre la IP y la máscara

	for contador := range ip {
		direccion_red[contador] = ip[contador] & mascara[contador]
	}

	//Paso 9, retornar

	return "Su dirección broadcast es: " + direccion_broadcast.String() + " Y su dirección de red es: " + direccion_red.String()

}
