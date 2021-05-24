package utils

import (
	"time"
	"unicode"
)

func CheckInteger(name string, value int, min int, max int) *Error {
	var err *Error
	err = nil
	if value < min || value > max {
		err = MakeError(100000, name, value, min, max)
	}
	return err
}

func CheckStringLength(name string, value string, minLen int, maxLen int) *Error {
	var err *Error
	err = nil
	if len(value) < minLen || len(value) > maxLen {
		err = MakeError(100001, name, value, minLen, maxLen)
	}
	return err
}

func CheckDateTimeString(name string, value string) *Error {
	_, e := time.Parse("2006-01-02 15:04:05", value)
	var err *Error
	err = nil
	if e != nil {
		err = MakeError(100002, name, value)
	}
	return err
}

func CheckMobilePhoneNumber(phoneNum string) *Error {
	isValid := true
	err := CheckStringLength("mobile phone num", phoneNum, 11, 11)
	if err != nil {
		isValid = false
	}
	if isValid {
		for _, r := range phoneNum {
			isValid = unicode.IsDigit(rune(r))
			if !isValid {
				break
			}
		}
	}
	if isValid {
		switch phoneNum[0:3] {
		case "134", "135", "136", "137", "138", "139", "150", "151", "157", "158", "159", "182", "187", "188":
			isValid = true
		case "130", "131", "132", "152", "155", "156", "185", "186":
			isValid = true
		case "133", "153", "180", "189":
			isValid = true
		default:
			isValid = false
		}
	}
	if isValid {
		return nil
	}
	err = MakeError(100010, phoneNum)
	return err
}

func CheckIdentityCardNumber(idCardNum string) *Error {
	isValid := false
	err := CheckStringLength("identity card number", idCardNum, 15, 15)
	if err == nil {
		isValid = true
	}
	if !isValid {
		err := CheckStringLength("identity card number", idCardNum, 18, 18)
		if err == nil {
			isValid = true
		}
	}
	if isValid {
		return nil
	}
	err = MakeError(100011, idCardNum)
	return err
}
