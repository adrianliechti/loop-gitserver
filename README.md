# Simple Git Server

## Usage

```bash
docker build . --tag loop-gitserver
docker run -it --rm -p 8080:8080 loop-gitserver
```

```bash
# create repository
curl -X POST -v http://localhost:8080/admin/repo/demo

# clone repository
git clone http://localhost:8080/demo

# add some content
echo Awesome Project > README.md

git add -A
git commit -m 'starting something new'
git push
```