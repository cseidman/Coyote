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

opendb("northwind","c:/data/sqlitedata/northwind.db")

select OrderID, ProductID, (UnitPrice*Quantity)-(UnitPrice*Quantity*Discount) as ProductTotal
FROM OrderDetail limit 10;
