package main
import "fmt"

func half(num int) (int,bool)  {
	return num/2, num%2==0
}
func main(){

	var num int
	fmt.Println("Enter a number : ")
	fmt.Scan(&num)
	fmt.Println(half(num))
}
