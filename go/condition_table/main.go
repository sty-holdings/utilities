package main

import (
	"fmt"
)

type ConditionTable map[any]FunctionTable
type FunctionTable map[any]func(any) any

func main() {

	a := NewFunctionTable()
	a["Cond_1"] = ConditionOne()
	a["Cond_2"] = ConditionTwo()

	fmt.Printf("Cond_1: %v\n", a["Cond_1"](2))
	fmt.Printf("Cond_2: %v\n", a["Cond_2"]("CHECK"))
	fmt.Printf("Cond_2: %v\n", a["Cond_2"]("check"))

	b := NewConditionTable()
	b["TR_IN"] = a

	fmt.Printf("TR_IN: %v\n", b["TR_IN"]["Cond_1"](2))
	fmt.Printf("TR_IN: %v\n", b["TR_IN"]["Cond_2"]("CHECK"))
	fmt.Printf("TR_IN: %v\n", b["TR_IN"]["Cond_2"]("check"))

	rows := []any{"TR_IN", "TR_OUT"}                         // List of look up values
	columns := []any{"Cond_1", "Cond_2", "Cond_3", "Cond_4"} // List of look up values
	functionList := []func(any) any{ConditionTwo(), ConditionOne(), ConditionOne(), ConditionOne(), ConditionOne(), ConditionTwo(), ConditionOne(), ConditionTwo()}

	fmt.Println("===========================")

	LoadConditions(rows, columns, functionList)

}

func LoadConditions(rows, columns []any, functionList []func(any) any) {

	count := 0
	ct := NewConditionTable()
	for _, i := range rows {
		ft := NewFunctionTable()
		for j := 0; j < len(columns); j++ {
			fmt.Println(functionList[count])
			ft[columns[j]] = functionList[count]
			count++
		}
		ct[i] = ft
		fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		fmt.Printf("ct[i]: %v\n", ct)
		fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	}

	fmt.Println("---------------------------")

	fmt.Println(ct["TR_IN"]["Cond_2"]("CHECK"))

	fmt.Println("---------------------------")
}

func NewConditionTable() ConditionTable {
	return ConditionTable{}
}

func NewFunctionTable() FunctionTable {
	return FunctionTable{}
}

func ConditionOne() func(nbrIn any) any {
	return func(nbrIn any) any { return nbrIn.(int) * 2 }
}

func ConditionTwo() func(stringIn any) any {
	return func(stringIn any) any {
		if stringIn == "CHECK" {
			return 1
		} else {
			return 2
		}
	}
}
