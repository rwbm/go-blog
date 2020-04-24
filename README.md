# Go Blog backend

This service "listen" for new blog posts templates files to be copied into a specific folder, parse them and save the details in a local SQLite database. 

There's also end endpoint to retrieve the blog posts, and filter by author, date, tags, etc.


## Compilation

In order to build the project you need to have Go 1.14 (or newer) installed. Just clone the project anywhere you wish and compile the project running `build.sh` script (on Linux or Mac; for Windows, create a similar script).

## Running

To run the service from the folder with the source code, you can use the script `run.sh`. That script compiles the project and start the service. Service binary is located in `./cmd/backend`.

Before running the service, you need to create the configuration file. There's a template named `config.local.yml` that can be used as a template to create a new one. Copy the file to a new one named `config.yml`. By default this file must be located in the same location where the binary is.

## Configuration file

Configuration files are located in `./cmd/backend`. Copy a new file from the template and then edit it to change the default values. 

**Configuration fields**

| Field                    | Description  |
|--------------------------|-----------------|
| server.name              | Server name, just for reference |
| server.port              | Port number where the HTTP server is going to serve |
| server.read_timeout      | HTTP reat timeout |
| server.write_timeout     | HTTP write timeout |
| database.filename        | db filename; placeholder `$APP_HOME` may be used to refer to the application location |
| template.base_location   | location where blog templates are stored; placeholder `$APP_HOME` may be used |
| template.processed_ok    | location where blog templates are stored after correctly processed; placeholder `$APP_HOME` may be used |
| template.processed_error | location where blog templates are stored after processed with errors; placeholder `$APP_HOME` may be used |
| template.check_cycle     | How many seconds to wait before checking for new templates in `template.base_location` |

If not defined, the service will assume some default values:

- server.port: `8080`
- server.read_timeout: `5 seconds`
- server.write_timeout: `2 seconds`
- database.filename = `$APP_HOME/blog.db`
- template.base_location = `$APP_HOME/templates`
- template.processed_ok = `$APP_HOME/templates/ok`
- template.processed_error = `$APP_HOME/templates/error`
- template.check_cycle = `30 seconds`

## Database

Database file is generated if not found in the location defined in the configuration setting _database.filename_. Before moving the application or makeing any change in the database, please consider making a backup.

The schema defined for the database is the following:

### post

Contains the main post info.

```[sql]
CREATE TABLE post (
    id_post           INTEGER       PRIMARY KEY AUTOINCREMENT,
    date_created      DATETIME      NOT NULL,
    date_updated      DATETIME      NOT NULL,
    title             VARCHAR (128) NOT NULL,
    author            VARCHAR (128) NOT NULL,
    content           TEXT          NOT NULL,
    original_filename VARCHAR (128) NOT NULL
);
```

### post_category

Lits of categories definedd for a post.

```[sql]
CREATE TABLE post_category (
    id_post INTEGER       NOT NULL,
    name    VARCHAR (128) NOT NULL
);
```

### post_tag

Lits of tags definedd for a post.

```[sql]
CREATE TABLE post_tag (
    id_post INTEGER       NOT NULL,
    name    VARCHAR (128) NOT NULL
);
```

Data strcuture is defined at the model base, in $PROJECT/pkg/util/model/post.go. GORM is used as the ORM to handle DB, so you can make the required changes here, move the old DB and start the service again. If you just want to make some minor change, like increase a field length, just make that change to the current DB with an external DB tool so you can keep the data.

## Get posts endpoint

Endpoint to the stored endpoints can be accesed by calling `/posts`; for example:

`curl http://127.0.0.1:8088/posts`

Additional filter parameters may be sent in order to get only the required results. Those parameter can be sent as query paramrters. The following list contains the list of valid parameters:

| Parameter   | Description                                                         |
|-------------|---------------------------------------------------------------------|
| id_post     | Internal ID of the post, extracted from the database                |
| author      | Name of the author of the post                                      |
| date-from   | Start creation date to filter, format `YYYY-MM-dd`                  |
| date-to     | End creation date to filter, format `YYYY-MM-dd`                    |
| categories  | Comma separated values with the list of categories to filter        |
| tags        | Comma separated values with the list of tags to filter              |
| page        | Indicates the page number, default value is `1`                     |
| page-size   | Indicates the max number of rows to retrieve; default value is `25` |

Example:

*Request:*

`curl http://127.0.0.1:8080/posts\?date-from\=2020-04-01\&date-to\=2020-04-20`

*Response:*

```
{
   "pagination":{
      "page":1,
      "page_size":25
   },
   "posts":[
      {
         "id_post":1,
         "date_created":"2020-04-15T12:09:57Z",
         "date_updated":"2020-04-15T12:09:57Z",
         "title":"My First Blog Post",
         "author":"John Doe",
         "content":"PGJvZHk+CiAgICA8aDE+TXkgRmlyc3QgQmxvZyBQb3N0PC9oMT4KICAgIDxwPlRoaXMgaXMgbXkgZmlyc3QgQmxvZyBwb3N0LCBqdXN0IHRvIHRyeSBpZiB0ZW1wbGF0ZXMgYXJlIHdvcmtpbmcgT0suPC9wPgoKPC9ib2R5Pg==",
         "categories":"Go Programming",
         "tags":"go,programming,web"
      },
      {
         "id_post":2,
         "date_created":"2020-04-20T20:35:17.508954-03:00",
         "date_updated":"2020-04-20T20:35:17.508954-03:00",
         "title":"My Second Blog Post",
         "author":"John Doe",
         "content":"PGJvZHk+CiAgICA8aDE+RGV2T3BzIEJsb2cgUG9zdDwvaDE+CiAgICA8cD5UaGlzIGlzIG15IHNlY29uZCBCbG9nIHBvc3QuIEhlcmUgd2UgZGVhbCB3aXRoIERldk9wcyByZWxhdGVkIHRvcGljcy48L3A+Cgo8L2JvZHk+",
         "categories":"AWS",
         "tags":"aws,devops"
      }
   ]
}
```
