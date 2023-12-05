# HO

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
```

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
