package main
import ( "fmt"
)
func max(num ...int) int  {
	max :=0
	for _,temp:=range num  {
		if max<temp{
			max = temp
		}
	}
	return max
}
func main(){

	list:=[]int{60,59,07,01,63,32,22,57,12,13,45}

	fmt.Println(max(list...))
}
