deploy:
	$(MAKE) -C src all
	$(MAKE) -C provision deploy
