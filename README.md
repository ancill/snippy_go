### COMANDS

`ss -tnlp` - show open ports

If you don’t like the term model, you might want to think of it as a service layer or data access layer instead.

Remember: The internal directory is being used to hold ancillary non-application-specific code,

## Disabling directory listings

Warning: http.ServeFile() does not automatically sanitize the file path. If you’re constructing a file path from untrusted user input, to avoid directory traversal attacks you must sanitize the input with filepath.Clean() before using it.

If you want to disable directory listings there are a few different approaches you can take.

The simplest way? Add a blank index.html file to the specific directory that you want to disable listings for. This will then be served instead of the directory listing, and the user will get a 200 OK response with no body. If you want to do this for all directories under ./ui/static you can use the command:

`find ./ui/static -type d -exec touch {}/index.html \;`

## Database

### Install to local machine

`brew install mysql` |
`sudo apt install mysql-server`

### Starting script for this app

```sql
-- Create a new UTF-8 `snippetbox` database.
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Switch to using the `snippetbox` database.
USE snippetbox;

-- Create a `snippets` table.
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets(created);

-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

-- create a sessions table in our MySQL database to hold the session data for our users.

USE snippetbox;

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);


USE snippetbox;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

### Usefull commads for testing app

```bash
# conect to database
mysql -D snippetbox -u web -p

# select snipptes
SELECT id, title, expires FROM snippets;

# generate tls cert
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

```

## Working with transactions

It’s important to realize that calls to Exec(), Query() and QueryRow() can use any connection from the sql.DB pool. Even if you have two calls to Exec() immediately next to each other in your code, there is no guarantee that they will use the same database connection.

Sometimes this isn’t acceptable. For instance, if you lock a table with MySQL’s LOCK TABLES command you must call UNLOCK TABLES on exactly the same connection to avoid a deadlock.

To guarantee that the same connection is used you can wrap multiple statements in a transaction. Here’s the basic pattern:

```go

type ExampleModel struct {
    DB *sql.DB
}

func (m *ExampleModel) ExampleTransaction() error {
    // Calling the Begin() method on the connection pool creates a new sql.Tx
    // object, which represents the in-progress database transaction.
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }

    // Defer a call to tx.Rollback() to ensure it is always called before the
    // function returns. If the transaction succeeds it will be already be
    // committed by the time tx.Rollback() is called, making tx.Rollback() a
    // no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
    // the changes before the function returns.
    defer tx.Rollback()

    // Call Exec() on the transaction, passing in your statement and any
    // parameters. It's important to notice that tx.Exec() is called on the
    // transaction object just created, NOT the connection pool. Although we're
    // using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
    // exactly the same way.
    _, err = tx.Exec("INSERT INTO ...")
    if err != nil {
        return err
    }

    // Carry out another transaction in exactly the same way.
    _, err = tx.Exec("UPDATE ...")
    if err != nil {
        return err
    }

    // If there are no errors, the statements in the transaction can be committed
    // to the database with the tx.Commit() method.
    err = tx.Commit()
    return err
}

