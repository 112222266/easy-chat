.PHONY: push

push:
	git add .
	git commit -m "$${MSG:-$(shell date +'%Y-%m-%d %H:%M:%S')}"
	git push origin main