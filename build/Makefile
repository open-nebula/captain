IMAGE := nebula-captain
URL := wss://88078771.ngrok.io/join

# Dynamic Port Exposure
# https://stackoverflow.com/questions/19897743/exposing-a-port-on-a-live-docker-container

# Or using host network (doesn't work for Mac because host is VM)

build:
	cd .. && docker build -t $(IMAGE) -f build/Dockerfile .

run: build
	docker run -it $(IMAGE) $(URL)
