/*create database GoGreen;*/
create database GoGreen;
use GoGreen;
create table `Users` (
	`ID` int not null auto_increment, 
    `Username` varchar(30) not null, 
    `Password` varchar(150) not null, 
    `Name` varchar(50) not null, 
    `Role` varchar(30) not null, 
    `Email` varchar(100) not null, 
    `Address` varchar(100) not null, 
    `Contact` varchar(8) not null,
    `Date_Joined` date not null,
    Primary Key (`ID`), 
    Unique Key (`Username`)
);

CREATE TABLE `Categories` (
  `ID` int NOT NULL auto_increment,
  `Name` varchar(200) NOT NULL,
  `Description` varchar(500),
  Primary Key (`ID`),
  Unique Key (`Name`)
);

CREATE TABLE `Brands` (
  `ID` int NOT NULL auto_increment,
  `Name` varchar(200) NOT NULL,
  `Description` varchar(500),
  Primary Key (`ID`),
  Unique Key (`Name`)
);

CREATE TABLE `Products` (
  `ID` int NOT NULL auto_increment,
  `Name` varchar(200) NOT NULL,
  `Image` varchar(100), /*just the path of the image location on the server*/
  `Desc_Short` varchar(500) NOT NULL,
  `Desc_Long` varchar(1000),
  `Date_Created` date not null,
  `Date_Modified` date not null,
  `Price` decimal(10,2) NOT NULL,
  `Quantity` int NOT NULL,
  `Quantity_Sold` int NOT NULL,
  `Condition` varchar(10) NOT NULL,
  `Category_ID` int NOT NULL,
  `Brand_ID` int NOT NULL,
  `Status` varchar(100) NOT NULL, /*Status can be: live, sold out, or discontinued*/
  Primary Key (`ID`),
  Foreign Key(`Category_ID`) references Categories(`ID`) on update cascade on delete cascade,
  Foreign Key(`Brand_ID`) references Brands(`ID`) on update cascade on delete cascade,
  Unique Key (`Name`)
);

CREATE TABLE `Customer_Orders` (
  `ID` int NOT NULL auto_increment,
  `Status` varchar(200) NOT NULL,
  `User_ID` int NOT NULL,
  `Amount` decimal(10,2) NOT NULL,
  `Order_Date` date not null,
  `Delivery_Date` date not null,
  Primary Key (`ID`),
  Foreign Key(`User_ID`) references Users(`ID`) on update cascade on delete cascade
);

CREATE TABLE `Product_Orders` (
  `ID` int NOT NULL auto_increment,
  `Status` varchar(200) NOT NULL,
  `Quantity` int not null,
  `Product_ID` int NOT NULL,
  `Customer_Order_ID` int NOT NULL,
  Primary Key (`ID`),
  Foreign Key(`Product_ID`) references Products(`ID`) on update cascade on delete cascade,
  Foreign Key(`Customer_Order_ID`) references Customer_Orders(`ID`) on update cascade on delete cascade
);