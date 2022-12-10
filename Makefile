.PHONY: deploy backend provision backend-arm destroy

all: backend provision clean

backend:
	$(MAKE) -C backend/ amd

backend-arm:
	$(MAKE) -C backend/ arm

provision:
	$(MAKE) -C provision/ deploy

clean:
	$(MAKE) -C backend/ clean

destroy:
	$(MAKE) -C provision/ destroy
	$(MAKE) -C backend/ clean
