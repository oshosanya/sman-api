# SMAN API
You need to have go installed to run this project as it is written in go, and there are no binaries on this repo.

To run the project, simply run `go run main.go` in the project directory.

An HTTP Server is started on port 1323

Make a post request to the server using form-data with the following items

```aidl
name: string
position: string
branch: string
id_number: string
passport: file
```

A file is generated and stored on a s3 bucket (Do not forget to create a .env file with your s3 bucket credentials), and the url of the file is sent in the response.
Feel free to modify the project and store the file instead on your local machine.

Disclaimer: This project is not beautiful, has poor error handling and will probably make you puke.
Read the source code at your own discretion. Thanks for understanding.