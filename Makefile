bump:
	goxc bump
	echo >> .goxc.json
	git commit -m 'Bump version [nostory]' -- .goxc.json
	goxc tag

release:
	goxc
