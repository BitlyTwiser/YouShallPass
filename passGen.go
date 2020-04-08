package main

import (
	mRand "math/rand"
	"http/net"
	//"flags"
	"crypto/rand"
	"strings"
	"strconv"
	"log"
)

type passGen struct{
	list []string
}

func genPassword(count int) passGen {

	var list []string
	baseList := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
	baseSpecial := []string{"!","@","#","$","%","^","&","*","(",")","_","=","+","-",">","<","?",",",".","/","|","\\", " "}

	
	for _, i := range baseList{
		list = append(list, i)
		list = append(list, strings.ToUpper(i))
	}
	
	for i:=0;i<11;i++{
		list = append(list, strconv.Itoa(i))
	}
	
	for _, v := range baseSpecial{
		list = append(list, v)
	}

	return passGen{list: genPass(list, count)}

}

//Pull details from body to determine length of password.
//Will eventually have flags regarding special chars etc..
func genPasswordWeb(w http.ResponseWriter, req *http.Request){

}

//Generates random slice of bytes up to the numerical digit presented.
// The length will be a flag for length of password.. after this is passed ito the slice generator the integers therein will be utilized to obtain the values randomly.
func randGen() []byte{
	by := make([]byte, 20)
	_, err := rand.Read(by)
	if err != nil{
		log.Printf("error: %v", err)
	}
	return by
}

func mathRand(list []string) string{
	mRand.Seed(int64(randGen()[0]))	
	return list[mRand.Intn(len(list))]
}

func (l *passGen) genPass(list []string, count int) string{
	var tempPass []string
	for i:=0;i<count;i++{
		tempPass = append(tempPass, mathRand(l.list))
	}
	return strings.Join(tempPass, "")
}