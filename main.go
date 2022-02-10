package main

import (
        "fmt"
        "sync"
)

type Selection struct {
        selections []int
        entryNum   int
}

type Selections struct {
        powerball  int
        selections Selection
}

type Resulted struct {
        entryNum int
        count    int
}

func resultToMap(result []int) map[int]bool {
        resultMap := make(map[int]bool)

        for _, res := range result {
                resultMap[res] = true
        }

        return resultMap
}

func countBalls(resmap map[int]bool, selectionz Selections, without chan Resulted, with chan Resulted, wg *sync.WaitGroup) {
        defer wg.Done()

        counter := 0
        for _, res := range selectionz.selections.selections {
                if _, found := resmap[res]; found == true {
                        counter++
                }
        }

        result := Resulted{
                entryNum: selectionz.selections.entryNum,
                count:    counter,
        }

        switch selectionz.powerball {
        case 0:
                without <- result

        case 1:
                with <- result
        }
}

func result(with chan Resulted, without chan Resulted, wg *sync.WaitGroup, done chan bool) {
        defer wg.Done()

        for {
                select {
                case a := <-with:
                        wg.Add(1)
                        if a.count >= 2 {
                                fmt.Printf("won with powerball on entry: %d with %d balls\n\n", a.entryNum, a.count)
                                wg.Done()
                        }
                case a := <-without:
                        wg.Add(1)
                        if a.count >= 5 {
                                fmt.Printf("won without powerball on entry: %d with %d balls\n\n", a.entryNum, a.count)
                                wg.Done()
                        }
                case _ = <-done:
                        wg.Add(1)
                        return
                }
        }

}

func main() {
        with := make(chan Resulted)
        without := make(chan Resulted)
        done := make(chan bool)
        var wg sync.WaitGroup

        selection := Selections{
                powerball: 1,
                selections: Selection{
                        selections: []int{3, 4, 5, 6, 7, 8, 9},
                        entryNum:   1,
                },
        }
        selectionArray := []Selections{selection}

        res := []int{1, 2, 3, 4, 5, 6, 7}
        resMap := resultToMap(res)

        for _, xD := range selectionArray {
                wg.Add(1)
                go countBalls(resMap, xD, with, without, &wg)
        }
        go result(with, without, &wg, done)
        wg.Wait()
        done <- true
}
