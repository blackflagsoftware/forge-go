package menu

import (
	"fmt"
	"os"

	m "github.com/blackflagsoftware/forge-go/internal/model"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func SetupMenu() (pwd string, err error) {
	pwd, err = os.Getwd()
	if err != nil {
		err = fmt.Errorf("Unable to get present working directory: %s", err)
		return
	}
	setupHeader()
	fmt.Printf("Your current directory is: \n\n%s\n\n", pwd)
	msg := fmt.Sprint("Is this the root folder of the project")
	useDir := util.AskYesOrNo(msg)
	if !useDir {
		err = fmt.Errorf("Current diretory, NOT chosen, good bye")
		return
	}
	return
}

func setupHeader() {
	util.ClearScreen()
	fmt.Printf("*** Welcome to forge ***\n\n")
	fmt.Printf("*forge* needs to initialize with a few questions\n\n")
}

// TODO: can these two functions, be combined with the ones in project.go... yep, send int a setupHeader as a argument for different printing headers, there will be other small
// differences but I think it is possible
func StorageInitMenu(p *m.ProjectFile) {
	setupHeader()
	mainMesssge := []string{"Storage Type", "This option will be saved to the project (this can be changed later)"}
	prompts := []string{"(s) SQL", "(f) File", "(m) MongoDB"}
	acceptablePrompts := []string{"s", "f", "m"}
	p.Storage = util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "", setupHeader)
	if p.Storage == "s" {
		mainMesssge = []string{"SQL Option", "Choice which SQL implementation"}
		prompts = []string{"(p) Postgres", "(m) Mysql", "(s) Sqlite"}
		acceptablePrompts = []string{"p", "m", "s"} // haha... get it?!??!!  pms... haha
		p.SqlStorage = util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "", setupHeader)
		// p.UseORM = util.AskYesOrNo("Would you like to use an ORM") only have this as an option in an "admin" screen
	}
}

func TagFormatInitMenu(p *m.ProjectFile) {
	setupHeader()
	mainMesssge := []string{"Tag Format", "What format you want your 'json' tags to be (this can be changed later)?"}
	prompts := []string{"(s) Snake Case (tag_format)", "(c) Camel Case (tagFormat)", "(p) Pascal Case (TagFormat)", "(k) Kebab Case (tag-format)", "(l) Lower Case (tag format)", "(u) Upper (TAG FORMAT)"}
	acceptablePrompts := []string{"s", "c", "p", "k", "l", "u"}
	tagFormat := util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "", setupHeader)
	switch tagFormat {
	case "s":
		p.TagFormat = "snakeCase"
	case "k":
		p.TagFormat = "kebabCase"
	case "c":
		p.TagFormat = "camelCase"
	case "p":
		p.TagFormat = "pascalCase"
	case "u":
		p.TagFormat = "upperCase"
	case "l":
		p.TagFormat = "lowerCase"
	}
}
