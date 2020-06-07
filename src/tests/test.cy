create table Person (
    first_name string ,
    last_name string ,
    age int not null
);

var ageNum = 27

insert into Person (first_name, last_name, age) values ('John','Smith',$ageNum);
insert into Person (first_name, last_name, age) values ('Mary','Jones',42);
insert into Person (first_name, last_name, age) values ('George','Carlin',66);

//var df table
var df = select first_name, last_name, age from Person ;
showdata(df)

create table Person2 as select * from Person where age <50;
select * from Person2;

opendb("northwind","../data/northwind.db")

//
select
a.OrderID, a.ProductID, b.ProductName, b.JJ, (a.UnitPrice*a.Quantity)-(a.UnitPrice*a.Quantity*a.Discount) as ProductTotal
FROM OrderDetail a join Product b on a.ProductId = b.Id
limit 10;
