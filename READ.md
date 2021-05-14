Project GoGreen
Product Management, Shopping Cart and Helpdesk/Enquiry
by Wong Joey and Joe Yen Yen
for Project Go Live 2021

Github Link: https://github.com/yy-joe/GoGreen/tree/plsbcool-patch-3

Running the implemented system
The following are the steps to run the system.
1.	Create a MySQL database with the SQL script GoGreenDB.sql in the database folder. The database server port number and password for the root user are set in the .env file in the root folder with the environment variables name DBPORT and DBPASSWORD, respectively.

2.	Run the main executable GoGreen.exe. The HTTP/TLS server at localhost will go live and listen to port 3000, i.e. the server address is https://localhost:3000.

3.	The URL of the main Products page for the admin (i.e. seller) is https://localhost:3000/products/all, while the URL of the main Products page for the non-admin user (i.e. shopper) is  https://localhost:3000/ .

Using the implemented system
Once the server is running, it can be tested by performing all the functions provided in the application website, with various possible inputs, either valid or invalid. Error is mainly printed to the terminal of both the server and the client. We also prepare sample data that can be used to populate the database by running the SQL script testdata.sql in the database folder.

Appendix: List of files in the application folder
Folder	Sub-folder	File	Description
GoGreen	--	.env	Environment file to store DB password and port number
	--	GoGreen.exe	Executable file for the application
	--	cert.pem	TLS certificate of the server
	--	key.pem	TLS key of the server
	--	server.go	Source file for HTTP/TLS server 
	database	GoGreenDB.sql	Database schema
		testdata.sql	Sample data to populate database
	products	client.go	Source file for client service application, mainly for admin UI
		clientUser.go	Source file for client service application, mainly for non-admin UI
		mergesort.go	Source file for merge sort algorithm to sort products by their attributes
		prodserver.go	Source file for API endpoint functions
		queries.go	Source file for DB query functions
		search.go	Source file for product search function 
		sendMail.go	Source file for mail sending function for user enquiry
	products/template	*.gohtml	Go HTML templates


