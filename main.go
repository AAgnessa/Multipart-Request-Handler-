package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func main(){
	go GoServer()
	GoServer1()
}

func GoServer() {
	http.HandleFunc("/",DownloadUpload)

	fmt.Println("listening 0.0.0")
	log.Fatal(http.ListenAndServe(":6600", nil))
}

func GoServer1() {
	http.HandleFunc("/status",StatusLoad)

	fmt.Println("listening 0.0.0")
	log.Fatal(http.ListenAndServe(":6601", nil))
}

func DownloadUpload (w http.ResponseWriter, r *http.Request){
	//По ключу находим  изображение
	FileImage,_,err:=r.FormFile("IMAGE")
	if err!=nil{
		fmt.Println("Error retrieving the image",err)
		http.Error(w,"Error retrieving the image",400)
		return
	}
	defer FileImage.Close()

	//По ключу находим json файл
	FileJson:=r.FormValue("JSON")
	fmt.Println(FileJson)
	

	//Создаем форму для отпраки на сервер
	url:="http://localhost:6601/status"

	// Буфер для хранения нашего тела запроса в виде байтов
	var b bytes.Buffer

	//NewWriter возвращает Writer, который генерирует multipart сообщение
	multiPartWriter:=multipart.NewWriter(&b)

	//Инициализируем поле
	//Создаем новый заголофок данных формы с указанием имени заголовка и именем файла
	fw,err:=multiPartWriter.CreateFormFile("FileImage","IMAGE.png")
	if err!=nil{
		fmt.Println("Error in CreateFormFile",err)
		http.Error(w,"Error in CreateFormFile",400)
		return
	}

	//Скопируем содержимое файла в поле
	_,err=io.Copy(fw,FileImage)// fm <- FileImage
	if err!=nil{
		fmt.Println("Error in copy image",err) 
		http.Error(w,"Error in copy image",400)
		return
	}

	//Заполняем остальные поля 
	fwj,err:=multiPartWriter.CreateFormField("FileJson")
	if err!=nil{
		fmt.Println("Error in CreateFormField",err)
		http.Error(w,"Error in CreateFormField",400)
		return
	}

	_, err = fwj.Write([]byte(FileJson))
	if err != nil {
		fmt.Println("Error writer json",err)
		http.Error(w,"Error writer json",400)
		return
	}

	//Закрываем запись данных
	multiPartWriter.Close()
	//Для отправуи формы используем дефолтного клиента
	client:=http.Client{}
	res,err:=client.Post(url,multiPartWriter.FormDataContentType(), &b)
	if err!=nil{
		fmt.Println("Error in client-post")
		http.Error(w,"Error in client-post",400)
		return
	}
	
	defer res.Body.Close()

}

func StatusLoad(w http.ResponseWriter, r *http.Request) {
	FileJson:=r.FormValue("FileJson")
	log.Println("StatusLoad",FileJson)

}