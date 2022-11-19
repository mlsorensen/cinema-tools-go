PLATFORMS := linux/amd64/lin64 linux/arm64/linarm windows/amd64/win64.exe darwin/arm64/macM1

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
alias = $(word 3, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o build/lumagen-monitor-$(alias) cmd/lumagen-monitor.go
	GOOS=$(os) GOARCH=$(arch) go build -o build/urtsi2-cmd-$(alias) cmd/urtsi2-cmd.go
	GOOS=$(os) GOARCH=$(arch) go build -o build/cinemacontrol-$(alias) cmd/cinemacontrol.go

clean:
	rm -rf ./build