## GoRM

GoRM is an ORM for Go. It lets you map Go `struct`s to tables in a database. It's intended to be very lightweight, doing very little beyond what you really want. For example, when fetching data, instead of re-inventing a query syntax, we just delegate your query to the underlying database, so you can write the "where" clause of your SQL statements directly. This allows you to have more flexibility while giving you a convenience layer. But GoRM also has some smart defaults, for those times when complex queries aren't necessary.

### How do we use it?

Open a database

	db, _ := OpenDB("test.db")
	db.Close() // you'll probably wanna close it at some point

Model a struct after a table in the db

	type Person struct {
		Id int
		Name string
		Age int
	}

Create an object and save it

	var someone Person
	someone.Name = "john"
	someone.Age = 20
	
	db.Save(&someone)

Fetch a single object

	var person1 Person
	db.Get(&person1, "id = ?", 3)

	var person2 Person
	db.Get(&person2, 3) // this is shorthand for the version above
	
	var person3 Person
	db.Get(&person3, "name = ?", "john") // more complex query
	
	var person4 Person
	db.Get(&person4, "name = ? and age < ?", "john", 88) // even more complex

Fetch multiple objects

	var bobs []Person
	err := db.GetAll(&bobs, "name = ?", "bob")

	var everyone []Person
	err := db.GetAll(&everyone, "") // use empty string to omit "where" clause

Saving new and existing objects

	person2.Name = "Jack" // an already-existing person in the database, from the example above
	db.Save(&person2)
	
	var newGuy Person
	newGuy.Name = "that new guy"
	newGuy.Age = 27
	
	db.Save(&newGuy)
	// newGuy.Id is suddenly valid, and he's in the database now.

### Installing GoRM

Obviously [Go](http://golang.org/) should be installed. [The official installation directions](http://golang.org/doc/install.html) are recommended, rather than installing it through a package (such as homebrew).

This package also requires [my version of the `sqlite` package](https://github.com/sdegutis/sqlite-go-wrapper). Clone it, then run `make install` and you'll be all set.

### Known bugs

Right now, it only interfaces with SQLite. The goal however is to add support for other databases in the future, including maybe PostgreSQL and CouchDB or NoSQL? Who knows.

Also, at the moment, relationship-support is in the works, but not yet implemented.

All in all, it's not entirely ready for advanced use yet, but it's getting there.

### Etc

The idea came about in #go-nuts on irc.freenode.net... Namegduf and wrtp were instrumental in helping solidify the main principles, and I think wrtp came up with the name.

Feel free to send pull requests with cool features added :)