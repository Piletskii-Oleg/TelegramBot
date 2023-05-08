package storage

const msgHelp = `I can save and keep your pages, which you can retrieve later.
To save the page, send me the link
To get a random page from my list, send '/rnd'
Caution! The link will be deleted after that command`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedPages   = "You have no saved pages"
	msgSaved          = "Saved!"
	msgAlreadyExists  = "This page already exists"
)
