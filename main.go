package main

import (
	// "github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"fmt"
)

type id_t = uint64

type Product struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
	Price uint64 `json:"price"`
}

type Alias struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
	ProductID id_t `json:"product_id"`
}

type Category struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
}

type Sale struct {
	ID id_t `json:"id"`
	ProductID id_t `json:"product_id"`
	Quantity int `json:"quantity"`
	Price uint64 `json:"price"`
}

type Ticket struct {
	ID id_t `json:"id"`
	SaleID []id_t `json:"sale_id"`
	Total uint64 `json:"total"`
}

type Database struct {
	db      *sql.DB
	parser  *SqlParser
}

func (d *Database) Execute(query string, params ...any) (r sql.Result, err error) {
	if d.parser == nil{
		err = fmt.Errorf("No parser found");
		return
	}

	q, ok := d.parser.formats[query]
	if !ok {
		err = fmt.Errorf("Query not found");
		return
	}
	r, err = d.db.Exec(q, params...);
	if err != nil {
		err = fmt.Errorf("could not query [%s] error: [%v]", query, err);
		return
	}
	return
}

func (d *Database)NewProduct(_name string, _price uint64) (r sql.Result, err error) {
	r, err = d.Execute("ADD-PRODUCT", _name, _price);
	return
}

func (d *Database)newAlias(_name string, _product_id id_t) (r sql.Result, err error) {
	r, err = d.Execute("ADD-ALIAS", _name, _product_id);
	return
}

func (d *Database)newCategory(_name string) (r sql.Result, err error) {
	r, err = d.Execute("ADD-CATEGORY", _name);
	return
}

func (d *Database)newSale(_product_id id_t, _quantity int, _price uint64) (r sql.Result, err error) {
	r, err = d.Execute("ADD-SALE", _product_id, _quantity, _price);
	return
}

func (d *Database)newTicket(_sale_id []id_t, _total uint64) (r sql.Result, err error) {
	r, err = d.Execute("ADD-TICKET", _sale_id, _total);
	return
}

func initDatabase(p *SqlParser) (database *Database, err error){
	db, err := sql.Open("sqlite3", "./database.db");

	database = &Database{
		parser: p,
	};
	if err != nil {
		panic(err);
	}

	database.db = db;

	// Create tables
	tables := []string {
		"PRODUCT",
		"ALIAS",
		"CATEGORY",
		"SALE",
		"TICKET",
	};
	for _, table := range tables {
		_, err = database.Execute("CREATE-" + table);
		if err != nil {
			fmt.Printf("COULD NOT CREATE TABLE: %s %v\n", table, err);
			return
		}
	}

	return
}

func seedDatabase(database *Database) (err error) {
	path := "./.ignore/seed.csv";
	data, err := os.ReadFile(path);
	if err != nil {
		err = fmt.Errorf("could not open file %s", path);
		return
	}

	data_str := string(data);
	lines := strings.Split(data_str, "\n");

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, ";");
		if len(parts) == 0 {
			err = fmt.Errorf("could not parse line %s", line);
			return
		}

		name := parts[0];
		price, err := strconv.ParseUint(parts[1], 10, 64);
		if err != nil {
			err = fmt.Errorf("could not parse price %s", parts[1]);
			return err
		}

		_, err = database.NewProduct(name, price);
	}
	return
}

func (p *Database) List() (r sql.Result, e error) {
	r, e = p.Execute("LIST-PRODUCTS");
	return
}

func main(){
	parser := newSqlParser();
	files := []string {
		"./sql/add.sql",
		"./sql/create_table.sql",
		"./sql/list.sql",
		"./sql/delete.sql",
	};
	for _, file := range files {
		err := parser.AddFromFile(file);
		if err != nil {
			fmt.Println(err);
			return
		}
	}

	database, err := initDatabase(parser);
	if err != nil {
		return
	}

	re, err := database.List();

	if err != nil {
		fmt.Println(err);
		return
	}

	fmt.Println(re);


	// e := echo.New()

	// e.Use(middleware.Logger())

	// e.GET("/api/v1/product/list", func (c echo.Context) error {
	// 	return c.String(200, "List of products")
	// });
}

