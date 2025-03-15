# HTTP-GO

## Parsing request body pattern

- Note, you are likely going to create a RequestStruct to ensure you only take in the right data.

```go
type parameters struct {
  Body string `json:"body"`
}
decoder := json.NewDecoder(r.Body)
params := parameters{}
err := decoder.Decode(&params)

if err != nil {
  log.Printf("Error decoding parameters: %s", err)
  w.WriteHeader(500)
  return
}
```

## Building response body pattern

- Note, you are likely going to create ResponseStructs or something similar to a DTO in NestJS.
- Ensures that your response body is structured with the expected data.s

```go
type responseVal struct {
  Valid string `json:"valid"`
}

respBody := responseVal{
  Valid: "true",
}

dat, err := json.Marshal(respBody)
if err != nil {
    log.Printf("Error marshalling JSON: %s", err)
    w.WriteHeader(500)
    return
}
w.WriteHeader(200)
w.Write(dat)
```
