package utils

import (
	"fmt"
	"testing"
)

func TestCheckInteger(t *testing.T) {
	InitGlobalError()
	var err *Error
	err = CheckInteger("test_int1", 10, 1, 2)
	fmt.Println(GetErrorString(err))
	err = CheckInteger("test_int2", 10, 1, 10)
	fmt.Println(GetErrorString(err))
	err = CheckInteger("test_int3", 10, 10, 20)
	fmt.Println(GetErrorString(err))
	err = CheckInteger("test_int4", 10, 20, 100)
	fmt.Println(GetErrorString(err))
}

func TestCheckStringLength(t *testing.T) {
	InitGlobalError()
	var err *Error
	err = CheckStringLength("test_str1", "hello world", 1, 2)
	fmt.Println(GetErrorString(err))
	err = CheckStringLength("test_str2", "hello world", 1, 11)
	fmt.Println(GetErrorString(err))
	err = CheckStringLength("test_str3", "hello world", 11, 20)
	fmt.Println(GetErrorString(err))
	err = CheckStringLength("test_str4", "hello world", 20, 100)
	fmt.Println(GetErrorString(err))
}

func TestCheckMobilePhoneNumber(t *testing.T) {
	InitGlobalError()
	var err *Error
	err = CheckMobilePhoneNumber("13112341234")
	fmt.Println(GetErrorString(err))
	err = CheckMobilePhoneNumber("1311234123")
	fmt.Println(GetErrorString(err))
	err = CheckMobilePhoneNumber("131123412345")
	fmt.Println(GetErrorString(err))
	err = CheckMobilePhoneNumber("23112341234")
	fmt.Println(GetErrorString(err))
}

func TestCheckIdentityCardNumber(t *testing.T) {
	InitGlobalError()
	var err *Error
	err = CheckIdentityCardNumber("411081199004235")
	fmt.Println(GetErrorString(err))
	err = CheckIdentityCardNumber("411081199004235955")
	fmt.Println(GetErrorString(err))
	err = CheckIdentityCardNumber("41108119900423")
	fmt.Println(GetErrorString(err))
	err = CheckIdentityCardNumber("41108119900423595")
	fmt.Println(GetErrorString(err))
	err = CheckIdentityCardNumber("4110811990042359555")
	fmt.Println(GetErrorString(err))
}
