deploy:
	$(MAKE) -C src/backend all
	$(MAKE) -C provision/ deploy

backend:
	$(MAKE) -C src/backend all

destroy:
	$(MAKE) -C src/backend clean
	$(MAKE) -C provision/ destroy
