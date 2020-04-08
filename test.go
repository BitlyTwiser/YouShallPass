package main

import (
	mRand "math/rand"
	//"http/net"
	//"flags"
	"crypto/rand"
	"strings"
	"strconv"
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"time"
)
var (
	w bool = true
	list []string
	bU []string
	baseList = []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
	baseSpecial = []string{"!","@","#","$","%","^","&","*","(",")","_","=","+","-",">","<","?",",",".","/","|","\\", " "}
	data []string
)
type passGen struct{
	list []string
	specialList []string
	baseUpper []string
}

func main(){
	if w{
		pass := wordlist()
		fmt.Println(pass.genPass(20))
	} else {
		pass := getChars()
		fmt.Println(pass.genPass(20))
	}			
}

func getChars() passGen {
	
	for _, i := range baseList{
		list = append(list, i)
		list = append(list, strings.ToUpper(i))
		bU = append(bU, strings.ToUpper(i))
	}
	
	for i:=0;i<11;i++{
		list = append(list, strconv.Itoa(i))
	}
	
	for _, v := range baseSpecial{
		list = append(list, v)
	}

	return passGen{list: list, specialList: baseSpecial, baseUpper: bU}

}

//Obtains words from wordlist and inserts into slice to obtain random values.
func wordlist() passGen {

	content, err := os.Open("/home/smallz/Documents/gitclones/YouShallPass/parsed.csv")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer content.Close()

	lines, err := csv.NewReader(content).ReadAll()
	if err != nil {
		log.Fatalf("Error: %v")
		os.Exit(0)
	}

	for _, line := range lines {
		data = append(data, line[0])
	}

	for _, char := range getChars().list{
		data = append(data, char)
	}

	return passGen{list: data, specialList: baseSpecial, baseUpper: bU}
	
}
//Pull details from body to determine length of password.
//Will eventually have flags regarding special chars etc..
// func genPasswordWeb(w http.ResponseWriter, req *http.Request){

// }
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

//Validates that a given password holds a letter, number, and special char.
func (l *passGen) passValidation(pass string, count int) string{
	mRand.Seed(time.Now().UnixNano())

	if len(pass) > count{
		pass = pass[0:count]
	}
	for _, val := range l.specialList{		
		if strings.Contains(pass, val){
			return pass
		} else {
			pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass) - 7)]), l.mathRand(l.specialList), -1)
			return pass
		}
	}
	for _, val := range l.baseUpper{		
		if strings.Contains(pass, val){
			return pass
		} else {
			pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass) - 7)]), l.mathRand(l.specialList), -1)
			return pass
		}
	}
	return pass
}

func (l *passGen) mathRand(list []string) string{
	mRand.Seed(int64(randGen()[0]))	
	return list[mRand.Intn(len(list))]
}

func (l *passGen) genPass(count int) string{
	var tempPass []string
	for i:=0;i<count;i++{
		tempPass = append(tempPass, l.mathRand(l.list))
	}

	return l.passValidation(strings.Join(tempPass, ""), count)
}