package main
import (
	"fmt"
	"math"
)

func main(){
	var a,b,rem float64
	fmt.Println("Enter a Smaller Value : ")
	fmt.Scan(&a)
	fmt.Println("Enter a Larger Value :")
	fmt.Scan(&b)
	rem= math.Mod(b,a)
	fmt.Println("The Remainder is" + rem)


}

