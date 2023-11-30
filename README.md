# HO

## Disabling directory listings

Warning: http.ServeFile() does not automatically sanitize the file path. If you’re constructing a file path from untrusted user input, to avoid directory traversal attacks you must sanitize the input with filepath.Clean() before using it.

If you want to disable directory listings there are a few different approaches you can take.

The simplest way? Add a blank index.html file to the specific directory that you want to disable listings for. This will then be served instead of the directory listing, and the user will get a 200 OK response with no body. If you want to do this for all directories under ./ui/static you can use the command:

`find ./ui/static -type d -exec touch {}/index.html \;`
