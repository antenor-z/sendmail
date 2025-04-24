all:
	go build
install: all
	mkdir ~/.config/a4sendmail
	cp config.toml.template ~/.config/a4sendmail/config.toml
	cp sendmail ~/.local/bin/sendmail
