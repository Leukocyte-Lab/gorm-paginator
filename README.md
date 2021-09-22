# gorm-paginator

## Usage

```bash
go get github.com/Leukocyte-Lab/gorm-paginator
```

```go
pgntr := paginator.New(page, orders, map[string]string{})

var models []*Model
result := pgntr.GenGormTransaction(dao.DB).Find(&models)
if result.Error != nil {
    // TODO: handle error
}

err := pgntr.CountPageTotal(result)
if result.Error != nil {
    // TODO: handle error
}

// return models, pgntr.Page
```

## License

[MIT](LICENSE)