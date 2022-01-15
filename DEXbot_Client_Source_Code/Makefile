
dexbot-linux: 

	#Build the dexbot-linux image
	docker build -t dexbot-linux-image --build-arg BUILDOS=linux --build-arg BUILDARCH=amd64 --build-arg BUILDNAME=dexbot-linux --no-cache .

	#Run the dexbot-linux-image in a new container called dexbot-linux
	docker run --name dexbot-linux  dexbot-linux-image

	#Copy the compiled dexbot binary to the host machine
	docker cp dexbot-linux:/app/dexbot-linux .
	
	#Remove the dexbot-linux container
	docker rm dexbot-linux

	#Remove the dexbot-linux-image
	docker image rm dexbot-linux-image



dexbot-windows: 

	#Build the dexbot-windows image
	docker build -t dexbot-windows-image --build-arg BUILDOS=windows --build-arg BUILDARCH=amd64 --build-arg BUILDNAME=dexbot-windows.exe --no-cache .

	#Run the dexbot-windows-image image in a new container called dexbot-windows
	docker run --name dexbot-windows  dexbot-windows-image

	#Copy the compiled dexbot binary to the host machine
	docker cp dexbot-windows:/app/dexbot-windows.exe .
	
	#Remove the dexbot-windows container
	docker rm dexbot-windows

	#Remove the dexbot-windows-image
	docker image rm dexbot-windows-image



dexbot-mac: 

	#Build the dexbot-mac image
	docker build -t dexbot-mac-image --build-arg BUILDOS=darwin --build-arg BUILDARCH=amd64 --build-arg BUILDNAME=dexbot-mac --no-cache .

	#Run the dexbot-mac-image in a new container called dexbot-mac
	docker run --name dexbot-mac  dexbot-mac-image


	#Copy the compiled dexbot binary to the host machine
	docker cp dexbot-mac:/app/dexbot-mac .
	
	#Remove the dexbot-mac container
	docker rm dexbot-mac

	#Remove the dexbot-mac-image
	docker image rm dexbot-mac-image

	

