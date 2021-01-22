## bscli squad cards

diplay active cards of the given squad (Only cards updated in the last month)

### Synopsis

display the users for the squad
	

```
bscli squad cards [flags]
```

### Examples

```
bscli squad --id <squad id> cards --filterEmail my@email.com --updatedSince "2020-01-31T22:37:22-03:00" 
```

### Options

```
      --filterEmail string       filter for cards for the email
      --filterPLabel string      filter for cards for the Primary Label
      --filterSLane string       filter for cards for the SwimLane
      --filterWorkState string   filter for cards for the WorkState Type
      --updatedSince string      filter for cards for the Primary Label (default "2020-12-22T13:10:50-03:00")

```

### Options inherited from parent commands

```
  -h, --help        display this help and exit
      --id string   squad id
```

### SEE ALSO

* [bscli squad](bscli_squad.md)	 - display information for a given squad

