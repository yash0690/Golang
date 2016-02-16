package main
import ( 
 "fmt"
 )
func main(){
	var name string = ""
	fmt.Println("May I Know Your Name Please: ")
	fmt.Scan(&name)
	fmt.Println("Hello "+ name+" Nice Meeting You ..!!")
}
