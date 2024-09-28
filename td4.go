package main
 
import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"os/exec"
	"strconv"
	"github.com/eiannone/keyboard"
	
)

func clearConsole() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func openBinary(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Открытие: ", err)
		return ""
	}
	defer file.Close()

	fileInfo, err := file.Stat() // для размера
	if err != nil {
		fmt.Println("Информация", err)
		return ""
	}

	fileSize := fileInfo.Size()
	data := make([]byte, fileSize)

	_, err = file.Read(data)
	if err != nil {
		fmt.Println("Чтение", err)
		return ""
	}
	// в двоичную строку
	var binaryString string
	for _, b := range data {
		binaryString += fmt.Sprintf("%08b", b)
	}

	return binaryString
}

func openTxt(filename string) string {
	f, err := os.Open(filename)
	if err != nil { // Ловим ошибку
		panic(err)
		return ""
	}
 
	c, err := ioutil.ReadAll(f) // Чтение в переменную c
	f.Close()

	if err != nil {
		panic(err)
		return ""
	}
	return string(c)
}
func main() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Клавиатура", err)
	}
	defer keyboard.Close()

	var inp = openBinary("input")
	//var inp = openTxt("input")

	if len(inp)%8 != 0 {
		fmt.Printf("Cant assemble!")
		return
	} else if len(inp) > 128 {
		inp = inp[0:128]
	}

	var A uint8 = 0
	var B uint8 = 0
	var C uint8 = 0
	var Input uint8 = 10
	var Output uint8 = 0
	var command string
	var counter int = -1
	var clock bool = false
	var auto = false

	// ПРИ ЗАПУСКЕ

	clearConsole()
	fmt.Printf("Program counter: 0000\n")
	fmt.Printf("Register A: %04b Register B: %04b C Flag: %04b\n", A, B, C)
	fmt.Printf("Input Port: %04b Output Port: %04b\n", Input, Output)
	fmt.Printf("For Clock/10Hz put P, then C, I - Input")
	// горутина
	go func() {
		for {
			char, _, err := keyboard.GetKey()
			if err != nil {
				fmt.Println("Клавиатура", err)
			}

			if char == 'p' || char == 'P' {
				if auto == false {
					auto = true
				} else { 
					auto = false
				}
			}
			if char == 'c' || char == 'C' {
				clock = true
			}
			if char == 'i' || char == 'I' {
				if Input < 15 { 
					Input++
				} else {
					Input = 0 
				}
				fmt.Printf("\r    Input порт увеличен: %04b                        ", Input)
			}
		}
	}()

	for {
		if auto {
			clock = true
			time.Sleep(100 * time.Millisecond) 
		}

		if clock {
			counter += 1
			if counter > 16 {
				counter = 0
			}

			clearConsole()
			fmt.Printf("Program counter: %04b\n", counter)
			fmt.Printf("Register A: %04b Register B: %04b C Flag: %04b\n", A, B, C)
			fmt.Printf("Input Port: %04b Output Port: %04b\n", Input, Output)
			fmt.Printf("For Clock/10Hz put P, then C, I - Input")

			if len(inp) >= counter*8+8 {
				command = inp[counter*8 : counter*8+4]
				arg, err := strconv.ParseInt(inp[counter*8+4:counter*8+8], 2, 8)
				if err == nil {
					switch command {
					case "0000": // ADD A Im
						A += uint8(arg) % 16
						if A > 15 {
							A = A % 16
							C = 1
						} else {
							C = 0
						}
					case "0101": // ADD B Im
						B += uint8(arg) % 16
						if B > 15 {
							B = B % 16
							C = 1
						} else {
							C = 0
						}
					case "0011": // MOV A Im
						A = uint8(arg) % 16
						C = 0
					case "0111": // MOV B Im
						B = uint8(arg) % 16
						C = 0
					case "0001": // MOV A B
						A = B
						C = 0
					case "0100": // MOV B A
						B = A
						C = 0
					case "1111": // JMP Im
						counter = int(arg%16) - 1
						C = 0
					case "1110": // JNC Im
						if C == 0 {
							counter = int(arg%16) - 1
						}
						C = 0
					case "0010": // In A
						A = Input
						C = 0
					case "0110": // In B
						B = Input
						C = 0
					case "1001": // OUT B
						Output = B
						C = 0
					case "1011": // OUT Im
						Output = uint8(arg)
						C = 0
					case "1000": // OUT A
						Output = A
						C = 0
					case "1010": // ADD A B
						A = A + B
						if A > 15 {
							A = A % 16
							C = 1
						} else {
							C = 0
						}
					case "1100": // ADD B A
						B = B + A
						if B > 15 {
							B = B % 16
							C = 1
						} else {
							C = 0
						}
					}
				}
			}

			clock = false
		}
	}
}