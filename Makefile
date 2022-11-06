deploy:
	$(MAKE) -C src/backend/scanner all
	$(MAKE) -C src/backend/api all
	$(MAKE) -C provision/ deploy

backend:
	$(MAKE) -C src/backend all

destroy:
	$(MAKE) -C provision/ destroy
	$(MAKE) -C src/backend clean
