# 将文件转为base64

> 在windows平台可以双击main.exe执行文件 自动生成对应的base64文件(.b64.txt)

```bash
$ ./main.exe # 生成当前目录中的(ttf,jpg,png,ico,txt)文件生成对应的base64文件(.b64.txt)

$ ./main.exe -R # 递归生成当前目录中的(ttf,jpg,png,ico和txt)文件生成对应的base64文件(.b64.txt)

$ ./main.exe -R --accept=jpg,png # 递归生成目录中的(jpg和png)文件生成对应的base64文件(.b64.txt)

$ ./main.exe -P # 生成当前目录中的(ttf,jpg,png,ico,txt)文件生成对应的base64文件(.b64.txt)并添加web前端需要的data:image/png;base64,内容

$ ./main.exe -P /home/dev/logo.png # 将logo.png生成对应的base64文件(.b64.txt)并添加web前端需要的data:image/png;base64,内容

$ ./main.exe -P -R c:/home/dev/pic # 递归生成指定目录中的(ttf,jpg,png,ico,txt)文件生成对应的base64文件(.b64.txt)并添加web前端需要的data:image/png;base64,内容
```