```

Important: You must always call either Rollback() or Commit() before your function returns. If you don’t the connection will stay open and not be returned to the connection pool. This can lead to hitting your maximum connection limit/running out of resources. The simplest way to avoid this is to use defer tx.Rollback() like we are in the example above.

---

_Transactions_ are also super-useful if you want to execute multiple SQL statements as a `single atomic action`. So long as you use the tx.Rollback() method in the event of any errors, the transaction ensures that either:

All statements are executed successfully; or
No statements are executed and the database remains unchanged.

## Verbosity

If you’re coming from Ruby, Python or PHP, the code for querying SQL databases may feel a bit verbose, especially if you’re used to dealing with an abstraction layer or ORM.

But the upside of the verbosity is that our code is non-magical; we can understand and control exactly what is going on. And with a bit of time, you’ll find that the patterns for making SQL queries become familiar and you can copy-and-paste from previous work.

If the verbosity really is starting to grate on you, you might want to consider trying the jmoiron/sqlx package. It’s well designed and provides some good extensions that make working with SQL queries quicker and easier. Another, newer, option you may want to consider is the blockloop/scan package.

## Managing null values

One thing that Go doesn’t do very well is managing NULL values in database records.

Let’s pretend that the title column in our snippets table contains a NULL value in a particular row. If we queried that row, then rows.Scan() would return an error because it can’t convert NULL into a string:

sql: Scan error on column index 1: unsupported Scan, storing driver.Value type
&lt;nil&gt; into type \*string
Very roughly, the fix for this is to change the field that you’re are scanning into from a string to a sql.NullString type. See this gist for a working example.

But, as a rule, the easiest thing to do is simply avoid NULL values altogether. Set NOT NULL constraints on all your database columns, like we have done in this book, along with sensible DEFAULT values as necessary.

## Packages

### Upgrading packages

Once a package has been downloaded and added to your go.mod file the package and version are ‘fixed’. But there are many reasons why you might want to upgrade to use a newer version of a package in the future.

To upgrade to latest available minor or patch release of a package, you can simply run go get with the -u flag like so:

$ go get -u github.com/foo/bar
Or alternatively, if you want to upgrade to a specific version then you should run the same command but with the appropriate @version suffix. For example:

$ go get -u github.com/foo/bar@v2.0.0

### Removing unused packages

Sometimes you might go get a package only to realize later that you don’t need it anymore. When this happens you’ve got two choices.

You could either run go get and postfix the package path with @none, like so:

$ go get github.com/foo/bar@none
Or if you’ve removed all references to the package in your code, you could run go mod tidy, which will automatically remove any unused packages from your go.mod and go.sum files.

$ go mod tidy -v

## Prepared statements

As I mentioned earlier, the Exec(), Query() and QueryRow() methods all use prepared statements behind the scenes to help prevent SQL injection attacks. They set up a prepared statement on the database connection, run it with the parameters provided, and then close the prepared statement.

This might feel rather inefficient because we are creating and recreating the same prepared statements every single time.

In theory, a better approach could be to make use of the DB.Prepare() method to create our own prepared statement once, and reuse that instead. This is particularly true for complex SQL statements (e.g. those which have multiple JOINS) and are repeated very often (e.g. a bulk insert of tens of thousands of records). In these instances, the cost of re-preparing statements may have a noticeable effect on run time.

Here’s the basic pattern for using your own prepared statement in a web application:

```go


// We need somewhere to store the prepared statement for the lifetime of our
// web application. A neat way is to embed in the model alongside the connection
// pool.
type ExampleModel struct {
    DB         *sql.DB
    InsertStmt *sql.Stmt
}

// Create a constructor for the model, in which we set up the prepared
// statement.
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
    // Use the Prepare method to create a new prepared statement for the
    // current connection pool. This returns a sql.Stmt object which represents
    // the prepared statement.
    insertStmt, err := db.Prepare("INSERT INTO ...")
    if err != nil {
        return nil, err
    }

    // Store it in our ExampleModel object, alongside the connection pool.
    return &ExampleModel{db, insertStmt}, nil
}

// Any methods implemented against the ExampleModel object will have access to
// the prepared statement.
func (m *ExampleModel) Insert(args...) error {
    // Notice how we call Exec directly against the prepared statement, rather
    // than against the connection pool? Prepared statements also support the
    // Query and QueryRow methods.
    _, err := m.InsertStmt.Exec(args...)

    return err
}

// In the web application's main function we will need to initialize a new
// ExampleModel struct using the constructor function.
func main() {
    db, err := sql.Open(...)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    // Create a new ExampleModel object, which includes the prepared statement.
    exampleModel, err := NewExampleModel(db)
    if err != nil {
       errorLog.Fatal(err)
    }

    // Defer a call to Close() on the prepared statement to ensure that it is
    // properly closed before our main function terminates.
    defer exampleModel.InsertStmt.Close()
}

