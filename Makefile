.PHONY: deploy backend provision backend-arm destroy

all: backend provision clean

backend:
	$(MAKE) -C src/backend/scanner amd
	$(MAKE) -C src/backend/api amd

backend-arm:
	$(MAKE) -C src/backend/scanner arm
	$(MAKE) -C src/backend/api arm

provision:
	$(MAKE) -C provision/ deploy

clean:
	$(MAKE) -C src/backend/scanner clean
	$(MAKE) -C src/backend/api clean

destroy:
	$(MAKE) -C provision/ destroy
	$(MAKE) -C src/backend/scanner clean
	$(MAKE) -C src/backend/api/clean
