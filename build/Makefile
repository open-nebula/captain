IMAGE := nebula-captain
REPO := codyperakslis/captain
URL := wss://412db34f.ngrok.io/join

# Dynamic Port Exposure
# https://stackoverflow.com/questions/19897743/exposing-a-port-on-a-live-docker-container

# Or using host network (doesn't work for Mac because host is VM)

build:
	cd .. && docker build -t $(IMAGE) -f build/Dockerfile .

run: build
	docker run -it -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE) $(URL) 
