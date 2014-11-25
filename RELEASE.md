# Release process

Dependencies:

  - [github.com/laher/goxc](https://github.com/laher/goxc)

```shell
make bump release
git push --all origin
```

Then, go to the releases page on github, and edit the release you just made by pushing a tag. Add the contents of the `dist/` directory as individual binary attachments to the release.
