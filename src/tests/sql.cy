create table Person
(
   name string,
   country string,
   age int
)

insert into Person (name, country, age) values ("Bob","USA",20)
insert into Person (name, country, age) values ("Bill","France",30)
insert into Person (name, country, age) values ("Mary","SPain",40)

select
    *
from Person
where age = 30