```

There are a few things to be wary of though.

Prepared statements exist on database connections. So, because Go uses a pool of many database connections, what actually happens is that the first time a prepared statement (i.e. the sql.Stmt object) is used it gets created on a particular database connection. The sql.Stmt object then remembers which connection in the pool was used. The next time, the sql.Stmt object will attempt to use the same database connection again. If that connection is closed or in use (i.e. not idle) the statement will be re-prepared on another connection.

Under heavy load, it’s possible that a large amount of prepared statements will be created on multiple connections. This can lead to statements being prepared and re-prepared more often than you would expect — or even running into server-side limits on the number of statements (in MySQL the default maximum is 16,382 prepared statements).

The code too is more complicated than not using prepared statements.

So, there is a trade-off to be made between performance and complexity. As with anything, you should measure the actual performance benefit of implementing your own prepared statements to determine if it’s worth doing. For most cases, I would suggest that using the regular Query(), QueryRow() and Exec() methods — without preparing statements yourself — is a reasonable starting point.

### SameSite ‘Strict’ setting

If you want, you can change the session cookie to use the SameSite=Strict setting instead of (the default) SameSite=Lax. Like this:

```go
sessionManager := scs.New()
sessionManager.Cookie.SameSite = http.SameSiteStrictMode
```

But it’s important to be aware that using SameSite=Strict will block the session cookie being sent by the user’s browser for all cross-site usage — including safe requests with HTTP methods like GET and HEAD.

While that might sound even safer (and it is!) the downside is that the session cookie won’t be sent when a user clicks on a link to your application from another website. In turn, that means that your application would initially treat the user as ‘not logged in’ even if they have an active session containing their "authenticatedUserID" value.

So if your application will potentially have other websites linking to it (or even links shared in emails or private messaging services), then SameSite=Lax is generally the more appropriate setting.

SameSite cookies and TLS 1.3
Earlier in this chapter I said that we can’t solely rely on the SameSite cookie attribute to prevent CSRF attacks, because it isn’t fully supported by all browsers.

But there is an exception to this rule, due to the fact that there is no browser that exists which supports TLS 1.3 and does not support SameSite cookies.

In other words, if you were to make TLS 1.3 the minimum supported version in the TLS config for your server, then all browsers able to use your application will support SameSite cookies.

```go
tlsConfig := &tls.Config{
MinVersion: tls.VersionTLS13,
}

```

So long as you only allow HTTPS requests to your application and enforce TLS 1.3 as the minimum TLS version, you don’t need to make any additional mitigation against CSRF attacks (like using the justinas/nosurf package). Just make sure that you always:

Set SameSite=Lax or SameSite=Strict on the session cookie; and
Use an ‘unsafe’ HTTP method (i.e. POST, PUT or DELETE) for any state-changing requests.

## The request context syntax

The basic code for adding information to a request’s context looks like this:

```go
// Where r is a \*http.Request...
ctx := r.Context()
ctx = context.WithValue(ctx, "isAuthenticated", true)
r = r.WithContext(ctx)
```

Let’s step through this line-by-line.

First, we use the r.Context() method to retrieve the existing context from a request and assign it to the ctx variable.
Then we use the context.WithValue() method to create a new copy of the existing context, containing the key "isAuthenticated" and a value of true.
Then finally we use the r.WithContext() method to create a copy of the request containing our new context.

## Profiling test coverage

```bash
go test -cover ./...

// cover file
go test -coverprofile=/tmp/profile.out ./...
go tool cover -func=/tmp/profile.out
go tool cover -html=/tmp/profile.out


```

Note: If you’re running some of your tests in parallel, you should use the -covermode=atomic flag (instead of -covermode=count) to ensure an accurate count.
