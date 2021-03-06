## bscli card update

update an existing card

### Synopsis

update an existing card
	

```
bscli card update [flags]
```

### Examples

```
bscli card update --card "09lfdk" -s "Default Swimlane" \
-w "Backlog" \
-t "My new card" \
-d "My new card<br>new line>" \
-p "Primary label1" -p "primary label2" \
-a "assingee@email.com" -a "assingee2@email.com" \
--dueDate "01/31/2021 15:00:00"

```

### Options

```
  -a, --assignees strings      card assignee emails
      --card string            card identifier
  -d, --description string     card description
      --dueDate string         card due date
  -p, --primarylabel strings   card primary label names
  -s, --swimlane string        swimlane name
  -t, --title string           card title
  -w, --workstate string       workstate name
```

### Options inherited from parent commands

```
  -h, --help   display this help and exit
```

### SEE ALSO

* [bscli card](bscli_card.md)	 - create, update or create comment for a given card

