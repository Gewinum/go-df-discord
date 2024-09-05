package utils

import (
    "strconv"
)

func ErrorPanic(err error) {
    if err != nil {
        panic(err)
    }
}

func GetNumberFirstDigits(number, digitsAmount int) int {
    firstDigits, err := strconv.Atoi(strconv.Itoa(number)[:digitsAmount])
    ErrorPanic(err)
    return firstDigits
}
