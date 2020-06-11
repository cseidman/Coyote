opendb("northwind","../data/northwind.db")

var country = "France"

select
a.ProductID,
b.ProductName,
sum(a.Quantity) as TotalQty
FROM OrderDetail a
    join Product b on a.ProductId = b.Id
    join "Order" c on c.Id = a.OrderId
where c.ShipCountry = $country
group by a.ProductID,
         b.ProductName;
