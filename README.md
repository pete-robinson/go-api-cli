# Fun evening CLI project
CLI to fetch data on "todos" from a given API in a number of ways. Utilising some basic concurrency via worker pools to speed up multiple API calls.

Tool is centred around 3 core commands with additional arguments.

## Docker Build
```
docker build . -t todo-cli
```

## Commands
### fetchall
Fetches entire list of results. Supports result limiting via `-count` argument.

#### Local
```
go run main.go -command=fetchall -count=10
```

#### Docker
```
docker run cli -command=fetchall -count=10
```

### fetchone
Fetches one result - requires `-id` argument.
#### Local
```
go run main.go -command=fetchone -id=20
```

#### Docker
```
docker run cli -command=fetchone -id=20
```

### fetchlist
Accepts an absolute file path (via `-list` argument). Uses worker pools to fetch results concurrently. `data/sample.csv` is the sample data.

Can be optimised a little further but is initially taking 2s to make 100 separate calls synchronously compared to ~500ms when fetched async using worker pools.
#### Local
```
go run main.go -command=fetchlist -list=/path/to/file.csv
```
#### Docker
```
docker run cli -command=fetchlist -list=/build/data/sample.csv
```

## Result format
Simple logging output. Example:
```
INFO[0000] Response received    completed=true id=22 title="distinctio vitae autem nihil ut molestias quo" userId=2
INFO[0000] Response received    completed=false id=53 title="qui labore est occaecati recusandae aliquid quam" userId=3
INFO[0000] Response received    completed=false id=31 title="repudiandae totam in est sint facere fuga" userId=2
INFO[0000] Response received    completed=false id=68 title="aut id perspiciatis voluptatem iusto" userId=4
```

## Error format
```
ERRO[0000] No results found for id: 2323    command=fetchone
```

## Architecture
A simplistic take on **Clean** without `infra` layer (this was an evening's efforts...).

```
/
│
└───cmd/cli/
│      main.go                          *main entry point*
│   
└───internal/
│   │   
│   └───domain/                         *domain structs*
│   │   │
│   │   └───model/
│   │   │       todo.go                 
│   │   │        
│   │   └───repository/
│   │           todo_repository.go      *repository structs w/ infra*
│   │          
│   └───handlers/cli/                   *handlers for entry point adapters*
│   │       cli.go                      
│   │     
│   └───services/                       *services distributing domain logic*
│           reader.go                   *reader - really this isn't a service and should be relocated*
│           todo.go                     *main todo service*
│   │         
│   └───utils/                       
│           request.go                  *simple HTTP request manager - no auth etc*
│   
│   Dockerfile
│   go.mod
│   go.sum
│   README.md
```