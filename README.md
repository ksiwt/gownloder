# Gownloder
Basically, it download file in single thread.  
If a URL supports http header - Accept-Ranges, it will be chunked and download it concurrently.

# Usage
```
./gownloader -u [URL] -d [directory]
```
- `-u` flag mean URL of download file.
- `-d` flag mean destination of save downloaded file.