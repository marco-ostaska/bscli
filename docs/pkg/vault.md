# vault
--
    import "github.com/marco-ostaska/bscli/cmd/vault"

Package vault is mainly a reference to cobra command vault

But it has the essentials to vault adminstration to be used throughout the
application.

## Usage

```go
const (
	Dir    = "bscli"               // Vault user dir
	File   = "bscli.vlt"           // vault usr file
	APIKey = "Bluesight-API-Token" // default bluesight token key
)
```
vault basic constants

```go
var Cmd = &cobra.Command{
	Use:   "vault",
	Short: "create or update vault credentials",
	Long:  `create or update vault credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v\n", cmd.Long)
		if err := cmd.Usage(); err != nil {
			log.Fatalln(err)
		}
	},
}
```
Cmd represents the vault command

```go
var Credential uvault.Credential
```
Credential is a reference to uvault.Credential

```go
var HTTP httpcalls.APIData
```
HTTP is a reference to httpcalls.APIData with uservault uploaded

#### func  ReadVault

```go
func ReadVault()
```
ReadVault reads the user vault contents
