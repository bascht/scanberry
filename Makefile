##
# Scan Builder
#
# @file
# @version 0.1

himbeer:
	GOOS=linux GOARCH=arm GOARM=7 go build -o himbeerscan
	scp himbeerscan himbeerkompott:~/bin/

# end
