/*create database GoGreen;*/
use GoGreen;
create table `Users` (
	`ID` int(11) not null auto_increment, 
    `Username` varchar(30), 
    `Password` varchar(20), 
    `Name` varchar(50) not null, 
    `Role` varchar(30) not null, 
    `Email` varchar(100) not null, 
    `Address` varchar(100), 
    `Contact` int(8),
    Primary Key (`ID`)
);


CREATE TABLE `Categories` (
  `ID` int(11) NOT NULL auto_increment,
  `Name` varchar(200) NOT NULL,
  `Description` varchar(500),
  `Number_Of_Products` int(6),
  Primary Key (`ID`)
) DEFAULT CHARSET=latin1;

CREATE TABLE `Brands` (
  `ID` int(11) NOT NULL auto_increment,
  `Name` varchar(200) NOT NULL,
  `Description` varchar(500),
  `Number_Of_Products` int(6),
  Primary Key (`ID`)
) DEFAULT CHARSET=latin1;

CREATE TABLE `Products` (
  `ID` int(11) NOT NULL auto_increment,
  `Name` varchar(200) NOT NULL,
  `Image` BLOB, 
  `Details` varchar(500) NOT NULL,
  `Date_Added` date DEFAULT '0000-00-00',
  `Price` decimal(15,2) NOT NULL,
  `Quantity` int(11) NOT NULL,
  `Category_ID` int(11) NOT NULL,
  `Brand_ID` int(11) NOT NULL,
  Primary Key (`ID`),
  Foreign Key(`Category_ID`) references Categories(`ID`),
  Foreign Key(`Brand_ID`) references Brands(`ID`)
) DEFAULT CHARSET=latin1;

CREATE TABLE `Customer_Orders` (
  `ID` int(11) NOT NULL auto_increment,
  `Status` varchar(200) NOT NULL,
  `User_ID` int(11) NOT NULL,
  `Amount` decimal(15,2) NOT NULL,
  `Order_Date` date DEFAULT '0000-00-00',
  `Delivery_Date` date DEFAULT '0000-00-00',
  Primary Key (`ID`),
  Foreign Key(`User_ID`) references Users(`ID`)
) DEFAULT CHARSET=latin1;

CREATE TABLE `Product_Orders` (
  `ID` int(11) NOT NULL auto_increment,
  `Status` varchar(200) NOT NULL,
  `Quantity` int(11),
  `Product_ID` int(11) NOT NULL,
  `Customer_Order_ID` int(11) NOT NULL,
  Primary Key (`ID`),
  Foreign Key(`Product_ID`) references Products(`ID`),
  Foreign Key(`Customer_Order_ID`) references Customer_Orders(`ID`)
) DEFAULT CHARSET=latin1;

