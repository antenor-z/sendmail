all:
	go build
install: all
	mkdir -p ~/.config/a4sendmail
	cp config.toml.template ~/.config/a4sendmail/config.toml
	mkdir -p ~/.local/bin
	cp sendmail ~/.local/bin/sendmail
