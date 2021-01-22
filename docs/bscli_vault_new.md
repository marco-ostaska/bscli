## bscli vault new

create new vault.

### Synopsis

create new vault.
For more information how to create a token please 
check the API section https://portal.bluesight.io/tutorial.html 


```
bscli vault new [flags]
```

### Examples

```

  Unix Based OS: (use single quotes)
      bscli vault new -k '<token>' --url 'https://www.bluesight.io/graphql'
  Windows: (use double quotes)
      bscli vault new -k "<token>" --url "https://www.bluesight.io/graphql"

```

### Options

```
  -k, --key string   API key value
      --url string   API URI
```

### Options inherited from parent commands

```
  -h, --help   display this help and exit
```

### SEE ALSO

* [bscli vault](bscli_vault.md)	 - create or update vault credentials

